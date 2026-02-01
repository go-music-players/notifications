# go-music/notifications

Universal D-Bus desktop notifications library for Go music players.

## Overview

This library provides a simple, clean interface for adding desktop notifications to any Go music player application. Notifications appear when tracks change, integrating seamlessly with Linux desktop environments.

## Features

- ✅ **Simple API** - One function call to show notifications
- ✅ **Smart deduplication** - Automatically prevents duplicate notifications
- ✅ **Configurable** - Control timeout, icon, behavior
- ✅ **Cross-platform build** - Compiles on all platforms (D-Bus only on Linux)
- ✅ **Lightweight** - Minimal dependencies, efficient
- ✅ **Well-tested** - Used in production by multiple projects

## Installation

```bash
go get github.com/go-music/notifications
```

## Quick Start

```go
package main

import (
    "github.com/go-music/notifications"
)

func main() {
    // Create notifier with default options
    opts := notifications.DefaultOptions("My Music Player")
    notifier, err := notifications.NewNotifier(opts)
    if err != nil {
        // Notifications not available (non-Linux or D-Bus issues)
        log.Printf("Notifications not available: %v", err)
        return
    }
    defer notifier.Close()

    // Show notification when track changes
    track := &notifications.TrackInfo{
        Title:  "Bohemian Rhapsody",
        Artist: "Queen",
        Album:  "A Night at the Opera",
    }

    err = notifier.Notify(track, notifications.StatePlaying)
    if err != nil {
        log.Printf("Failed to show notification: %v", err)
    }
}
```

## Usage

### Basic Notifications

The simplest way to use this library:

```go
// Create notifier
opts := notifications.DefaultOptions("myapp")
notifier, err := notifications.NewNotifier(opts)
if err != nil {
    // Handle error
}
defer notifier.Close()

// Notify on track change
track := &notifications.TrackInfo{
    Title:  "Song Title",
    Artist: "Artist Name",
    Album:  "Album Name",
}
notifier.Notify(track, notifications.StatePlaying)
```

### Custom Options

Customize notification behavior:

```go
opts := notifications.Options{
    AppName:         "My Player",
    Icon:            "media-playback-start",
    Timeout:         3000,  // 3 seconds
    NotifyOnPause:   true,  // Show notification when paused
    ReplaceExisting: true,  // Replace previous notification
}

notifier, err := notifications.NewNotifier(opts)
```

### Automatic Deduplication

`Notify()` automatically deduplicates notifications:

```go
// First call - shows notification
notifier.Notify(track, notifications.StatePlaying)

// Second call with same track - no notification shown
notifier.Notify(track, notifications.StatePlaying)

// Different track - shows notification
track2 := &notifications.TrackInfo{Title: "Different Song"}
notifier.Notify(track2, notifications.StatePlaying)
```

### Force Notification

Use `NotifyNow()` to bypass deduplication:

```go
// Always shows notification, even if track hasn't changed
notifier.NotifyNow(track, notifications.StatePlaying)
```

### Radio Stations

For radio stations, use the `Station` field:

```go
track := &notifications.TrackInfo{
    Title:   "Current Song",
    Artist:  "Artist Name",
    Station: "KEXP 90.3 FM",
}
notifier.Notify(track, notifications.StatePlaying)
// Shows: "Current Song" with "Artist Name\nKEXP 90.3 FM" as body
```

### Check Capabilities

Query what the notification daemon supports:

```go
caps, err := notifier.GetCapabilities()
if err == nil {
    for _, cap := range caps {
        fmt.Println("Capability:", cap)
    }
}
// Common capabilities: "actions", "body", "body-markup", "icon-static", etc.
```

## API Reference

### Types

#### TrackInfo

```go
type TrackInfo struct {
    Title    string        // Track title
    Artist   string        // Artist name
    Album    string        // Album name
    Station  string        // Station name (for radio/streaming)
    ImageURL string        // Album art URL (future use)
    Duration time.Duration // Track duration (future use)
}
```

#### PlaybackState

```go
type PlaybackState string

const (
    StatePlaying PlaybackState = "Playing"
    StatePaused  PlaybackState = "Paused"
    StateStopped PlaybackState = "Stopped"
)
```

#### Options

```go
type Options struct {
    AppName         string // Application name (default: "Music Player")
    Icon            string // Icon name (default: "media-playback-start")
    Timeout         int32  // Milliseconds (default: 5000)
    NotifyOnPause   bool   // Show on pause (default: false)
    ReplaceExisting bool   // Replace vs stack (default: true)
}
```

### Functions

#### NewNotifier

```go
func NewNotifier(options Options) (*Notifier, error)
```

Creates a new notification service. Returns error if D-Bus is unavailable.

#### DefaultOptions

```go
func DefaultOptions(appName string) Options
```

Returns sensible default options.

### Methods

#### Notify

```go
func (n *Notifier) Notify(track *TrackInfo, state PlaybackState) error
```

Shows a notification if the track has changed. Automatically deduplicates.

#### NotifyNow

```go
func (n *Notifier) NotifyNow(track *TrackInfo, state PlaybackState) error
```

Shows a notification immediately without deduplication.

#### Close

```go
func (n *Notifier) Close() error
```

Closes the D-Bus connection. Should be called when done.

#### GetCapabilities

```go
func (n *Notifier) GetCapabilities() ([]string, error)
```

Returns capabilities supported by the notification daemon.

## Examples

### Example: Event-Driven Player

```go
func main() {
    player := NewMyPlayer()

    opts := notifications.DefaultOptions("myplayer")
    notifier, err := notifications.NewNotifier(opts)
    if err != nil {
        log.Printf("Notifications not available: %v", err)
        // Continue without notifications
    } else {
        defer notifier.Close()
    }

    // Listen for track changes
    player.OnTrackChange(func(track TrackMetadata, state PlayState) {
        if notifier != nil {
            info := &notifications.TrackInfo{
                Title:  track.Title,
                Artist: track.Artist,
                Album:  track.Album,
            }

            var nState notifications.PlaybackState
            switch state {
            case PlayStatePlaying:
                nState = notifications.StatePlaying
            case PlayStatePaused:
                nState = notifications.StatePaused
            default:
                nState = notifications.StateStopped
            }

            notifier.Notify(info, nState)
        }
    })

    player.Run()
}
```

### Example: Polling-Based Player

```go
func main() {
    player := NewMyPlayer()

    opts := notifications.DefaultOptions("myplayer")
    notifier, _ := notifications.NewNotifier(opts)
    defer notifier.Close()

    ticker := time.NewTicker(5 * time.Second)
    for range ticker.C {
        track := player.GetCurrentTrack()
        state := player.GetPlaybackState()

        if track != nil {
            info := &notifications.TrackInfo{
                Title:  track.Title,
                Artist: track.Artist,
                Album:  track.Album,
            }
            notifier.Notify(info, state)
        }
    }
}
```

### Example: Custom Notification Style

```go
func main() {
    opts := notifications.Options{
        AppName:         "🎵 My Music Player",
        Icon:            "audio-headphones",
        Timeout:         10000, // 10 seconds
        NotifyOnPause:   true,  // Show "Paused" notifications
        ReplaceExisting: false, // Stack notifications
    }

    notifier, _ := notifications.NewNotifier(opts)
    defer notifier.Close()

    // Your player code...
}
```

## Testing

Test notifications manually:

```bash
# Using the library in a test program
go run examples/test.go

# Or test with dbus-send directly
dbus-send --session \
  --dest=org.freedesktop.Notifications \
  /org/freedesktop/Notifications \
  org.freedesktop.Notifications.Notify \
  string:"Test App" \
  uint32:0 \
  string:"media-playback-start" \
  string:"Test Title" \
  string:"Test Body" \
  array:string: \
  dict:string:string: \
  int32:5000
```

Check notification daemon capabilities:

```bash
dbus-send --session --print-reply \
  --dest=org.freedesktop.Notifications \
  /org/freedesktop/Notifications \
  org.freedesktop.Notifications.GetCapabilities
```

## Platform Support

| Platform | Support | Notes |
|----------|---------|-------|
| Linux | ✅ Full | Requires D-Bus session bus and notification daemon |
| macOS | ⚠️ Compiles | Notifications not functional (no D-Bus) |
| Windows | ⚠️ Compiles | Notifications not functional (no D-Bus) |

The library compiles on all platforms but notification functionality is Linux-only (D-Bus requirement).

## Desktop Environment Support

Works with any desktop environment that supports the freedesktop.org notification specification:

- ✅ GNOME
- ✅ KDE Plasma
- ✅ XFCE
- ✅ Cinnamon
- ✅ MATE
- ✅ Budgie
- ✅ Most other Linux DEs

## Projects Using This Library

- **[HeosTUI](https://github.com/stashluk/heostui)** - Terminal UI for HEOS-enabled devices
- **Your project here!** - Open a PR to add your project

## Comparison with Other Libraries

### Why not use godbus directly?

godbus is excellent but low-level. This library provides:
- Higher-level abstraction (one function call vs managing D-Bus details)
- Automatic deduplication
- Sensible defaults
- Type-safe track metadata
- Tested across multiple desktop environments

### Why not use other notification libraries?

Most existing Go notification libraries are:
- Abandoned or unmaintained
- Missing features (deduplication, replacement)
- Not music-focused
- Overly complex for simple use cases

This library is actively maintained by the go-music community and designed specifically for music players.

## Contributing

Contributions are welcome! Areas of interest:

- **Additional platforms** - Investigate alternatives for macOS/Windows
- **Album art** - Support image-data D-Bus hint
- **Actions** - Add action buttons (skip, pause, etc.)
- **Testing** - More comprehensive test coverage
- **Documentation** - More examples and use cases

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - See [LICENSE](LICENSE) for details.

## Related Projects

- **[go-music/mpris](https://github.com/go-music/mpris)** - MPRIS D-Bus interface
- **[go-music/backend](https://github.com/go-music/backend)** - Universal backend interface

## Specifications

This library implements:
- [Desktop Notifications Specification](https://specifications.freedesktop.org/notification-spec/notification-spec-latest.html)
- [D-Bus Specification](https://dbus.freedesktop.org/doc/dbus-specification.html)

## Support

- **Issues**: https://github.com/go-music/notifications/issues
- **Discussions**: https://github.com/go-music/notifications/discussions
- **Matrix**: #go-music:matrix.org

## Acknowledgments

- Extracted from [HeosTUI](https://github.com/stashluk/heostui)
- Based on notification implementations from stmp and other projects
- Built on [godbus](https://github.com/godbus/dbus)
