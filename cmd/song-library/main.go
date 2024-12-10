package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/foreground-eclipse/song-library/internal/config"
	"github.com/foreground-eclipse/song-library/internal/logger"
	"github.com/foreground-eclipse/song-library/internal/storage/postgres"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	logLevel := "INFO"
	log, err := logger.NewLogger(logLevel)
	if err != nil {
		fmt.Println(err)
		return
	}

	if log == nil {
		fmt.Println("Logger is not properly initialized")
		return
	}

	initLogger(log)

	db, err := postgres.NewDatabase()
	if err != nil {
		panic(err)
	}

	postgres.Migrate(db)

	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		log.LogInfo("Received GET request to /ping")
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	if err := router.Run(":8080"); err != nil {
		log.LogError("Failed to start server", zap.Error(err))
		if err := log.Sync(); err != nil {
			fmt.Println(err)
		}
		os.Exit(1)
	}

	if err := log.Sync(); err != nil {
		fmt.Println(err)
	}

}

func initLogger(log *logger.Logger) {
	if log == nil {
		fmt.Println("Logger is not properly initialized")
		return
	}

	log.LogInfo("This is an info message")
	log.LogError("This is an error message")

}
