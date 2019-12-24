package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/urfave/cli.v1"
)

func init() {
	// Help text template
	cli.AppHelpTemplate = `{{.Name}} v{{.Version}}

DESCRIPTION:
	{{.Description}}

USAGE:
	{{.HelpName}} [options] path-to-database-binary

OPTIONS:
	{{range .Flags}}{{ . }}
	{{end}}
`

	// Custom Version printer
	cli.VersionPrinter = func(ctx *cli.Context) {
		fmt.Printf("%s (%s)\n", ctx.App.HelpName, ctx.App.Name)
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Operating System: %s\n", runtime.GOOS)
		fmt.Printf("Architecture: %s\n", runtime.GOARCH)
		fmt.Printf("Go Version: %s\n", runtime.Version())
	}
}

var (
	// List and define option flags
	flags = []cli.Flag{
		listenPortFlag,
		verbosityFlag,
	}
	// TODO: replace cli.v1 with cli.v2 when it's out of beta. cli.v2 handles flag aliases in a cleaner way
	// For now, first one in each flag name alias list is assumed to be the full name.
	listenPortFlag = cli.IntFlag{
		Name:  "port, p",
		Usage: "Network listening port",
		Value: 8080,
	}
	verbosityFlag = cli.IntFlag{
		Name:  "verbosity",
		Usage: "Logging verbosity: 0=panic, 1=fatal, 2=error, 3=warn, 4=info, 5=debug, 6=trace",
		Value: 4,
	}
)

// newApp creates an app with sane defaults.
func newApp(desc string) *cli.App {
	app := cli.NewApp()
	app.Name = name
	app.HelpName = filepath.Base(os.Args[0])
	app.Author = author
	app.Version = version
	app.Description = desc
	app.Writer = os.Stdout
	return app
}

func printHelpOnError(ctx *cli.Context, err error) {
	fmt.Printf("%s %s\n\n", "Incorrect Usage.", err)
	_ = cli.ShowAppHelp(ctx)
}
