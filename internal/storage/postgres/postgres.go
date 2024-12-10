package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Song struct {
	gorm.Model
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func NewDatabase() (*gorm.DB, error) {
	dsn := "host=localhost user=myuser password=mypassword dbname=mydb port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Song{})
}
