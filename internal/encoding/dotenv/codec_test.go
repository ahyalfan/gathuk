// Package dotenv
package dotenv

import (
	"fmt"
	"reflect"
	"testing"

	customtests "github.com/ahyalfan/gathuk/internal/utils/custom-test"
	"github.com/ahyalfan/gathuk/option"
)

type Example struct {
	Hello string
	Holla string
}

func TestCodec(t *testing.T) {
	cdc := Codec[Example]{}
	cdc.ApplyDecodeOption(&option.DecodeOption{
		AutomaticEnv: false,
	})
	got, err := cdc.Decode([]byte(
		`HELLO=apa
		 HOLLA=1a`))

	customtests.OK(t, err)
	fmt.Println(got)

	t.Run("Test 2: Encode", func(t *testing.T) {
		cdc := Codec[Example]{}
		got, err := cdc.Encode(Example{
			Hello: "hello_world",
			Holla: "asdasdas",
		})

		customtests.OK(t, err)
		// customtests.Equals(t, got, want)
		fmt.Print(string(got))
	})
}

func BenchmarkCodec(b *testing.B) {
	b.Run("Benchmarking 1: Decode", func(b *testing.B) {
		for b.Loop() {
			cdc := Codec[Example]{}
			cdc.ApplyDecodeOption(&option.DecodeOption{
				AutomaticEnv: false,
			})
			_, err := cdc.Decode([]byte(
				`HELLO=apa
		 HOLLA=1a`))

			customtests.OK(b, err)
		}
	})
	b.Run("Benchmarking 2: Encode ", func(b *testing.B) {
		for b.Loop() {
			cdc := Codec[Example]{}
			_, err := cdc.Encode(Example{
				Hello: "hello_world",
				Holla: "asdasdas",
			})

			customtests.OK(b, err)
		}
	})
}

func BenchmarkSetValue(b *testing.B) {
	// Contoh struct
	type MyStruct struct {
		Name   string
		Age    int
		Price  float64
		Active bool
	}

	// Contoh objek
	obj := MyStruct{}
	v := reflect.ValueOf(&obj).Elem() // Dapatkan value dari pointer struct

	// Benchmarking untuk set field "Name"
	b.Run("SetName", func(b *testing.B) {
		for b.Loop() {
			setValue(v.FieldByName("Name"), "Alice")
		}
	})

	// Benchmarking untuk set field "Age"
	b.Run("SetAge", func(b *testing.B) {
		for b.Loop() {
			setValue(v.FieldByName("Age"), "30")
		}
	})

	// Benchmarking untuk set field "Price"
	b.Run("SetPrice", func(b *testing.B) {
		for b.Loop() {
			setValue(v.FieldByName("Price"), "100.50")
		}
	})

	// Benchmarking untuk set field "Active"
	b.Run("SetActive", func(b *testing.B) {
		for b.Loop() {
			setValue(v.FieldByName("Active"), "true")
		}
	})
}

func BenchmarkSetValueV2(b *testing.B) {
	// Contoh struct
	type MyStruct struct {
		Name   string
		Age    int
		Price  float64
		Active bool
	}

	// Contoh objek
	obj := MyStruct{}
	v := reflect.ValueOf(&obj).Elem() // Dapatkan value dari pointer struct

	// Benchmarking untuk set field "Name"
	b.Run("SetName", func(b *testing.B) {
		for b.Loop() {
			setValueAny(v.FieldByName("Name"), "Alice")
		}
	})

	// Benchmarking untuk set field "Age"
	b.Run("SetAge", func(b *testing.B) {
		for b.Loop() {
			setValueAny(v.FieldByName("Age"), 30)
		}
	})

	// Benchmarking untuk set field "Price"
	b.Run("SetPrice", func(b *testing.B) {
		for b.Loop() {
			setValueAny(v.FieldByName("Price"), 100.50)
		}
	})

	// Benchmarking untuk set field "Active"
	b.Run("SetActive", func(b *testing.B) {
		for b.Loop() {
			setValueAny(v.FieldByName("Active"), true)
		}
	})
}
