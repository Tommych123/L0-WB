package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Структура для считывания переменных окружения из env
type Config struct {
	DBHost      string
	DBPort      int
	DBUser      string
	DBPassword  string
	DBName      string
	KafkaBroker string
	ServicePort int
}

//Функция загрузки переменных окружения из env
func Load() *Config {
	godotenv.Load()
	return &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnvAsInt("DB_PORT", 5432),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "secret"),
		DBName:      getEnv("DB_NAME", "orders_db"),
		KafkaBroker: getEnv("KAFKA_BROKER", "localhost:9092"),
		ServicePort: getEnvAsInt("SERVICE_PORT", 8080),
	}
}

// Вспомогательная функция для получения строки из env или задания дефолтного значения
func getEnv(key string, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists && val != "" {
		return val
	}
	return defaultVal
}

// Вспомогательная функция для получения int из env или задания дефолтного значения
func getEnvAsInt(name string, defaultVal int) int {
	valStr := os.Getenv(name)
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}
