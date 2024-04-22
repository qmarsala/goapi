package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var testDb *gorm.DB

func seedDB(posts []post) []post {
	createdPosts := []post{}
	for _, p := range posts {
		tx := testDb.Model(post{}).Create(&p)
		createdPosts = append(createdPosts, p)
		if tx.Error != nil {
			fmt.Println(tx.Error)
		}
	}
	return createdPosts
}

func cleanupSeedDB(posts []post) {
	for _, p := range posts {
		tx := testDb.Model(p).Delete(p)
		if tx.Error != nil {
			fmt.Println(tx.Error)
		}
	}
}

func TestMain(t *testing.M) {
	testDb = connectDB("test")
	posts := []post{
		{ID: 1, Message: "Hello!"},
		{ID: 2, Message: "Hello, Go!"},
		{ID: 3, Message: "Hello, World!"},
	}
	insertedPosts := seedDB(posts)
	code := t.Run()
	cleanupSeedDB(insertedPosts)
	os.Exit(code)
}

func TestGetPosts(t *testing.T) {
	rPath := "/posts"
	router := gin.Default()
	router.GET(rPath, makeHandler(testDb, getPosts))

	req, _ := http.NewRequest("GET", rPath, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	t.Run("Returns 200 status code", func(t *testing.T) {
		if recorder.Code != 200 {
			t.Error("Expected 200, got ", recorder.Code)
		}
	})

	t.Run("Returns list of posts", func(t *testing.T) {
		var body PostsResponse
		json.Unmarshal(recorder.Body.Bytes(), &body)
		if len(body.Posts) < 1 {
			t.Error("Expected at least 1 post, got 0 ", body.Posts)
		}
	})
}

func TestGetPost(t *testing.T) {
	rPath := "/posts/:id"
	router := gin.Default()
	router.GET(rPath, makeHandler(testDb, getPost))

	req, _ := http.NewRequest("GET", fmt.Sprintf("/posts/%d", 1), nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	t.Run("Returns 200 status code", func(t *testing.T) {
		if recorder.Code != 200 {
			t.Error("Expected 200, got ", recorder.Code)
		}
	})

	t.Run("Returns post", func(t *testing.T) {
		var body post
		json.Unmarshal(recorder.Body.Bytes(), &body)
		if len(body.Message) < 1 {
			t.Error("Expected post with a message, message is empty ", body.Message)
		}
		if body.ID < 1 {
			t.Error("Expected post with an ID, ID is 0 ", body.Message)
		}
	})
}

func TestCreatePost(t *testing.T) {
	rPath := "/posts"
	router := gin.Default()
	router.POST(rPath, makeHandler(testDb, createPost))
	rPost := post{
		Message: "Testing Create Post",
	}
	bodyBytes, _ := json.Marshal(rPost)
	req, _ := http.NewRequest("POST", rPath, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)
	var body post
	json.Unmarshal(recorder.Body.Bytes(), &body)
	defer testDb.Model(body).Delete(body)

	t.Run("Returns 201 status code", func(t *testing.T) {
		if recorder.Code != 201 {
			t.Error("Expected 201, got ", recorder.Code)
		}
	})
	t.Run("Returns post", func(t *testing.T) {
		if len(body.Message) < 1 {
			t.Error("Expected post with a message, message is empty ", body.Message)
		}
		if body.ID < 1 {
			t.Error("Expected post with an ID, ID is 0 ", body.ID)
		}
	})
}

func TestDeletePost(t *testing.T) {
	rPath := "/posts/:id"
	router := gin.Default()
	router.DELETE(rPath, makeHandler(testDb, deletePost))
	testId := uint(1000)
	testPost := post{
		ID:      testId,
		Message: "To Be Deleted",
	}
	testDb.Model(testPost).Create(&testPost)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/posts/%d", testId), nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	t.Run("Returns 204 status code", func(t *testing.T) {
		if recorder.Code != 204 {
			t.Error("Expected 204, got ", recorder.Code)
		}
	})
	t.Run("Post is deleted", func(t *testing.T) {
		tx := testDb.Model(testPost).Find(&testPost)
		if tx.RowsAffected > 0 {
			t.Error("expected not to find post, but found post in db ")
			testDb.Model(testPost).Delete(testPost)
		}
	})
}
