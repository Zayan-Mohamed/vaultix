package cli

import (
	"fmt"
	"sync"
	"time"
)

type ProgressSpinner struct {
	mu      sync.Mutex
	done    chan struct{}
	action  string
	message string
	current int
	total   int
	frames  []string
	idx     int
}

func NewProgressSpinner(action string) *ProgressSpinner {
	return &ProgressSpinner{
		done:   make(chan struct{}),
		action: action,
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	}
}

func (s *ProgressSpinner) Start() {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-s.done:
				// Clear the line on done
				fmt.Printf("\r\033[K")
				return
			case <-ticker.C:
				s.render()
			}
		}
	}()
}

func (s *ProgressSpinner) Update(current, total int, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.current = current
	s.total = total

	// Truncate message if it's too long
	if len(message) > 40 {
		message = message[:37] + "..."
	}
	s.message = message
}

func (s *ProgressSpinner) Stop() {
	close(s.done)
}

func (s *ProgressSpinner) render() {
	s.mu.Lock()
	defer s.mu.Unlock()

	frame := s.frames[s.idx]
	s.idx = (s.idx + 1) % len(s.frames)

	// Ensure we clear the line
	fmt.Printf("\r\033[K")

	if s.total > 0 {
		fmt.Printf("%s %s [%d/%d] %s", frame, s.action, s.current, s.total, s.message)
	} else if s.message != "" {
		fmt.Printf("%s %s %s", frame, s.action, s.message)
	} else {
		fmt.Printf("%s %s...", frame, s.action)
	}
}
