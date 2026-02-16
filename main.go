package main

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-toast/toast"
)

// Notification represents a notification with type, message, and options
type Notification struct {
	Type      string
	Title     string
	Message   string
	Timeout   int
	AutoClose bool
}

// Icon data for each notification type (colored circle icons)
var iconData = map[string]struct {
	Color    color.RGBA
	Symbol   string
}{
	"success": {Color: color.RGBA{R: 46, G: 204, B: 113, A: 255}, Symbol: "✓"},   // Green
	"error":   {Color: color.RGBA{R: 231, G: 76, B: 60, A: 255}, Symbol: "✗"},    // Red
	"info":    {Color: color.RGBA{R: 52, G: 152, B: 219, A: 255}, Symbol: "ℹ"},   // Blue
	"warning": {Color: color.RGBA{R: 241, G: 196, B: 15, A: 255}, Symbol: "⚠"},   // Yellow
}

func main() {
	args := os.Args[1:]

	// Default values
	notificationType := "info"
	timeout := 5
	autoClose := true
	customTitle := ""
	var message string

	// Parse arguments
	i := 0
	for i < len(args) {
		arg := args[i]

		if arg == "--help" || arg == "-help" || arg == "-h" {
			showHelp()
			os.Exit(0)
		}

		if strings.HasPrefix(arg, "--type=") {
			notificationType = strings.TrimPrefix(arg, "--type=")
			i++
			continue
		}

		if arg == "--type" || arg == "-type" {
			if i+1 < len(args) {
				notificationType = args[i+1]
				i += 2
				continue
			}
		}

		if strings.HasPrefix(arg, "--title=") {
			customTitle = strings.TrimPrefix(arg, "--title=")
			i++
			continue
		}

		if arg == "--title" || arg == "-title" {
			if i+1 < len(args) {
				customTitle = args[i+1]
				i += 2
				continue
			}
		}

		if strings.HasPrefix(arg, "--timeout=") {
			if val, err := strconv.Atoi(strings.TrimPrefix(arg, "--timeout=")); err == nil {
				timeout = val
			}
			i++
			continue
		}

		if arg == "--timeout" || arg == "-timeout" {
			if i+1 < len(args) {
				if val, err := strconv.Atoi(args[i+1]); err == nil {
					timeout = val
				}
				i += 2
				continue
			}
		}

		if strings.HasPrefix(arg, "--autoclose=") {
			autoClose = parseBool(strings.TrimPrefix(arg, "--autoclose="))
			i++
			continue
		}

		if arg == "--autoclose" || arg == "-autoclose" {
			if i+1 < len(args) {
				autoClose = parseBool(args[i+1])
				i += 2
				continue
			}
		}

		if !strings.HasPrefix(arg, "-") {
			message = arg
			i++
			continue
		}

		i++
	}

	if message == "" {
		fmt.Println("Message is required as a positional argument")
		showHelp()
		os.Exit(1)
	}

	// Validate notification type
	validTypes := []string{"success", "error", "info", "warning"}
	isValidType := false
	for _, t := range validTypes {
		if notificationType == t {
			isValidType = true
			break
		}
	}

	if !isValidType {
		fmt.Printf("Invalid notification type: %s. Valid types are: success, error, info, warning\n", notificationType)
		os.Exit(1)
	}

	// Determine title
	title := customTitle
	if title == "" {
		title = strings.Title(notificationType)
	}

	// Create notification
	notification := &Notification{
		Type:      notificationType,
		Title:     title,
		Message:   message,
		Timeout:   timeout,
		AutoClose: autoClose,
	}

	// Display the notification
	if err := displayNotification(notification); err != nil {
		fmt.Printf("Error displaying notification: %v\n", err)
		os.Exit(1)
	}
}

func parseBool(s string) bool {
	return strings.ToLower(s) == "true"
}

func showHelp() {
	fmt.Println(`notify - A CLI notification utility

Usage:
  notify MESSAGE [OPTIONS]

Arguments:
  MESSAGE             The notification message (positional argument)

Options:
  --title TITLE       Custom title for the notification (default: based on type)
  --type TYPE         Type of notification: success, error, info, warning (default: info)
  --timeout SECONDS   Timeout in seconds (default: 5)
  --autoclose BOOLEAN Auto close after timeout (default: true)
  --help              Show this help message

Examples:
  notify "Operation completed successfully" --type success
  notify "An error occurred" --type error --timeout 10
  notify "Build done" --title "My App" --type success
  notify "Download started" --title "Downloader" --type info --autoclose false
`)
}

// createIcon creates a colored icon PNG and returns the path
func createIcon(nType string) (string, error) {
	data, ok := iconData[nType]
	if !ok {
		data = iconData["info"]
	}

	// Create a 64x64 image
	size := 64
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Draw a filled circle with the color
	center := size / 2
	radius := size / 2 - 4

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			// Calculate distance from center
			dx := x - center
			dy := y - center
			distance := dx*dx + dy*dy

			if distance <= radius*radius {
				// Inside circle - draw color
				img.Set(x, y, data.Color)
			} else {
				// Outside circle - transparent
				img.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 0})
			}
		}
	}

	// Get temp directory
	tempDir := os.TempDir()
	iconPath := filepath.Join(tempDir, fmt.Sprintf("notify_icon_%s.png", nType))

	// Create the file
	file, err := os.Create(iconPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Encode as PNG
	err = png.Encode(file, img)
	if err != nil {
		return "", err
	}

	return iconPath, nil
}

// getIconPath returns the path to an icon file for the notification type
func getIconPath(nType string) (string, error) {
	// Try to create icon in temp directory
	iconPath, err := createIcon(nType)
	if err != nil {
		return "", err
	}
	return iconPath, nil
}

func displayNotification(n *Notification) error {
	// Create icon for this notification type
	iconPath, err := getIconPath(n.Type)
	if err != nil {
		// Continue without icon if there's an error
		iconPath = ""
	}

	// Build toast notification
	notification := toast.Notification{
		AppID:               "Notify CLI",
		Title:               n.Title,
		Message:             n.Message,
		Icon:                iconPath,
		Duration:            toast.Short,
		ActivationType:      "protocol",
		ActivationArguments: "dismiss",
	}

	// We can't easily add custom XML attributes through the go-toast library
	// The library generates PowerShell code that creates the toast
	// By default, clicking a toast dismisses it from the Action Center

	// Set audio based on type
	switch n.Type {
	case "success", "error", "warning":
		notification.Audio = toast.Default
	default:
		notification.Audio = toast.Silent
	}

	if !n.AutoClose {
		notification.Duration = toast.Long
	}

	// Show the notification - it will dismiss when clicked
	err = notification.Push()
	if err != nil {
		return err
	}

	// Small delay to ensure notification is sent before program exits
	time.Sleep(500 * time.Millisecond)

	// Clean up icon file
	if iconPath != "" {
		os.Remove(iconPath)
	}

	return nil
}

// Embedded icon as base64 (fallback)
var embeddedIcons = map[string]string{
	"success": "iVBORw0KGgoAAAANSUhEUgAAAEAAAABACAYAAACqaXHeAAABhUlEQVR4Xu2YsU7DMBCGv0lZQAKI7xCU2aJsACZgDZYNGMFrG4AN2AAbgA3gAV2oSRqFpCRqivYKvED8SK3s+mb+xf/dyA/yjN//B/gHHoD7+/v7ASb4+/t7QPVtAHP+PxUA+QXw9/f3B3jfFJCWwD+A/v6+wI7n5+fnBPB9fqT+Avi/v7//AHB/f2cB6TsB5QU0wP39/X8B9/f3F5C8Afj+/v4C0PQEwAfwYQvIAugF7u/vLwD9BQH39/cXkCZJAPXzAVL1BQXQAHB/f38BKVoKQP0ZgPj+BiYg+YMAykvIBJAXeT8BvAYkCygeQOqC4h6oAXi/v7+ANHd+fn4CyAJ4BbQC+xMBvEcAyQJIAugFaAXsTwTwEJC8gfASkNdgfhvAdHuLXsBdRNMXQH4D6Q3MvweQ/18Dkn1A8QfEL8D9/f0FJHkB94WlB9A+wPsEYCggNICyA+J7uL+/vwDS+gC1VUAd6gEIDaB8AM0DwBsEpCkgRKEB5QfQugAdo1oA/wLqfQDjAnwGeL2/v7+A0gC6BmhGkK4BdArgNqDZAnoBaQOdALYG0QHKANICSCsgKYBMALcG0gzIDkgpIE0ArQkYXYDeBJYgAPd/+2KOfwCbAONbQHkE5QAAAABJRU5ErkJggg==",
}

// getEmbeddedIconPath extracts embedded icon and returns path
func getEmbeddedIconPath(nType string) (string, error) {
	data, ok := embeddedIcons[nType]
	if !ok {
		data = embeddedIcons["success"]
	}

	// Decode base64
	imgData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	// Write to temp file
	tempDir := os.TempDir()
	iconPath := filepath.Join(tempDir, fmt.Sprintf("notify_icon_%s_embed.png", nType))

	err = os.WriteFile(iconPath, imgData, 0644)
	if err != nil {
		return "", err
	}

	return iconPath, nil
}