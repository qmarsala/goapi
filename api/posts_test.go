package main

// functions to create reqs with content type set correctly

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
var api *gin.Engine

func seedDB(posts []Post) []Post {
	createdPosts := []Post{}
	for _, p := range posts {
		tx := testDb.Model(Post{}).Create(&p)
		createdPosts = append(createdPosts, p)
		if tx.Error != nil {
			fmt.Println(tx.Error)
		}
	}
	return createdPosts
}

func cleanupSeedDB(posts []Post) {
	for _, p := range posts {
		tx := testDb.Model(p).Delete(p)
		if tx.Error != nil {
			fmt.Println(tx.Error)
		}
	}
}

func createJsonRequest(method string, path string, requestObj interface{}) (*http.Request, error) {
	bodyBytes, err := json.Marshal(requestObj)
	if err != nil {
		return nil, err
	}

	if req, err := http.NewRequest(method, path, bytes.NewBuffer(bodyBytes)); err != nil {
		return nil, err
	} else {
		req.Header.Set("Content-Type", "application/json")
		return req, nil
	}
}

func readResponseBody[T Post | PostsResponse](bytes []byte) *T {
	var responseBody T
	json.Unmarshal(bytes, &responseBody)
	return &responseBody
}

func TestMain(t *testing.M) {
	testDb = connectDB("test")
	api = setupRoutes(testDb)
	posts := []Post{
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
	req, _ := http.NewRequest("GET", "/api/posts", nil)
	recorder := httptest.NewRecorder()
	api.ServeHTTP(recorder, req)

	t.Run("Returns 200 status code", func(t *testing.T) {
		if recorder.Code != 200 {
			t.Error("Expected 200, got ", recorder.Code)
		}
	})

	t.Run("Returns list of posts", func(t *testing.T) {
		postsResponse := readResponseBody[PostsResponse](recorder.Body.Bytes())
		if len(postsResponse.Posts) < 1 {
			t.Error("Expected at least 1 post, got 0 ", postsResponse.Posts)
		}
	})
}

func TestGetPost(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/posts/%d", 1), nil)
	recorder := httptest.NewRecorder()
	api.ServeHTTP(recorder, req)

	t.Run("Returns 200 status code", func(t *testing.T) {
		if recorder.Code != 200 {
			t.Error("Expected 200, got ", recorder.Code)
		}
	})

	t.Run("Returns post", func(t *testing.T) {
		post := readResponseBody[Post](recorder.Body.Bytes())
		if len(post.Message) < 1 {
			t.Error("Expected post with a message, message is empty ", post.Message)
		}
		if post.ID < 1 {
			t.Error("Expected post with an ID, ID is 0 ", post.Message)
		}
	})
}

func TestCreatePost(t *testing.T) {
	rPost := Post{
		Message: "Testing Create Post",
	}
	req, _ := createJsonRequest("POST", "/api/posts", rPost)
	recorder := httptest.NewRecorder()

	api.ServeHTTP(recorder, req)
	post := readResponseBody[Post](recorder.Body.Bytes())
	defer testDb.Model(post).Delete(post)

	t.Run("Returns 201 status code", func(t *testing.T) {
		if recorder.Code != 201 {
			t.Error("Expected 201, got ", recorder.Code)
		}
	})
	t.Run("Returns post", func(t *testing.T) {
		if post.Message != rPost.Message {
			t.Error("Expected message to match request, got ", post.Message)
		}
		if post.ID < 1 {
			t.Error("Expected post with an ID, ID is 0 ", post.ID)
		}
	})
}

func TestDeletePost(t *testing.T) {
	testId := uint(1000)
	testPost := Post{
		ID:      testId,
		Message: "To Be Deleted",
	}
	testDb.Model(testPost).Create(&testPost)
	delReq, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/posts/%d", testId), nil)
	delRecorder := httptest.NewRecorder()

	api.ServeHTTP(delRecorder, delReq)

	t.Run("Returns 204 status code", func(t *testing.T) {
		if delRecorder.Code != 204 {
			t.Error("Expected 204, got ", delRecorder.Code)
		}
	})
	t.Run("Post is deleted", func(t *testing.T) {
		getReq, _ := http.NewRequest("GET", fmt.Sprintf("/api/posts/%d", testId), nil)
		getRecorder := httptest.NewRecorder()
		api.ServeHTTP(getRecorder, getReq)
		if getRecorder.Code != 404 {
			t.Error("Expected 404, got ", getRecorder.Code)
		}
	})
}

func TestUpdatePost(t *testing.T) {
	testId := uint(2000)
	testPost := Post{
		ID:      testId,
		Message: "To Be Updated",
	}
	testDb.Model(testPost).Create(&testPost)
	defer testDb.Model(testPost).Delete(testPost)

	updateMessage := "I am updated!"
	updatedPost := Post{
		ID:      testId,
		Message: updateMessage,
	}
	updateReq, _ := createJsonRequest("PUT", fmt.Sprintf("/api/posts/%d", testId), updatedPost)
	updateRecorder := httptest.NewRecorder()
	api.ServeHTTP(updateRecorder, updateReq)

	t.Run("Returns 200 status code", func(t *testing.T) {
		if updateRecorder.Code != 200 {
			t.Error("Expected 200, got ", updateRecorder.Code)
		}
	})
	t.Run("updated post is returned", func(t *testing.T) {
		responsePost := readResponseBody[Post](updateRecorder.Body.Bytes())
		if responsePost.Message != updateMessage {
			t.Error("expected message to be updated in database, got ", responsePost.Message)
		}
	})
	t.Run("Post is updated", func(t *testing.T) {
		getReq, _ := http.NewRequest("GET", fmt.Sprintf("/api/posts/%d", testId), nil)
		getRecorder := httptest.NewRecorder()
		api.ServeHTTP(getRecorder, getReq)
		response := readResponseBody[Post](getRecorder.Body.Bytes())
		if response.Message != updateMessage {
			t.Error("expected message to be updated in database, got ", response.Message)
		}
		if response.ID != testId {
			t.Error("expected post to have correct ID, got ", response.ID)
		}
	})
}
