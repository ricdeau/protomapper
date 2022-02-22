package protomapper

// Config protomapper configuration.
type Config struct {
	// AppName - app name to include in generated comment.
	AppName string
	// TypesDir - folder, where types will be generated.
	TypesDir string
	// TypesGoPackage - go package name for generated types.
	TypesGoPackage string
	// ConvertersDir - folder, where type converters will be generated.
	ConvertersDir string
	// ConvertersGoPkg - go package for converters.
	ConvertersGoPkg string
	// TypesImport - import types package in converters.
	TypesImport string
	// PbImport - import protobuf generated package in converters.
	PbImport   string
	GenHelpers bool
}
