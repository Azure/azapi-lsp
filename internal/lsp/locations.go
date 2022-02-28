package lsp

import (
	"path/filepath"

	lsp "github.com/Azure/azapi-lsp/internal/protocol"
	"github.com/Azure/azapi-lsp/internal/uri"
	"github.com/hashicorp/hcl-lang/decoder"
)

func RefOriginsToLocations(origins decoder.ReferenceOrigins) []lsp.Location {
	locations := make([]lsp.Location, len(origins))

	for i, origin := range origins {
		originUri := uri.FromPath(filepath.Join(origin.Path.Path, origin.Range.Filename))
		locations[i] = lsp.Location{
			URI:   lsp.DocumentURI(originUri),
			Range: HCLRangeToLSP(origin.Range),
		}
	}

	return locations
}
