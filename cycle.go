package akashi

import "github.com/seiyab/akashi/internal/doc"

type cycle struct{}

func (c cycle) docs() []doc.Doc {
	return []doc.Doc{
		doc.Inline("<circular reference>"),
	}
}

func (c cycle) loss() float64 {
	return 0
}
