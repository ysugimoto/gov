package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

// Execute git commands
func git(version Version) error {
	// First, confirm commit target files
	status, err := exec.Command("git", "status", "-s").Output()
	if err != nil {
		fmt.Println("git status error")
		return err
	}

	// If commit files isn't only .versions file, abort process
	files := bytes.Split(bytes.Trim(status, "\r\n"), []byte("\n"))
	if len(files) != 1 {
		return fmt.Errorf("Uncommitted file found. Please stash or remove.")
	} else if !bytes.HasSuffix(files[0], []byte(".versions")) {
		return fmt.Errorf("Counldn't listed .versions file")
	}

	vs := version.VersionString()
	// Make commit
	commit := exec.Command("git", "commit", "-a", "-m", vs)
	if err := commit.Run(); err != nil {
		fmt.Println("git commit error")
		return err
	}
	// Make tag
	tag := exec.Command("git", "tag", vs)
	if err := tag.Run(); err != nil {
		fmt.Println("git tag error")
		return err
	}
	return nil
}
