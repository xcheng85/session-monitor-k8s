package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xcheng85/session-monitor-k8s/internal/logger"
)

func TestViperConfig_Get(t *testing.T) {
	t.Setenv("DEFAULT_PATH", "../../cmds/session-monitor/dummy.yaml")
	t.Setenv("CONFIG_PATH", "../../cmds/session-monitor/config.yaml")
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})
	viper, err := NewViperConfig(os.Getenv("DEFAULT_PATH"), []string{os.Getenv("CONFIG_PATH")}, logger)
	assert.Nil(t, err, "no error to create Viper Config")
	assert.NotNil(t, viper, "Viper Config is defined")
	assert.Equal(t, "evd-cia3dviz", viper.Get("app.pod_namespace"))
	assert.Equal(t, "127.0.0.1:6379", viper.Get("app.redis_address"))
	assert.Equal(t, false, viper.Get("app.redis_mock"))
	assert.Equal(t, "enqueue_session_test", viper.Get("app.enqueue_session_stream_key"))
}
