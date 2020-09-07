package env

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../..")
)

type Config struct {
	WalletAddressChecksum int
}

func New() *Config {
	return &Config{
		WalletAddressChecksum: getEnvAsInt("WALLET_ADDRESS_CHECKSUM", 1),
	}
}

func GetEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(Root + "/.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := GetEnvVariable(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
