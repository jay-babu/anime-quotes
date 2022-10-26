package web

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Main struct {
	router *gin.Engine
	log    *zap.SugaredLogger
}

func (m Main) Start() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	m.router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "The available groups are [...]")
	})
	if err := m.router.Run(":" + port); err != nil {
		m.log.Panicf("error: %s", err)
	}
}
