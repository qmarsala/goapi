package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db := initializeDB()
	api := gin.Default()
	api.Use(cors.Default())
	api.GET("/api/posts", makeHandler(db, getPosts))
	api.POST("/api/posts", makeHandler(db, createPost))
	api.GET("/api/posts/:id", makeHandler(db, getPost))
	api.PUT("/api/posts/:id", makeHandler(db, updatePost))
	api.DELETE("/api/posts/:id", makeHandler(db, deletePost))
	log.Fatal(api.Run())
}

func initializeDB() *gorm.DB {
	db := connectDB("api")
	db.AutoMigrate(post{})
	return db
}

func connectDB(dbName string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s.db", dbName)), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
