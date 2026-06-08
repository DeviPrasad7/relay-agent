// Package orchestration scheduler runs the pipeline on a ticker.
//
// Responsibilities:
//   - Start background ticker
//   - Handle graceful shutdown
//   - Call pipeline and update state
//
// Last updated: 2026-06-09
package orchestration

import (
	"context"
	"log"
	"time"
)

// Scheduler runs the pipeline periodically.
type Scheduler struct {
	pipeline      *Pipeline
	stateStore    *StateStore
	tickerInterval time.Duration
	stopChan      chan struct{}
	running       bool
}

// NewScheduler creates a scheduler.
func NewScheduler(pipeline *Pipeline, stateStore *StateStore, intervalSeconds int) *Scheduler {
	return &Scheduler{
		pipeline:       pipeline,
		stateStore:     stateStore,
		tickerInterval: time.Duration(intervalSeconds) * time.Second,
		stopChan:       make(chan struct{}),
	}
}

// Start begins the periodic execution.
func (s *Scheduler) Start(ctx context.Context) {
	if s.running {
		return
	}
	s.running = true
	go s.run(ctx)
}

// Stop signals the scheduler to stop after current tick.
func (s *Scheduler) Stop() {
	if !s.running {
		return
	}
	close(s.stopChan)
	s.running = false
}

func (s *Scheduler) run(ctx context.Context) {
	ticker := time.NewTicker(s.tickerInterval)
	defer ticker.Stop()
	// Run immediately on start
	s.tick(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.tick(ctx)
		}
	}
}

func (s *Scheduler) tick(ctx context.Context) {
	lastTs, err := s.stateStore.GetLastProcessed(ctx)
	if err != nil {
		log.Printf("failed to get last processed timestamp: %v", err)
		return
	}
	newTs, err := s.pipeline.Run(ctx, lastTs)
	if err != nil {
		log.Printf("pipeline run error: %v", err)
		return
	}
	if err := s.stateStore.SetLastProcessed(ctx, newTs); err != nil {
		log.Printf("failed to update last processed timestamp: %v", err)
	}
	log.Printf("tick completed: last=%d, new=%d", lastTs, newTs)
}
