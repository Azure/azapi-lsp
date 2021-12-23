package cmd

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/filesystem"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver/handlers/complete"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/logging"
	ilsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/lsp"
	lsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/protocol"
)

type CompletionCommand struct {
	Ui cli.Ui

	atPos string
}

func (c *CompletionCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("completion")

	fs.StringVar(&c.atPos, "at-pos", "", "at zero-indexed position (line:col)")

	fs.Usage = func() { c.Ui.Error(c.Help()) }

	return fs
}

func (c *CompletionCommand) Run(args []string) int {
	f := c.flags()
	if err := f.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}

	if f.NArg() != 1 {
		c.Ui.Output(fmt.Sprintf("args is %q", c.flags().Args()))
		return 1
	}

	path := f.Arg(0)

	path, err := filepath.Abs(path)
	if err != nil {
		c.Ui.Output(err.Error())
		return 1
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("reading file at %q failed: %s", path, err))
		return 1
	}

	fh := ilsp.FileHandlerFromPath(path)

	parts := strings.Split(c.atPos, ":")
	if len(parts) != 2 {
		c.Ui.Error(fmt.Sprintf("Error parsing at-pos argument: %q (expected line:col format)", c.atPos))
		return 1
	}
	line, err := strconv.Atoi(parts[0])
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing line: %s (expected number)", err))
		return 1
	}
	col, err := strconv.Atoi(parts[1])
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing column: %s (expected number)", err))
		return 1
	}
	lspPos := lsp.Position{Line: uint32(line), Character: uint32(col)}

	logger := logging.NewLogger(os.Stderr)

	fs := filesystem.NewFilesystem()
	fs.SetLogger(logger)
	fs.CreateAndOpenDocument(fh, "terraform", content)

	doc, err := fs.GetDocument(fh)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	fPos, err := ilsp.FilePositionFromDocumentPosition(lsp.TextDocumentPositionParams{
		TextDocument: lsp.TextDocumentIdentifier{
			URI: fh.DocumentURI(),
		},
		Position: lspPos,
	}, doc)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	pos := fPos.Position()

	src, err := doc.Text()
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}
	candidates := complete.CandidatesAtPos(src, doc.Filename(), pos, logger)

	c.Ui.Output(fmt.Sprintf("%#v", candidates))
	return 0
}

func (c *CompletionCommand) Help() string {
	helpText := `
Usage: terraform-ls completion [options] [path]

` + c.Synopsis() + "\n\n" + helpForFlags(c.flags())
	return strings.TrimSpace(helpText)
}

func (c *CompletionCommand) Synopsis() string {
	return "Lists available completion items"
}
