package couplet

import (
	"net/http"

	"github.com/foreground-eclipse/song-library/internal/lib/api/response"
	"github.com/foreground-eclipse/song-library/internal/logger"
	"github.com/foreground-eclipse/song-library/internal/storage/postgres"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Request struct {
	Group       string `json:"group" validate:"required"`
	Song        string `json:"song" validate:"required"`
	ReleaseDate string `json:"release_date,omitempty"`
	Text        string `json:"text,omitempty"`
	Link        string `json:"link,omitempty"`
	Page        int    `json:"page,omitempty"`
}

type CoupletGetter interface {
	GetCouplet(filter postgres.Song, page int) (string, error)
}

/**
 * New returns a couplet with given attributes
 * New godoc
 * @Summary Gets a couplet for a song
 * @Tags couplet
 * @Description Gets a couplet for a song from the database.
 * @Param group query string true "The group of the song"
 * @Param song query string true "The name of the song"
 * @Param page query integer false "The page number of the couplet"
 * @Param request body Request true "Request body"
 * @Success 200 {string} "The couplet for the song"
 * @Failure 400 {object} response "Bad request"
 * @Failure 500 {object} response "Internal server error"
 * @Router /api/v1/song/couplet [get]
 */
func New(log *logger.Logger, coupletGetter CoupletGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.couplet_get.New"

		var req Request

		if req.Page < 0 {
			req.Page = 0
		}

		log.LogDebug("received request", zap.String("op", op),
			zap.String("group", req.Group),
			zap.String("song", req.Song),
			zap.String("release_date", req.ReleaseDate),
			zap.String("text", req.Text),
			zap.String("link", req.Link))

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, response.Error(err))
			log.LogError("error happened at handler",
				zap.String("op:", op),
				zap.Error(err))

			return
		}

		log.LogInfo("got new request at", zap.String("op", op),
			zap.String("group", req.Group),
			zap.String("song", req.Song))

		var filter postgres.Song
		filter.Group = req.Group
		filter.Song = req.Song
		filter.ReleaseDate = req.ReleaseDate
		filter.Text = req.Text
		filter.Link = req.Link
		couplet, err := coupletGetter.GetCouplet(filter, req.Page)

		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(err))
			log.LogError("error getting the song details at ", zap.String("op", op),
				zap.Error(err))

			return
		}

		c.JSON(http.StatusOK, response.OK(couplet))
	}
}
