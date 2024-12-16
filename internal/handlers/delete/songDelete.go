package songdelete

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

// New deletes the song with given group and song
func New(log *logger.Logger, songDeleter SongDeleter) gin.HandlerFunc {
	/**
	 * @Summary Deletes a song
	 * @Description Deletes a song from the database by its group and name.
	 * @Param Request body required true "The song attributes to delete"
	 * @Success 200 {object} response "Song deleted successfully"
	 * @Failure 400 {object} response "Bad request"
	 * @Failure 500 {object} response "Internal server error"
	 * @Router /song [delete]
	 */
	return func(c *gin.Context) {
		const op = "handlers.song_delete.New"
		var req Request
		log.LogDebug("received request", zap.String("op", op),
			zap.String("group", req.Group),
			zap.String("song", req.Song),
		)

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

		err := songDeleter.DeleteSong(req.Group, req.Song)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(err))
			log.LogError("error deleting the song at ", zap.String("op", op),
				zap.Error(err))

			return
		}
		c.JSON(http.StatusOK, response.OK(nil))

	}
}
