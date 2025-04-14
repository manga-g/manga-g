package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("manga.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	// Auto-migrate the Manga model
	err = DB.AutoMigrate(&Manga{})
	if err != nil {
		return err
	}
	return nil
}

type Manga struct {
	ID    string `gorm:"primaryKey"`
	Title string
	// Add more fields as needed
}

func GetMangaFromDB(query string) []Manga {
	if DB == nil {
		fmt.Println("Error: Database not initialized")
		return []Manga{}
	}
	var mangas []Manga
	DB.Where("title LIKE ?", "%"+query+"%").Find(&mangas)
	return mangas
}

func GetMangaDetailFromDB(id string) Manga {
	if DB == nil {
		fmt.Println("Error: Database not initialized")
		return Manga{}
	}
	var manga Manga
	DB.Where("id = ?", id).First(&manga)
	return manga
}
