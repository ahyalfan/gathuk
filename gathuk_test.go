// Package gathuk
package gathuk

import (
	"testing"

	customtests "github.com/ahyalfan/gathuk/internal/utils/custom-test"
)

type Simple struct {
	SimpleC string
	SimpleE int
}

type Simple2 struct {
	Simplee  int      `config:"simple_e"`
	Debug    bool     `config:"debug_c"`
	Database Database `nested:"db"`
}

type Database struct {
	User       string
	Server     string `config:"server_port"`
	PoolingMax int    `config:"poling_max_pool"`
}

func TestGathuk(t *testing.T) {
	t.Run("Test 1 : Simple Load Gathuk config", func(t *testing.T) {
		gt := NewGathuk[Simple]()

		err := gt.LoadConfigFiles(".example.env")
		customtests.OK(t, err)
		customtests.Equals(t, "hore", gt.GetConfig().SimpleC)
		customtests.Equals(t, 2, gt.GetConfig().SimpleE)
	})
	t.Run("Test 1 : Nested Struct Load Gathuk config", func(t *testing.T) {
		gt := NewGathuk[Simple2]()

		err := gt.LoadConfigFiles(".example.env")
		customtests.OK(t, err)

		customtests.Equals(t, 2, gt.GetConfig().Simplee)
		customtests.Equals(t, true, gt.GetConfig().Debug)
		customtests.Equals(t, "dbtest", gt.GetConfig().Database.User)
		customtests.Equals(t, "halo", gt.GetConfig().Database.Server)
	})
}
