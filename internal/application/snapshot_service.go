package application

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/example/habitica-monitor/internal/domain/stats"
	"github.com/example/habitica-monitor/internal/domain/user"
)

// SnapshotService периодически снимает статы всех пользователей и пишет
// каждую выборку отдельной записью в историю.
type SnapshotService struct {
	users    user.Repository
	history  stats.Repository
	fetcher  StatsFetcher
	interval time.Duration
	now      func() time.Time
}

// NewSnapshotService создаёт сервис. interval — настраиваемый период (по умолчанию 1 час).
func NewSnapshotService(
	users user.Repository,
	history stats.Repository,
	fetcher StatsFetcher,
	interval time.Duration,
) *SnapshotService {
	if interval <= 0 {
		interval = time.Hour
	}
	return &SnapshotService{
		users:    users,
		history:  history,
		fetcher:  fetcher,
		interval: interval,
		now:      time.Now,
	}
}

// Run запускает планировщик до отмены ctx. Делает первый прогон сразу,
// затем повторяет каждые interval. Блокирующий — запускайте в горутине.
func (s *SnapshotService) Run(ctx context.Context) {
	log.Printf("snapshot: запуск, интервал %s", s.interval)
	s.runOnce(ctx)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("snapshot: остановка")
			return
		case <-ticker.C:
			s.runOnce(ctx)
		}
	}
}

// runOnce снимает статы по всем пользователям. Каждый пользователь
// обрабатывается в ОТДЕЛЬНОЙ горутине; WaitGroup ждёт завершения цикла.
func (s *SnapshotService) runOnce(ctx context.Context) {
	users, err := s.users.List(ctx)
	if err != nil {
		log.Printf("snapshot: не удалось получить список пользователей: %v", err)
		return
	}

	var wg sync.WaitGroup
	for _, u := range users {
		wg.Add(1)
		go func(u *user.User) {
			defer wg.Done()
			s.snapshotUser(ctx, u)
		}(u)
	}
	wg.Wait()
}

func (s *SnapshotService) snapshotUser(ctx context.Context, u *user.User) {
	st, err := s.fetcher.GetUserStats(ctx, u.ID, u.APIKey)
	if err != nil {
		log.Printf("snapshot: пользователь %s (%s): ошибка получения статов: %v", u.Name, u.ID, err)
		return
	}

	snap, err := stats.NewSnapshot(u.ID, st.HP, st.MP, st.Exp, st.GP, st.Lvl, s.now())
	if err != nil {
		log.Printf("snapshot: пользователь %s: невалидный снапшот: %v", u.Name, err)
		return
	}
	if err := s.history.Add(ctx, snap); err != nil {
		log.Printf("snapshot: пользователь %s: не удалось сохранить: %v", u.Name, err)
		return
	}
	log.Printf("snapshot: пользователь %s — hp=%.1f mp=%.1f exp=%.1f gp=%.1f lvl=%d",
		u.Name, st.HP, st.MP, st.Exp, st.GP, st.Lvl)
}
