package renderer

type Config interface {
	GetAppName() string
	GetTypesDir() string
	GetTypesGoPackage() string
	GetConvertersDir() string
	GetConvertersGoPkg() string
	GetTypesImport() string
	GetPbImport() string
	GetGenHelpers() bool
	GetAddComment() bool
}
