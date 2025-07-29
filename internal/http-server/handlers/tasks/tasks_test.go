package taskshandler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	taskshandler "github.com/Swiffy12/taskify/internal/http-server/handlers/tasks"
	"github.com/Swiffy12/taskify/internal/http-server/handlers/tasks/mocks"
	"github.com/Swiffy12/taskify/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTask(t *testing.T) {
	//Создаем моковые зависимости handlera
	taskServiceMock := mocks.NewMockTaskService(t)
	logger := slogdiscard.NewDiscardLogger()

	handler := taskshandler.New(logger, nil, taskServiceMock)

	cases := []struct {
		name           string
		reqBody        string
		setupMock      func()
		expectedStatus int
		expectedResp   func(t *testing.T, resp map[string]any)
	}{
		{
			name:    "Success",
			reqBody: `{"title": "Что то", "description":"Описание чего то"}`,
			setupMock: func() {
				taskServiceMock.EXPECT().
					CreateTask(mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(123, nil).
					Once()
			},
			expectedStatus: http.StatusOK,
			expectedResp: func(t *testing.T, resp map[string]any) {
				assert.Equal(t, "OK", resp["result"])
				assert.Equal(t, 123, int(resp["data"].(map[string]any)["id"].(float64)))
			},
		},
		{
			name:           "Empty request",
			reqBody:        ``,
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedResp: func(t *testing.T, resp map[string]any) {
				assert.Equal(t, "error", resp["result"])
				assert.Equal(t, "empty request", resp["error"])
			},
		},
		{
			name:           "Invalid request",
			reqBody:        `{"title": 123, "description":"любая строка"}`,
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedResp: func(t *testing.T, resp map[string]any) {
				assert.Equal(t, "error", resp["result"])
				assert.Equal(t, "failed to decode request body", resp["error"])
			},
		},
		{
			name:           "Validation error",
			reqBody:        `{"title": "", "description":"любая строка"}`,
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedResp: func(t *testing.T, resp map[string]any) {
				assert.Equal(t, "error", resp["result"])
				assert.Equal(t, "field Title is a required field", resp["error"])
			},
		},
		{
			name:    "Mock error",
			reqBody: `{"title": "нормальная строка", "description":"любая строка"}`,
			setupMock: func() {
				taskServiceMock.EXPECT().
					CreateTask(mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(1, errors.New("DB is broken")).
					Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp: func(t *testing.T, resp map[string]any) {
				assert.Equal(t, "error", resp["result"])
				assert.Equal(t, "failed to create task", resp["error"])
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			req := httptest.NewRequest("POST", "/tasks", strings.NewReader(tc.reqBody))
			w := httptest.NewRecorder()

			handler.CreateTask(w, req)

			resp := w.Result()
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			var response map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if tc.expectedResp != nil {
				tc.expectedResp(t, response)
			}

		})
	}

}
