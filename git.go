package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

var LF = []byte("\n")

func git(version Version) error {
	status, err := exec.Command("git", "status", "-s").Output()
	if err != nil {
		fmt.Println("git status error")
		return err
	}

	files := bytes.Split(bytes.Trim(status, "\r\n"), LF)
	if len(files) != 1 {
		return fmt.Errorf("Uncommitted file found. Please stash or remove.")
	} else if !bytes.HasSuffix(files[0], []byte(".versions")) {
		return fmt.Errorf("Counldn't listed .versions file")
	}
	vs := version.VersionString()
	commit := exec.Command("git", "commit", "-a", "-m", vs)
	if err := commit.Run(); err != nil {
		fmt.Println("git commit error")
		return err
	}
	tag := exec.Command("git", "tag", vs)
	if err := tag.Run(); err != nil {
		fmt.Println("git tag error")
		return err
	}
	return nil
}
