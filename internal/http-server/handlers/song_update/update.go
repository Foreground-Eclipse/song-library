package updatesong

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

type SongUpdater interface {
	UpdateSong(song postgres.Song) error
}

// New updates the song and sets the given attributes
func New(log *logger.Logger, songUpdate SongUpdater) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.update.UpdateSong"

		var req Request
		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(err))
			log.LogError("error binding JSON at ", zap.String("op", op),
				zap.Error(err))

			return
		}

		log.LogInfo("got new request at", zap.String("op", op),
			zap.String("group", req.Group),
			zap.String("song", req.Song))
		var song postgres.Song
		err = songUpdate.UpdateSong(song)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(err))
			log.LogError("error getting the song details at ", zap.String("op", op),
				zap.Error(err))
			// TODO: ADD DEBUG LOG

			return
		}
		c.JSON(http.StatusOK, response.OK(song))
	}
}
