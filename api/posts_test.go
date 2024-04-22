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

func createJsonRequest(method string, path string, requestObj interface{}) (*http.Request, error) {
	bodyBytes, _ := json.Marshal(requestObj)
	if req, err := http.NewRequest(method, path, bytes.NewBuffer(bodyBytes)); err != nil {
		return nil, err
	} else {
		req.Header.Set("Content-Type", "application/json")
		return req, nil
	}
}

func readPostsResponse(bytes []byte) *PostsResponse {
	responseBody := &PostsResponse{}
	json.Unmarshal(bytes, responseBody)
	return responseBody
}

func readPost(bytes []byte) *post {
	responseBody := &post{}
	json.Unmarshal(bytes, responseBody)
	return responseBody
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
		postsResponse := readPostsResponse(recorder.Body.Bytes())
		if len(postsResponse.Posts) < 1 {
			t.Error("Expected at least 1 post, got 0 ", postsResponse.Posts)
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
		post := readPost(recorder.Body.Bytes())
		if len(post.Message) < 1 {
			t.Error("Expected post with a message, message is empty ", post.Message)
		}
		if post.ID < 1 {
			t.Error("Expected post with an ID, ID is 0 ", post.Message)
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
	req, _ := createJsonRequest("POST", rPath, rPost)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)
	post := readPost(recorder.Body.Bytes())
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

func TestUpdatePost(t *testing.T) {
	rPath := "/posts/:id"
	router := gin.Default()
	router.PUT(rPath, makeHandler(testDb, updatePost))
	testId := uint(2000)
	testPost := post{
		ID:      testId,
		Message: "To Be Updated",
	}
	testDb.Model(testPost).Create(&testPost)
	defer testDb.Model(testPost).Delete(testPost)

	updateMessage := "I am updated!"
	updatedPost := post{
		ID:      testId,
		Message: updateMessage,
	}
	req, _ := createJsonRequest("PUT", fmt.Sprintf("/posts/%d", testId), updatedPost)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	t.Run("Returns 200 status code", func(t *testing.T) {
		if recorder.Code != 200 {
			t.Error("Expected 200, got ", recorder.Code)
		}
	})
	t.Run("updated post is returned", func(t *testing.T) {
		responsePost := readPost(recorder.Body.Bytes())
		if responsePost.Message != updateMessage {
			t.Error("expected message to be updated in database, got ", responsePost.Message)
		}
	})
	t.Run("Post is updated", func(t *testing.T) {
		toBeUpdated := post{ID: testId}
		testDb.Model(post{}).Find(&toBeUpdated)
		if toBeUpdated.Message != updateMessage {
			t.Error("expected message to be updated in database, got ", toBeUpdated.Message)
		}
	})
}
