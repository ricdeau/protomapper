package renderer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ricdeau/enki"
	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/helpers"
	"github.com/ricdeau/protomapper/internal/mappers"
	"github.com/ricdeau/protomapper/internal/registry"
	"github.com/ricdeau/protomapper/internal/types"
)

type TypeRenderer struct {
	app           string
	dir           string
	pkg           string
	genComment    bool
	typeResolver  types.TypeResolver
	fileResolver  types.TypeResolver
	fieldResolver types.FieldResolver
	typeMapper    *mappers.TypeMapper
	dryRun        bool
}

func NewTypeRenderer(cfg Config, typeMapper *mappers.TypeMapper) *TypeRenderer {
	return &TypeRenderer{
		app:           cfg.GetAppName(),
		pkg:           cfg.GetTypesGoPackage(),
		dir:           cfg.GetTypesDir(),
		genComment:    cfg.GetAddComment(),
		typeResolver:  helpers.CamelCaseName,
		fileResolver:  helpers.SnakeCaseGoTypeFile,
		fieldResolver: helpers.StandardFieldResolver,
		typeMapper:    typeMapper,
	}
}

func (c *TypeRenderer) SetTypeResolver(f func(r types.TypeResolver) types.TypeResolver) {
	c.typeResolver = f(c.typeResolver)
}

func (c *TypeRenderer) SetFileResolver(f func(r types.TypeResolver) types.TypeResolver) {
	c.fileResolver = f(c.fileResolver)
}

func (c *TypeRenderer) SetFieldResolver(f func(r types.FieldResolver) types.FieldResolver) {
	c.fieldResolver = f(c.fieldResolver)
}

func (r *TypeRenderer) DryRun() *TypeRenderer {
	r.dryRun = true
	return r
}

func (r *TypeRenderer) Render(pbTyp ast.Type) (err error) {
	typ := registry.Types.GetType(pbTyp)
	if typ == nil {
		typ, err = r.typeMapper.FromProtoType(pbTyp)
		if err != nil {
			return fmt.Errorf("from proto type %s: %v", pbTyp, err)
		}
	}

	if _, ok := typ.(*types.Struct); !ok {
		return nil
	}
	fileName := r.fileResolver.Resolve(typ, pbTyp)
	if fileName == "" {
		return nil
	}

	var out io.Writer = os.Stdout
	if !r.dryRun {
		f, err := os.Create(filepath.Join(r.dir, fileName))
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	}

	file := enki.NewFile()
	file.Package(r.pkg)
	if r.genComment {
		file.GeneratedBy(r.app)
	}
	file.NewLine()

	fields := typ.GetFields()
	fieldStatements := make([]enki.Statement, 0, len(fields))

	for _, f := range fields {
		fieldType := f.GetType()
		pbType := registry.Types.GetPbType(fieldType)
		st := enki.Field(r.fieldResolver.Resolve(f, pbType))
		fieldStatements = append(fieldStatements, st)

		err = r.Render(pbType)
		if err != nil {
			return err
		}
	}

	typeName := typ.GetName()
	comment := typ.GetComment()
	if len(comment) >= 1 {
		comment[0] = strings.TrimSpace(comment[0])
		firstWordIdx := strings.Index(comment[0], " ")
		if firstWordIdx >= 0 {
			comment[0] = comment[0][firstWordIdx:]
		}
		comment[0] = typeName + comment[0]
	}
	for _, line := range comment {
		file.Line("// @1", strings.TrimSpace(line))
	}
	file.Add(enki.T(typeName).Struct(fieldStatements...))

	return file.Write(out)
}
