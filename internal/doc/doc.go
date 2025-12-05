package doc

import (
	"strings"
)

const (
	indentSize = 2
	sightWidth = 2
)

type lineKind int

const (
	lineContext lineKind = iota
	lineLeft
	lineRight
)

type ColorProfile struct {
	Left    string
	Right   string
	Context string
	Reset   string
}

const defaultReset = "\x1b[0m"

type Doc interface {
	print(depth int) virtualLines
	Left() Doc
	Right() Doc
	AddPrefix(prefix string) Doc
	AddSuffix(suffix string) Doc
}

var _ Doc = block{}
var _ Doc = inline{}

func PrintDoc(ds []Doc, profile ...ColorProfile) string {
	var p ColorProfile
	if len(profile) > 0 {
		p = profile[0]
	}
	p = normalizeProfile(p)
	vls, shouldPrint := renderLines(ds)

	var lines []string
	for i, vl := range vls {
		if !shouldPrint[i] {
			continue
		}
		if i > 0 && !shouldPrint[i-1] {
			lines = append(lines, ":")
		}
		lines = append(lines, paintLine(vl, p))
	}

	return strings.Join(lines, "\n")
}

type virtualLine struct {
	isContext bool
	isDiff    bool
	kind      lineKind
	text      string
}

type virtualLines []virtualLine

func renderLines(ds []Doc) (virtualLines, []bool) {
	vls := virtualLines{}
	for _, d := range ds {
		vls = append(vls, d.print(0)...)
	}

	shouldPrint := make([]bool, len(vls))
	for i, vl := range vls {
		if vl.isContext {
			shouldPrint[i] = true
		}
		if vl.isDiff {
			for j := i - sightWidth; j <= i+sightWidth; j++ {
				if j >= 0 && j < len(vls) {
					shouldPrint[j] = true
				}
			}
			if k := i - sightWidth - 2; k >= 0 && shouldPrint[k] {
				// prevent skipping just a single line
				shouldPrint[k+1] = true
			}
		}
	}

	return vls, shouldPrint
}

func paintLine(vl virtualLine, profile ColorProfile) string {
	var color string
	switch vl.kind {
	case lineLeft:
		color = profile.Left
	case lineRight:
		color = profile.Right
	case lineContext:
		color = profile.Context
	}

	if color == "" {
		return vl.text
	}
	return color + vl.text + profile.Reset
}

func normalizeProfile(p ColorProfile) ColorProfile {
	if p.Reset == "" {
		p.Reset = defaultReset
	}
	return p
}

func (v virtualLines) asContext() virtualLines {
	var out virtualLines
	for _, l := range v {
		l.isContext = true
		out = append(out, l)
	}
	return out
}

func (v virtualLines) hasDiff() bool {
	for _, l := range v {
		if l.isDiff {
			return true
		}
	}
	return false
}
