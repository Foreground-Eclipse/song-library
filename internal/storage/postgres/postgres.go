package postgres

import (
	"database/sql"
	"fmt"

	"github.com/foreground-eclipse/song-library/internal/config"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type Storage struct {
	db *sql.DB
}

type Song struct {
	gorm.Model
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// New creates new instance of database
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

// func (s *Storage) Init() error {
// 	const op = "storage.postgres.Init"

// 	if err := s.CreateSongsTable(); err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}
// 	return nil
// }

// func (s *Storage) CreateSongsTable() error {
// 	const op = "storage.postgres.createSongsTable"

// 	_, err := s.db.Exec(`
// 		CREATE TABLE IF NOT EXISTS songs (
// 			id SERIAL PRIMARY KEY,
// 			group VARCHAR(255) NOT NULL,
// 			song VARCHAR(255) NOT NULL,
// 			release_date VARCHAR(255) NOT NULL,
// 			text TEXT NOT NULL,
// 			link VARCHAR(255) NOT NULL
// 		);
// 	`)

// 	return fmt.Errorf("%s: %w", op, err)
// }

func (s *Storage) InsertSong(song Song) error {
	const op = "storage.postgres.insertSong"

	_, err := s.db.Exec(`insert into songs (group, song, release_date, text, link)
	values ($1, $2, $3, $4, $5);`,
		song.Group, song.Song, song.ReleaseDate, song.Text, song.Link)

	return fmt.Errorf("%s: %w", op, err)
}

func (s *Storage) GetSongs() ([]Song, error) {
	const op = "storage.postgres.GetSongs"

	rows, err := s.db.Query("SELECT * FROM songs")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var songs []Song
	for rows.Next() {
		var song Song
		err = rows.Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		songs = append(songs, song)
	}
	return songs, nil
}

func (s *Storage) UpdateSong(song Song) error {
	const op = "storage.postgres.UpdateSong"

	_, err := s.db.Exec(`
		UPDATE songs
		SET group = $1, song = $2, release_date = $3, text = $4, link = $5
		WHERE id = $6;
	`, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link, song.ID)
	return fmt.Errorf("%s: %w", op, err)
}

func (s *Storage) DeleteSong(id int) error {
	const op = "storage.postgres.DeleteSong"

	_, err := s.db.Exec("DELETE FROM songs WHERE id = $1", id)
	return fmt.Errorf("%s: %w", op, err)
}
