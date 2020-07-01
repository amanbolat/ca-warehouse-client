// +build mage

package main

import (
	"github.com/magefile/mage/sh"
)

func GoMod() error {
	return sh.Run("go", "mod", "download")
}

func Run() error {
	if err := GoMod(); err != nil {
		return err
	}

	return sh.Run("go", "run", "./cmd/main.go")
}

func Mod() error {
	return sh.Run("go", "mod", "download")
}
