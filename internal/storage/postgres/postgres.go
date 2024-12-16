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
	Group       string `json:"group" db:"group"`
	Song        string `json:"song" db:"song"`
	ReleaseDate string `json:"release_date" db:"release_date"`
	Text        string `json:"text" db:"text"`
	Link        string `json:"link" db:"link"`
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
	const op = "storage.postgres.AddSong"

	_, err := s.db.Exec(`insert into songs ("group", song, release_date, text, link)
	values ($1, $2, $3, $4, $5);`,
		song.Group, song.Song, song.ReleaseDate, song.Text, song.Link)
	fmt.Println(song.Group)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// GetSongs gets all the songs from database with given filter and page
func (s *Storage) GetSongs(filter Song, page int) (Song, error) {
	const op = "storage.postgres.GetSongs"

	// Build the base query
	query := "SELECT \"group\", song, release_date, text, link FROM songs WHERE "
	params := make([]interface{}, 0)

	rv := reflect.ValueOf(filter)
	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		value := rv.Field(i).Interface()

		if value != "" {
			fmt.Println(field.Name)
			query += fmt.Sprintf("\"%s\" = $%d AND ", strings.ToLower(field.Tag.Get("db")), len(params)+1)
			params = append(params, value)
		}
	}

	// Удалить последний " AND "
	if len(params) > 0 {
		query = query[:len(query)-5]
	}

	// Добавить пагинацию
	query += fmt.Sprintf(" LIMIT 1 OFFSET %d", page-1)

	// Execute the query
	rows, err := s.db.Query(query, params...)
	var song Song
	if err != nil {
		return song, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	// Scan the results and return them
	for rows.Next() {
		err = rows.Scan(&song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link)
		if err != nil {
			return song, fmt.Errorf("%s: %w", op, err)
		}
	}

	return song, nil
}

// GetCouplet gets all couplets from the song with given filter
func (s *Storage) GetCouplet(filter Song, page int) (string, error) {
	const op = "storage.postgres.GetSongs"

	query := "SELECT  text FROM songs WHERE "
	params := make([]interface{}, 0)

	rv := reflect.ValueOf(filter)
	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		value := rv.Field(i).Interface()

		if value != "" {
			fmt.Println(field.Name)
			query += fmt.Sprintf("\"%s\" = $%d AND ", strings.ToLower(field.Tag.Get("db")), len(params)+1)
			params = append(params, value)
		}
	}

	// Удалить последний " AND "
	if len(params) > 0 {
		query = query[:len(query)-5]
	}

	// Добавить пагинацию
	query += fmt.Sprintf(" LIMIT 1")

	// Execute the query
	rows, err := s.db.Query(query, params...)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	// Scan the results and return them
	var couplet string
	for rows.Next() {
		err = rows.Scan(&couplet)
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}
	}

	part := strings.Split(couplet, "\n")
	fmt.Println(part[1])
	fmt.Println(part[2])

	return part[page-1], nil
}
func (s *Storage) UpdateSong(song Song) error {
	const op = "storage.postgres.UpdateSong"

	_, err := s.db.Exec(`
		UPDATE songs
		 SET "group" = $1, song = $2, release_date = $3, text = $4, link = $5
		  WHERE "group" = $6 and song = $7;
	`, song.Group,
		song.Song,
		song.ReleaseDate,
		song.Text,
		song.Link,
		song.Group,
		song.Song)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) DeleteSong(group, song string) error {
	const op = "storage.postgres.DeleteSong"

	_, err := s.db.Exec("DELETE FROM songs WHERE \"group\" = $1 and song = $2", group, song)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
