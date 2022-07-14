package renderer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/ricdeau/enki"
	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/helpers"
	"github.com/ricdeau/protomapper/internal/mappers"
	"github.com/ricdeau/protomapper/internal/registry"
	"github.com/ricdeau/protomapper/internal/types"
)

type ConvertersRenderer struct {
	app             string
	dir             string
	pkg             string
	protoPkg        string
	typesPkg        string
	typeDict        *registry.TypeDict
	fileResolver    types.TypeResolver
	typeResolver    types.TypeResolver
	importsResolver types.ImportsResolver
	typeMapper      *mappers.TypeMapper
	dryRun          bool
	genHelpers      bool
	genComment      bool
	helpersDone     int32
}

func NewConvertersRenderer(cfg Config, typeMapper *mappers.TypeMapper) *ConvertersRenderer {
	dir := cfg.GetConvertersDir()
	return &ConvertersRenderer{
		app:             cfg.GetAppName(),
		dir:             dir,
		pkg:             helpers.PkgFromDir(dir),
		protoPkg:        cfg.GetPbImport(),
		typesPkg:        cfg.GetTypesImport(),
		typeDict:        registry.Types,
		fileResolver:    helpers.SnakeCaseGoTypeFile,
		typeResolver:    helpers.CamelCaseName,
		importsResolver: helpers.DefaultImportsResolver,
		typeMapper:      typeMapper,
		genComment:      cfg.GetAddComment(),
		genHelpers:      cfg.GetGenHelpers(),
	}
}

func (c *ConvertersRenderer) SetTypeResolver(f func(r types.TypeResolver) types.TypeResolver) {
	c.typeResolver = f(c.typeResolver)
}

func (c *ConvertersRenderer) SetFileResolver(f func(r types.TypeResolver) types.TypeResolver) {
	c.fileResolver = f(c.fileResolver)
}

func (c *ConvertersRenderer) SetImportsResolver(f func(r types.ImportsResolver) types.ImportsResolver) {
	c.importsResolver = f(c.importsResolver)
}

func (r *ConvertersRenderer) DryRun() *ConvertersRenderer {
	r.dryRun = true
	return r
}

func (c *ConvertersRenderer) renderHelpers() error {
	if !c.genHelpers {
		return nil
	}

	var out io.Writer = os.Stdout
	if !c.dryRun {
		f, err := os.Create(filepath.Join(c.dir, "helpers.go"))
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	}

	file := enki.NewFile()
	file.Package(c.pkg)
	file.GeneratedBy(c.app)
	file.NewLine()

	file.Add(enki.F("Map[T1, T2 any]").Params("src []T1", "mapFunc func(x T1) T2").Returns("[]T2").
		Body(enki.Stmt().
			Line(`result := make([]T2, len(src))`).
			Line(`for i, x := range src {`).
			Line(`	result[i] = mapFunc(x)`).
			Line(`}`).
			NewLine().
			Line(`return result`),
		))

	return file.Write(out)
}

func (c *ConvertersRenderer) Render(pbTyp ast.Type) (err error) {
	if atomic.CompareAndSwapInt32(&c.helpersDone, 0, 1) {
		err := c.renderHelpers()
		if err != nil {
			return fmt.Errorf("render helpers: %v", err)
		}
	}

	typ := registry.Types.GetType(pbTyp)
	if typ == nil {
		typ, err = c.typeMapper.FromProtoType(pbTyp)
		if err != nil {
			return fmt.Errorf("from proto type %s: %v", pbTyp, err)
		}
	}

	msg, ok := pbTyp.(*ast.Message)
	if !ok {
		return nil
	}

	var out io.Writer = os.Stdout
	if !c.dryRun {
		fileName := c.fileResolver.Resolve(typ, pbTyp)
		f, err := os.Create(filepath.Join(c.dir, fileName))
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	}

	typeName := typ.GetName()
	fromPbName := typeName + "FromPb"
	toPbName := typeName + "ToPb"
	pbTypeName := helpers.PbGoStructName(msg)

	file := enki.NewFile()
	file.Package(c.pkg)
	if c.genComment {
		file.GeneratedBy(c.app)
	}
	file.NewLine()
	file.Import("pb", c.protoPkg)
	file.Import("types", c.typesPkg)

	for _, imprt := range c.importsResolver.Resolve(typ, pbTyp) {
		parts := strings.Split(imprt, " ")
		if len(parts) > 1 {
			file.Import(parts[0], parts[1])
		} else {
			file.Import("", imprt)
		}
	}

	file.NewLine()

	fields := msg.Fields
	fromPbFields := make([]enki.Statement, 0, len(fields)+2)
	toPbFields := make([]enki.Statement, 0, len(fields)+2)

	ifNil := enki.Stmt().
		Line(`if src == nil {`).
		Line(` return nil`).
		Line(`}`).
		NewLine()
	initType := enki.Stmt().Line("result := new(types.@1)", typeName)
	initPb := enki.Stmt().Line("result := new(pb.@1)", pbTypeName)

	fromPbFields = append(fromPbFields, ifNil, initType)
	toPbFields = append(toPbFields, ifNil, initPb)
	for _, field := range fields {
		mapper, err := mappers.GetMapper(field)
		if err != nil {
			return fmt.Errorf("get mapper: %v", err)
		}

		fieldName := helpers.PbGoFieldName(field)
		fromProtoField := enki.Stmt().Line("result.@1 = @2", fieldName, mapper.FromPb(fieldName)("src"))
		fromPbFields = append(fromPbFields, fromProtoField)
		toProtoField := enki.Stmt().Line("result.@1 = @2", fieldName, mapper.ToPb(fieldName)("src"))
		toPbFields = append(toPbFields, toProtoField)
	}
	fromPbFields = append(fromPbFields, enki.Stmt().Line("return result"))
	toPbFields = append(toPbFields, enki.Stmt().Line("return result"))

	file.Add(enki.F(fromPbName).Params("src *pb." + pbTypeName).Returns("*types." + typeName).
		Body(fromPbFields...))

	file.NewLine()

	file.Add(enki.F(toPbName).Params("src *types." + typeName).Returns("*pb." + pbTypeName).
		Body(toPbFields...))

	return file.GoFmt(true).Write(out)
}
