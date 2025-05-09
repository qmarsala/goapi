package main

import (
	"github.com/gin-gonic/gin"
)

type LabelsResponse struct {
	Labels []Label `json:"labels"`
}

type GetLabelRequest struct {
	ID uint `uri:"id" binding:"required"`
}

type UpdateLabelRequest struct {
	ID     uint   `uri:"id" binding:"required"`
	Text   string `json:"text" binding:"required"`
	Target string `json:"target" binding:"required"`
}

type CreateLabelRequest struct {
	Text   string `json:"text"`
	Target string `json:"target" binding:"required"`
}

type Label struct {
	ID     uint   `json:"id" gorm:"primarykey"`
	Text   string `json:"text"`
	Target string `json:"target"`
}

func getPostById(db Database, id uint) (*Label, error) {
	p := Label{}
	if tx := db.Limit(1).Find(&p, id); tx.Error != nil {
		return nil, tx.Error
	} else if tx.RowsAffected < 1 {
		return nil, nil
	}
	return &p, nil
}

func getLabels(db Database, c *gin.Context) {
	posts := []Label{}
	if tx := db.Limit(25).Find(&posts); tx.Error != nil {
		c.Status(500)
	} else {
		c.JSON(200, &LabelsResponse{Labels: posts})
	}
}

func getLabel(db Database, c *gin.Context) {
	var getPostRequest GetLabelRequest
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

func createLabel(db Database, c *gin.Context) {
	var createPostRequest CreateLabelRequest
	if err := c.ShouldBind(&createPostRequest); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}

	newLabel := Label{
		Text:   createPostRequest.Text,
		Target: createPostRequest.Target,
	}
	if tx := db.Model(Label{}).Create(&newLabel); tx.Error == nil {
		c.JSON(201, newLabel)
	} else {
		c.Status(500)
	}
}

func updatePost(db Database, c *gin.Context) {
	var update UpdateLabelRequest
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

func deletePost(db Database, c *gin.Context) {
	var getPostRequest GetLabelRequest
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
