package main

import (
	"bytes"
	"fmt"
	"strings"

	"os/exec"
)

// Execute git commands
func git(version Version) error {
	// Confirm you are in master branch
	branch, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return fmt.Errorf("[Git] Finding branch error")
	} else if branch != "master" {
		return fmt.Errorf("[Git] You aren't master branch. You need to checkout master branch")
	}
	// Confirm commit target files
	status, err := exec.Command("git", "status", "-s").Output()
	if err != nil {
		return fmt.Errorf("[Git] Fetch status error")
	}

	// If commit files isn't only .versions file, abort process
	files := bytes.Split(bytes.Trim(status, "\r\n"), []byte("\n"))
	if len(files) != 1 {
		return fmt.Errorf("[Git] Uncommitted file found. Please stash or remove.\n%s", string.Join(files, "\n"))
	} else if !bytes.HasSuffix(files[0], []byte(".versions")) {
		return fmt.Errorf("[Git] Counldn't find .versions file in commit files")
	}

	vs := version.VersionString()
	// Make commit
	commit := exec.Command("git", "commit", "-a", "-m", vs)
	if err := commit.Run(); err != nil {
		fmt.Println("[Git] Failed to commit version file")
		return err
	}
	// Make tag
	tag := exec.Command("git", "tag", "-a", vs, "-m", "")
	if err := tag.Run(); err != nil {
		fmt.Println("[Git] Failed to create version tag")
		return err
	}
	return nil
}
