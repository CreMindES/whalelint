package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
	log "github.com/sirupsen/logrus"
	"robpike.io/filter"

	Linter "github.com/cremindes/whalelint/linter"
	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
	Report "github.com/cremindes/whalelint/report"
	Utils "github.com/cremindes/whalelint/utils"
)

/*
CLI plan so far:

Commands:
  help [automatic]
  lsp
	--port
    -c, --config
  lint [default]
    - | --output [textsummary, json]
    --short-summary? [mutually exclusive with output?]
    --return-value [program,lint]
    file list
	-c, --config
  version
*/

type WhaleLintCLI struct {
	Lint    LintCommand    `kong:"cmd,help='run linter.'"`
	Version VersionCommand `kong:"cmd,help='show version.'"`

	Config  string         `kong:"help='config file path'"` // nolint:gofmt,gofumpt,goimports
}

func (*WhaleLintCLI) Options() []kong.Option {
	return []kong.Option{
		kong.Name("whalelint"),
		kong.Description("WhaleLint is a Dockerfile linter. It can function as\n" +
			"  - CLI linter as default\n" +
			"  - Pre-commit hook with option --return-value=\"errnum\"\n" +
			"  - CI linter with options like --format=json and/or --return-value=\"errnum\"\n" +
			"  - Language server for plugins\n" +
			" (- GitHub Action - in the roadmap)" + "\n" +
			"More documentation at https://github.com/CreMindES/whalelint"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{ // nolint:exhaustivestruct
			Compact: true,
		}),
	}
}

type LintCommand struct { // nolint:maligned
	Format       string   `kong:"help='TODO',default:'summary',enum:'json,summary'"`
	NoColor      bool     `kong:"help='No color output'"`
	Paths        []string `kong:"arg,required,help:'Paths to remove.',type:'path'"`
	ReturnValue  string   `kong:"help='Return value TODO'"`
	ShortSummary bool     `kong:"help='Print only a short summary.'"`
}

func (lintCommand *LintCommand) Run() error {
	log.Println("Running linter... TODO", lintCommand)

	if len(lintCommand.Paths) > 1 {
		log.Warning("Although it's planned, for now only one Dockerfile can be validated at a time.")
	}

	filePath := lintCommand.Paths[0]

	// TODO: move to parser
	/* Parse Dockerfile */
	stageList, metaArgs, err := Utils.GetDockerfileAst(filePath)
	if err != nil {
		return fmt.Errorf("linter | %w", err)
	}

	if metaArgs != nil {
		log.Debug("metaArgs |", metaArgs)
	}

	// Run Linter
	linter := Linter.Linter{}
	ruleValidationResultArray := linter.Run(stageList)
	violations := filter.Choose(ruleValidationResultArray,
		func(x RuleSet.RuleValidationResult) bool {
			return x.IsViolated()
		},
	)

	/* Print result | TODO: cli dependent output */
	Report.PrintResultAsJSON(violations)

	return nil
}

type VersionCommand struct{}

func (versionCommand *VersionCommand) Run(k *kong.Context) error {
	version := "v0.0.1"
	k.Printf("%s", version)

	return nil
}

func (t WhaleLintCLI) ApplyDefaultCommand(err error, args *[]string) {
	var parseError *kong.ParseError
	if errors.As(err, &parseError) {
		if !strings.HasPrefix(err.Error(), "unknown flag") {
			// insert the default command in case there is only one argument
			// supporting case of ./whalelint [Dockerfile path] -> ./whalelint lint [Dockerfile path]
			*args = append([]string{"lint"}, *args...)
		}
	}
}
