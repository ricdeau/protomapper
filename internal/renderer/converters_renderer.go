package renderer

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/ricdeau/enki"
	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/dicts"
	"github.com/ricdeau/protomapper/internal/mappers"
	"github.com/ricdeau/protomapper/internal/types"
)

type ConvertersRenderer struct {
	app              string
	dir              string
	pkg              string
	protoPkg         string
	typesPkg         string
	typeDict         *dicts.TypeDict
	mappers          *mappers.FieldMappers
	typeNameResolver Resolver
	fileNameResolver Resolver
	dryRun           bool
}

func NewConvertersRenderer(app, dir, pbPkg, typesPkg string, typeDict *dicts.TypeDict, fieldDict *dicts.FieldDict) *ConvertersRenderer {
	return &ConvertersRenderer{
		app:              app,
		dir:              dir,
		pkg:              pkgFromDir(dir),
		protoPkg:         pbPkg,
		typesPkg:         typesPkg,
		typeDict:         typeDict,
		typeNameResolver: CamelCaseName,
		fileNameResolver: SnakeCaseGoTypeFile,
		mappers:          mappers.NewFieldMappers(dicts.NewMappersDict(), fieldDict),
	}
}

func (r *ConvertersRenderer) DryRun() *ConvertersRenderer {
	r.dryRun = true
	return r
}

func (c *ConvertersRenderer) Render(t types.Type) error {
	protoType, _ := c.typeDict.GetByName(t.GetName())
	if protoType == nil {
		return errors.Errorf("type %T not registered", t)
	}

	msg, ok := protoType.(*ast.Message)
	if !ok {
		return nil
	}

	var out io.Writer = os.Stdout
	if !c.dryRun {
		fileName := c.fileNameResolver(t)
		f, err := os.Create(filepath.Join(c.dir, fileName))
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	}

	typeName := c.typeNameResolver(t)
	fromPbName := typeName + "FromPb"
	toPbName := typeName + "ToPb"
	pbTypeName := goStructName(msg)

	file := enki.NewFile()
	file.Package(c.pkg)
	file.GeneratedBy(c.app)
	file.NewLine()
	file.Import("pb", c.protoPkg)
	file.Import("types", c.typesPkg)
	file.NewLine()

	fields := msg.Fields
	fromPbFields := make([]enki.Statement, 0, len(fields)+2)
	toPbFields := make([]enki.Statement, 0, len(fields)+2)

	initType := enki.Stmt().Line("result := new(types.@1)", typeName)
	initPb := enki.Stmt().Line("result := new(pb.@1)", pbTypeName)

	fromPbFields = append(fromPbFields, initType)
	toPbFields = append(toPbFields, initPb)
	for _, field := range fields {
		mapper := c.mappers.GetMapper(field)
		fromProtoField := enki.Stmt().Line(mapper.GetFromProto("src", "result"))
		fromPbFields = append(fromPbFields, fromProtoField)
		toProtoField := enki.Stmt().Line(mapper.GetToProto("src", "result"))
		toPbFields = append(toPbFields, toProtoField)
	}
	fromPbFields = append(fromPbFields, enki.Stmt().Line("return result"))
	toPbFields = append(toPbFields, enki.Stmt().Line("return result"))

	file.Add(enki.F(fromPbName).Params("src *pb." + pbTypeName).Returns("*types." + typeName).
		Body(fromPbFields...))

	file.NewLine()

	file.Add(enki.F(toPbName).Params("src *types." + typeName).Returns("*pb." + pbTypeName).
		Body(toPbFields...))

	return file.Write(out)
}
