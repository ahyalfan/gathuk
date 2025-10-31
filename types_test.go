// Package gathuk_test
package gathuk_test

import (
	"testing"

	"github.com/ahyalfan/gathuk"
	customtests "github.com/ahyalfan/gathuk/internal/utils/custom-test"
)

func TestTypes(t *testing.T) {
	t.Run("Test 1: Create variable use Tag type", func(t *testing.T) {
		var example gathuk.Tag = "default"
		example.Set("change")
		customtests.Equals(t, gathuk.Tag("change"), example.Get())
	})
}
