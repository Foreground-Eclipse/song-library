package getcouplet

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
	GetCouplet(filter postgres.Song, page int) ([]postgres.Song, error)
}

func New(log *logger.Logger, coupletGetter CoupletGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.couplet_get.New"

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
		var filter postgres.Song
		filter.Group = req.Group
		filter.Song = req.Song
		filter.ReleaseDate = req.ReleaseDate
		filter.Text = req.Text
		filter.Link = req.Link
		couplet, err := coupletGetter.GetCouplet(filter, req.Page)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(err))
			log.LogError("error getting the song details at ", zap.String("op", op),
				zap.Error(err))
			// TODO: ADD DEBUG LOG

			return
		}
		c.JSON(http.StatusOK, response.OK(couplet))
	}
}
