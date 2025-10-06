package main

import (
	"log"
	"os/exec"
	"strings"

	"github.com/pkg/browser"
)

// Opens the specified URL in the default browser of the user.
func OpenURL(url string) {
	var err error

	if isWSL() {
		// Open default browser within Windows instead of the Linux subsystem
		err = exec.Command("cmd.exe", "/c", "", "start", url).Start()
	} else {
		err = browser.OpenURL(url)
	}

	if err != nil {
		log.Fatal(err)
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
