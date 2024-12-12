// BEGIN: 7f3a1g2h3i4j
package middleware

import (
	"ginapp/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahuigo/gohttptool/httpreq"
	"github.com/gin-gonic/gin"
)

func TestCors(t *testing.T) {
	req,_ := httpreq.R().SetFormData(map[string]string{
		"username": "alex",
	}).SetReq("POST","http://m/api/v1/mauth/login").ToRequest()

	resp, ctx := test.CreateTestCtx(req)
	CORS(ctx)

	// Check the response headers
	responseHeaders := resp.Header()
	allowedHeaders := responseHeaders.Get("Access-Control-Allow-Headers")

	if allowedHeaders != "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Traceparent" {
		t.Errorf("Expected Access-Control-Allow-Headers header to be 'Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Traceparent', got '%s'", allowedHeaders)
	}
}

func TestCORSWithServer(t *testing.T) {
	// Create a new Gin router
	router := gin.New()

	// Use the CORS middleware
	router.Use(CORS)

	// Define a test route
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Test route")
	})

	// Create a test request
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the Origin header
	req.Header.Set("Origin", "http://example.com")

	// Create a test response recorder
	recorder := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(recorder, req)

	// Check the response headers
	responseHeaders := recorder.Header()
	allowedOrigin := responseHeaders.Get("Access-Control-Allow-Origin")
	allowedMethods := responseHeaders.Get("Access-Control-Allow-Methods")
	allowedHeaders := responseHeaders.Get("Access-Control-Allow-Headers")
	allowCredentials := responseHeaders.Get("Access-Control-Allow-Credentials")

	// Assert the expected values
	if allowedOrigin != "http://example.com" {
		t.Errorf("Expected Access-Control-Allow-Origin header to be 'http://example.com', got '%s'", allowedOrigin)
	}

	if allowedMethods != "POST, OPTIONS, GET, PUT, DELETE, PATCH" {
		t.Errorf("Expected Access-Control-Allow-Methods header to be 'POST, OPTIONS, GET, PUT, DELETE, PATCH', got '%s'", allowedMethods)
	}

	if allowedHeaders != "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Traceparent" {
		t.Errorf("Expected Access-Control-Allow-Headers header to be 'Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Traceparent', got '%s'", allowedHeaders)
	}

	if allowCredentials != "true" {
		t.Errorf("Expected Access-Control-Allow-Credentials header to be 'true', got '%s'", allowCredentials)
	}

	// Check the response status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	// Check the response body
	expectedBody := "Test route"
	if recorder.Body.String() != expectedBody {
		t.Errorf("Expected response body to be '%s', got '%s'", expectedBody, recorder.Body.String())
	}
}

// END: 7f3a1g2h3i4j
