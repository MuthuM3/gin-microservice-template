package app

import (
	"log"
	"net/http"

	"github.com/MuthuM3/gin-microservice-template/internal/config"
)

type App struct {
	config *config.Config
	logger *log.Logger
	server *http.Server
}
