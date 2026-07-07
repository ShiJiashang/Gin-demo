package test

import (
	"encoding/json"
	"gin_gorm_demo/database"
	"gin_gorm_demo/dto"
	"gin_gorm_demo/handler"
	"gin_gorm_demo/middleware"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

var token string

func TestMain(m *testing.M) {
	if err := database.Init("../gin-demo.db"); err != nil {
		log.Fatalf("failed to init database: %v", err)
	}
	os.Exit(m.Run())
}

func TestLogin(t *testing.T) {
	// 使用httptest测试登陆接口
	r := gin.Default()
	r.POST("/login", handler.Login)

	payload := dto.LoginRequest{
		Name:     "test",
		Password: "123456",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	req := httptest.NewRequest("POST", "/login", strings.NewReader(string(jsonPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}

	var wrap apiResponse
	if err := json.Unmarshal(w.Body.Bytes(), &wrap); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}
	var resp dto.LoginResponse
	if err := json.Unmarshal(wrap.Data, &resp); err != nil {
		t.Fatalf("Failed to parse response data: %v", err)
	}
	if resp.Token == "" {
		t.Fatalf("Expected non-empty token")
	}
	token = resp.Token
}

type apiResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"` // ⭐ 关键
}

func TestGetUsers(t *testing.T) {
	// 使用httptest测试获取用户列表接口
	r := gin.Default()
	r.Use(middleware.AuthMiddleware())
	r.GET("/users", handler.GetAuthorizedUsers)

	req := httptest.NewRequest("GET", "/users", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}

	var wrap apiResponse
	w.Body.Bytes()
	if err := json.Unmarshal(w.Body.Bytes(), &wrap); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}
	var users []dto.UserResponse
	if err := json.Unmarshal(wrap.Data, &users); err != nil {
		t.Fatalf("Failed to parse response data: %v", err)
	}

	found := false
	for _, it := range users {
		t.Logf("User %s found", it.Name)
		if it.Name == "test" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("User %s not found", "test")
	}

}
