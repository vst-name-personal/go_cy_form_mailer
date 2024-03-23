package env_var

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Singleton struct {
	Data map[string]string
}

var instance *Singleton

func init() {
	instance = &Singleton{
		Data: make(map[string]string),
	}
	// Load environment variables into singleton data
	instance.LoadEnv()
	log.Println("loading env vars")
}

func GetInstance() *Singleton {
	return instance
}

// LoadEnv loads environment variables into the singleton
func (s *Singleton) LoadEnv() {
	// Required environment variables
	requiredVars := []string{
		"mail_server_domain",
		"mail_server_sender",
		"mail_server_receiver",
		"mail_server_passwd",
	}

	// Iterate over required environment variables and check their presence
	for _, v := range requiredVars {
		if value := os.Getenv(v); value != "" {
			s.Data[v] = value
		}
	}

	// If not set, load from .env
	if len(s.Data) != len(requiredVars) {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		log.Println("Running undockerized/dev")
	}
}
