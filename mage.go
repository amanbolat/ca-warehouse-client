// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"gopkg.in/errgo.v2/errors"
	"strings"
	"time"
)

const repoName = "ca-warehouse-client"
const dockerRegistryName = "087613087242.dkr.ecr.us-west-2.amazonaws.com/ca-warehouse-client"

func Mod() error {
	return sh.RunV("go", "mod", "download")
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
	return sh.RunV("go", "build", "-ldflags", "-X main.GitCommit=$GIT_COMMIT", "-s -w", "-o", "./dist/whclient", "./cmd/main.go")
}

func TagPush() error {
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
