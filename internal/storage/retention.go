package storage

import (
	"context"
	"log"
	"time"
)

type RetentionManager struct {
	logRepo   LogRepository
	traceRepo TraceRepository
	eventRepo EventRepository
	cacheRepo CacheRepository
	retentionDays int
	interval      time.Duration
	stopChan      chan struct{}
}

func NewRetentionManager(logRepo LogRepository, traceRepo TraceRepository, eventRepo EventRepository, cacheRepo CacheRepository, retentionDays int, interval time.Duration) *RetentionManager {
	return &RetentionManager{
		logRepo:      logRepo,
		traceRepo:    traceRepo,
		eventRepo:    eventRepo,
		cacheRepo:    cacheRepo,
		retentionDays: retentionDays,
		interval:     interval,
		stopChan:     make(chan struct{}),
	}
}

func (r *RetentionManager) Start(ctx context.Context) {
	go r.run(ctx)
}

func (r *RetentionManager) Stop() {
	close(r.stopChan)
}

func (r *RetentionManager) run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-r.stopChan:
			return
		case <-ticker.C:
			r.cleanup(ctx)
		}
	}
}

func (r *RetentionManager) cleanup(ctx context.Context) {
	cutoff := time.Now().AddDate(0, 0, -r.retentionDays).Unix()
	if err := r.logRepo.DeleteOlderThan(ctx, cutoff); err != nil {
		log.Printf("Retention: failed to delete logs: %v", err)
	}
	if err := r.traceRepo.DeleteOlderThan(ctx, cutoff); err != nil {
		log.Printf("Retention: failed to delete traces: %v", err)
	}
	if err := r.eventRepo.DeleteOlderThan(ctx, cutoff); err != nil {
		log.Printf("Retention: failed to delete events: %v", err)
	}
	if err := r.cacheRepo.CleanExpired(ctx); err != nil {
		log.Printf("Retention: failed to clean cache: %v", err)
	}
	log.Println("Retention cleanup completed")
}
