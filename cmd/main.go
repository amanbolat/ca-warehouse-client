package main

import (
	"fmt"
	"github.com/amanbolat/ca-warehouse-client/common"
	"github.com/amanbolat/ca-warehouse-client/config"
	"github.com/amanbolat/ca-warehouse-client/server"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

var conf config.Config
var logger *logrus.Logger

// var version string

func main() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})
	fmt.Println("VERSION:", common.Version)

	app := &cli.App{
		Name:      fmt.Sprintf("CrossAsia warehouse client service [%s]", common.Version),
		Version:   common.Version,
		Writer:    os.Stderr,
		ErrWriter: os.Stderr,
		Before: func(context *cli.Context) error {
			v := context.Bool("verbose")
			if v {
				logger.Info("debug mode on")
				logger.SetLevel(logrus.DebugLevel)
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "path to config(dotenv file) file",
			},
			&cli.BoolFlag{
				Name:  "verbose, v",
				Usage: "run in debug mode",
			},
		},
		Commands: []*cli.Command{
			{
				Name: "version",
				Action: func(context *cli.Context) error {
					cli.ShowVersion(context)
					return nil
				},
			},
			{
				Name:  "run",
				Usage: "run warehouse client",
				Action: func(context *cli.Context) error {
					configFile := context.String("config")

					err := godotenv.Load(configFile)
					if err != nil {
						logger.Fatalf("could not load env file: %v", err)
					}

					err = envconfig.Process("", &conf)
					if err != nil {
						logger.Fatalf("could not parse env vars: %v", err)
					}

					s, err := server.NewServer(conf, logger)
					if err != nil {
						return err
					}
					s.Start(conf.Port)

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
	}
}
