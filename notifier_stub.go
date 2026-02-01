//go:build !linux

package notifications

import "fmt"

// Notifier stub for non-Linux platforms
type Notifier struct {
	options Options
}

// NewNotifier returns an error on non-Linux platforms
func NewNotifier(options Options) (*Notifier, error) {
	return nil, fmt.Errorf("D-Bus notifications are only available on Linux")
}

// Close is a no-op on non-Linux platforms
func (n *Notifier) Close() error {
	return nil
}

// Notify is a no-op on non-Linux platforms
func (n *Notifier) Notify(track *TrackInfo, state PlaybackState) error {
	return nil
}

// NotifyNow is a no-op on non-Linux platforms
func (n *Notifier) NotifyNow(track *TrackInfo, state PlaybackState) error {
	return nil
}

// GetCapabilities returns empty on non-Linux platforms
func (n *Notifier) GetCapabilities() ([]string, error) {
	return []string{}, nil
}
