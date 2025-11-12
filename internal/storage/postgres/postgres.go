package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/MuthuM3/gin-microservice-template/internal/config"
)

type Store struct {
	db        *sql.DB
	authStore *AuthStore
	todoStore *TodoStore
	config    *config.DatabaseConfig
	logger    *log.Logger

	// Connection Monitoring
	mu              sync.RWMutex
	lastHealthCheck time.Time
	isHealthy       bool
	stats           ConnectionStats
}

type ConnectionStats struct {
	OpenConnections       int
	InUseConnections      int
	IdleConnection        int
	WaitCount             int
	WaitDuration          time.Duration
	MaxIdleClosed         int64
	MaxIdleTimeClosed     int64
	MaxIdleLifeTimeClosed int64
}

func newStore(connectionsString string, cfg *config.DatabaseConfig, logger *log.Logger) (*Store, error) {
	db, err := sql.Open("postgres", connectionsString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Apply configuration settings
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Test the connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()

		return nil, fmt.Errorf("failed to pink database: %w", err)
	}

	store := &Store{
		db:              db,
		authStore:       NewAuthStore(db),
		todoStore:       newTodoStore(db),
		config:          cfg,
		logger:          logger,
		isHealthy:       true,
		lastHealthCheck: time.Now(),
	}

	return store, nil
}

func (s *Store) startConnectionMonitoring() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {}
	}

}

func (s *Store) monitorConnections() {
	stats := s.GetStats()

	s.logger.Printf("DB stats: Open=%d, InUse=%d, Idle=%d, WaitCount=%d, WaitDuration=%v", stats.OpenConnections, stats.InUseConnections, stats.IdleConnection, stats.WaitCount, stats.WaitDuration)
}

func (s *Store) GetStats() ConnectionStats {
	s.mu.RLock()
	defer s.mu.Unlock()

	dbStats := s.db.Stats()

	return ConnectionStats{
		OpenConnections:       dbStats.OpenConnections,
		InUseConnections:      dbStats.InUse,
		IdleConnection:        dbStats.Idle,
		WaitCount:             int(dbStats.WaitCount),
		WaitDuration:          dbStats.WaitDuration,
		MaxIdleClosed:         dbStats.MaxIdleClosed,
		MaxIdleTimeClosed:     dbStats.MaxIdleTimeClosed,
		MaxIdleLifeTimeClosed: dbStats.MaxLifetimeClosed,
	}
}
