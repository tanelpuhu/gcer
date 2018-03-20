package main

import (
	"os"
	"testing"
)

func TestFileExists(t *testing.T) {
	if fileExists("main.go") == false {
		t.Errorf("main.go does not exist?")
	}
}

func TestGetDirSize(t *testing.T) {
	res := getDirSize(".")
	if res == 0 {
		t.Errorf("getDirSize returned %v?", res)
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Invalid path did not cause panic")
		}
	}()
	getDirSize("/path/to/no/where")
}

func TestChdir(t *testing.T) {
	wd, _ := os.Getwd()
	defer chdir(wd)
	chdir("..")
}

func TestChdirPanic(t *testing.T) {
	wd, _ := os.Getwd()
	defer chdir(wd)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Invalid path did not cause panic")
		}
	}()
	chdir("/path/to/no/where")
}

func TestRunGC(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Invalid path did not cause panic")
		}
	}()
	runGC("/path/to/no/where")
}

func isEqual(t *testing.T, size int64, expect string) {
	res := fmtInt(size)
	if res != expect {
		t.Errorf("FmtInt result wrong, got '%v', expected '%v'", res, expect)
	}
}

func TestFmtInt(t *testing.T) {
	isEqual(t, 10, "10b")
	isEqual(t, 1023, "1023b")
	isEqual(t, 1024, "1Kb")
	isEqual(t, 1024*1.5, "1Kb")
	isEqual(t, 1024*2, "2Kb")
	isEqual(t, 1024*1023, "1023Kb")
	isEqual(t, 1024*1024, "1Mb")
	isEqual(t, 1024*1024*1024, "1Gb")
	isEqual(t, 1024*1024*1024*1024, "1024Gb")
}
