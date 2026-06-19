// Package configs — загрузка конфигурации из переменных окружения.
package configs

import (
	"os"
	"time"
)

// Config — настройки приложения.
type Config struct {
	HTTPAddr         string        // адрес HTTP-сервера, напр. ":8080"
	DBPath           string        // путь к файлу SQLite
	XClient          string        // заголовок x-client: "<your-user-id>-<appname>"
	SnapshotInterval time.Duration // период снятия снапшотов (настраиваемый)
	CORSOrigin       string        // разрешённый origin для Web UI ("*" по умолчанию)
}

// Load читает конфиг из окружения с разумными значениями по умолчанию.
func Load() Config {
	return Config{
		HTTPAddr:         getEnv("HTTP_ADDR", ":8080"),
		DBPath:           getEnv("DB_PATH", "monitor.db"),
		XClient:          getEnv("HABITICA_X_CLIENT", "REPLACE-WITH-YOUR-USER-ID-TeamMonitor"),
		SnapshotInterval: getEnvDuration("SNAPSHOT_INTERVAL", time.Hour),
		CORSOrigin:       getEnv("CORS_ORIGIN", "*"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
