package application

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/example/habitica-monitor/internal/domain/stats"
	"github.com/example/habitica-monitor/internal/domain/user"
	"github.com/example/habitica-monitor/internal/infrastructure/habitica"
)

// --- моки доменных портов ---

type fakeUserRepo struct{ users []*user.User }

func (f *fakeUserRepo) Save(context.Context, *user.User) error { return nil }
func (f *fakeUserRepo) GetByID(context.Context, string) (*user.User, error) {
	return nil, user.ErrNotFound
}
func (f *fakeUserRepo) List(context.Context) ([]*user.User, error) { return f.users, nil }

type fakeStatsRepo struct {
	mu    sync.Mutex
	added []*stats.Snapshot
}

func (f *fakeStatsRepo) Add(_ context.Context, s *stats.Snapshot) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.added = append(f.added, s)
	return nil
}
func (f *fakeStatsRepo) ListByUser(context.Context, string, int) ([]*stats.Snapshot, error) {
	return f.added, nil
}

type fakeFetcher struct{}

func (fakeFetcher) GetUserStats(context.Context, string, string) (habitica.Stats, error) {
	return habitica.Stats{HP: 42, MP: 10, Exp: 100, GP: 50, Lvl: 7}, nil
}

// runOnce должен сохранить по одному снапшоту на каждого пользователя.
func TestSnapshotRunOnce(t *testing.T) {
	users := &fakeUserRepo{users: []*user.User{
		{ID: "u1", APIKey: "k1", Name: "alice"},
		{ID: "u2", APIKey: "k2", Name: "bob"},
	}}
	history := &fakeStatsRepo{}

	svc := NewSnapshotService(users, history, fakeFetcher{}, time.Hour)
	svc.runOnce(context.Background())

	if len(history.added) != 2 {
		t.Fatalf("ожидалось 2 снапшота, получено %d", len(history.added))
	}
	for _, s := range history.added {
		if s.HP != 42 || s.Lvl != 7 {
			t.Fatalf("неверные данные снапшота: %+v", s)
		}
	}
}
