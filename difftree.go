package akashi

import (
	"fmt"

	"github.com/seiyab/akashi/internal/doc"
)

type root struct {
	inner   diffTree
	printer func([]doc.Doc) string
}

type diffTree interface {
	docs() []doc.Doc
	loss() float64
}

var _ diffTree = split{}
var _ diffTree = mixed{}
var _ diffTree = cycle{}
var _ diffTree = nilNode{}
var _ diffTree = format1{}
var _ diffTree = format2{}
var _ diffTree = fail{}

func (d root) format() string {
	o := d.inner.docs()
	if d.printer != nil {
		return d.printer(o)
	}
	return doc.PrintDoc(o)
}

func quote(s string) string {
	return fmt.Sprintf("%q", s)
}
