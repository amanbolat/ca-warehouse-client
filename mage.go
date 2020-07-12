// +build mage

package main

import (
	"github.com/magefile/mage/sh"
)

func Mod() error {
	return sh.Run("go", "mod", "download")
}

func Generate() error {
	return sh.Run("go", "generate", "./...")
}

func Run() error {
	if err := Mod(); err != nil {
		return err
	}

	return sh.Run("go", "run", "./cmd/main.go")
}

func ClearDist() error {
	return sh.Run("rm", "-rf", "./dist")
}

func Build() error {
	ClearDist()
	return sh.Run("go", "build", "-ldflags", "-s -w", "-o", "./dist/whclient", "./cmd/main.go")
}
