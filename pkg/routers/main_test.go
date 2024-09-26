package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupRouter initializes the Gin router and sets up routes for testing.
func setupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/api/v1/packages/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
	})
	return router
}

// TestInitRouter tests the routes defined in setupRouter.
func TestInitRouter(t *testing.T) {
	router := setupRouter()

	testCases := []struct {
		description    string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{"BaseRoute", "GET", "/", http.StatusNotFound, ""},
		{"APIv1Route", "GET", "/api/v1", http.StatusNotFound, ""},
		{"PackagesRoute", "GET", "/api/v1/packages/", http.StatusOK, `{"data":[]}`},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.path, nil)
			router.ServeHTTP(recorder, req)
			assert.Equal(t, tc.expectedStatus, recorder.Code)

			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, recorder.Body.String())
			}
		})
	}
}
