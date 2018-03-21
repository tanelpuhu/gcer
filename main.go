package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const gcerVersion string = "0.0.5"

var flagVersion bool
var flagAgressive bool

func fileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func getDirSize(path string) int64 {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	if err != nil {
		panic(err)
	}
	return size

}

func chdir(path string) {
	if err := os.Chdir(path); err != nil {
		panic(err)
	}
}

func runGC(path string) time.Duration {
	wd, _ := os.Getwd()
	defer chdir(wd)
	chdir(path)
	start := time.Now()
	lastArg := "--auto"
	if flagAgressive {
		lastArg = "--aggressive"
	}
	_, err := exec.Command("git", "gc", lastArg).Output()
	if err != nil {
		panic(err)
	}
	return time.Now().Sub(start)
}

func fmtInt(size int64) string {
	unit := "b"
	// i know, if if if'y
	if size >= 1024 {
		unit, size = "Kb", size/1024
	}
	if size >= 1024 {
		unit, size = "Mb", size/1024
	}
	if size >= 1024 {
		unit, size = "Gb", size/1024
	}
	result := fmt.Sprintf("%d%s", size, unit)
	return result
}

func sizeAndRunGC(path string) {
	sizeBefore := getDirSize(path)
	fmt.Printf("%-64s %11s -> ", path, fmtInt(sizeBefore))
	elapsed := runGC(path)
	sizeAfter := getDirSize(path)
	fmt.Printf("%-14s\t%v%%\t%.2fs\n", fmtInt(sizeAfter), 100*sizeAfter/sizeBefore, elapsed.Seconds())
}

func walkCallback(path string, info os.FileInfo, err error) error {
	if err != nil {
		panic(err)
	}
	if info.IsDir() && info.Name() == ".git" {
		basepath, err := filepath.Abs(filepath.Dir(path))
		if err != nil {
			panic(err)
		}
		if fileExists(filepath.Join(path, "HEAD")) && fileExists(filepath.Join(path, "refs")) {
			sizeAndRunGC(basepath)
		}
		return filepath.SkipDir
	}
	return nil
}

func init() {
	flag.BoolVar(&flagVersion, "V", false, "Print version")
	flag.BoolVar(&flagAgressive, "a", false, "use --aggressive")
	flag.Parse()
}

func main() {
	if flagVersion {
		fmt.Printf("gcer %v\n", gcerVersion)
		return
	}
	root := "."
	if len(flag.Args()) > 1 {
		root = flag.Args()[1]
	}
	filepath.Walk(root, walkCallback)
}
