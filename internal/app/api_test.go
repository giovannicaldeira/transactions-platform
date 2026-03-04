package app

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	tests := []struct {
		name         string
		portEnv      string
		expectedPort string
	}{
		{
			name:         "default port when no env set",
			portEnv:      "",
			expectedPort: "8080",
		},
		{
			name:         "custom port from env",
			portEnv:      "3000",
			expectedPort: "3000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.portEnv != "" {
				os.Setenv("PORT", tt.portEnv)
				defer os.Unsetenv("PORT")
			} else {
				os.Unsetenv("PORT")
			}

			ctx := context.Background()
			api, err := Build(ctx)

			require.NoError(t, err)
			require.NotNil(t, api)
			assert.Equal(t, tt.expectedPort, api.port)
			assert.NotNil(t, api.server)
			assert.NotNil(t, api.router)

			// Clean up
			api.Close(ctx)
		})
	}
}

func TestAPI_Close(t *testing.T) {
	ctx := context.Background()
	api, err := Build(ctx)
	require.NoError(t, err)

	err = api.Close(ctx)
	assert.NoError(t, err)
}

func TestAPI_ServerConfiguration(t *testing.T) {
	ctx := context.Background()
	api, err := Build(ctx)
	require.NoError(t, err)
	defer api.Close(ctx)

	assert.Equal(t, 15*time.Second, api.server.ReadTimeout)
	assert.Equal(t, 15*time.Second, api.server.WriteTimeout)
	assert.Equal(t, 60*time.Second, api.server.IdleTimeout)
}

func TestAPI_RunAndShutdown(t *testing.T) {
	os.Setenv("PORT", "0") // Use random available port
	defer os.Unsetenv("PORT")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	api, err := Build(ctx)
	require.NoError(t, err)
	defer api.Close(ctx)

	// Run server in goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- api.Run(ctx)
	}()

	// Wait for context timeout (simulates graceful shutdown)
	select {
	case err := <-errChan:
		assert.NoError(t, err)
	case <-time.After(3 * time.Second):
		t.Fatal("test timeout")
	}
}

func TestAPI_HealthEndpoint(t *testing.T) {
	ctx := context.Background()
	api, err := Build(ctx)
	require.NoError(t, err)
	defer api.Close(ctx)

	// Verify the router is configured with routes
	assert.NotNil(t, api.router)
	assert.NotEmpty(t, api.router.Routes())
}
