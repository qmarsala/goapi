package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// todo: add error handling back
func main() {
	db := setupDB()

	api := gin.Default()
	api.GET("/posts", makeHandler(db, getPosts))
	api.POST("/posts", makeHandler(db, createPost))
	api.GET("/posts/:id", makeHandler(db, getPost))
	api.PUT("/posts/:id", makeHandler(db, updatePost))
	api.DELETE("/posts/:id", makeHandler(db, deletePost))
	log.Fatal(api.Run())
}

func setupDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&post{})

	if tx := db.Limit(3).Find(&post{}); tx.RowsAffected < 1 {
		for _, p := range []post{
			{Message: "Hello!"},
			{Message: "Hello, Go!"},
			{Message: "Hello, World!"},
		} {
			tx := db.Model(&post{}).Create(&p)
			fmt.Println(tx.Error)
		}
	}
	return db
}
