package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	log "github.com/sirupsen/logrus"

	Linter "github.com/cremindes/whalelint/linter"
	Lsp "github.com/cremindes/whalelint/lsp"
	Parser "github.com/cremindes/whalelint/parser"
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
    --format [json, summary]
    --return-value [app, bool, num]
    --verbosity [short, normal, high]
    file list
	-c, --config
  version
*/

type WhaleLintCLI struct {
	Lint    LintCommand    `kong:"cmd,help='run linter.'"`
	Lsp     LspCommand     `kong:"cmd,help='run language server'"`
	Version VersionCommand `kong:"cmd,help='show version.'"`

	Config  string         `kong:"help='config file path NOTIMPLEMENTED'"` // nolint:gofmt,gofumpt,goimports
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

type LintCommand struct {
	Format      string   `kong:"help='Report format [${enum}].',default='summary',enum='json, summary'"`
	NoColor     bool     `kong:"help='No color output'"`
	Paths       []string `kong:"arg,required,help='Path to Dockerfile.',type:'path'"`
	ReturnValue string   `kong:"help='Set return value to one of [${enum}] NOT_IMPLEMENTED.',default='app',enum='app, bool, num'"` // nolint:lll
	Verbosity   string   `kong:"help='Verbosity level [${enum}].',default='normal',enum='normal, short'"`
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

	fileContent, err := Utils.ReadFileContents(filePath)
	if err != nil {
		// this is virtually unreachable, so no test case for this branch
		// could put before the AST parsing, but it would only server test coverage percentage fetish,
		// while hurting readability marginally.
		// TODO: rethink error handling here.
		return fmt.Errorf("linter | %w", err)
	}

	Parser.RawParser.UpdateRawStr(fileContent)

	// Run Linter
	linter := Linter.Linter{}
	ruleValidationResultArray := linter.Run(stageList)

	switch lintCommand.Format {
	case "json":
		Report.PrintResultAsJSON(ruleValidationResultArray, os.Stdout)
	case "summary":
		var verbosity Report.VerbosityLevel

		switch lintCommand.Verbosity {
		case "high":
			verbosity = Report.VerbosityHigh
		case "normal":
			verbosity = Report.VerbosityNormal
		case "short":
			verbosity = Report.VerbosityShort
		}

		options := Report.SummaryOption{
			NoColor:   lintCommand.NoColor,
			Verbosity: verbosity,
		}
		Report.PrintSummary(ruleValidationResultArray, os.Stdout, options)
	}

	return nil
}

type LspCommand struct {
	Port int `help:"Port number" default:"18888"`
}

// Run starts the Language Server.
func (lspCommand *LspCommand) Run() error {
	err := Lsp.Serve(lspCommand.Port)

	return fmt.Errorf("%w", err)
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
