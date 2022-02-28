package lsp

import (
	lsp "github.com/Azure/azapi-lsp/internal/protocol"
	"github.com/hashicorp/hcl-lang/lang"
)

func Links(links []lang.Link, caps lsp.DocumentLinkClientCapabilities) []lsp.DocumentLink {
	docLinks := make([]lsp.DocumentLink, len(links))

	for i, link := range links {
		tooltip := ""
		if caps.TooltipSupport {
			tooltip = link.Tooltip
		}
		docLinks[i] = lsp.DocumentLink{
			Range:   HCLRangeToLSP(link.Range),
			Target:  link.URI,
			Tooltip: tooltip,
		}
	}

	return docLinks
}
