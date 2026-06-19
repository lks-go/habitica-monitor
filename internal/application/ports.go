// Package application — слой сценариев использования (use cases).
// Оркеструет доменные порты и внешний Habitica-клиент.
package application

import (
	"context"

	"github.com/example/habitica-monitor/internal/infrastructure/habitica"
)

// StatsFetcher — выходной порт получения статов из Habitica.
// Реализуется *habitica.Client; вынесен в интерфейс для тестируемости.
type StatsFetcher interface {
	GetUserStats(ctx context.Context, apiUser, apiKey string) (habitica.Stats, error)
}
