//go:build linux

package notifications

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

const (
	notificationsInterface = "org.freedesktop.Notifications"
	notificationsPath      = "/org/freedesktop/Notifications"
)

// Notifier sends desktop notifications via D-Bus
type Notifier struct {
	conn      *dbus.Conn
	options   Options
	lastID    string // Track ID to detect changes
	replaceID uint32 // Replace previous notification
}

// NewNotifier creates a new D-Bus notification service
func NewNotifier(options Options) (*Notifier, error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to session bus: %w", err)
	}

	// Test that notifications are available
	obj := conn.Object(notificationsInterface, notificationsPath)
	call := obj.Call(notificationsInterface+".GetCapabilities", 0)
	if call.Err != nil {
		conn.Close()
		return nil, fmt.Errorf("D-Bus notifications not available: %w", call.Err)
	}

	return &Notifier{
		conn:      conn,
		options:   options,
		replaceID: 0,
	}, nil
}

// Close closes the D-Bus connection
func (n *Notifier) Close() error {
	if n.conn != nil {
		return n.conn.Close()
	}
	return nil
}

// Notify shows a notification for a track
// Only notifies if the track has changed (based on title/artist/album)
func (n *Notifier) Notify(track *TrackInfo, state PlaybackState) error {
	if track == nil {
		return nil
	}

	// Don't notify if nothing is playing
	if track.Title == "" && track.Artist == "" {
		return nil
	}

	// Don't notify on pause unless configured to do so
	if state == StatePaused && !n.options.NotifyOnPause {
		return nil
	}

	// Check if track has changed
	currentID := fmt.Sprintf("%s-%s-%s", track.Title, track.Artist, track.Album)
	if currentID == n.lastID {
		return nil // Same track, don't notify again
	}

	// Update last track
	n.lastID = currentID

	// Show notification
	return n.showNotification(track, state)
}

// NotifyNow shows a notification immediately without deduplication
func (n *Notifier) NotifyNow(track *TrackInfo, state PlaybackState) error {
	if track == nil {
		return nil
	}
	return n.showNotification(track, state)
}

// showNotification displays a desktop notification
func (n *Notifier) showNotification(track *TrackInfo, state PlaybackState) error {
	obj := n.conn.Object(notificationsInterface, notificationsPath)

	// Build notification body
	var body string
	if track.Artist != "" && track.Album != "" {
		body = fmt.Sprintf("%s\n%s", track.Artist, track.Album)
	} else if track.Artist != "" {
		body = track.Artist
	} else if track.Station != "" {
		body = track.Station
	} else {
		body = "Now Playing"
	}

	// Add state indicator if paused
	if state == StatePaused {
		body = "⏸ " + body
	}

	// Notification summary (title)
	summary := track.Title
	if summary == "" {
		summary = "Now Playing"
	}

	// Application name
	appName := n.options.AppName
	if appName == "" {
		appName = "Music Player"
	}

	// Icon
	icon := n.options.Icon
	if icon == "" {
		icon = "media-playback-start"
	}

	// Actions (empty for now - could add skip/pause buttons)
	actions := []string{}

	// Hints (could add album art via image-data hint)
	hints := map[string]dbus.Variant{}

	// Determine replace ID
	replaceID := n.replaceID
	if !n.options.ReplaceExisting {
		replaceID = 0 // Always create new notification
	}

	// Call Notify
	call := obj.Call(
		notificationsInterface+".Notify",
		0,
		appName,           // app_name
		replaceID,         // replaces_id (0 = new notification, >0 = replace)
		icon,              // app_icon
		summary,           // summary
		body,              // body
		actions,           // actions
		hints,             // hints
		n.options.Timeout, // expire_timeout (-1 = default, 0 = never, >0 = milliseconds)
	)

	if call.Err != nil {
		return fmt.Errorf("failed to show notification: %w", call.Err)
	}

	// Store the notification ID so we can replace it next time
	if n.options.ReplaceExisting && len(call.Body) > 0 {
		if id, ok := call.Body[0].(uint32); ok {
			n.replaceID = id
		}
	}

	return nil
}

// GetCapabilities returns the capabilities supported by the notification daemon
func (n *Notifier) GetCapabilities() ([]string, error) {
	obj := n.conn.Object(notificationsInterface, notificationsPath)
	call := obj.Call(notificationsInterface+".GetCapabilities", 0)

	if call.Err != nil {
		return nil, fmt.Errorf("failed to get capabilities: %w", call.Err)
	}

	if len(call.Body) > 0 {
		if caps, ok := call.Body[0].([]string); ok {
			return caps, nil
		}
	}

	return []string{}, nil
}
