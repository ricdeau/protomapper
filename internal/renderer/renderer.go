package renderer

import (
	"path/filepath"

	"github.com/iancoleman/strcase"
	"github.com/ricdeau/protoast/ast"
)

func goStructName(msg *ast.Message) string {
	if msg.GetParentMsg() == nil {
		return strcase.ToCamel(msg.GetName())
	}

	return goStructName(msg.GetParentMsg()) + "_" + msg.GetName()
}

func pkgFromDir(dir string) string {
	return filepath.Base(dir)
}
