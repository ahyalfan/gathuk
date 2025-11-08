// Package gathuk
package gathuk

import (
	"os"
	"testing"

	customtests "github.com/ahyalfan/gathuk/internal/utils/custom-test"
)

var (
	EXAMPLE_ENV_FILE   string = "./example/dotenv/.example.env"
	EXAMPLE_2_ENV_file string = "./example/dotenv/.example_2.env"
	EXAMPLE_1_ENV_file string = "./example/dotenv/.example_1.env"
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

type Simple3 struct {
	Simplee     int      `config:"simple_e"`
	Debug       bool     `config:"debug_c"`
	Database    Database `nested:"db"`
	ExampleType string
	User        string
	Editor      string
}

func TestGathukLoad(t *testing.T) {
	t.Run("Test 1 : Simple Load Gathuk config", func(t *testing.T) {
		gt := NewGathuk[Simple]()

		err := gt.LoadConfigFiles(EXAMPLE_ENV_FILE)
		customtests.OK(t, err)
		customtests.Equals(t, "hore", gt.GetConfig().SimpleC)
		customtests.Equals(t, 2, gt.GetConfig().SimpleE)
	})
	t.Run("Test 2 : Nested Struct Load Gathuk config", func(t *testing.T) {
		gt := NewGathuk[Simple2]()

		err := gt.LoadConfigFiles(EXAMPLE_ENV_FILE)
		customtests.OK(t, err)

		customtests.Equals(t, 2, gt.GetConfig().Simplee)
		customtests.Equals(t, true, gt.GetConfig().Debug)
		customtests.Equals(t, "dbtest", gt.GetConfig().Database.User)
		customtests.Equals(t, "halo", gt.GetConfig().Database.Server)
	})

	t.Run("Test 3 : Nested Struct Load Gathuk config", func(t *testing.T) {
		gt := NewGathuk[Simple2]()

		gt.SetConfigFiles(EXAMPLE_ENV_FILE)
		err := gt.LoadConfigFiles(EXAMPLE_1_ENV_file)
		// err := gt.LoadConfigFiles(EXAMPLE_ENV_FILE, EXAMPLE_1_ENV_file) // still works
		customtests.OK(t, err)

		customtests.Equals(t, 2, gt.GetConfig().Simplee)
		customtests.Equals(t, true, gt.GetConfig().Debug)
		customtests.Equals(t, "dbtest", gt.GetConfig().Database.User)
		customtests.Equals(t, "halo", gt.GetConfig().Database.Server)
		customtests.Equals(t, 200, gt.GetConfig().Database.PoolingMax)
		customtests.Equals(t, "senin", gt.GetConfig().ExampleType)
	})

	t.Run("Test 4 : Load Gathuk config and option global", func(t *testing.T) {
		t.Run("Test 4.1: option global default", func(t *testing.T) {
			gt := NewGathuk[Simple3]()

			err := gt.LoadConfigFiles(EXAMPLE_2_ENV_file)
			customtests.OK(t, err)

			customtests.Equals(t, 200, gt.GetConfig().Database.PoolingMax)
			customtests.Equals(t, "senin", gt.GetConfig().ExampleType)
			customtests.Equals(t, "bukan_ahyalfan", gt.GetConfig().User)
		})
		t.Run("Test 4.2: option global with automaticenv", func(t *testing.T) {
			gt := NewGathuk[Simple3]()

			gt.globalDecodeOpt.AutomaticEnv = true

			err := gt.LoadConfigFiles(EXAMPLE_2_ENV_file)
			customtests.OK(t, err)

			customtests.Equals(t, 200, gt.GetConfig().Database.PoolingMax)
			customtests.Equals(t, "senin", gt.GetConfig().ExampleType)
			customtests.Equals(t, "ahyalfan", gt.GetConfig().User)
			customtests.Equals(t, "nvim", gt.GetConfig().Editor)
		})

		t.Run("Test 4.3: option global with automaticenv but file priority", func(t *testing.T) {
			gt := NewGathuk[Simple3]()

			gt.globalDecodeOpt.AutomaticEnv = true
			gt.globalDecodeOpt.PreferFileOverEnv = true

			err := gt.LoadConfigFiles(EXAMPLE_2_ENV_file)
			customtests.OK(t, err)

			customtests.Equals(t, 200, gt.GetConfig().Database.PoolingMax)
			customtests.Equals(t, "senin", gt.GetConfig().ExampleType)
			customtests.Equals(t, "bukan_ahyalfan", gt.GetConfig().User)
			customtests.Equals(t, "nvim", gt.GetConfig().Editor)
		})

		t.Run("Test 4.4: option global with set in os env", func(t *testing.T) {
			gt := NewGathuk[Simple3]()

			gt.globalDecodeOpt.AutomaticEnv = true
			gt.globalDecodeOpt.PreferFileOverEnv = true
			gt.globalDecodeOpt.PersistToOSEnv = true

			err := gt.LoadConfigFiles(EXAMPLE_ENV_FILE)
			customtests.OK(t, err)

			customtests.Equals(t, "hore", os.Getenv("SIMPLE_C"))
		})
	})
}

func TestGathukWrite(t *testing.T) {
	t.Run("Test 1: simple write config", func(t *testing.T) {
		gt := NewGathuk[Simple]()

		err := gt.writeFile("example/dotenv/example_12.env", 0, Simple{
			SimpleC: "hore",
			SimpleE: 100,
		})
		customtests.OK(t, err)
	})
}

func BenchmarkGathuk(b *testing.B) {
	b.Run("Benchmark 1 : Simple Load Gathuk config", func(b *testing.B) {
		for b.Loop() {
			gt := NewGathuk[Simple]()

			err := gt.LoadConfigFiles(EXAMPLE_ENV_FILE)
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

			err := gt.LoadConfigFiles(EXAMPLE_ENV_FILE)
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

			err := gt.LoadConfigFiles(EXAMPLE_ENV_FILE, EXAMPLE_1_ENV_file)
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
