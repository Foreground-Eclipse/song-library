package addsong

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/foreground-eclipse/song-library/internal/lib/api/response"
	"github.com/foreground-eclipse/song-library/internal/logger"
	"github.com/foreground-eclipse/song-library/internal/storage/postgres"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Request struct {
	Group string `json:"group" validate:"required"`
	Song  string `json:"song" validate:"required"`
}

type SongDetail struct {
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type SongAdder interface {
	AddSong(song postgres.Song) error
}

// New adds the song in database
func New(log *logger.Logger, songAdder SongAdder) gin.HandlerFunc {

	/**
	 * @BasePath /api/v1
	 * @Summary Adds a new song
	 * @Description Adds a new song to the database with the given attributes.
	 * @Tag Song
	 * @OperationId addSong
	 * @Param Request body required true "The song attributes to add"
	 * @Success 200 {object} Song "The added song"
	 * @Failure 400 {object} response "Bad request"
	 * @Failure 500 {object} response "Internal server error"
	 * @Router /api/v1/song/add [post]
	 */
	return func(c *gin.Context) {
		const op = "handlers.song_add.New"

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
		song := postgres.Song{
			Group: req.Group,
			Song:  req.Song,
		}
		info, err := GetInfo(log, req.Group, req.Song)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(err))
			log.LogError("error at handling at ", zap.String("op", op),
				zap.Error(err))
			return
		}
		song.Link = info.Link
		song.ReleaseDate = info.ReleaseDate
		song.Text = info.Text

		err = songAdder.AddSong(song)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(err))
			log.LogError("error at handling at ", zap.String("op", op),
				zap.Error(err))
			return
		}

		c.JSON(http.StatusOK, response.OK(song))
	}
}

// GetInfo gets full information about the song from API
func GetInfo(log *logger.Logger, group, song string) (SongDetail, error) {
	const op = "handlers.song_add.GetInfo"

	log.LogInfo("trying to get full info", zap.String("op", op),
		zap.String("group", group),
		zap.String("song", song))

	client := &http.Client{}

	url := "http://localhost:8080/info?group=" + group + "&song=" + song

	resp, err := client.Get(url)
	if err != nil {
		return SongDetail{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return SongDetail{}, errors.New("internal error")
	}

	if resp.Body == nil {
		return SongDetail{}, errors.New("response body is nil")
	}

	var info SongDetail
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return SongDetail{}, err
	}

	log.LogDebug("got the full information", zap.String("op", op),
		zap.String("group", group),
		zap.String("song", song))

	return info, nil

}
