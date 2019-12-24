package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/SimonSK/geoipserver/pkg/webapi"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
)

const logTimestampFormat = time.RFC3339

var (
	log *logrus.Logger
	app *cli.App
)

func init() {
	log = logrus.New()
	log.Out = os.Stderr
	log.SetFormatter(&nested.Formatter{
		TimestampFormat: logTimestampFormat,
		HideKeys:        false,
		NoColors:        true,
		ShowFullLevel:   false,
		TrimMessages:    false,
	})
}

func setLogLevel(ctx *cli.Context) error {
	verbosity := ctx.GlobalInt(strings.Split(verbosityFlag.Name, ",")[0])
	log.SetLevel(logrus.Level(verbosity))
	return nil
}

func makeServerConfig(ctx *cli.Context) (*webapi.Config, error) {
	dbBinaryFullpath, err := filepath.Abs(ctx.Args()[0])
	if err != nil {
		return nil, err
	}
	listenPort := uint16(ctx.GlobalInt(strings.Split(listenPortFlag.Name, ",")[0]))
	log.Debugf("[databaseBinaryFilepath=%s listenPort=%d] server configs", dbBinaryFullpath, listenPort)
	return &webapi.Config{
		nameWithVersion,
		description,
		log,
		dbBinaryFullpath,
		listenPort,
	}, err
}

func start(ctx *cli.Context) error {
	log.Debugf("starting %s", nameWithVersion)

	// Check argument
	if !ctx.Args().Present() {
		err := fmt.Errorf("database binary filename required")
		printHelpOnError(ctx, err)
		return err
	}

	// Get server configs
	config, err := makeServerConfig(ctx)
	if err != nil {
		return err
	}

	// Create server
	s := &webapi.Server{Config: *config}

	// Start server
	return s.Start()
}

func main() {
	app = newApp(description)
	app.Before = setLogLevel
	app.Action = start
	app.Flags = append(app.Flags, flags...)
	sort.Sort(cli.FlagsByName(app.Flags))
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
