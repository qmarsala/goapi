package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostsResponse struct {
	Posts []Post `json:"posts"`
}

type GetPostRequest struct {
	ID uint `uri:"id" binding:"required"`
}

type UpdatePostRequest struct {
	ID      uint   `uri:"id" binding:"required"`
	Message string `json:"message"`
}

type Post struct {
	ID      uint   `json:"id" gorm:"primarykey"`
	Message string `json:"message"`
}

func getPostById(db *gorm.DB, id uint) (*Post, error) {
	p := &Post{}
	if tx := db.Find(&p, "ID = ?", uint(id)); tx.Error != nil {
		return nil, tx.Error
	} else if tx.RowsAffected < 1 {
		return nil, nil
	}
	return p, nil
}

func getPosts(db *gorm.DB, c *gin.Context) {
	posts := []Post{}
	if tx := db.Limit(25).Find(&posts); tx.Error != nil {
		c.Status(500)
	} else {
		c.JSON(200, &PostsResponse{Posts: posts})
	}
}

func getPost(db *gorm.DB, c *gin.Context) {
	var getPostRequest GetPostRequest
	if err := c.ShouldBindUri(&getPostRequest); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}

	p, err := getPostById(db, getPostRequest.ID)
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
	newPost := Post{}
	c.Bind(&newPost)
	if tx := db.Model(&Post{}).Create(&newPost); tx.Error == nil {
		c.JSON(201, newPost)
	} else {
		c.Status(500)
	}
}

func updatePost(db *gorm.DB, c *gin.Context) {
	var update UpdatePostRequest
	if err := c.ShouldBind(&update); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	p, err := getPostById(db, update.ID)
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
	var getPostRequest GetPostRequest
	if err := c.ShouldBindUri(&getPostRequest); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}

	p, _ := getPostById(db, getPostRequest.ID)
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
