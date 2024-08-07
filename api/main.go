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
	db := initializeDB[Label]("api")
	api := setupRoutes(db)
	api.Use(cors.Default())
	api.Use(gin.Recovery())
	log.Fatal(api.Run())
}

func initializeDB[t interface{}](dbName string) *gorm.DB {
	db := connectDB(dbName)
	db.AutoMigrate(new(t))
	return db
}

func connectDB(dbName string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s.db", dbName)), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func setupRoutes(db *gorm.DB) *gin.Engine {
	api := gin.Default()
	api.GET("/api/labels", makeHandler(db, getLabels))
	api.POST("/api/labels", makeHandler(db, createLabel))
	api.GET("/api/labels/:id", makeHandler(db, getLabel))
	api.PUT("/api/posts/:id", makeHandler(db, updatePost))
	api.DELETE("/api/posts/:id", makeHandler(db, deletePost))
	return api
}

func makeHandler(db *gorm.DB, fn func(*gorm.DB, *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn(db, c)
	}
}
