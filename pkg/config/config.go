package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func Init() {
	// viper
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	// get root path using go.mod as marker
	rootPath, err := findRootPath()
	if err != nil {
		log.Fatalf("Error finding root path: %s", err)
	}

	// add config path
	viper.AddConfigPath(rootPath)

	// read config file
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// db config
	err = viper.Unmarshal(&dbConfig)
	if err != nil {
		log.Fatalf("Error unmarshalling config, %s", err)
	}
}

// findRootPath finds the project root directory by looking for go.mod file
func findRootPath() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	// Keep going up until we find go.mod or hit the root
	for {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			return currentDir, nil
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			return "", fmt.Errorf("could not find go.mod in any parent directory")
		}
		currentDir = parent
	}
}
