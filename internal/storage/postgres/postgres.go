package postgres

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/foreground-eclipse/song-library/internal/config"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

type Song struct {
	Group       string   `json:"group"`
	Song        string   `json:"song"`
	ReleaseDate string   `json:"release_date"`
	Text        string   `json:"text"`
	Link        string   `json:"link"`
	Couplet     []string `json:"couplet"`
}

// New initializing new database connection
func New(cfg *config.Config) (*Storage, error) {
	const op = "storage.postgres.New"

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) AddSong(song Song) error {
	const op = "storage.postgres.insertSong"

	_, err := s.db.Exec(`insert into songs (group, song, release_date, text, link)
	values ($1, $2, $3, $4, $5);`,
		song.Group, song.Song, song.ReleaseDate, song.Text, song.Link)

	return fmt.Errorf("%s: %w", op, err)
}

// GetSongs gets all the songs from database with given filter and page
func (s *Storage) GetSongs(filter Song, page int) ([]Song, error) {
	const op = "storage.postgres.GetSongs"

	var query string

	query = "select * from songs where "

	value := reflect.ValueOf(filter)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.Kind() == reflect.String && field.Len() != 0 {
			query += fmt.Sprintf("%s = '%s' AND", value.Type().Field(i).Name, field.String())
		}
	}

	query = query[:len(query)-5]

	query += " LIMIT? OFFSET?"

	rows, err := s.db.Query(query, 1, (page - 1))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var songs []Song

	for rows.Next() {
		var song Song
		err = rows.Scan(&song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		songs = append(songs, song)
	}
	return songs, nil
}

// GetCouplet gets all couplets from the song with given filter
func (s *Storage) GetCouplet(filter Song, page int) ([]Song, error) {
	const op = "storage.postgres.GetSongs"

	var query string

	query = "SELECT id, group, song, release_date, text, link FROM songs WHERE "

	value := reflect.ValueOf(filter)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.Kind() == reflect.String && field.Len() != 0 {
			query += fmt.Sprintf("%s = '%s' AND", value.Type().Field(i).Name, field.String())
		}
	}

	query = query[:len(query)-5]

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", 1, (page - 1))

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var songs []Song

	for rows.Next() {
		var song Song
		err = rows.Scan(&song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		song.Couplet = strings.Split(song.Text, "\n\n")

		songs = append(songs, song)
	}
	return songs, nil
}
func (s *Storage) UpdateSong(song Song) error {
	const op = "storage.postgres.UpdateSong"

	_, err := s.db.Exec(`
		UPDATE songs
		SET group = $1, song = $2, release_date = $3, text = $4, link = $5
		WHERE group = $6 and song = $7;
	`, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link, song.Group, song.Song)
	return fmt.Errorf("%s: %w", op, err)
}

func (s *Storage) DeleteSong(group, song string) error {
	const op = "storage.postgres.DeleteSong"

	_, err := s.db.Exec("DELETE FROM songs WHERE group = $1 and song = $2", group, song)
	return fmt.Errorf("%s: %w", op, err)
}
