package templates_test

import (
	_ "embed"
	"os"
	"testing"

	"github.com/lorenzogood/x/quotestack/internal/templates"
	"github.com/lorenzogood/x/quotestack/internal/testutil"
	"github.com/stretchr/testify/require"
)

//go:embed test.tar
var testDir []byte

// Ensure that our template renderer can parse and execute templates in a non-flat directory structure.
func TestTemplate_Renderer_Non_Flat(t *testing.T) {
	testutil.TSetup(t)

	tempDir := testutil.ExtractTar(t, testDir)
	f := os.DirFS(tempDir)

	tp, err := templates.New(f, "templ")
	require.NoError(t, err)

	require.Equal(t, tp.List(), []string{"", "hi/hi.tmpl.html", "index.tmpl.html"})

	err = tp.Run(os.Stderr, "hi/hi.tmpl.html", nil)
	require.NoError(t, err)
}
