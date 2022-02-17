module github.com/ricdeau/protomapper

go 1.16

replace github.com/ricdeau/protoast => ../../proj/protoast

require (
	github.com/iancoleman/strcase v0.2.0
	github.com/pkg/errors v0.9.1
	github.com/ricdeau/enki v1.0.6
	github.com/ricdeau/protoast v0.30.2
	github.com/stretchr/testify v1.7.0
	google.golang.org/protobuf v1.27.1
)
