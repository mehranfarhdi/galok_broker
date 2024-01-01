package enviroment

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

func Load(key string) (string, error) {
	// Try to load environment variables from the .env file
	err := godotenv.Load(".env")
	if err != nil {
		// .env file not found, log a message (optional)
		log.Println(".env file not found, reading from system environment variables.")
	}

	// Iterate over the system environment variables
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		envKey := pair[0]
		envValue := pair[1]

		// If the key matches, return the corresponding value
		if envKey == key {
			return envValue, nil
		}
	}

	// Key not found
	return "", fmt.Errorf("key not found: %s", key)
}
