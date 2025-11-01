// Package dotenv
package dotenv

import (
	"fmt"
	"reflect"
	"testing"
)

type Example struct {
	Hello string
	Holla string
}

func TestCodec(t *testing.T) {
	cdc := Codec[Example]{}
	got, _ := cdc.Decode([]byte(
		`HELLO=apa
		 HOLLA=1a`))
	fmt.Println(got)
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
