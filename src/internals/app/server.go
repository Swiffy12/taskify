package app

import (
	"context"
	"net/http"
	"time"

	"github.com/Swiffy12/taskify/src/api"
	"github.com/Swiffy12/taskify/src/internals/app/handlers"
	"github.com/Swiffy12/taskify/src/internals/app/services"
	"github.com/Swiffy12/taskify/src/internals/app/storages"
	"github.com/Swiffy12/taskify/src/internals/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type Server struct {
	config config.Config
	ctx    context.Context
	srv    *http.Server
	db     *pgxpool.Pool
}

func NewServer(config config.Config, ctx context.Context) *Server {
	server := new(Server)
	server.config = config
	server.ctx = ctx
	return server
}

func (server *Server) Listen() {
	logrus.Println("Starting server")
	var err error

	server.db, err = pgxpool.Connect(server.ctx, server.config.GetStringDatabaseConnection())
	if err != nil {
		logrus.Fatalln(err)
	}

	usersStorage := storages.NewUsersStorage(server.db)
	tasksStorage := storages.NewTasksStorage(server.db)

	tasksService := services.NewTasksService(tasksStorage)
	usersService := services.NewUsersService(usersStorage)

	tasksHandler := handlers.NewTasksHandler(tasksService)
	usersHandler := handlers.NewUsersHandler(usersService)

	routes := api.CreateRoutes(tasksHandler, usersHandler)

	server.srv = &http.Server{
		Addr:    ":" + server.config.Port,
		Handler: routes,
	}

	logrus.Println("Server started")

	err = server.srv.ListenAndServe()
	if err != nil {
		logrus.Fatalln(err)
	}
}

func (server *Server) Shutdown() {
	logrus.Println("Server stopped")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	server.db.Close()
	defer func() {
		cancel()
	}()

	if err := server.srv.Shutdown(ctxShutdown); err != nil {
		logrus.Fatal("server Shutdown failed: ", err)
	}

	logrus.Println("Server exited properly")
}
