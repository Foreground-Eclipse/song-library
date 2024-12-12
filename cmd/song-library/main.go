package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/foreground-eclipse/song-library/cmd/migrator"
	"github.com/foreground-eclipse/song-library/internal/config"
	addsong "github.com/foreground-eclipse/song-library/internal/http-server/handlers/song_add"
	"github.com/foreground-eclipse/song-library/internal/logger"
	"github.com/foreground-eclipse/song-library/internal/storage/postgres"
	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	logLevel := "DEBUG"
	log, err := logger.NewLogger(logLevel)
	if err != nil {
		fmt.Println(err)
		return
	}

	initLogger(log)

	if log == nil {
		fmt.Println("Logger is not properly initialized")
		return
	}

	storage, err := postgres.New(cfg)
	if err != nil {
		log.LogError("failed to create database connection", zap.Error(err))
		os.Exit(1)
	}

	err = migrator.Migrate(log, cfg)
	if err != nil {
		log.LogError("failed to apply migrations", zap.Error(err))
	}

	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		log.LogInfo("Received GET request to /ping")
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.POST("/songs/add", addsong.New(log, storage))

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

	log.LogInfo("logger initialized")

}
