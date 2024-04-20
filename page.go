package main

import (
	"gorm.io/gorm"
)

type Page struct {
	gorm.Model
	Title   string
	Content string
}

func loadPage(db *gorm.DB, title string) (*Page, error) {
	page := Page{}
	tx := db.First(&page, "Title = ?", title)
	return &page, tx.Error
}

func createPage(db *gorm.DB, title string, content string) (*Page, error) {
	p := &Page{Title: title, Content: content}
	tx := db.Create(p)
	return p, tx.Error
}

func updatePage(db *gorm.DB, title string, content string) (*Page, error) {
	p, err := loadPage(db, title)
	if err != nil {
		return nil, err
	}
	tx := db.Model(p).Update("Content", content)
	return &Page{p.Model, p.Title, content}, tx.Error
}

func deletePage(db *gorm.DB, title string) error {
	p, err := loadPage(db, title)
	if err != nil {
		return err
	}
	db.Delete(&p)
	return nil
}
