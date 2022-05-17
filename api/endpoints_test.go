package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/GlobalWebIndex/platform2.0-go-challenge/db"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		panic("Failed to connect to database")
	}
	db.Migrate(database)
	return database
}

func performRequest(r http.Handler, method, path string, body io.ReadCloser) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestGetUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database := initDB()
	// Skip the auth process
	router := CreateTestEngine(database, false)

	w := performRequest(router, "GET", "/api/v1/users", nil)

	assert.Equal(t, 404, w.Code)
	db.FillDB(database)

	w = performRequest(router, "GET", "/api/v1/users", nil)
	assert.Equal(t, 200, w.Code)

	var got []gin.H
	want := []gin.H{{
		"id":       float64(1),
		"username": "test",
	}}
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(got)
	assert.Equal(t, want, got)
}

func TestPostUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database := initDB()
	// Skip the auth process
	router := CreateTestEngine(database, false)

	data, err := json.Marshal(gin.H{
		"username": "test_username",
		"password": "test_password",
	})
	assert.Equal(t, err, nil)
	body := io.NopCloser(bytes.NewBuffer(data))
	w := performRequest(router, "POST", "/api/v1/users", body)
	assert.Equal(t, 201, w.Code)

	w = performRequest(router, "GET", "/api/v1/users", nil)
	assert.Equal(t, 200, w.Code)

	var got []gin.H
	want := []gin.H{{
		"id":       float64(1),
		"username": "test_username",
	}}
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(got)
	assert.Equal(t, want, got)
}

func TestGetUserByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database := initDB()
	// Skip the auth process
	router := CreateTestEngine(database, false)

	data, err := json.Marshal(gin.H{
		"username": "test_username",
		"password": "test_password",
	})
	assert.Equal(t, err, nil)
	body := io.NopCloser(bytes.NewBuffer(data))
	w := performRequest(router, "POST", "/api/v1/users", body)
	assert.Equal(t, 201, w.Code)
	assert.Equal(t, 1, len(w.Result().Header["Location"]))

	// Verify the user details
	w = performRequest(router, "GET", w.Result().Header["Location"][0], nil)
	assert.Equal(t, 200, w.Code)

	var got gin.H
	want := gin.H{
		"id":       float64(1),
		"username": "test_username",
	}
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(got)
	assert.Equal(t, want, got)
}

// TODO: test all endpoints
