package rest

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/xcheng85/session-monitor-k8s/internal/test"
	"github.com/xcheng85/session-monitor-k8s/k8s/internal/handler"
)

func TestNewK8sRouter_RegisterLivenessProbe(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	mux := chi.NewRouter()
	mockK8sHandler := &handler.MockIK8sHandler{}
	mockK8sHandler.On("GetLivenessProbe", mock.Anything, mock.Anything).Return().Once()
	k8sHandler := NewK8sRouter(mockK8sHandler, ctx, mux)
	k8sHandler.Register()
	ts := httptest.NewServer(mux)
	defer ts.Close()
	if _, body := test.TestRequest(t, ts, "GET", "/livenessProbe", nil); body != "" {
		t.Fatalf(body)
	}
}

func TestNewK8sRouter_RegisterReadinessProbe(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	mux := chi.NewRouter()
	mockK8sHandler := &handler.MockIK8sHandler{}
	mockK8sHandler.On("GetReadinessProbe", mock.Anything, mock.Anything).Return().Once()
	k8sHandler := NewK8sRouter(mockK8sHandler, ctx, mux)
	k8sHandler.Register()
	ts := httptest.NewServer(mux)
	defer ts.Close()
	if _, body := test.TestRequest(t, ts, "GET", "/readinessProbe", nil); body != "" {
		t.Fatalf(body)
	}
}
