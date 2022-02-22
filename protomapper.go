package protomapper

import (
	"fmt"

	"github.com/ricdeau/protoast"
	"github.com/ricdeau/protoast/ast"
	"github.com/ricdeau/protomapper/internal/helpers"
	"github.com/ricdeau/protomapper/internal/mappers"
	"github.com/ricdeau/protomapper/internal/registry"
	"github.com/ricdeau/protomapper/internal/renderer"
)

var excludeNone = func(field ProtoField) bool {
	return false
}

type ProtoMapper struct {
	astBuilder          *protoast.Builder
	typeMapper          *mappers.TypeMapper
	typeRenderer        *renderer.TypeRenderer
	converterRenderer   *renderer.ConvertersRenderer
	excludeMessageField func(field ProtoField) bool
}

// NewProtoMapper setup new protomapper.
func NewProtoMapper(cfg *Config) *ProtoMapper {
	result := &ProtoMapper{}

	result.excludeMessageField = excludeNone
	typeMapper := mappers.NewTypeMapper(result.excludeMessageField)
	result.typeMapper = typeMapper
	result.typeRenderer = renderer.NewTypeRenderer(
		cfg.AppName, cfg.TypesDir, cfg.TypesGoPackage, typeMapper,
	)
	result.converterRenderer = renderer.NewConvertersRenderer(
		cfg.AppName, cfg.ConvertersDir, cfg.PbImport, cfg.TypesImport, cfg.GenHelpers, typeMapper,
	)

	return result
}

// WithOptions add options for protomapper.
func (p *ProtoMapper) WithOptions(opts ...Option) *ProtoMapper {
	for _, opt := range opts {
		if opt != nil {
			opt.apply(p)
		}
	}

	return p
}

// ResolveTypes resolver *.proto files and collect defined types.
func (p *ProtoMapper) ResolveTypes(resolver FileResolver, fileNames ...string) error {
	p.astBuilder = protoast.NewBuilder(protoast.NewFilesViaResolver(resolver), func(err error) {
		panic(err)
	})

	for _, fName := range fileNames {
		file, err := p.astBuilder.AST(fName)
		if err != nil {
			return fmt.Errorf("get AST for file %q: %v", fName, err)
		}
		for _, t := range file.Types {
			pbType, err := p.typeMapper.FromProtoType(t)
			if err != nil {
				return fmt.Errorf("map type %T: %v", t, err)
			}
			registry.Types.Put(t, pbType)
		}
	}

	return nil
}

// TypeMapper get type mapper.
func (p *ProtoMapper) TypeMapper() TypeMapper {
	return p.typeMapper
}

// Types get types.
func (p *ProtoMapper) Types() *registry.TypeDict {
	return registry.Types
}

// Fields get fields.
func (p *ProtoMapper) Fields() *registry.FieldDict {
	return registry.Fields
}

// TypeRenderer get type renderer.
func (p *ProtoMapper) TypeRenderer() Renderer {
	return p.typeRenderer
}

// ConvertersRenderer get converters renderer.
func (p *ProtoMapper) ConvertersRenderer() Renderer {
	return p.converterRenderer
}

func (p *ProtoMapper) ASTBuilder() *protoast.Builder {
	return p.astBuilder
}

func (p *ProtoMapper) SetTypeResolver(resolver func(r TypeResolver) TypeResolver) {
	p.typeRenderer.SetTypeResolver(resolver)
	p.converterRenderer.SetTypeResolver(resolver)
	p.typeMapper.SetTypeResolver(resolver)
}

func (p *ProtoMapper) SetFileResolver(resolver func(r TypeResolver) TypeResolver) {
	p.typeRenderer.SetFileResolver(resolver)
	p.converterRenderer.SetFileResolver(resolver)
}

func (p *ProtoMapper) SetFieldResolver(resolver func(r FieldResolver) FieldResolver) {
	p.typeRenderer.SetFieldResolver(resolver)
}

func (p *ProtoMapper) SetImportsResolver(resolver func(r ImportsResolver) ImportsResolver) {
	p.converterRenderer.SetImportsResolver(resolver)
}

func AddMapper(pbTypeName string, mapper FieldMapper) {
	registry.Mappers.Put(pbTypeName, mapper)
}

func GoTypeName(typ ast.Type) string {
	return helpers.GoTypeName(typ)
}
