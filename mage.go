// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"gopkg.in/errgo.v2/errors"
	"strings"
)

const repoName = "ca-warehouse-client"
const dockerRegistryName = "087613087242.dkr.ecr.us-west-2.amazonaws.com/ca-warehouse-client"

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

func BuildImage() error {
	out, _ := sh.Output("git", "status", "-s")
	if strings.TrimSpace(out) != "" {
		return errors.New("Some files are not committed, can't build docker image")
	}

	gitTag, err := sh.Output("git", "rev-parse", " --short", "HEAD")
	if err != nil {
		return err
	}

	taggedRepoName := fmt.Sprintf("%s:%s", repoName, gitTag)

	err = sh.Run("docker", "build", "-t", taggedRepoName, ".")
	if err != nil {
		return err
	}

	taggedRegistryRepoName := fmt.Sprintf("%s:%s", dockerRegistryName, gitTag)
	err = sh.Run("docker", "tag", taggedRepoName, taggedRegistryRepoName)
	if err != nil {
		return err
	}

	return sh.Run("docker", "push", taggedRegistryRepoName)
}

func GitStatus() error {
	out, _ := sh.Output("git", "status", "-s")
	fmt.Printf("[%v]", strings.TrimSpace(out))

	return nil
}
