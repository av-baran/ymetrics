package config

import "time"

type StorageConfig struct {
	DatabaseDSN  string
	QueryTimeout time.Duration
}
