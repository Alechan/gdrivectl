package app

import "time"

type Config struct {
	GcloudBin string
	Timeout   time.Duration
	JSON      bool
	Debug     bool
}

func NewConfig(gcloudBin string, timeout time.Duration, json, debug bool) Config {
	return Config{GcloudBin: gcloudBin, Timeout: timeout, JSON: json, Debug: debug}
}
