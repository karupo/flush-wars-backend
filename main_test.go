package main

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// TestInitApp tests the initialization of the Fiber app.
func TestInitApp(t *testing.T) {
	// Call InitApp to initialize the app
	app := InitApp()

	// Assert that the app is not nil and is of type *fiber.App
	assert.NotNil(t, app)
	assert.IsType(t, &fiber.App{}, app)
}

// TestHealthCheckRoute tests the health check route.
func TestHealthCheckRoute(t *testing.T) {
	// Initialize the app
	app := InitApp()

	// Perform a request to the /healthz route
	req := httptest.NewRequest("GET", "/healthz", nil)
	resp, err := app.Test(req)

	// Assert that the request was successful (200 OK)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

}
