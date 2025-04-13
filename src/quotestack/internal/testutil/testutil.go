package testutil

import (
	"archive/tar"
	"bytes"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/neilotoole/slogt"
	"github.com/stretchr/testify/require"
)

// Extract a embedded test directory tarball into a temporary directory on the host.
func ExtractTar(t testing.TB, data []byte) string {
	tempDir, err := os.MkdirTemp("", "qtest-*")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	tr := tar.NewReader(bytes.NewReader(data))
	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)

		fname := filepath.Join(tempDir, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			err := os.MkdirAll(fname, 0o755)
			require.NoError(t, err)
		case tar.TypeReg:
			f, err := os.OpenFile(fname, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			require.NoError(t, err)

			_, err = io.Copy(f, tr)
			f.Close()
			require.NoError(t, err)
		}
	}

	return tempDir
}

func TSetup(t *testing.T) {
	t.Helper()
	t.Parallel()

	logger := slogt.New(t, slogt.JSON())
	slog.SetDefault(logger)
}
