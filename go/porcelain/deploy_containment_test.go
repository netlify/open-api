package porcelain

import (
	gocontext "context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The deploy walk and the function/edge bundlers read every file through an
// os.Root over the directory they were given, so nothing outside that directory
// is ever hashed or uploaded. These tests assert that contract directly.

// manifestFunctionRel is the guard that keeps a customer-controlled functions
// manifest from naming a path outside the functions directory.
func TestManifestFunctionRel(t *testing.T) {
	root := filepath.Join("srv", "functions")

	cases := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{name: "plain file", path: "fn.zip", want: "fn.zip"},
		{name: "nested file", path: filepath.Join("sub", "fn.zip"), want: filepath.Join("sub", "fn.zip")},
		{name: "cleaned but in-root", path: filepath.Join("sub", "..", "fn.zip"), want: "fn.zip"},
		{name: "absolute in-root", path: filepath.Join(mustAbs(t, root), "fn.zip"), want: "fn.zip"},
		{name: "parent escape", path: filepath.Join("..", "escape"), wantErr: true},
		{name: "deep escape", path: filepath.Join("sub", "..", "..", "escape"), wantErr: true},
		{name: "bare parent", path: "..", wantErr: true},
		{name: "absolute escape", path: filepath.Join(mustAbs(t, "elsewhere"), "escape"), wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := manifestFunctionRel(root, tc.path)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func mustAbs(t *testing.T, p string) string {
	t.Helper()
	abs, err := filepath.Abs(p)
	require.NoError(t, err)
	return abs
}

// bundle rejects a functions manifest whose path resolves outside the functions
// directory, and never reads the out-of-tree file it names.
func TestBundle_RejectsManifestPathEscapingRoot(t *testing.T) {
	root := t.TempDir()
	functionsDir := filepath.Join(root, "functions")
	require.NoError(t, os.Mkdir(functionsDir, 0o755))

	// A file outside the functions directory that the manifest tries to reach.
	outside := filepath.Join(root, "secret.zip")
	require.NoError(t, os.WriteFile(outside, []byte("out-of-tree"), 0o644))

	manifest := map[string]any{
		"version": 1,
		"functions": []map[string]any{
			{"name": "escape", "path": filepath.Join("..", "secret.zip"), "runtime": "go"},
		},
	}
	body, err := json.Marshal(manifest)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(functionsDir, "manifest.json"), body, 0o644))

	_, _, _, err = bundle(gocontext.Background(), testDir(t, functionsDir), newTestTempDir(t), mockObserver{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "outside the functions directory")
}

// A symlink in the deploy directory is not followed: the walk records regular
// files only, so a link pointing outside the tree contributes nothing.
func TestWalk_DoesNotFollowSymlinkOutOfRoot(t *testing.T) {
	root := t.TempDir()
	deployDir := filepath.Join(root, "publish")
	require.NoError(t, os.Mkdir(deployDir, 0o755))

	require.NoError(t, os.WriteFile(filepath.Join(deployDir, "index.html"), []byte("ok"), 0o644))

	outside := filepath.Join(root, "secret")
	require.NoError(t, os.WriteFile(outside, []byte("out-of-tree"), 0o644))
	if err := os.Symlink(outside, filepath.Join(deployDir, "leak")); err != nil {
		t.Skipf("symlinks unavailable on this platform: %v", err)
	}

	files, err := walk(testDir(t, deployDir), mockObserver{}, false, false)
	require.NoError(t, err)

	assert.NotNil(t, files.Files["index.html"], "regular file should be deployed")
	assert.Nil(t, files.Files["leak"], "symlink should not be deployed")
	for name := range files.Files {
		assert.NotEqual(t, "out-of-tree", contentOf(t, deployDir, name))
	}
}

func contentOf(t *testing.T, dir, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(name)))
	require.NoError(t, err)
	return string(b)
}
