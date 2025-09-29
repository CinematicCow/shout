package clipboard

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

func CopyFileToClipboard(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	uri := (&url.URL{Scheme: "file", Path: absPath}).String()
	gnomePayload := fmt.Sprintf("copy\n%s\n", uri)

	if isWayland() && commandExists("wl-copy") {
		if err := runClipboardCommand("wl-copy", bytes.NewBufferString(uri+"\n"), "--type", "text/uri-list"); err == nil {
			return nil
		}
		if err := runClipboardCommand("wl-copy", bytes.NewBufferString(gnomePayload)); err == nil {
			return nil
		}
	}

	if commandExists("xclip") {
		if err := runClipboardCommand("xclip", bytes.NewBufferString(gnomePayload), "-selection", "clipboard", "-t", "x-special/gnome-copied-files"); err == nil {
			return nil
		}
	}

	if commandExists("xsel") {
		if err := runClipboardCommand("xsel", bytes.NewBufferString(gnomePayload), "--clipboard", "--input", "--mime-type", "x-special/gnome-copied-files"); err == nil {
			return nil
		}
	}

	return fmt.Errorf("no compatible clipboard utility found")
}

func runClipboardCommand(name string, input *bytes.Buffer, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = input
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run clipboard command %s: %w", name, err)
	}
	return nil
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func isWayland() bool {
	return os.Getenv("WAYLAND_DISPLAY") != ""
}
