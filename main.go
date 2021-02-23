package main

import (
	"os"

	"github.com/alecthomas/kong"
	log "github.com/sirupsen/logrus"

	CLI "github.com/cremindes/whalelint/cli"
	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debug("We have", RuleSet.Get().Count(), "ruleset.")

	// Get arguments
	args := os.Args[1:]

	cli := CLI.WhaleLintCLI{}

	// Create our CLI
	parser := kong.Must(&cli, cli.Options()...)

	// If no argument is given, show help/usage
	if len(args) == 0 {
		args = []string{"--help"}
	}

	// Parse arguments
	ctx, err := parser.Parse(args)
	// Use lint as default command if none is given
	if err != nil {
		cli.ApplyDefaultCommand(err, &args)

		ctx, err = parser.Parse(args)
		if err != nil {
			// log.Error(err)
			parser.FatalIfErrorf(err)
			os.Exit(1)
		}
	}

	// Run command selected by CLI
	err = ctx.Run(
		kong.Name("WhaleLint"),
	)

	if err != nil {
		log.Error(err)
	}
}
