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
	Simplee     int      `config:"simple_e"`
	Debug       bool     `config:"debug_c"`
	Database    Database `nested:"db"`
	ExampleType string
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
	t.Run("Test 2 : Nested Struct Load Gathuk config", func(t *testing.T) {
		gt := NewGathuk[Simple2]()

		err := gt.LoadConfigFiles(".example.env")
		customtests.OK(t, err)

		customtests.Equals(t, 2, gt.GetConfig().Simplee)
		customtests.Equals(t, true, gt.GetConfig().Debug)
		customtests.Equals(t, "dbtest", gt.GetConfig().Database.User)
		customtests.Equals(t, "halo", gt.GetConfig().Database.Server)
	})

	t.Run("Test 3 : Nested Struct Load Gathuk config", func(t *testing.T) {
		gt := NewGathuk[Simple2]()

		err := gt.LoadConfigFiles(".example.env", ".example_1.env")
		customtests.OK(t, err)

		customtests.Equals(t, 2, gt.GetConfig().Simplee)
		customtests.Equals(t, true, gt.GetConfig().Debug)
		customtests.Equals(t, "dbtest", gt.GetConfig().Database.User)
		customtests.Equals(t, "halo", gt.GetConfig().Database.Server)
		customtests.Equals(t, 200, gt.GetConfig().Database.PoolingMax)
		customtests.Equals(t, "senin", gt.GetConfig().ExampleType)
	})
}

func BenchmarkGathuk(b *testing.B) {
	b.Run("Benchmark 1 : Simple Load Gathuk config", func(b *testing.B) {
		for b.Loop() {
			gt := NewGathuk[Simple]()

			err := gt.LoadConfigFiles(".example.env")
			if err != nil {
				b.Fatalf("Failed to load config: %v", err)
			}
			_ = gt.GetConfig().SimpleC
			_ = gt.GetConfig().SimpleE
		}
	})

	b.Run("Benchmark 2 : Nested Struct Load Gathuk config", func(b *testing.B) {
		for b.Loop() {
			gt := NewGathuk[Simple2]()

			err := gt.LoadConfigFiles(".example.env")
			if err != nil {
				b.Fatalf("Failed to load config: %v", err)
			}
			_ = gt.GetConfig().Simplee
			_ = gt.GetConfig().Debug
			_ = gt.GetConfig().Database.User
			_ = gt.GetConfig().Database.Server
		}
	})

	b.Run("Benchmark 3 : Nested Struct Load Gathuk config with multiple files", func(b *testing.B) {
		for b.Loop() {
			gt := NewGathuk[Simple2]()

			err := gt.LoadConfigFiles(".example.env", ".example_1.env")
			if err != nil {
				b.Fatalf("Failed to load config: %v", err)
			}
			_ = gt.GetConfig().Simplee
			_ = gt.GetConfig().Debug
			_ = gt.GetConfig().Database.User
			_ = gt.GetConfig().Database.Server
			_ = gt.GetConfig().Database.PoolingMax
			_ = gt.GetConfig().ExampleType
		}
	})
}
