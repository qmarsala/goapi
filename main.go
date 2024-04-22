package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// todo: add error handling back
func main() {
	db := initializeDB()
	api := gin.Default()
	api.GET("/posts", makeHandler(db, getPosts))
	api.POST("/posts", makeHandler(db, createPost))
	api.GET("/posts/:id", makeHandler(db, getPost))
	api.PUT("/posts/:id", makeHandler(db, updatePost))
	api.DELETE("/posts/:id", makeHandler(db, deletePost))
	log.Fatal(api.Run())
}

func initializeDB() *gorm.DB {
	db := connectDB()
	db.AutoMigrate(&post{})
	return db
}

func connectDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
