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
	api := gin.Default()
	api.Use(cors.Default())
	setupRoutes(api, db)
	log.Fatal(api.Run())
}

type Database struct {
	*gorm.DB
}

func initializeDB[t interface{}](dbName string) Database {
	db := connectDB(dbName)
	db.AutoMigrate(new(t))
	return db
}

func connectDB(dbName string) Database {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s.db", dbName)), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return Database{db}
}

func setupRoutes(api *gin.Engine, db Database) {
	api.GET("/api/labels", makeHandler(db, getLabels))
	api.POST("/api/labels", makeHandler(db, createLabel))
	api.GET("/api/labels/:id", makeHandler(db, getLabel))
	api.PUT("/api/labels/:id", makeHandler(db, updatePost))
	api.DELETE("/api/labels/:id", makeHandler(db, deletePost))
}

func makeHandler(db Database, fn func(Database, *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn(db, c)
	}
}
