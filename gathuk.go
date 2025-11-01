// Package gathuk
package gathuk

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// Gathuk a configuration
type Gathuk[T any] struct {
	PersistInEnvirontment   bool   // set env in os or no, jika tidak maka taruh di memory atau return aja
	OverrideWithEnvironment bool   // true untuk override environment, false untuk string, true jika config file lebih di utamakan
	Mode                    string // dev, staging, production. mungkin set modenya di taruh di flag pas jalanin binary
	// mode file example dev.env,stag.env,dev.json

	// value
	value T

	// codec interface
	codecRegistry CodecRegistry[T]

	// logger
	logger *slog.Logger
}

type Option[T any] interface {
	apply(g *Gathuk[T])
}

type optionFunc[T any] func(g *Gathuk[T])

func (fn optionFunc[T]) apply(g *Gathuk[T]) {
	fn(g)
}

func NewGathuk[T any]() *Gathuk[T] {
	g := &Gathuk[T]{}
	g.codecRegistry = NewDefaultCodecRegister[T]()
	g.logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // default slog
	return g
}

func (g *Gathuk[T]) LoadConfigFiles(srcFiles ...string) error {
	resolveFilenames(srcFiles...)
	for _, filename := range srcFiles {
		return g.loadFile(filename)
	}
	return nil
}

func (g *Gathuk[T]) LoadConfig(src io.Reader, format string) error {
	return g.load(src, format)
}

func (g *Gathuk[T]) WriteConfigFile(dst string, config T) error {
	return nil
}

func (g *Gathuk[T]) WriteConfig(out io.Writer, config T) error {
	return nil
}

func (g *Gathuk[T]) loadFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	ext := filepath.Ext(filename)

	ext = strings.Trim(ext, ".")

	return g.load(f, ext)
}

func (g *Gathuk[T]) load(src io.Reader, format string) error {
	var buf bytes.Buffer

	_, err := io.Copy(&buf, src)
	if err != nil {
		return err
	}

	by := buf.Bytes()

	dc, err := g.codecRegistry.Decoder(format)
	if err != nil {
		return err
	}

	v, err := dc.Decode(by)
	if err != nil {
		return err
	}

	g.value = v

	return nil
}

func (g *Gathuk[T]) GetConfig() T {
	return g.value
}
