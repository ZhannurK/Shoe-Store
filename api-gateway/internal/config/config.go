package config

import "os"

type Config struct {
	Port           string
	TransactionURL string
	AuthURL        string
	InventoryURL   string
	JwtSecret      string
}

func Load() Config {
	cfg := Config{
		Port:           getEnv("GATEWAY_PORT", "8181"),
		TransactionURL: getEnv("TRANSACTION_SERVICE_URL", "http://transaction-service:8088"),
		AuthURL:        getEnv("AUTH_SERVICE_URL", "http://auth-service:8087"),
		InventoryURL:   getEnv("INVENTORY_SERVICE_URL", "http://inventory:8082"),
		JwtSecret:      getEnv("JWT_SECRET", "supersecret"),
	}
	return cfg
}

const (
	AuthServiceURL = "http://auth-service:8081"
)

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
