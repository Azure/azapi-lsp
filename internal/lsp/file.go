package lsp

import (
	lsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/protocol"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/source"
)

type File interface {
	URI() string
	FullPath() string
	Dir() string
	Filename() string
	Lines() source.Lines
	LanguageID() string
}

type file struct {
	fh         *fileHandler
	ls         source.Lines
	text       []byte
	version    int
	languageID string
}

func (f *file) URI() string {
	return f.fh.URI()
}

func (f *file) FullPath() string {
	return f.fh.FullPath()
}

func (f *file) Dir() string {
	return f.fh.Dir()
}

func (f *file) Filename() string {
	return f.fh.Filename()
}

func (f *file) Text() []byte {
	return f.text
}

func (f *file) Lines() source.Lines {
	return f.lines()
}

func (f *file) LanguageID() string {
	return f.languageID
}

func (f *file) lines() source.Lines {
	if f.ls == nil {
		f.ls = source.MakeSourceLines(f.fh.Filename(), f.text)
	}
	return f.ls
}

func (f *file) Version() int {
	return f.version
}

func FileFromDocumentItem(doc lsp.TextDocumentItem) *file {
	return &file{
		fh:         FileHandlerFromDocumentURI(doc.URI),
		text:       []byte(doc.Text),
		version:    int(doc.Version),
		languageID: doc.LanguageID,
	}
}
