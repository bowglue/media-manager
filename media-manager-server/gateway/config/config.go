package config

// Config stores configuration data (e.g., ports, service addresses)
type Config struct {
	HTTPPort int
	GRPCPort int
	UserServiceURL string
	MediaServiceURL string
}

// LoadConfig loads configuration (this could be from a file or environment variables)
func LoadConfig() *Config {
	// In this example, we're using hardcoded values.
	return &Config{
		HTTPPort: 8080,
		GRPCPort: 9090,
		UserServiceURL: "localhost:9091", // Example URL for user service
		MediaServiceURL: "localhost:9092", // Example URL for media service
	}
}
