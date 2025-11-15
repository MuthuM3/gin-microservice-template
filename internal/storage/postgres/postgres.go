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

	// Lifecycle management
	ctx    context.Context
	cancel context.CancelFunc
}

type ConnectionStats struct {
	OpenConnections   int
	InUseConnections  int
	IdleConnection    int
	WaitCount         int
	WaitDuration      time.Duration
	MaxIdleClosed     int64
	MaxIdleTimeClosed int64
	MaxLifeTimeClosed int64
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
	pintCtx, pingCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer pingCancel()

	if err := db.PingContext(pintCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to pink database: %w", err)
	}

	// Create context for lifecycle management
	ctx, cancel := context.WithCancel(context.Background())

	store := &Store{
		db:              db,
		config:          cfg,
		logger:          logger,
		isHealthy:       true,
		lastHealthCheck: time.Now(),
		ctx:             ctx,
		cancel:          cancel,
	}

	store.authStore = NewAuthStore(db, store)
	store.todoStore = newTodoStore(db, store)

	// Start connection monitoring
	go store.startConnectionMonitoring()
	logger.Printf("Database connection established with %d max open connections", cfg.MaxOpenConns)

	return store, nil
}

func (s *Store) startConnectionMonitoring() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.monitorConnections()
		}
	}
}

func (s *Store) monitorConnections() {
	stats := s.GetStats()

	s.logger.Printf("DB stats: Open=%d, InUse=%d, Idle=%d, WaitCount=%d, WaitDuration=%v",
		stats.OpenConnections,
		stats.InUseConnections,
		stats.IdleConnection,
		stats.WaitCount,
		stats.WaitDuration,
	)

	// Warn if connection usage is high
	maxConns := s.config.MaxOpenConns

	if stats.OpenConnections > int(float64(maxConns)*0.8) {
		s.logger.Printf(
			"WARNING: High connection usage: %d/%d (%.1f%%)",
			stats.OpenConnections,
			maxConns,
			float64(stats.OpenConnections)/float64(maxConns)*100,
		)
	}

	// Warn if wait times are high
	if stats.WaitDuration > time.Second {
		s.logger.Printf("WARNING: High connection wait time: %v", stats.WaitDuration)
	}

	// Perform periodic health check
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.HealthCheck(ctx); err != nil {
		s.logger.Printf("Periodic health check failed: %v", err)
	}
}

func (s *Store) GetStats() ConnectionStats {
	s.mu.RLock()
	defer s.mu.Unlock()

	dbStats := s.db.Stats()

	return ConnectionStats{
		OpenConnections:   dbStats.OpenConnections,
		InUseConnections:  dbStats.InUse,
		IdleConnection:    dbStats.Idle,
		WaitCount:         int(dbStats.WaitCount),
		WaitDuration:      dbStats.WaitDuration,
		MaxIdleClosed:     dbStats.MaxIdleClosed,
		MaxIdleTimeClosed: dbStats.MaxIdleTimeClosed,
		MaxLifeTimeClosed: dbStats.MaxLifetimeClosed,
	}
}

func (s *Store) HealthCheck(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	start := time.Now()
	err := s.db.PingContext(ctx)
	duration := time.Since(start)

	s.lastHealthCheck = time.Now()
	s.isHealthy = err == nil

	if err != nil {
		s.logger.Printf("Database health check failed (took %v): %v", duration, err)
		return fmt.Errorf("database health check failed: %w", err)
	}

	s.logger.Printf("Database health check passed (took %v)", duration)
	return nil
}

// IsHealthy returns the current health status
func (s *Store) IsHealthy() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isHealthy
}

// Close closes the database connection
func (s *Store) Close() error {
	s.logger.Println("Closing database connection...")
	
	// Cancel monitoring goroutine
	if s.cancel != nil {
		s.cancel()
	}

	return s.db.Close()
}

// DB returns the underlying database connection (for migrations, etc..)
func (s *Store) DB() *sql.DB {
	return s.db
}

// ExecuteWithRetry execute a function with retry logic for database operations
func (s *Store) ExecuteWithRetry(ctx context.Context, opertion func() error, maxRetries int) error {
	var lastErr error

	for attempt := 1; attempt < maxRetries; attempt++ {
		if err := opertion(); err != nil {
			lastErr = err
			s.logger.Printf("Database operation attemps %d failed: %v", attempt, err)

			if attempt < maxRetries {
				// Exponential backoff
				backOff := time.Duration(attempt) * 100 * time.Millisecond
				select {
				case <-time.After(backOff):
					continue
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		} else {
			if attempt > 1 {
				s.logger.Printf("Database operation succeeded on attempt %d", attempt)
			}
			return nil
		}
	}

	return fmt.Errorf("database operation failed after %d attempts: %w", maxRetries, lastErr)
}
