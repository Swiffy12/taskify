package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Swiffy12/taskify/src/internals/app/models"
	"github.com/Swiffy12/taskify/src/internals/app/services"
	"github.com/Swiffy12/taskify/src/internals/constants"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

type TasksHandler struct {
	service *services.TasksService
}

func NewTasksHandler(service *services.TasksService) *TasksHandler {
	tasksHandler := new(TasksHandler)
	tasksHandler.service = service
	return tasksHandler
}

func (tasksHandler *TasksHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask models.CreateTaskRequestDTO

	userId := r.Context().Value(constants.UserIdKey).(string)
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		WrapErrorBadRequest(w, errors.New("ошибка в параметрах запроса"))
		return
	}

	validation, err := govalidator.ValidateStruct(newTask)
	if !validation {
		WrapErrorBadRequest(w, errors.New("недопустимая форма ввода"))
		return
	}
	if err != nil {
		logrus.Errorln(err)
		WrapErrorInternalServerError(w)
		return
	}

	createdUser, err := tasksHandler.service.CreateOneTask(userId, newTask)
	if err != nil {
		WrapErrorInternalServerError(w)
		return
	}

	WrapOK(w, createdUser)
}

func (tasksHandler *TasksHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	var id, creatorId, assignedId uint64
	var err error

	queryId := vars.Get("id")
	if queryId != "" {
		id, err = strconv.ParseUint(queryId, 10, 64)
		if err != nil {
			WrapErrorBadRequest(w, errors.New("недопустимый формат ввода"))
			return
		}
	}

	queryCreatorId := vars.Get("creator_id")
	if queryCreatorId != "" {
		creatorId, err = strconv.ParseUint(queryCreatorId, 10, 64)
		if err != nil {
			WrapErrorBadRequest(w, errors.New("недопустимый формат ввода"))
			return
		}
	}

	queryAssignedId := vars.Get("assigned_id")
	if queryAssignedId != "" {
		assignedId, err = strconv.ParseUint(queryAssignedId, 10, 64)
		if err != nil {
			WrapErrorBadRequest(w, errors.New("недопустимый формат ввода"))
			return
		}
	}

	queryParams := models.GetTasksRequestDTO{
		Id:         id,
		Title:      vars.Get("title"),
		CreatorId:  creatorId,
		AssignedId: assignedId,
	}

	tasks, err := tasksHandler.service.GetTasksWithFilter(queryParams)
	if err != nil {
		WrapErrorInternalServerError(w)
		return
	}

	WrapOK(w, tasks)
}

func (tasksHandler *TasksHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["id"] == "" {
		WrapErrorBadRequest(w, errors.New("пропущено id задачи"))
		return
	}

	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		WrapErrorBadRequest(w, errors.New("недопустимый формат ввода id задачи"))
		return
	}

	task, err := tasksHandler.service.GetOneTask(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			WrapErrorNotFound(w, errors.New("не удалось найти данную задачу"))
			return
		}
		WrapErrorInternalServerError(w)
		return
	}
	WrapOK(w, task)
}

func (tasksHandler *TasksHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["id"] == "" {
		WrapErrorBadRequest(w, errors.New("пропущено id задачи"))
		return
	}

	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		WrapErrorBadRequest(w, errors.New("недопустимый формат ввода id задачи"))
		return
	}

	err = tasksHandler.service.DeleteOneTask(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			WrapErrorNotFound(w, errors.New("не удалось найти данную задачу"))
			return
		}
		WrapErrorInternalServerError(w)
		return
	}
	WrapOK(w, nil)
}

func (tasksHandler *TasksHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["id"] == "" {
		WrapErrorBadRequest(w, errors.New("пропущено id задачи"))
		return
	}

	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		WrapErrorBadRequest(w, errors.New("недопустимый формат ввода id задачи"))
		return
	}

	var taskBody models.UpdateTaskRequestDTO

	err = json.NewDecoder(r.Body).Decode(&taskBody)
	if err != nil {
		WrapErrorBadRequest(w, errors.New("ошибка в параметрах запроса"))
		return
	}

	validation, err := govalidator.ValidateStruct(taskBody)
	if !validation {
		WrapErrorBadRequest(w, errors.New("недопустимая форма ввода"))
		return
	}
	if err != nil {
		logrus.Errorln(err)
		WrapErrorInternalServerError(w)
		return
	}

	updatedTask, err := tasksHandler.service.UpdateOneTask(id, taskBody)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			WrapErrorNotFound(w, errors.New("не удалось найти данную задачу"))
			return
		}
		WrapErrorInternalServerError(w)
		return
	}
	WrapOK(w, updatedTask)
}
