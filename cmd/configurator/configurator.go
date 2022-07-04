package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	logger := log.StandardLogger()
	fakeConfig := []byte(`
# mock config already converted
remote_write:
- name: test
  url: https://fake-endpoint:8000
`)
	if err := os.WriteFile("./etc/prometheus/config/config.yaml", fakeConfig, 0444); err != nil {
		os.Exit(1)
		logger.Fatalf("fail to create config file: %s", err)
	}

	logger.Printf("config created")
}
