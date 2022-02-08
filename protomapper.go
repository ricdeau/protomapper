package protomapper

import (
	"github.com/pkg/errors"
	"github.com/ricdeau/protoast"
	"github.com/ricdeau/protomapper/internal/dicts"
	"github.com/ricdeau/protomapper/internal/mappers"
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
		cfg.AppName, cfg.TypesDir, cfg.TypesGoPackage,
	)
	result.converterRenderer = renderer.NewConvertersRenderer(
		cfg.AppName, cfg.ConvertersDir, cfg.PbImport, cfg.TypesImport, typeMapper.Types, typeMapper.Fields,
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
			return errors.Wrapf(err, "get AST for file %q", fName)
		}
		for _, t := range file.Types {
			p.typeMapper.Types.Put(t, nil)
		}
	}

	return nil
}

// TypeMapper get type mapper.
func (p *ProtoMapper) TypeMapper() TypeMapper {
	return p.typeMapper
}

// Types get types.
func (p *ProtoMapper) Types() *dicts.TypeDict {
	return p.typeMapper.Types
}

// Fields get fields.
func (p *ProtoMapper) Fields() *dicts.FieldDict {
	return p.typeMapper.Fields
}

// TypeRenderer get type renderer.
func (p *ProtoMapper) TypeRenderer() Renderer {
	return p.typeRenderer
}

// ConvertersRenderer get converters renderer.
func (p *ProtoMapper) ConvertersRenderer() Renderer {
	return p.converterRenderer
}
