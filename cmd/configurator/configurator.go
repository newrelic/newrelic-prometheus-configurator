package main

import (
	"log"
	"os"
)

func main() {
	fakeConfig := []byte(`
# mock config already converted 
remote_write:
- name: test
  url: https://fake-endpoint:8000	
`)
	if err := os.WriteFile("./etc/prometheus/config/config.yaml", fakeConfig, 444); err != nil {
		os.Exit(1)
		log.Fatalf("fail to create config file: %s", err)
	}

	log.Printf("config created")
}
