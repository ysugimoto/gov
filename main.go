package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/ysugimoto/go-args"
	"io/ioutil"
	"path/filepath"
)

var sectionRegex = regexp.MustCompile("\\[v?([0-9]+)\\.?([0-9]+)?\\.?([0-9]+)?\\]([^\\[]+)")

type Version struct {
	major   int
	minor   int
	patch   int
	message string
}

func (v Version) String() string {
	var m string
	if v.message != "" {
		m = "\n" + v.message
	}
	return fmt.Sprintf("[%d.%d.%d]%s\n", v.major, v.minor, v.patch, m)
}
func (v Version) VersionString() string {
	return fmt.Sprintf("v%d.%d.%d", v.major, v.minor, v.patch)
}
func (v Version) BumpPatch(msg string) Version {
	return Version{
		major:   v.major,
		minor:   v.minor,
		patch:   v.patch + 1,
		message: msg,
	}
}
func (v Version) BumpMinor(msg string) Version {
	return Version{
		major:   v.major,
		minor:   v.minor + 1,
		message: msg,
	}
}
func (v Version) BumpMajor(msg string) Version {
	return Version{
		major:   v.major + 1,
		message: msg,
	}
}

func createVersion(major, minor, patch, message string) Version {
	_major, _ := strconv.Atoi(major)
	_minor, _ := strconv.Atoi(minor)
	_patch, _ := strconv.Atoi(patch)
	return Version{
		major:   _major,
		minor:   _minor,
		patch:   _patch,
		message: strings.Trim(message, "\r\n"),
	}
}

type Versions []Version

func (v Versions) Len() int {
	return len(v)
}
func (v Versions) Less(i, j int) bool {
	if v[i].major < v[j].major {
		return false
	} else if v[i].major > v[j].major {
		return true
	}

	if v[i].minor < v[j].minor {
		return false
	} else if v[i].minor > v[j].minor {
		return true
	}
	if v[i].patch < v[j].patch {
		return false
	} else if v[i].patch > v[j].patch {
		return true
	}
	return false
}
func (v Versions) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

var vs Versions
var versionFile string

func findup(dir string) (path string) {
	for {
		path = filepath.Join(dir, "versions")
		if _, err := os.Stat(path); err == nil {
			break
		} else if dir == "/" {
			return ""
		}
		dir = filepath.Dir(dir)
	}
	return filepath.Join(dir, "versions")
}

func setup() {
	vs = Versions{}
	pwd, _ := os.Getwd()
	versionFile = findup(pwd)
	if _, err := os.Stat(versionFile); err != nil {
		fmt.Println("Couldn't find version file.")
		os.Exit(1)
	}
	buf, err := ioutil.ReadFile(versionFile)
	if err != nil {
		fmt.Println("Couln't read Version file. Check file permission.")
		fmt.Println(err)
		os.Exit(1)
	}
	ret := sectionRegex.FindAllStringSubmatch(string(buf), -1)
	for _, v := range ret {
		version := createVersion(v[1], v[2], v[3], v[4])
		vs = append(vs, version)
	}
	sort.Sort(vs)
}

func main() {
	ctx := args.New().Alias("m", "message", "").Parse(os.Args[1:])
	subCommand := ctx.At(0)
	if subCommand == "init" {
		pwd, _ := os.Getwd()
		versionFile = findup(pwd)
		if versionFile != "" {
			fmt.Println("Version file has already exists.")
			os.Exit(1)
		}
		file, _ := os.OpenFile(filepath.Join(pwd, "versions"), os.O_CREATE|os.O_WRONLY, 0666)
		defer file.Close()
		file.WriteString("[0.0.1]\nFirst Version\n")
		fmt.Println("Version file created successfully!")
		os.Exit(0)
	}
	setup()
	if len(vs) == 0 {
		fmt.Println("Couldn't find any versions at versions file.")
		os.Exit(1)
	}
	latest := vs[0]
	if ctx.Len() == 0 {
		fmt.Println(latest.VersionString())
		os.Exit(0)
	}

	var newVersion Version
	switch ctx.At(0) {
	case "patch":
		newVersion = latest.BumpPatch(ctx.String("message"))
		writeVersion(newVersion)
		fmt.Println(newVersion.VersionString())
	case "minor":
		newVersion = latest.BumpMinor(ctx.String("message"))
		writeVersion(newVersion)
		fmt.Println(newVersion.VersionString())
	case "major":
		newVersion = latest.BumpMajor(ctx.String("message"))
		writeVersion(newVersion)
		fmt.Println(newVersion.VersionString())
	default:
		fmt.Println(latest.VersionString())
	}
}

func writeVersion(newVersion Version) {
	versions := Versions{newVersion}
	versions = append(versions, vs...)
	out := []string{}
	for _, v := range versions {
		out = append(out, v.String())
	}
	if err := ioutil.WriteFile(versionFile, []byte(strings.Join(out, "\n")), 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
