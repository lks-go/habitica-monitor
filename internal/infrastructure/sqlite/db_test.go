package sqlite

import (
	"context"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/example/habitica-monitor/internal/domain/stats"
	"github.com/example/habitica-monitor/internal/domain/user"
)

// TestConcurrentWritesNoLock воспроизводит проблему "database is locked":
// много горутин одновременно пишут снапшоты. С WAL + busy_timeout +
// SetMaxOpenConns(1) все записи должны пройти без ошибок SQLITE_BUSY.
func TestConcurrentWritesNoLock(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	userRepo := NewUserRepository(db)
	statsRepo := NewStatsRepository(db)

	// Готовим пользователей (внешний ключ требует существующих записей в user).
	const n = 20
	for i := 0; i < n; i++ {
		u := &user.User{ID: id(i), APIKey: "k", Name: "u"}
		if err := userRepo.Save(ctx, u); err != nil {
			t.Fatalf("save user: %v", err)
		}
	}

	// Параллельно пишем снапшоты — как делает SnapshotService.runOnce.
	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
	)
	now := time.Now()
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			snap, _ := stats.NewSnapshot(id(i), 1, 2, 3, 4, 5, now)
			if err := statsRepo.Add(ctx, snap); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()

	if len(errs) > 0 {
		t.Fatalf("ожидалось 0 ошибок при параллельной записи, получено %d, первая: %v", len(errs), errs[0])
	}

	// Все записи на месте.
	for i := 0; i < n; i++ {
		got, err := statsRepo.ListByUser(ctx, id(i), 0)
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("user %s: ожидался 1 снапшот, получено %d", id(i), len(got))
		}
	}
}

func id(i int) string {
	return "u-" + string(rune('a'+i%26)) + string(rune('0'+i/26))
}
