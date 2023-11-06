package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestNewZapLogger(t *testing.T) {
	scenarios := []struct {
		desc                      string
		environment               string
		expectedDebugLevelEnabled bool
	}{
		{
			desc:                      "Zap Logger for Development",
			environment:               "development",
			expectedDebugLevelEnabled: true,
		},
		{
			desc:                      "Zap Logger for Production",
			environment:               "production",
			expectedDebugLevelEnabled: false,
		},
	}

	for _, s := range scenarios {
		scenario := s
		t.Run(scenario.desc, func(t *testing.T) {
			logger := NewZapLogger(LogConfig{
				Environment: scenario.environment,
			})
			assert.Equal(t, scenario.expectedDebugLevelEnabled, logger.Core().Enabled(zapcore.DebugLevel), "log level should match")
		})
	}
}
