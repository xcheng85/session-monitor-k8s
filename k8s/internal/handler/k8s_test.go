package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestK8sHandlerGetLivenessProbe(t *testing.T) {
	k8sHandler := NewK8sHandler()
	// perform request
	request, err := http.NewRequest("GET", "/livenessProbe", nil)
	require.NoError(t, err)
	response := httptest.NewRecorder()
	k8sHandler.GetLivenessProbe(response, request)
	require.Equal(t, 200, response.Code)
	payload, _ := io.ReadAll(response.Body)
	assert.Equal(t, "{\"status\":\"livenessProbe passes\"}\n", string(payload), "happy path")
}

func TestK8sHandlerGetReadinessProbe(t *testing.T) {
	k8sHandler := NewK8sHandler()
	// perform request
	request, err := http.NewRequest("GET", "/readinessProbe", nil)
	require.NoError(t, err)
	response := httptest.NewRecorder()
	k8sHandler.GetReadinessProbe(response, request)
	require.Equal(t, 200, response.Code)
	payload, _ := io.ReadAll(response.Body)
	assert.Equal(t, "{\"status\":\"readinessProbe passes\"}\n", string(payload), "happy path")
}
