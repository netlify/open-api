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

func TestAddWithAssetManagement(t *testing.T) {
	files := newDeployFiles()
	tests := []struct {
		rel string
		sum string
	}{
		{"foo.jpg", "sum1"},
		{"bar.jpg", "sum2"},
		{"baz.jpg", "sum3:originalsha"},
	}

	for _, test := range tests {
		file := &FileBundle{}
		file.Sum = test.sum
		files.Add(test.rel, file)
	}

	out := files.Hashed["sum3"]
	if len(out) != 1 {
		t.Fatalf("expected `%d`, got `%d`", 1, len(out))
	}
	out2 := files.Sums["baz.jpg"]
	if out2 != "sum3:originalsha" {
		t.Fatalf("expected `%v`, got `%v`", "sum3:originalsha", out2)
	}
}
