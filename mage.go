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
	return sh.RunV("go", "build", "-ldflags", "-s -w", "-o", "./dist/whclient", "./cmd/main.go")
}

func BuildImage() error {
	out, _ := sh.Output("git", "status", "-s")
	if strings.TrimSpace(out) != "" {
		return errors.New("Some files are not committed, can't build docker image")
	}

	gitTag, err := sh.Output("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return err
	}

	taggedRepoName := fmt.Sprintf("%s:%s", repoName, gitTag)

	err = sh.RunV("docker", "build", "-t", taggedRepoName, ".")
	if err != nil {
		return err
	}

	taggedRegistryRepoName := fmt.Sprintf("%s:%s", dockerRegistryName, gitTag)
	err = sh.RunV("docker", "tag", taggedRepoName, taggedRegistryRepoName)
	if err != nil {
		return err
	}

	return sh.RunV("docker", "push", taggedRegistryRepoName)
}

func TagPush() error {
	out, _ := sh.Output("git", "status", "-s")
	if strings.TrimSpace(out) != "" {
		return errors.New("Some files are not committed, can't build docker image")
	}

	version := fmt.Sprintf("%d.%d%d.%d", time.Now().Year(), time.Now().Month(), time.Now().YearDay(), time.Now().Unix())
	err := sh.RunV("git", "tag", "-a", version)
	if err != nil {
		return err
	}

	err = sh.RunV("git", "push", version)
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
