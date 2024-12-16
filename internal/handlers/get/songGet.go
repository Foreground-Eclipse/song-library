package songget

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
	Page        int    `json:"page" validate:"required"`
	ReleaseDate string `json:"release_date,omitempty"`
	Text        string `json:"text,omitempty"`
	Link        string `json:"link,omitempty"`
}

type SongGetter interface {
	GetSongs(filter postgres.Song, page int) (postgres.Song, error)
}

// New is a handler for getting all songs with given filter
func New(log *logger.Logger, songGetter SongGetter) gin.HandlerFunc {
	/**
	 * @Summary Gets a song by attributes
	 * @Description Gets a song from the database by its attributes.
	 * @Param Request body required true "The song attributes to search"
	 * @Success 200 {object} Song "The found song"
	 * @Failure 400 {object} response "Bad request"
	 * @Failure 500 {object} response "Internal server error"
	 * @Router /song [get]
	 */
	return func(c *gin.Context) {
		const op = "handlers.song_get.New"

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
			zap.String("song", req.Song),
			zap.Int("page", req.Page))
		var filter postgres.Song
		filter.Group = req.Group
		filter.Song = req.Song
		filter.ReleaseDate = req.ReleaseDate
		filter.Text = req.Text
		filter.Link = req.Link

		song, err := songGetter.GetSongs(filter, req.Page)
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
