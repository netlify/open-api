package porcelain

import (
	"strings"
	"testing"
)

func TestGetAssetManagementSha(t *testing.T) {
	tests := []struct {
		contents string
		length   int
	}{
		{"Not a pointer file", 0},
		{"version https://git-lfs.github.com/spec/v1\n" +
			"oid sha256:7e56e498ccb4cbb9c672e1aed6710fb91b2fd314394a666c11c33b2059ea3d71\n" +
			"size 1743570", 64},
	}

	for _, test := range tests {
		file := strings.NewReader(test.contents)
		out := getAssetManagementSha(file)
		if len(out) != test.length {
			t.Fatalf("expected `%d`, got `%d`", test.length, len(out))
		}
	}
}
