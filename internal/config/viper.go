package config

import (
	"path/filepath" // go 1.21.3+
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type ViperConfig struct {
}

func (s *ViperConfig) Get(key string) any {
	return viper.Get(key)
}

func (s *ViperConfig) Set(key string, value any) {
	viper.Set(key, value)
}

func extractPath(path string) (dir string, filename string, filetype string) {
	dir, file := filepath.Split(path)
	splitR := strings.Split(file, ".")
	if len(splitR) == 1 {
		filetype = "json"
	} else {
		filetype = splitR[1]
	}
	filename = splitR[0]
	return dir, filename, filetype
}

// interface and implementation
func NewViperConfig(defautPath string, extraPaths []string, logger *zap.Logger) (cfg IConfig, err error) {
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()
	fdir, fname, ftype := extractPath(defautPath)
	viper.SetConfigType(ftype)
	viper.SetConfigName(fname)
	viper.AddConfigPath(fdir)
	viper.SetDefault("port", 8080)
	err = viper.ReadInConfig()
	viper.WatchConfig()
	if err != nil {
		logger.Sugar().Panicw("fatal error config file: %w", err)
	}
	for _, p := range extraPaths {
		fdir, fname, ftype := extractPath(p)
		v := viper.New()
		v.SetConfigName(fname)
		v.SetConfigType(ftype)
		v.AddConfigPath(fdir)
		err := v.ReadInConfig()
		if err != nil {
			logger.Sugar().Panicw("fatal error config file: %w", err)
		}
		viper.MergeConfigMap(v.AllSettings())
	}
	logger.Sugar().Info(viper.AllKeys())
	return &ViperConfig{}, nil
}
