// Package shared
package shared

import (
	"testing"

	customtests "github.com/ahyalfan/gathuk/internal/utils/custom-test"
)

func TestTypes(t *testing.T) {
	t.Run("Test 1: Create variable use Tag type", func(t *testing.T) {
		var example Tag = "default"
		example.Set("change")
		customtests.Equals(t, Tag("change"), example.Get())
	})
}
