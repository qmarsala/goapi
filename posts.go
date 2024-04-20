package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type post struct {
	ID      uint   `json:"id" gorm:"primarykey"`
	Message string `json:"message"`
}

func getPosts(db *gorm.DB, c *gin.Context) {
	posts := []post{}
	db.Limit(10).Find(&posts)
	fmt.Println(posts)
	c.JSON(200, posts)
}

func createPost(db *gorm.DB, c *gin.Context) {
	p := &post{Message: "test"}
	c.JSON(201, p)
}

func updatePost(db *gorm.DB, c *gin.Context) {
	p := &post{Message: "test"}
	c.JSON(200, p)
}

func deletePost(db *gorm.DB, c *gin.Context) {
	c.Status(204)
}

func makeHandler(db *gorm.DB, fn func(*gorm.DB, *gin.Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		fn(db, c)
	}
}
