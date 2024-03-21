package env_var

import (
	"fmt"
	"os"
	"sync"
)

// Singleton struct to hold global data
type Singleton struct {
	Data map[string]string
}

var instance *Singleton
var once sync.Once

// GetInstance returns the singleton instance
func GetInstance() *Singleton {
	once.Do(func() {
		instance = &Singleton{
			Data: make(map[string]string),
		}

		// Load environment variables into singleton data
		instance.LoadEnv()
	})
	return instance
}

// LoadEnv loads environment variables into singleton data
func (s *Singleton) LoadEnv() {
	// Load environment variables into singleton data
	s.Data["ENV_VAR_NAME"] = os.Getenv("ENV_VAR_NAME")
	// Add more environment variables as needed
}

func env_var() {
	// Get singleton instance
	s := GetInstance()

	// Access data from singleton
	fmt.Println("ENV_VAR_NAME:", s.Data["ENV_VAR_NAME"])
}
