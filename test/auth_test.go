package test

import (
	"encoding/json"
	"fmt"
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

func TestMain(m *testing.M) {
	dbPath := fmt.Sprintf("%s/gin-gorm-demo-test.db", os.TempDir())
	_ = os.Remove(dbPath)

	if err := database.Init(dbPath); err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	code := m.Run()
	_ = os.Remove(dbPath)
	os.Exit(code)
}

func TestLogin(t *testing.T) {
	// 使用httptest测试登陆接口
	r := gin.Default()
	r.POST("/users", handler.CreateUser)
	r.POST("/login", handler.Login)

	createTestUser(t, r, "login_test_user", "123456")
	token := loginAndGetToken(t, r, "login_test_user", "123456")
	if token == "" {
		t.Fatal("Expected non-empty token")
	}
}

type apiResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"` // ⭐ 关键
}

func TestGetUsers(t *testing.T) {
	// 使用httptest测试获取用户列表接口
	r := gin.Default()
	r.POST("/users", handler.CreateUser)
	r.POST("/login", handler.Login)

	createTestUser(t, r, "list_test_user", "123456")
	token := loginAndGetToken(t, r, "list_test_user", "123456")

	r.Use(middleware.AuthMiddleware())
	r.GET("/users", handler.GetAuthorizedUsers)

	req := httptest.NewRequest("GET", "/users?page=1&size=10&sort=id&order=desc", nil)
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
	var userList dto.UserListResponse
	if err := json.Unmarshal(wrap.Data, &userList); err != nil {
		t.Fatalf("Failed to parse response data: %v", err)
	}
	if userList.Page != 1 || userList.Size != 10 || userList.Sort != "id" || userList.Order != "desc" {
		t.Fatalf("Unexpected pagination data: %+v", userList)
	}

	found := false
	for _, it := range userList.Items {
		t.Logf("User %s found", it.Name)
		if it.Name == "list_test_user" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("User %s not found", "list_test_user")
	}

}

func createTestUser(t *testing.T, r *gin.Engine, name string, password string) {
	t.Helper()

	payload := dto.CreateUserRequest{
		Name:     name,
		Age:      18,
		Password: password,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal create user payload: %v", err)
	}

	req := httptest.NewRequest("POST", "/users", strings.NewReader(string(jsonPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected create user status code %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}
}

func loginAndGetToken(t *testing.T, r *gin.Engine, name string, password string) string {
	t.Helper()

	payload := dto.LoginRequest{
		Name:     name,
		Password: password,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal login payload: %v", err)
	}

	req := httptest.NewRequest("POST", "/login", strings.NewReader(string(jsonPayload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected login status code %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}

	var wrap apiResponse
	if err := json.Unmarshal(w.Body.Bytes(), &wrap); err != nil {
		t.Fatalf("Failed to parse login response body: %v", err)
	}

	var resp dto.LoginResponse
	if err := json.Unmarshal(wrap.Data, &resp); err != nil {
		t.Fatalf("Failed to parse login response data: %v", err)
	}

	return resp.Token
}
