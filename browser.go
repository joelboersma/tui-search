package main

import (
	"os/exec"
	"strings"

	"github.com/pkg/browser"
)

// Opens the specified URL in the default browser of the user.
func OpenURL(url string) error {
	if isWSL() {
		// Open default browser within Windows instead of the Linux subsystem
		return exec.Command("cmd.exe", "/c", "", "start", url).Start()
	} else {
		return browser.OpenURL(url)
	}
}

// Checks if the Go program is running inside Windows Subsystem for Linux
func isWSL() bool {
	releaseData, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(releaseData)), "microsoft")
}
