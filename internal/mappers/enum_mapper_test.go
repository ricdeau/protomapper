package mappers

import (
	"testing"

	"github.com/ricdeau/protoast/ast"
	"github.com/stretchr/testify/require"
)

func TestEnumMapper(t *testing.T) {
	file := &ast.File{
		Package: "api",
	}
	enum := &ast.Enum{
		File: file,
		ParentMsg: &ast.Message{
			File: file,
			Name: "Request",
		},
		Name: "Status",
	}
	m := EnumMapper{
		typ: enum,
	}

	got := m.FromProto("ProjectStatus")("src")
	want := "src.ProjectStatus.String()"

	require.Equal(t, want, got)
	t.Log(got)

	got = m.ToProto("ProjectStatus")("src")
	want = "pb.Request_Status(pb.Request_Status_value[src.ProjectStatus])"

	require.Equal(t, got, want)
	t.Log(got)
}
