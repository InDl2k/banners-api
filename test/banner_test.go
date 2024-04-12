package test

import (
	"banners/internal/database/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCode200UserBanner(t *testing.T) {
	admin := getToken("admin")
	user := getToken("user")
	createBanner(t, "/banner", admin)
	Code200UserBanner(t, "/user_banner?tag_id=1&feature_id=1", user)
	clearTable()
}

func Code200UserBanner(t *testing.T, url string, token string) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", token)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if len(m) == 0 {
		t.Errorf("Expected the 'json object', but Got 'empty'")
	}
	clearTable()
}

func TestCode400UserBanner(t *testing.T) {
	admin := getToken("admin")
	user := getToken("user")

	Code400UserBanner(t, "/user_banner", user)

	createBanner(t, "/banner", admin)

	Code400UserBanner(t, "/user_banner?tag_id=1", user)
	Code400UserBanner(t, "/user_banner?feature_id=1", user)
	Code400UserBanner(t, "/user_banner?tag_id=f&feature_id=1", user)
	Code400UserBanner(t, "/user_banner?tag_id=1&feature_id=f", user)
	clearTable()
}

func Code400UserBanner(t *testing.T, url string, token string) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", token)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Некорректные данные" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Некорректные данные'. Got '%s'", m["error"])
	}
	clearTable()
}

func TestCode401UserBanner(t *testing.T) {
	admin := getToken("admin")
	none := ""

	createBanner(t, "/banner", admin)
	Code401UserBanner(t, "/user_banner?tag_id=1&feature_id=1", none)
	clearTable()
}

func Code401UserBanner(t *testing.T, url string, token string) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", token)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
	clearTable()
}

func TestCode404UserBanner(t *testing.T) {
	admin := getToken("admin")
	user := getToken("user")

	createBanner(t, "/banner", admin)

	Code404UserBanner(t, "/user_banner?tag_id=1&feature_id=2", user)
	Code404UserBanner(t, "/user_banner?tag_id=2&feature_id=1", user)
	clearTable()
}

func Code404UserBanner(t *testing.T, url string, token string) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", token)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
	clearTable()
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func clearTable() {
	models.GetDB().Exec("DELETE FROM banners")
}

func getToken(role string) string {
	req, _ := http.NewRequest("GET", "/token/"+role, nil)
	response := executeRequest(req)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	return m["token"]
}

func createBanner(t *testing.T, url string, token string) {

	var json = []byte(`{
			"tag_ids": [
				1
			],
			"feature_id": 1,
			"content": {"title": "some_title", "text": "some_text", "url": "some_url"},
			"is_active": true
	}`)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
}
