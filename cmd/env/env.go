package env_var

import (
	"log"
	"os"
	"sync"
)

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
		log.Println("loading env vars")

	})
	return instance
}

// LoadEnv loads environment variables into singleton data
func (s *Singleton) LoadEnv() {

	s.Data["mail_server_domain"] = os.Getenv("mail_server_domain")
	s.Data["mail_server_sender"] = os.Getenv("mail_server_sender")
	s.Data["mail_server_receiver"] = os.Getenv("mail_server_receiver")
	s.Data["mail_server_passwd"] = os.Getenv("mail_server_passwd")


}
