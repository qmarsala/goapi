package main

import (
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
	//todo: validate request
	newPost := &post{}
	c.Bind(newPost)
	tx := db.Model(&post{}).Create(newPost)
	if tx.Error == nil {
		c.JSON(201, newPost)
		return
	}
	c.Status(500)
}

func updatePost(db *gorm.DB, c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(err)
		return
	}
	//todo: validate request
	update := &post{}
	c.Bind(update)
	p, err := getPostById(db, uint(id))
	switch {
	case p != nil:
		if tx := db.Model(p).UpdateColumns(update); tx.Error == nil {
			c.JSON(200, p)
		} else {
			c.Status(500)
		}
	case err != nil:
		c.Status(500)
	default:
		c.Status(404)
	}
}

func deletePost(db *gorm.DB, c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	p, _ := getPostById(db, uint(id))
	switch {
	case p != nil:
		db.Model(p).Delete(p)
		c.Status(204)
	default:
		c.Status(404)
	}
}

func makeHandler(db *gorm.DB, fn func(*gorm.DB, *gin.Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		fn(db, c)
	}
}
