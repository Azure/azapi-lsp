package common

import "github.com/hashicorp/hcl/v2"

func RangeOver(a hcl.Range, b hcl.Range) hcl.Range {
	if a.Empty() {
		return b
	}
	if b.Empty() {
		return a
	}

	var start, end hcl.Pos
	if a.Start.Line < b.Start.Line || a.Start.Line == b.Start.Line && a.Start.Column < b.Start.Column {
		start = a.Start
	} else {
		start = b.Start
	}
	if a.End.Line > b.End.Line || a.End.Line == b.End.Line && a.End.Column > b.End.Column {
		end = a.End
	} else {
		end = b.End
	}
	return hcl.Range{
		Filename: a.Filename,
		Start:    start,
		End:      end,
	}
}
