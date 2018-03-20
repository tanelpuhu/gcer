package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

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

func runGC(path string) {
	wd, _ := os.Getwd()
	defer chdir(wd)
	chdir(path)
	_, err := exec.Command("git", "gc", "--aggressive").Output()
	if err != nil {
		panic(err)
	}
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
	fmt.Printf("%-54s %11s -> ", path, fmtInt(sizeBefore))
	runGC(path)
	sizeAfter := getDirSize(path)
	fmt.Printf("%-14s\t%v%%\n", fmtInt(sizeAfter), 100*sizeAfter/sizeBefore)
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
		if fileExists(filepath.Join(path, "HEAD")) {
			sizeAndRunGC(basepath)
		}
		return filepath.SkipDir
	}
	return nil
}

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	filepath.Walk(root, walkCallback)
}
