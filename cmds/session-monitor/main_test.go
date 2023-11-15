package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func Test_PublicEndpoint_Integration(t *testing.T) {
	t.Setenv("CONFIG_PATH", "./config.yaml")
	ioc, err := newIocContainer()
	assert.Nil(t, err, "cannot create ioc container")
	err = ioc.container.Invoke(func(root *CompositionRoot) error {
		mux := root.mux
		ts := httptest.NewServer(mux)
		defer ts.Close()

		t.Run("it should return 200 for k8s livenessProbe", func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("%s/livenessProbe", ts.URL))

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			assert.Equal(t, 200, resp.StatusCode)
		})

		t.Run("it should return 200 for k8s readinessProbe", func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("%s/readinessProbe", ts.URL))

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			assert.Equal(t, 200, resp.StatusCode)
		})

		return nil
	})
}
