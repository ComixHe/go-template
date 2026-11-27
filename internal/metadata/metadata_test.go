package metadata_test

import (
	"strings"
	"testing"

	"github.com/ComixHe/go-template/internal/metadata"
)

func TestString(t *testing.T) {
	t.Parallel()

	got := metadata.String()
	for _, want := range []string{
		metadata.AppName,
		metadata.Version,
		metadata.Commit,
		metadata.BuildDate,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("String() = %q, want it to contain %q", got, want)
		}
	}
}
