package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"mangahub/internal/api-server/handlers"
	"mangahub/internal/api-server/mangas"
	"mangahub/internal/api-server/routes"
	"mangahub/internal/api-server/users"
	"mangahub/internal/auth"
	"mangahub/pkg/utils/database"

	"mangahub/pkg/models"

	"github.com/gin-gonic/gin"
)

// ------------------------------
// MOCK CONFIG STRUCT
// ------------------------------
type MockConfig struct {
	API_PORT   string
	DB_PATH    string
	JWT_SECRET string
}

// Returns a new mock config with in-memory DB per test
func NewMockConfig() *MockConfig {
	return &MockConfig{
		API_PORT:   "8081",
		DB_PATH:    ":memory:", // in-memory DB
		JWT_SECRET: "test-secret",
	}
}

// ------------------------------
// SEED MANGA
// ------------------------------
func seedManga(t *testing.T, db *database.Database) {
	var count int64
	if err := db.DB().Model(&models.Manga{}).Where("id = ?", "one_piece").Count(&count).Error; err != nil {
		t.Fatalf("Failed to check manga existence: %v", err)
	}

	if count > 0 {
		// already seeded, skip
		return
	}

	m := models.Manga{
		ID:            "one_piece",
		Title:         "One Piece",
		Author:        "Eiichiro Oda",
		Genres:        "Action Adventure",
		Status:        "ONGOING",
		TotalChapters: 1100,
		Description:   "A young pirate's adventure...",
	}
	if err := db.DB().Create(&m).Error; err != nil {
		t.Fatalf("Failed to seed manga: %v", err)
	}
}

// ------------------------------
// TEST ROUTER SETUP
// ------------------------------
func SetupTestRouter(t *testing.T, cfg *MockConfig) (*gin.Engine, *database.Database) {
	gin.SetMode(gin.TestMode)

	// ------------------------------
	// 1. Initialize in-memory test database
	// ------------------------------
	dbConnector, err := database.NewDatabaseConnection(cfg.DB_PATH)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	db := dbConnector.(*database.Database)
	sqlConn := db.SQLDB()
	t.Cleanup(func() {
		sqlConn.Close()
	})

	// ------------------------------
	// 2. Repositories
	// ------------------------------
	userRepo := users.NewRepository(db)
	mangaRepo := mangas.NewRepository(db)

	// ------------------------------
	// 3. Utilities
	// ------------------------------
	jwtManager, err := auth.NewJWTManager(cfg.JWT_SECRET, "1h")
	if err != nil {
		t.Fatalf("Failed to create JWT manager: %v", err)
	}
	passwordHasher := auth.NewPasswordUtility()

	// ------------------------------
	// 4. Services
	// ------------------------------
	authService := auth.NewService(userRepo, jwtManager, passwordHasher)
	userService := users.NewService(userRepo)
	mangaService := mangas.NewService(mangaRepo)

	// ------------------------------
	// 5. Handlers
	// ------------------------------
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	mangaHandler := handlers.NewMangaHandler(mangaService)

	// ------------------------------
	// 6. Router setup
	// ------------------------------
	router := gin.Default()
	apiV1 := router.Group("/api/v1")

	apiV1.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "good"})
	})

	routes.AuthRoutesV1(apiV1.Group("/auth"), authHandler)
	routes.MangaRoutesV1(apiV1.Group("/manga"), mangaHandler)

	userGroup := apiV1.Group("/users")
	userGroup.Use(auth.AuthMiddleware(jwtManager))
	routes.UserRoutesV1(userGroup, userHandler)

	// Auto-migrate models
	if err := db.DB().AutoMigrate(&models.User{}, &models.Manga{}, &models.UserLibrary{}); err != nil {
		t.Fatalf("Failed to migrate test DB: %v", err)
	}

	// Seed manga table for library tests
	seedManga(t, db)

	return router, db
}

// ------------------------------
// HELPER: REGISTER AND LOGIN
// ------------------------------
func registerAndLogin(t *testing.T, router *gin.Engine) string {
	// Generate a unique username each time
	username := fmt.Sprintf("user_%d", time.Now().UnixNano())
	password := "password123"

	body := fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 201 && w.Code != 200 {
		t.Fatalf("Register failed: %v", w.Code)
	}

	// Login
	req2 := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer([]byte(body)))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != 200 {
		t.Fatalf("Login failed: %v", w2.Code)
	}

	var resp map[string]interface{}
	_ = json.Unmarshal(w2.Body.Bytes(), &resp)

	return resp["access_token"].(string)
}

// ------------------------------
// AUTH TEST CASES
// ------------------------------
func TestRegisterSuccess(t *testing.T) {
	router, _ := SetupTestRouter(t, NewMockConfig())

	body := `{"username":"john","password":"password123"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 201 && w.Code != 200 {
		t.Errorf("Expected 201, got %d", w.Code)
	}
}

func TestRegisterDuplicate(t *testing.T) {
	router, _ := SetupTestRouter(t, NewMockConfig())

	body := `{"username":"dup","password":"password123"}`

	// first
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// second (should fail)
	req2 := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer([]byte(body)))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != 409 { // match actual API response
		t.Errorf("Expected 409 for duplicate registration, got %d", w2.Code)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	router, _ := SetupTestRouter(t, NewMockConfig())

	body := `{"username":"userx","password":"password123"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("POST", "/api/v1/auth/login",
		bytes.NewBuffer([]byte(`{"username":"userx","password":"wrong"}`)))
	req2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w2, req2)

	if w2.Code != 400 {
		t.Errorf("Expected 400 invalid credentials, got %d", w2.Code)
	}
}

// ------------------------------
// JWT MIDDLEWARE TEST CASES
// ------------------------------
func TestProtectedRouteSuccess(t *testing.T) {
	router, _ := SetupTestRouter(t, NewMockConfig())
	token := registerAndLogin(t, router)

	req := httptest.NewRequest("GET", "/api/v1/users/library", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestProtectedRouteUnauthorized(t *testing.T) {
	router, _ := SetupTestRouter(t, NewMockConfig())

	req := httptest.NewRequest("GET", "/api/v1/users/library", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("Expected 401, got %d", w.Code)
	}
}

// ------------------------------
// MANGA TEST CASES
// ------------------------------
func TestGetMangaList(t *testing.T) {
	router, _ := SetupTestRouter(t, NewMockConfig())

	req := httptest.NewRequest("GET", "/api/v1/manga", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestGetMangaDetailNotFound(t *testing.T) {
	router, _ := SetupTestRouter(t, NewMockConfig())

	req := httptest.NewRequest("GET", "/api/v1/manga/unknown_id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Expected 404, got %d", w.Code)
	}
}

// ------------------------------
// LIBRARY TEST CASES
// ------------------------------
func TestAddToLibrary(t *testing.T) {
	router, db := SetupTestRouter(t, NewMockConfig())
	_ = db // db already seeded with one_piece
	token := registerAndLogin(t, router)

	body := `{"manga_id":"one_piece","current_chapter":10}`
	req := httptest.NewRequest("POST", "/api/v1/users/library", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestUpdateLibraryProgress(t *testing.T) {
	router, db := SetupTestRouter(t, NewMockConfig())
	_ = db
	token := registerAndLogin(t, router)

	body := `{"manga_id":"one_piece","current_chapter":1}`
	req1 := httptest.NewRequest("POST", "/api/v1/users/library", bytes.NewBuffer([]byte(body)))
	req1.Header.Set("Authorization", "Bearer "+token)
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	update := `{"manga_id":"one_piece","current_chapter":5}`
	req := httptest.NewRequest("PUT", "/api/v1/users/library", bytes.NewBuffer([]byte(update)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}
