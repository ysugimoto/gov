package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"

	"io/ioutil"
	"path/filepath"

	"github.com/ysugimoto/go-args"
)

// Section split regex. Format is:

// ```
// [(major).(minor).(patch)]
// (description...)
//
// ---
// ...
var sectionRegex = regexp.MustCompile("\\[v?([0-9]+)\\.?([0-9]+)?\\.?([0-9]+)?\\]([^\\-]{4})")
var vs Versions
var versionFile string

// Find .versions file up to root
func findup(dir string) (path string) {
	for {
		path = filepath.Join(dir, ".versions")
		if _, err := os.Stat(path); err == nil {
			break
		} else if dir == "/" {
			return ""
		}
		dir = filepath.Dir(dir)
	}
	return
}

// Read version file and make structs
func setup() error {
	vs = Versions{}
	pwd, _ := os.Getwd()
	versionFile = findup(pwd)
	if _, err := os.Stat(versionFile); err != nil {
		return fmt.Errorf("Couldn't find version file.")
	}
	buf, err := ioutil.ReadFile(versionFile)
	if err != nil {
		return fmt.Errorf("Couln't read Version file. Check file permission.")
	}
	ret := sectionRegex.FindAllStringSubmatch(string(buf), -1)
	for _, v := range ret {
		version := createVersion(v[1], v[2], v[3], v[4])
		vs = append(vs, version)
	}
	sort.Sort(vs)
	if len(vs) == 0 {
		vs = append(vs, Version{})
	}
	return nil
}

// Main
func main() {
	var newVersion Version
	var err error

	defer func() {
		if err != nil {
			fmt.Println(err)
		}
	}()

	ctx := args.New().Alias("message", "m", "").Parse(os.Args[1:])

	switch ctx.At(0) {
	case "init":
		pwd, _ := os.Getwd()
		versionFile = findup(pwd)
		if versionFile != "" {
			fmt.Println("Version file has already exists.")
			os.Exit(1)
		}
		file, _ := os.OpenFile(filepath.Join(pwd, ".versions"), os.O_CREATE|os.O_WRONLY, 0666)
		defer file.Close()
		file.WriteString("[0.0.1]\nFirst Version\n")
		fmt.Println("Version file created successfully!")
	case "patch":
		if err = setup(); err != nil {
			return
		}
		newVersion = vs[0].BumpPatch(ctx.String("message"))
		if err = writeVersion(newVersion); err == nil {
			fmt.Println(newVersion.VersionString())
		}
	case "minor":
		if err = setup(); err != nil {
			return
		}
		newVersion = vs[0].BumpMinor(ctx.String("message"))
		if err = writeVersion(newVersion); err == nil {
			fmt.Println(newVersion.VersionString())
		}
		writeVersion(newVersion)
		fmt.Println(newVersion.VersionString())
	case "major":
		if err = setup(); err != nil {
			return
		}
		newVersion = vs[0].BumpMajor(ctx.String("message"))
		if err = writeVersion(newVersion); err == nil {
			fmt.Println(newVersion.VersionString())
		}
	default:
		if err = setup(); err != nil {
			return
		}
		fmt.Println(vs[0].VersionString())
	}
}

// Write new version
func writeVersion(newVersion Version) error {
	versions := Versions{newVersion}
	versions = append(versions, vs...)
	fp, _ := os.OpenFile(versionFile, os.O_RDWR, 0755)
	defer fp.Close()
	for i, v := range versions {
		fp.WriteString(v.String())
		if i != len(versions)-1 {
			fp.WriteString("----\n\n")
		}
	}
	if err := git(newVersion); err != nil {
		return err
	}
	return nil
}
