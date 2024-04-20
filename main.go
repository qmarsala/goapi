package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"regexp"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type TestError struct {
	Message string
}

func (e TestError) Error() string {
	return e.Message
}

type Page struct {
	gorm.Model
	Title string
	Body  []byte
}

func (p *Page) save(db *gorm.DB) {
	db.Model(p).Update("Title", p.Title)
	db.Model(p).Update("Body", p.Body)
}

func loadPage(db *gorm.DB, title string) *Page {
	page := Page{}
	db.First(&page, "Title = ?", title)
	return &page
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}
	return m[2], nil
}

func viewHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request, title string) {
	p := loadPage(db, title)
	renderTemplate(w, "view", p)
}

func editHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request, title string) {
	p := loadPage(db, title)
	renderTemplate(w, "edit", p)
}

func saveHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := loadPage(db, title)
	p.Body = []byte(body)
	p.save(db)
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func makeHandler(db *gorm.DB, fn func(*gorm.DB, http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title, err := getTitle(w, r)
		if err != nil {
			return
		}
		fn(db, w, r, title)
	}
}

// todo: add error handling back
func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Page{})

	db.Create(&Page{
		Title: "test",
		Body:  []byte("This is a test page!"),
	})

	http.HandleFunc("/view/", makeHandler(db, viewHandler))
	http.HandleFunc("/edit/", makeHandler(db, editHandler))
	http.HandleFunc("/save/", makeHandler(db, saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
