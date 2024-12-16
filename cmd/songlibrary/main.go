package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/foreground-eclipse/song-library/cmd/migrator"
	docs "github.com/foreground-eclipse/song-library/cmd/songlibrary/docs"

	"github.com/foreground-eclipse/song-library/internal/config"
	addsong "github.com/foreground-eclipse/song-library/internal/handlers/add"
	"github.com/foreground-eclipse/song-library/internal/handlers/couplet"

	songdelete "github.com/foreground-eclipse/song-library/internal/handlers/delete"
	songget "github.com/foreground-eclipse/song-library/internal/handlers/get"
	"github.com/foreground-eclipse/song-library/internal/handlers/update"

	"github.com/foreground-eclipse/song-library/internal/logger"
	"github.com/foreground-eclipse/song-library/internal/storage/postgres"
	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go.uber.org/zap"
)

// @title Song library test api
// @version 1
// @Description simple rest api representing a song library
// @host localhost:8080
// @BasePath /api/v1

func init() {

}

func main() {
	cfg := config.MustLoad()
	docs.SwaggerInfo.BasePath = "/api/v1"
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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.POST("/api/v1/song/add", addsong.New(log, storage))
	router.GET("/api/v1/song/get", songget.New(log, storage))
	router.GET("/api/v1/song/couplet", couplet.New(log, storage))
	router.DELETE("/api/v1/song/delete", songdelete.New(log, storage))
	router.POST("/api/v1/song/update", update.New(log, storage))

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
