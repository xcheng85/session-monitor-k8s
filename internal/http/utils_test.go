package http

import (
	"errors"
	"net/http"
	"testing"

	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
)

func TestResponseRenderer(t *testing.T) {
	scenarios := []struct {
		desc               string
		responseRenderer   render.Renderer
		expectedStatusCode int
		expectedError      error
	}{
		{
			desc:               "400: bad request",
			responseRenderer:   ErrBadRequest(errors.New("caller error")),
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      nil,
		},
		{
			desc:               "401: unauthorized",
			responseRenderer:   ErrUnauthorized(errors.New("caller error")),
			expectedStatusCode: http.StatusUnauthorized,
			expectedError:      nil,
		},
		{
			desc:               "500: server internal error",
			responseRenderer:   ErrServerInternal(errors.New("caller error")),
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      nil,
		},
		{
			desc:               "200: ok",
			responseRenderer:   TextOkRender("caller ok"),
			expectedStatusCode: http.StatusOK,
			expectedError:      nil,
		},
	}

	for _, s := range scenarios {
		scenario := s
		t.Run(scenario.desc, func(t *testing.T) {
			errRenderable := scenario.responseRenderer.(*HttpResponse)
			assert.NotNil(t, errRenderable, "Should be struct: HttpResponse")
			assert.Equal(t, scenario.expectedStatusCode, errRenderable.HTTPStatusCode, "status code should be 400")
		})
	}
}
