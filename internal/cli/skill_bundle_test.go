package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// readPrintingPressSkillBundle reads the printing-press skill as one bundle:
// the router SKILL.md plus its references/ and scripts/ files. Since the
// progressive-disclosure split, per-phase content lives in references/ and the
// setup contract in scripts/preflight.sh, so assertions about skill text must
// look at the whole bundle.
func readPrintingPressSkillBundle(t *testing.T) string {
	t.Helper()

	skillDir := filepath.Join("..", "..", "skills", "printing-press")
	data, err := os.ReadFile(filepath.Join(skillDir, "SKILL.md"))
	require.NoError(t, err)

	var content strings.Builder
	content.Write(data)
	for _, glob := range []string{
		filepath.Join(skillDir, "references", "*.md"),
		filepath.Join(skillDir, "scripts", "*.sh"),
	} {
		paths, err := filepath.Glob(glob)
		require.NoError(t, err)
		for _, path := range paths {
			extra, err := os.ReadFile(path)
			require.NoError(t, err)
			content.WriteByte('\n')
			content.Write(extra)
		}
	}
	return content.String()
}
