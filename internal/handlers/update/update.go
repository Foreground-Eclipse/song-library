package update

import (
	"errors"
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
	ReleaseDate string `json:"release_date" validate:"required"`
	Text        string `json:"text" validate:"required"`
	Link        string `json:"link" validate:"required"`
}

type SongUpdater interface {
	UpdateSong(song postgres.Song) error
}

// New updates the song and sets the given attributes
func New(log *logger.Logger, songUpdate SongUpdater) gin.HandlerFunc {
	/**
	 * @Summary Updates a song
	 * @Description Updates an existing song in the database with the given attributes.
	 * @Param Request body required true "The updated song attributes"
	 * @Success 200 {object} Song "The updated song"
	 * @Failure 400 {object} response "Bad request"
	 * @Failure 500 {object} response "Internal server error"
	 * @Router /song [put]
	 */
	return func(c *gin.Context) {
		const op = "handlers.update.UpdateSong"

		var req Request
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

		if req.Group == "" || req.Song == "" || req.ReleaseDate == "" || req.Text == "" || req.Link == "" {
			c.JSON(http.StatusBadRequest, response.Error(errors.New("missing required fields")))
			return
		}

		song := postgres.Song{
			Group:       req.Group,
			Song:        req.Song,
			ReleaseDate: req.ReleaseDate,
			Text:        req.Text,
			Link:        req.Link,
		}

		err := songUpdate.UpdateSong(song)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(err))
			log.LogError("error getting the song details at ", zap.String("op", op),
				zap.Error(err))

			return
		}

		c.JSON(http.StatusOK, response.OK(song))
	}
}
