package config

//go:generate mockery --name IConfig
type IConfig interface {
	Get(key string) any
}
