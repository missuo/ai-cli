//go:build !windows

package main

import (
	"os"
	"syscall"
)

func runProvider(path string, argv []string) error {
	return syscall.Exec(path, argv, os.Environ())
}
