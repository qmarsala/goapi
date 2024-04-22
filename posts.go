package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostsResponse struct {
	Posts []post `json:"posts"`
}

type post struct {
	ID      uint   `json:"id" gorm:"primarykey"`
	Message string `json:"message"`
}

func parsePostId(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(err)
		return 0, err
	}
	return uint(id), nil
}

func getPostById(db *gorm.DB, id uint) (*post, error) {
	p := &post{}
	if tx := db.Find(&p, "ID = ?", uint(id)); tx.Error != nil {
		return nil, tx.Error
	} else if tx.RowsAffected < 1 {
		return nil, nil
	}
	return p, nil
}

func getPosts(db *gorm.DB, c *gin.Context) {
	posts := []post{}
	if tx := db.Limit(25).Find(&posts); tx.Error != nil {
		c.Status(500)
	} else {
		c.JSON(200, &PostsResponse{Posts: posts})
	}
}

func getPost(db *gorm.DB, c *gin.Context) {
	id, err := parsePostId(c)
	if err != nil {
		//todo: how to return better response
		c.Status(400)
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
	newPost := post{}
	c.Bind(&newPost)
	if tx := db.Model(&post{}).Create(&newPost); tx.Error == nil {
		c.JSON(201, newPost)
	} else {
		c.Status(500)
	}
}

func updatePost(db *gorm.DB, c *gin.Context) {
	id, err := parsePostId(c)
	if err != nil {
		c.Status(400)
		return
	}
	//todo: validate request
	update := &post{}
	c.Bind(update)
	p, err := getPostById(db, uint(id))
	switch {
	case p != nil:
		if tx := db.Model(p).UpdateColumns(update); tx.Error == nil {
			c.JSON(200, update)
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
	id, err := parsePostId(c)
	if err != nil {
		c.Status(404)
		return
	}

	p, _ := getPostById(db, id)
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
