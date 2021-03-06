package backend

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Server struct {
	UnimplementedBackendServer
	storage *PostgreStorage
}

type ServerConfig struct {
	ListenAddr  string `yaml:"listen_on"`
	DatabaseURI string `yaml:"db_uri"`

	MaxMessageSizeBytes int `yaml:"max_message_size_bytes"`

	// SentryDSN string `yaml:"sentry_dsn"`

	ServerCert string `yaml:"server_cert"`
	ServerKey  string `yaml:"server_key"`
	CACert     string `yaml:"ca_cert"`
}

func LoadConfig() ServerConfig {
	configPath := os.Getenv("BACKEND_CONFIG")
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf(
			"Error loading config from BACKEND_CONFIG=%s: %v",
			configPath,
			err,
		)
	}

	config := ServerConfig{}
	err = yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}
	return config
}

func NewServerFromConfig(cfg ServerConfig) *Server {
	storage, err := NewPostgreStorage(cfg.DatabaseURI)
	if err != nil {
		log.Fatalf("Could not create storage: %v", err)
	}

	return &Server{
		storage: storage,
	}
}

func (s *Server) Cleanup() {
	s.storage.db.Close()
}
