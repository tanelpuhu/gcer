package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func FileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetDirSize(path string) int64 {
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

func Chdir(path string) {
	if err := os.Chdir(path); err != nil {
		panic(err)
	}
}

func RunGC(path string) {
	_, err := exec.Command("git", "gc", "--aggressive").Output()
	if err != nil {
		panic(err)
	}
}

func FmtInt(size int64) string {
	unit := "b"
	// i know, if if if'y
	if size > 1024 {
		unit, size = "Kb", size/1024
	}
	if size > 1024 {
		unit, size = "Mb", size/1024
	}
	if size > 1024 {
		unit, size = "Gb", size/1024
	}
	result := fmt.Sprintf("%d%s", size, unit)
	return result
}

func SizeAndRunGC(path string) {
	size_before := GetDirSize(path)
	fmt.Printf("%-54s %11s -> ", path, FmtInt(size_before))
	wd, _ := os.Getwd()
	Chdir(path)
	RunGC(path)
	Chdir(wd)
	size_after := GetDirSize(path)
	fmt.Printf("%-14s\t%v%%\n", FmtInt(size_after), 100*size_after/size_before)
}

func WalkCallback(path string, info os.FileInfo, err error) error {
	if err != nil {
		panic(err)
	}
	if info.IsDir() && info.Name() == ".git" {
		basepath, err := filepath.Abs(filepath.Dir(path))
		if err != nil {
			panic(err)
		}
		if FileExists(filepath.Join(path, "HEAD")) {
			SizeAndRunGC(basepath)
		}
		return filepath.SkipDir
	}
	return nil
}

func main() {
	filepath.Walk(".", WalkCallback)
}
