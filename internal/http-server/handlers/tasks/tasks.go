package taskshandler

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	kafkaclient "github.com/Swiffy12/taskify/internal/clients/mail-notifier/kafka-client"
	"github.com/Swiffy12/taskify/internal/http-server/models"
	resp "github.com/Swiffy12/taskify/internal/lib/api/response"
	"github.com/Swiffy12/taskify/internal/lib/logger/sl"
	"github.com/Swiffy12/taskify/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type TaskHandler struct {
	log          *slog.Logger
	mailNotifier *kafkaclient.Producer
	service      TaskService
}

//go:generate mockery --r --name TaskService --output ./mocks --case underscore
type TaskService interface {
	CreateTask(title, description string) (int, error)
	GetAllTasks(title string) ([]models.Task, error)
	GetTask(id int) (models.Task, error)
	DeleteTask(id int) (int, error)
	UpdateTask(id int, req models.UpdateTaskRequest) (models.Task, error)
}

func New(log *slog.Logger, mailNotifier *kafkaclient.Producer, service TaskService) *TaskHandler {
	return &TaskHandler{log: log, mailNotifier: mailNotifier, service: service}
}

func (t *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.tasks.CreateTask"

	log := t.log.With(
		slog.String("op", op),
	)

	var req models.CreateTaskRequest
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty", sl.Err(err))
			resp.Error(w, r, http.StatusBadRequest, "empty request")
			return
		}
		log.Error("failed to decode request body", sl.Err(err))
		resp.Error(w, r, http.StatusBadRequest, "failed to decode request body")
		return
	}

	if err := validator.New().Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)
		log.Error("invalid request", sl.Err(err))
		resp.ValidationError(w, r, http.StatusBadRequest, validateErr)
		return
	}

	log.Info("creating task")
	id, err := t.service.CreateTask(req.Title, req.Description)
	if err != nil {
		log.Error("failed to create task", sl.Err(err))
		resp.Error(w, r, http.StatusInternalServerError, "failed to create task")
		return
	}

	err = t.mailNotifier.Send(fmt.Sprintf("task created with id: %d", id))
	if err != nil {
		log.Error("failed to send mail notification", sl.Err(err))
	}

	res := models.TaskIdResponse{
		Id: id,
	}

	log.Info("task created", slog.Int("id", id))
	resp.OK(w, r, http.StatusOK, res)
}

func (t *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.tasks.GetAllTasks"

	log := t.log.With(
		slog.String("op", op),
	)

	title := r.URL.Query().Get("title")

	tasks, err := t.service.GetAllTasks(title)
	if err != nil {
		log.Error("failed to get tasks", sl.Err(err))
		resp.Error(w, r, http.StatusInternalServerError, "failed to get tasks")
		return
	}

	resp.OK(w, r, http.StatusOK, tasks)
}

func (t *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.tasks.GetTask"

	log := t.log.With(
		slog.String("op", op),
	)

	strId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Error("failed to parse id", sl.Err(err))
		resp.Error(w, r, http.StatusBadRequest, "failed to parse id")
		return
	}

	task, err := t.service.GetTask(id)
	if err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			log.Error("task not found", sl.Err(err))
			w.WriteHeader(http.StatusNotFound)
			resp.Error(w, r, http.StatusNotFound, "task not found")
			return
		}
		log.Error("failed to get task", sl.Err(err))
		resp.Error(w, r, http.StatusInternalServerError, "failed to get task")
		return
	}

	resp.OK(w, r, http.StatusOK, task)
}

func (t *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.tasks.DeleteTask"

	log := t.log.With(
		slog.String("op", op),
	)

	strId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Error("failed to parse id", sl.Err(err))
		resp.Error(w, r, http.StatusBadRequest, "failed to parse id")
		return
	}

	log.Info("deleting task", slog.Int("id", id))
	id, err = t.service.DeleteTask(id)
	if err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			log.Error("task not found", sl.Err(err))
			w.WriteHeader(http.StatusNotFound)
			resp.Error(w, r, http.StatusNotFound, "task not found")
			return
		}
		log.Error("failed to delete task", sl.Err(err))
		resp.Error(w, r, http.StatusInternalServerError, "failed to delete task")
		return
	}

	res := models.TaskIdResponse{
		Id: id,
	}

	log.Info("task deleted", slog.Int("id", id))
	resp.OK(w, r, http.StatusOK, res)
}

func (t *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.tasks.UpdateTask"

	log := t.log.With(
		slog.String("op", op),
	)

	strId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Error("failed to parse id", sl.Err(err))
		resp.Error(w, r, http.StatusBadRequest, "failed to parse id")
		return
	}

	var req models.UpdateTaskRequest
	err = render.DecodeJSON(r.Body, &req)
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty", sl.Err(err))
			resp.Error(w, r, http.StatusBadRequest, "empty request")
			return
		}
		log.Error("failed to decode request body", sl.Err(err))
		resp.Error(w, r, http.StatusBadRequest, "failed to decode request body")
		return
	}

	log.Info("updating task", slog.Int("id", id))
	task, err := t.service.UpdateTask(id, req)
	if err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			log.Error("task not found", sl.Err(err))
			w.WriteHeader(http.StatusNotFound)
			resp.Error(w, r, http.StatusNotFound, "task not found")
			return
		}
		log.Error("failed to update task", sl.Err(err))
		resp.Error(w, r, http.StatusInternalServerError, "failed to update task")
		return
	}

	log.Info("task updated", slog.Int("id", int(task.Id)))
	resp.OK(w, r, http.StatusOK, task)
}
