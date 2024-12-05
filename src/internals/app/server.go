package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Swiffy12/taskify/src/api"
	"github.com/Swiffy12/taskify/src/internals/app/handlers"
	"github.com/Swiffy12/taskify/src/internals/app/services"
	"github.com/Swiffy12/taskify/src/internals/config"
)

type Server struct {
	config config.Config
	ctx    context.Context
	srv    *http.Server
	db     string
}

func NewServer(config config.Config, ctx context.Context) *Server {
	server := new(Server)
	server.config = config
	server.ctx = ctx
	return server
}

func (server *Server) Listen() {
	log.Println("Starting server")

	taskService := services.NewTasksService()

	taskHandler := handlers.NewTaskHandler(*taskService)

	routes := api.CreateRoutes(taskHandler)

	server.srv = &http.Server{
		Addr:    ":" + server.config.Port,
		Handler: routes,
	}

	log.Println("Server started")

	err := server.srv.ListenAndServe()

	if err != nil {
		log.Fatalln(err)
	}
}

func (server *Server) Shutdown() {
	log.Println("Server stopped")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		cancel()
	}()

	if err := server.srv.Shutdown(ctxShutdown); err != nil {
		log.Fatal("server Shutdown failed: ", err)
	}

	log.Println("Server exited properly")
}
