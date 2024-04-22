package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetPosts(t *testing.T) {
	db := setupDB()
	rPath := "/posts"
	router := gin.Default()
	router.GET(rPath, makeHandler(db, getPosts))

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
	db := setupDB()
	rPath := "/posts/:id"
	router := gin.Default()
	router.GET(rPath, makeHandler(db, getPost))

	req, _ := http.NewRequest("GET", "/posts/1", nil)
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
