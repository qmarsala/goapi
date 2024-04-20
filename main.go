package main

import (
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// todo: add error handling back
func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Page{})

	db.Create(&Page{
		Title:   "test",
		Content: "This is a test page!",
	})

	http.HandleFunc("/view/", makeHandler(db, viewHandler))
	http.HandleFunc("/edit/", makeHandler(db, editHandler))
	http.HandleFunc("/save/", makeHandler(db, saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
