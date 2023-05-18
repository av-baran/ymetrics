package config

import "time"

type StorageConfig struct {
	DatabaseDSN    string
	RequestTimeout time.Duration
}
