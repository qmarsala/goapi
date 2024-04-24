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
var api *gin.Engine

func seedDB(posts []Label) []Label {
	createdLabels := []Label{}
	for _, l := range posts {
		tx := testDb.Model(Label{}).Create(&l)
		createdLabels = append(createdLabels, l)
		if tx.Error != nil {
			fmt.Println(tx.Error)
		}
	}
	return createdLabels
}

func cleanupSeedDB(posts []Label) {
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

func readResponseBody[T Label | LabelsResponse](bytes []byte) *T {
	var responseBody T
	json.Unmarshal(bytes, &responseBody)
	return &responseBody
}

func TestMain(t *testing.M) {
	testDb = initializeDB[Label]("test")
	api = setupRoutes(testDb)
	posts := []Label{
		{ID: 1, Text: "Hello!", Target: "dog"},
		{ID: 2, Text: "Hello, Go!", Target: "cat"},
		{ID: 3, Text: "Hello, World!", Target: "bird"},
	}
	insertedPosts := seedDB(posts)
	code := t.Run()
	cleanupSeedDB(insertedPosts)
	os.Exit(code)
}

func TestGetLabels(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/labels", nil)
	recorder := httptest.NewRecorder()
	api.ServeHTTP(recorder, req)

	t.Run("Returns 200 status code", func(t *testing.T) {
		if recorder.Code != 200 {
			t.Error("Expected 200, got ", recorder.Code)
		}
	})

	t.Run("Returns list of posts", func(t *testing.T) {
		postsResponse := readResponseBody[LabelsResponse](recorder.Body.Bytes())
		if len(postsResponse.Labels) < 1 {
			t.Error("Expected at least 1 post, got 0 ", postsResponse.Labels)
		}
	})
}

func TestGetLabel(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/labels/%d", 1), nil)
	recorder := httptest.NewRecorder()
	api.ServeHTTP(recorder, req)

	t.Run("Returns 200 status code", func(t *testing.T) {
		if recorder.Code != 200 {
			t.Error("Expected 200, got ", recorder.Code)
		}
	})

	t.Run("Returns label", func(t *testing.T) {
		label := readResponseBody[Label](recorder.Body.Bytes())
		if len(label.Text) < 1 {
			t.Error("Expected label with a text, text is empty ", label.Text)
		}
		if len(label.Target) < 1 {
			t.Error("Expected label with a target, target is empty ", label.Target)
		}
		if label.ID < 1 {
			t.Error("Expected label with an ID, ID is 0 ", label.Text)
		}
	})
}

func TestCreateLabel(t *testing.T) {
	createLabelRequest := CreateLabelRequest{
		Text:   "Testing Create Label",
		Target: "building",
	}
	req, _ := createJsonRequest("POST", "/api/labels", createLabelRequest)
	recorder := httptest.NewRecorder()

	api.ServeHTTP(recorder, req)
	label := readResponseBody[Label](recorder.Body.Bytes())
	defer testDb.Model(label).Delete(label)

	t.Run("Returns 201 status code", func(t *testing.T) {
		if recorder.Code != 201 {
			t.Error("Expected 201, got ", recorder.Code)
		}
	})
	t.Run("Returns post", func(t *testing.T) {
		if label.Text != createLabelRequest.Text {
			t.Error("Expected text to match request, got ", label.Text)
		}
		if label.Target != createLabelRequest.Target {
			t.Error("Expected target to match request, got ", label.Target)
		}
		if label.ID < 1 {
			t.Error("Expected post with an ID, ID is 0 ", label.ID)
		}
	})
}

func TestCreateLabelWithBadInput(t *testing.T) {
	createLabelRequest := CreateLabelRequest{
		Text: "Testing Create Label",
	}
	req, _ := createJsonRequest("POST", "/api/labels", createLabelRequest)
	recorder := httptest.NewRecorder()

	api.ServeHTTP(recorder, req)
	label := readResponseBody[Label](recorder.Body.Bytes())
	defer testDb.Model(label).Delete(label)

	t.Run("Returns 400 status code", func(t *testing.T) {
		if recorder.Code != 400 {
			t.Error("Expected 400, got ", recorder.Code)
		}
	})
}

// func TestDeletePost(t *testing.T) {
// 	testId := uint(1000)
// 	testPost := Label{
// 		ID:   testId,
// 		Text: "To Be Deleted",
// 	}
// 	testDb.Model(testPost).Create(&testPost)
// 	delReq, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/posts/%d", testId), nil)
// 	delRecorder := httptest.NewRecorder()

// 	api.ServeHTTP(delRecorder, delReq)

// 	t.Run("Returns 204 status code", func(t *testing.T) {
// 		if delRecorder.Code != 204 {
// 			t.Error("Expected 204, got ", delRecorder.Code)
// 		}
// 	})
// 	t.Run("Post is deleted", func(t *testing.T) {
// 		getReq, _ := http.NewRequest("GET", fmt.Sprintf("/api/posts/%d", testId), nil)
// 		getRecorder := httptest.NewRecorder()
// 		api.ServeHTTP(getRecorder, getReq)
// 		if getRecorder.Code != 404 {
// 			t.Error("Expected 404, got ", getRecorder.Code)
// 		}
// 	})
// }

// func TestUpdatePost(t *testing.T) {
// 	testId := uint(2000)
// 	testPost := Label{
// 		ID:   testId,
// 		Text: "To Be Updated",
// 	}
// 	testDb.Model(testPost).Create(&testPost)
// 	defer testDb.Model(testPost).Delete(testPost)

// 	updateMessage := "I am updated!"
// 	updatedPost := Label{
// 		ID:   testId,
// 		Text: updateMessage,
// 	}
// 	updateReq, _ := createJsonRequest("PUT", fmt.Sprintf("/api/posts/%d", testId), updatedPost)
// 	updateRecorder := httptest.NewRecorder()
// 	api.ServeHTTP(updateRecorder, updateReq)

// 	t.Run("Returns 200 status code", func(t *testing.T) {
// 		if updateRecorder.Code != 200 {
// 			t.Error("Expected 200, got ", updateRecorder.Code)
// 		}
// 	})
// 	t.Run("updated post is returned", func(t *testing.T) {
// 		responsePost := readResponseBody[Label](updateRecorder.Body.Bytes())
// 		if responsePost.Text != updateMessage {
// 			t.Error("expected message to be updated in database, got ", responsePost.Text)
// 		}
// 	})
// 	t.Run("Post is updated", func(t *testing.T) {
// 		getReq, _ := http.NewRequest("GET", fmt.Sprintf("/api/posts/%d", testId), nil)
// 		getRecorder := httptest.NewRecorder()
// 		api.ServeHTTP(getRecorder, getReq)
// 		response := readResponseBody[Label](getRecorder.Body.Bytes())
// 		if response.Text != updateMessage {
// 			t.Error("expected message to be updated in database, got ", response.Text)
// 		}
// 		if response.ID != testId {
// 			t.Error("expected post to have correct ID, got ", response.ID)
// 		}
// 	})
// }
