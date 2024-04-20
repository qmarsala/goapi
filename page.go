package main

import (
	"gorm.io/gorm"
)

type Page struct {
	gorm.Model
	Title   string
	Content string
}

// todo: handle errors
func loadPage(db *gorm.DB, title string) *Page {
	page := Page{}
	db.First(&page, "Title = ?", title)
	return &page
}

func updatePage(db *gorm.DB, title string, content string) {
	p := loadPage(db, title)
	db.Model(p).Update("Content", content)
}
