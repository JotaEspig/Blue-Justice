package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sigma/config"
	"sigma/models/classroom"
	"sigma/server"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddClassroom(t *testing.T) {
	router := server.CreateTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login",
		bytes.NewBuffer([]byte(`username=admin&password=admin`)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// Gets json from writer
	jsonResp := map[string]string{}
	json.Unmarshal(w.Body.Bytes(), &jsonResp)

	token, ok := jsonResp["token"]
	assert.Equal(t, true, ok)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/admin/tools/classroom/add",
		bytes.NewBuffer([]byte(`{"name": "test", "year": 2022}`)))
	req.AddCookie(&http.Cookie{
		Name:  "auth",
		Value: token,
	})

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	err := config.DB.Unscoped().Delete(&classroom.Classroom{}, "name = ?", "test").Error
	assert.Equal(t, nil, err)
}

func TestGetClassroom(t *testing.T) {
	_ = server.CreateRouter()

	// Test GetClassroom endpoint with http.NewRequest
	t.Skip("Skipping GetClassroom endpoint test, because it's not implemented yet")
}

func TestGetAllClassrooms(t *testing.T) {
	_ = server.CreateRouter()

	// Test GetAllClassrooms endpoint with http.NewRequest
	t.Skip("Skipping GetAllClassrooms endpoint test, because it's not implemented yet")
}
