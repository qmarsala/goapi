package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type post struct {
	ID      uint   `json:"id" gorm:"primarykey"`
	Message string `json:"message"`
}

func getPosts(db *gorm.DB, c *gin.Context) {
	posts := []post{}
	db.Limit(25).Find(&posts)
	c.JSON(200, posts)
}

func getPostById(db *gorm.DB, id uint) (*post, error) {
	p := &post{}
	tx := db.Find(&p, "ID = ?", uint(id))
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected < 1 {
		return nil, nil
	}
	return p, nil
}

func getPost(db *gorm.DB, c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	p, err := getPostById(db, uint(id))
	switch {
	case err != nil:
		c.Status(500)
	case p != nil:
		c.JSON(200, p)
	default:
		c.Status(404)
	}
}

func createPost(db *gorm.DB, c *gin.Context) {
	p := &post{Message: "test"}
	c.JSON(201, p)
}

func updatePost(db *gorm.DB, c *gin.Context) {
	fmt.Println(c.Request)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	p, err := getPostById(db, uint(id))
	if p != nil {
		update := &post{}
		c.Bind(update)
		//todo: validate request
		tx := db.Model(p).UpdateColumns(update)
		if tx.Error == nil {
			c.JSON(200, p)
		} else {
			c.Status(500)
		}
	}
	if err != nil {
		c.Status(500)
		return
	}
	c.Status(404)
}

func deletePost(db *gorm.DB, c *gin.Context) {
	c.Status(204)
}

func makeHandler(db *gorm.DB, fn func(*gorm.DB, *gin.Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		fn(db, c)
	}
}
