package porcelain

import "testing"

func TestForceSlashSeparators(t *testing.T) {
	out := forceSlashSeparators("foo\\bar\\baz.js")
	if out != "foo/bar/baz.js" {
		t.Fatalf("expected `foo/bar/baz.js`, got `%s`", out)
	}
}
