// Package gathuk
package gathuk

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
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
		var (
			partial T
			err     error
		)
		partial, err = g.loadFile(filename)
		if err != nil {
			return err
		}
		if err := g.mergeStruct(&g.value, &partial); err != nil {
			return err
		}
	}
	return nil
}

func (g *Gathuk[T]) LoadConfig(src io.Reader, format string) error {
	var (
		partial T
		err     error
	)
	partial, err = g.load(src, format)
	if err != nil {
		return err
	}
	if err := g.mergeStruct(&g.value, &partial); err != nil {
		return err
	}

	return nil
}

func (g *Gathuk[T]) WriteConfigFile(dst string, config T) error {
	return nil
}

func (g *Gathuk[T]) WriteConfig(out io.Writer, config T) error {
	return nil
}

func (g *Gathuk[T]) loadFile(filename string) (T, error) {
	f, err := os.Open(filename)
	if err != nil {
		var zeroValue T
		return zeroValue, err
	}

	defer f.Close()

	ext := strings.Trim(filepath.Ext(filename), ".")

	return g.load(f, ext)
}

func (g *Gathuk[T]) load(src io.Reader, format string) (T, error) {
	var buf bytes.Buffer

	_, err := io.Copy(&buf, src)
	if err != nil {
		var zeroValue T
		return zeroValue, err
	}

	by := buf.Bytes()

	dc, err := g.codecRegistry.Decoder(format)
	if err != nil {
		var zeroValue T
		return zeroValue, err
	}

	v, err := dc.Decode(by)
	if err != nil {
		return v, err
	}

	return v, nil
}

func (g *Gathuk[T]) GetConfig() T {
	return g.value
}

func (g *Gathuk[T]) mergeStruct(dst, src any) error {
	dv := reflect.ValueOf(dst).Elem()
	sv := reflect.ValueOf(src).Elem()

	for i := 0; i < dv.NumField(); i++ {
		df := dv.Field(i)
		sf := sv.Field(i)

		if !df.CanSet() {
			continue
		}

		switch df.Kind() {
		case reflect.Struct:
			if err := g.mergeStruct(df.Addr().Interface(), sf.Addr().Interface()); err != nil {
				return err
			}
		default:
			if !isZeroValue(sf) {
				df.Set(sf)
			}
		}
	}

	return nil
}
