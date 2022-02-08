package renderer

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/ricdeau/enki"
	"github.com/ricdeau/protomapper/internal/types"
)

var CamelCaseName Resolver = func(t types.Type) string {
	return strcase.ToCamel(t.GetName())
}

var SnakeCaseGoTypeFile = func(t types.Type) string {
	return strcase.ToSnake(t.GetName()) + ".go"
}

type Resolver func(t types.Type) string

type TypeRenderer struct {
	app              string
	dir              string
	pkg              string
	typeNameResolver Resolver
	fileNameResolver Resolver
	dryRun           bool
}

func NewTypeRenderer(app, dir, pkg string) *TypeRenderer {
	return &TypeRenderer{
		app:              app,
		pkg:              pkg,
		dir:              dir,
		typeNameResolver: CamelCaseName,
		fileNameResolver: SnakeCaseGoTypeFile,
	}
}

func (c *TypeRenderer) SetTypeNameResolver(resolver func(t types.Type) string) {
	c.typeNameResolver = resolver
}

func (r *TypeRenderer) DryRun() *TypeRenderer {
	r.dryRun = true
	return r
}

func (r *TypeRenderer) Render(t types.Type) (err error) {
	if _, ok := t.(*types.Struct); !ok {
		return nil
	}

	var out io.Writer = os.Stdout
	if !r.dryRun {
		fileName := r.fileNameResolver(t)
		f, err := os.Create(filepath.Join(r.dir, fileName))
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	}

	file := enki.NewFile()
	file.Package(r.pkg)
	file.GeneratedBy(r.app)
	file.NewLine()

	fields := t.GetFields()
	fieldStatements := make([]enki.Statement, 0, len(fields))

	for _, f := range fields {
		st := enki.Field("@1 @2", f.GetName(), f.GetType().GetName())
		fieldStatements = append(fieldStatements, st)
		switch v := f.GetType().(type) {
		case *types.Struct:
			err = r.Render(v)
		case *types.Array:
			err = r.Render(v.Elem)
		case *types.Map:
			err = r.Render(v.Val)
		}

		if err != nil {
			return err
		}
	}

	typeName := r.typeNameResolver(t)
	comment := t.GetComment()
	if len(comment) >= 1 {
		comment[0] = strings.TrimPrefix(strings.TrimSpace(comment[0]), t.GetName())
		comment[0] = typeName + " " + comment[0]
	}
	for _, line := range comment {
		file.Line("// @1", strings.TrimSpace(line))
	}
	file.Add(enki.T(t.GetName()).Struct(fieldStatements...))

	return file.Write(out)
}
