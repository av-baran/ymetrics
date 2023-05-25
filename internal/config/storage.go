package config

import "time"

type StorageConfig struct {
	DatabaseDSN          string
	SingleRequestTimeout time.Duration
	BatchRequestTimeout  time.Duration
}
