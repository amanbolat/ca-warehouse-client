// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"gopkg.in/errgo.v2/errors"
	"strings"
	"time"
)

const packageName = "github.com/amanbolat/ca-warehouse-client"

func Mod() error {
	err := sh.RunV("go", "mod", "download")
	if err != nil {
		return err
	}

	return sh.RunV("go", "mod", "tidy")
}

func Generate() error {
	return sh.RunV("go", "generate", "./...")
}

func Run() error {
	if err := Mod(); err != nil {
		return err
	}
	return sh.RunV("go", "run", "./cmd/main.go", "-c", "docker.env", "run")
}

func ClearDist() error {
	return sh.RunV("rm", "-rf", "./dist")
}

func Build() error {
	ClearDist()
	tag, _ := sh.Output("git", "describe", "--tags", "--abbrev=0")
	return sh.RunV("go", "build", "-ldflags", "-s -w -X '"+packageName+"/common.Version="+tag+"'", "-o", "./dist/whclient", "./cmd/main.go")
}

func TagPush() error {
	err := Mod()
	if err != nil {
		return err
	}
	err = Generate()
	if err != nil {
		return err
	}
	out, _ := sh.Output("git", "status", "-s")
	if strings.TrimSpace(out) != "" {
		return errors.New("Some files are not committed, can't build docker image")
	}

	commitMsg, err := sh.Output("git", "log", "-1", "--pretty=%B")
	if err != nil {
		return err
	}

	version := fmt.Sprintf("%d.%d%d.%d", time.Now().Year(), time.Now().Month(), time.Now().YearDay(), time.Now().Unix())
	err = sh.RunV("git", "tag", "-a", version, "-m", commitMsg)
	if err != nil {
		return err
	}

	err = sh.RunV("git", "push", "origin", "master", "--tags")
	if err != nil {
		return err
	}

	return nil
}

func GitStatus() error {
	out, _ := sh.Output("git", "status", "-s")
	fmt.Printf("[%v]", strings.TrimSpace(out))

	return nil
}

func LatestTag() error {
	out, _ := sh.Output("git", "describe", "--tags", "--abbrev=0")
	fmt.Printf("%v", strings.TrimSpace(out))

	return nil
}
