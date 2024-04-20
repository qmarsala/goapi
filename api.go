package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"regexp"

	"gorm.io/gorm"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func viewHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request, title string) {
	p, _ := loadPage(db, title)
	renderTemplate(w, "view", p)
}

func editHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request, title string) {
	p, _ := loadPage(db, title)
	renderTemplate(w, "edit", p)
}

func saveHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request, title string) {
	content := r.FormValue("content")
	updatePage(db, title, content)
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}
	return m[2], nil
}

func makeHandler(db *gorm.DB, fn func(*gorm.DB, http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title, err := getTitle(w, r)
		if err != nil {
			return
		}
		fmt.Printf("requested: %s\n", r.URL.Path)
		fn(db, w, r, title)
		fmt.Println("done")
	}
}
