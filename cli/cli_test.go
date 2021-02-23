package cli_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	log "github.com/sirupsen/logrus"
	"gotest.tools/assert"

	"github.com/cremindes/whalelint/cli"
	TestHelper "github.com/cremindes/whalelint/testhelper"
)

type StdBuffer struct {
	stdOut *bytes.Buffer
	stdErr *bytes.Buffer
}

func generateStdBuffer() StdBuffer {
	return StdBuffer{
		stdOut: &bytes.Buffer{},
		stdErr: &bytes.Buffer{},
	}
}

func generateCLI(args []string) (*kong.Context, StdBuffer, error) {
	cli := cli.WhaleLintCLI{}
	parser := kong.Must(&cli, cli.Options()...)
	stdBuffer := generateStdBuffer()
	parser.Stdout, parser.Stderr = stdBuffer.stdOut, stdBuffer.stdErr

	ctx, err := parser.Parse(args)

	return ctx, stdBuffer, err // nolint:wrapcheck
}

func generateKongParseErrorUnknownFlag(t *testing.T) *kong.ParseError {
	t.Helper()

	args := []string{"--nonExistingFlag"}
	_, _, err := generateCLI(args)

	var kongParseError *kong.ParseError
	if !errors.As(err, &kongParseError) {
		t.FailNow()
	}

	return kongParseError
}

func TestVersionCommand_Run(t *testing.T) {
	t.Parallel()

	args := []string{"version"}

	ctx, stdBuffer, err := generateCLI(args)
	assert.NilError(t, err)

	err = ctx.Run()
	assert.NilError(t, err)

	assert.Equal(t, "whalelint: v0.0.1\n", stdBuffer.stdOut.String())
	assert.Equal(t, "", stdBuffer.stdErr.String())
}

func TestCliType_ApplyDefaultCommand(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name     string
		Args     []string
		Err      *kong.ParseError
		Expected []string
	}{
		{
			Name:     "Call default lint with path",
			Args:     []string{"pathToDockerfile"},
			Err:      generateKongParseErrorUnknownFlag(t),
			Expected: []string{"lint", "pathToDockerfile"},
		},
		{
			Name:     "Call explicit lint with path",
			Args:     []string{"lint", "pathToDockerfile"},
			Err:      nil,
			Expected: []string{"lint", "pathToDockerfile"},
		},
		{
			Name:     "Call version",
			Args:     []string{"version"},
			Err:      nil,
			Expected: []string{"version"},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			// prepare
			cli := cli.WhaleLintCLI{}
			parser := kong.Must(&cli, cli.Options()...)
			_, err := parser.Parse(testCase.Args)
			args := testCase.Args

			// test target
			cli.ApplyDefaultCommand(err, &args)

			// check
			assert.DeepEqual(t, testCase.Expected, args)
		})
	}
}

// nolint:funlen,paralleltest
func TestLintCommand_Run(t *testing.T) {
	testCases := []struct {
		Name           string
		TmpFileContent []string
		Expected       error
		ExpectedErrStr string
		ExpectedStdout string
	}{
		{
			Name:           "Call lint with path of simple Dockerfile.",
			TmpFileContent: []string{"FROM golang:1.16"},
			Expected:       nil,
			ExpectedErrStr: "",
			ExpectedStdout: "",
		},
		{
			Name:           "Call lint with path of non-existing file",
			TmpFileContent: []string{""},
			Expected:       syscall.ENOENT,
			ExpectedErrStr: "no such file",
			ExpectedStdout: "",
		},
		{
			Name:           "Call lint with path empty file",
			TmpFileContent: []string{" "},
			Expected:       &parser.ErrorLocation{},
			ExpectedErrStr: "file with no instructions",
			ExpectedStdout: "",
		},
		{
			Name:           "Call lint with path of simple Dockerfile with args.",
			TmpFileContent: []string{"ARG from=\"golang 1.16\"\nFROM ${from}"},
			Expected:       nil,
			ExpectedErrStr: "",
			ExpectedStdout: "",
		},
		{
			Name:           "Call lint with 2 paths of simple Dockerfiles",
			TmpFileContent: []string{"FROM golang:1.16", "FROM golang:1.16"},
			Expected:       nil,
			ExpectedErrStr: "",
			ExpectedStdout: "Although it's planned, for now only one Dockerfile can be validated at a time.",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.Name, func(t *testing.T) {
			// t.Parallel() - Linter is not thread safe!

			stdOut := &strings.Builder{}
			log.SetOutput(stdOut)

			tmpFileSlice := make([]*os.File, len(testCase.TmpFileContent))
			os.Args = make([]string, 1, len(tmpFileSlice))
			os.Args[0] = "lint"

			for i, fileContent := range testCase.TmpFileContent {
				if len(fileContent) > 0 {
					tmpFileSlice[i], _ = ioutil.TempFile("", "mock-dockerfile.*")

					_, errTmpFile := tmpFileSlice[i].WriteString(fileContent)
					assert.NilError(t, errTmpFile)
					errTmpFile = tmpFileSlice[i].Sync()
					assert.NilError(t, errTmpFile)

					os.Args = append(os.Args, tmpFileSlice[i].Name())
				}
			}

			if len(os.Args) == 1 {
				os.Args = append(os.Args, "bogusPath")
			}

			ctx, _, err := generateCLI(os.Args)
			assert.NilError(t, err)

			err = ctx.Run()

			isSameErr := TestHelper.CheckForErrorRecursively(t, err, testCase.Expected)
			assert.Equal(t, true, isSameErr)
			if testCase.Expected != nil {
				assert.ErrorContains(t, err, testCase.ExpectedErrStr)
			}
			println(stdOut.String())
			if len(testCase.ExpectedStdout) > 0 {
				assert.Equal(t, true, strings.Contains(stdOut.String(), testCase.ExpectedStdout))
			}

			for _, tmpFile := range tmpFileSlice {
				if tmpFile != nil {
					os.Remove(tmpFile.Name())
				}
			}
		})
	}
}
