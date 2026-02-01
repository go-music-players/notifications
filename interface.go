package notifications

import "time"

// TrackInfo represents track metadata for notifications
type TrackInfo struct {
	Title    string        // Track title
	Artist   string        // Artist name
	Album    string        // Album name
	Station  string        // Station name (for radio/streaming)
	ImageURL string        // Album art or station logo URL
	Duration time.Duration // Total track duration (0 if unknown)
}

// PlaybackState represents the current playback state
type PlaybackState string

const (
	StatePlaying PlaybackState = "Playing"
	StatePaused  PlaybackState = "Paused"
	StateStopped PlaybackState = "Stopped"
)

// Options configures notification behavior
type Options struct {
	AppName         string // Application name shown in notifications
	Icon            string // Icon name (defaults to "media-playback-start")
	Timeout         int32  // Notification timeout in milliseconds (default: 5000)
	NotifyOnPause   bool   // Show notification when paused (default: false)
	ReplaceExisting bool   // Replace previous notification instead of stacking (default: true)
}

// DefaultOptions returns sensible defaults
func DefaultOptions(appName string) Options {
	return Options{
		AppName:         appName,
		Icon:            "media-playback-start",
		Timeout:         5000,
		NotifyOnPause:   false,
		ReplaceExisting: true,
	}
}
