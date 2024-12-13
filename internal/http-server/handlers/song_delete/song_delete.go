package deletesong

import (
	"net/http"

	"github.com/foreground-eclipse/song-library/internal/lib/api/response"
	"github.com/foreground-eclipse/song-library/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Request struct {
	Group string `json:"group" validate:"required"`
	Song  string `json:"song" validate:"required"`
}

type SongDeleter interface {
	DeleteSong(group, song string) error
}

// TODO: VALIDATE EVERYTHING I MADE THAT AFTER SURGERY I DONT HAVE CLEAR MIND RN
// New deletes the song with given group and song
func New(log *logger.Logger, songDeleter SongDeleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.song_delete.New"
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

		err = songDeleter.DeleteSong(req.Group, req.Song)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(err))
			log.LogError("error deleting the song at ", zap.String("op", op),
				zap.Error(err))
			// TODO: ADD DEBUG LOG

			return
		}
		c.JSON(http.StatusOK, response.OK(nil))

	}
}
