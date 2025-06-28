package taskhandler

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Swiffy12/taskify/internal/http-server/models"
	taskservice "github.com/Swiffy12/taskify/internal/http-server/services/task"
	resp "github.com/Swiffy12/taskify/internal/lib/api/response"
	"github.com/Swiffy12/taskify/internal/lib/logger/sl"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type TaskHandler struct {
	log     *slog.Logger
	service *taskservice.TaskService
}

func New(log *slog.Logger, service *taskservice.TaskService) *TaskHandler {
	return &TaskHandler{log: log, service: service}
}

func (t *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.task.CreateTask"

	log := t.log.With(
		slog.String("op", op),
	)

	var req models.CreateTaskRequest
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty", sl.Err(err))
			render.JSON(w, r, resp.Error("empty request"))
			return
		}
		log.Error("failed to decode request body", sl.Err(err))
		render.JSON(w, r, "failed to decode request body")
		return
	}

	if err := validator.New().Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)
		log.Error("invalid request", sl.Err(err))
		render.JSON(w, r, resp.ValidationError(validateErr))
		return
	}

	log.Info("creating task")
	id, err := t.service.CreateTask(req.Title, req.Description)
	if err != nil {
		// Ошибку storage
		log.Error("failed to create task", sl.Err(err))
		render.JSON(w, r, resp.Error("failed to create task"))
		return
	}

	res := models.TaskIdResponse{
		Id: id,
	}

	log.Info("task created", slog.Int("id", id))
	render.JSON(w, r, resp.OK(res))
}

func (t *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.task.GetAllTasks"

	log := t.log.With(
		slog.String("op", op),
	)

	title := r.URL.Query().Get("title")

	tasks, err := t.service.GetAllTasks(title)
	if err != nil {
		// Ошибку storage
		log.Error("failed to get tasks", sl.Err(err))
		render.JSON(w, r, resp.Error("failed to get tasks"))
		return
	}

	render.JSON(w, r, resp.OK(tasks))
}

func (t *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.task.GetTask"

	log := t.log.With(
		slog.String("op", op),
	)

	strId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Error("failed to parse id", sl.Err(err))
		render.JSON(w, r, resp.Error("failed to parse id"))
		return
	}

	task, err := t.service.GetTask(id)
	if err != nil {
		// Ошибку storage not found
		log.Error("failed to get task", sl.Err(err))
		render.JSON(w, r, resp.Error("failed to get task"))
		return
	}

	render.JSON(w, r, resp.OK(task))
}

func (t *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.task.DeleteTask"

	log := t.log.With(
		slog.String("op", op),
	)

	strId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Error("failed to parse id", sl.Err(err))
		render.JSON(w, r, resp.Error("failed to parse id"))
		return
	}

	log.Info("deleting task", slog.Int("id", id))
	id, err = t.service.DeleteTask(id)
	if err != nil {
		// if errors.Is(err) {
		// 	log.Error("task not found")
		// }
		log.Error("failed to delete task", sl.Err(err))
		render.JSON(w, r, resp.Error("failed to delete task"))
		return
	}

	res := models.TaskIdResponse{
		Id: id,
	}

	log.Info("task deleted", slog.Int("id", id))
	render.JSON(w, r, resp.OK(res))
}

func (t *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.task.UpdateTask"

	log := t.log.With(
		slog.String("op", op),
	)

	strId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Error("failed to parse id", sl.Err(err))
		render.JSON(w, r, resp.Error("failed to parse id"))
		return
	}

	var req models.UpdateTaskRequest
	err = render.DecodeJSON(r.Body, &req)
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty", sl.Err(err))
			render.JSON(w, r, resp.Error("empty request"))
			return
		}
		log.Error("failed to decode request body", sl.Err(err))
		render.JSON(w, r, "failed to decode request body")
		return
	}

	log.Info("updating task", slog.Int("id", id))
	task, err := t.service.UpdateTask(id, req)
	if err != nil {
		// Ошибка storage not found
		log.Error("failed to update task", sl.Err(err))
		render.JSON(w, r, resp.Error("failed to update task"))
		return
	}

	log.Info("task updated", slog.Int("id", int(task.Id)))
	render.JSON(w, r, resp.OK(task))
}
