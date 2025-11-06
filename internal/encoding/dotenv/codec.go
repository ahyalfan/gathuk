// Package dotenv
package dotenv

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	utility "github.com/ahyalfan/gathuk/internal/utils"
	"github.com/ahyalfan/gathuk/option"
	"github.com/ahyalfan/gathuk/shared"
)

type Codec[T any] struct {
	option.DefaultCodec[T]

	do   *option.DecodeOption
	eo   *option.EncodeOption
	temp map[string][]byte
}

func (c *Codec[T]) ApplyEncodeOption(eo *option.EncodeOption) {
	c.eo = eo
}

func (c *Codec[T]) CheckEncodeOption() bool {
	return c.eo != nil
}

func (c *Codec[T]) Encode(val T) ([]byte, error) {
	if c.temp == nil {
		c.temp = make(map[string][]byte)
	}

	c.flattenWithNestedPrefix(val)
	// var build strings.Builder
	// for k, v := range c.temp {
	// 	build.WriteString(k)
	// 	build.WriteRune('=')
	// 	build.Write(v)
	// 	build.WriteRune('\n')
	// }
	//
	// return []byte(build.String()), nil

	var build []byte
	for k, v := range c.temp {
		build = append(build, []byte(k)...)
		build = append(build, '=')
		build = append(build, v...)
		build = append(build, '\n')
	}
	return build, nil
}

func (c *Codec[T]) ApplyDecodeOption(do *option.DecodeOption) {
	c.do = do
}

func (c *Codec[T]) CheckDecodeOption() bool {
	return c.do != nil
}

func (c *Codec[T]) Decode(buf []byte) (T, error) {
	var value T

	if c.temp == nil {
		c.temp = make(map[string][]byte)
	}

	lines := bytes.SplitSeq(buf, []byte{'\n'})

	for line := range lines {

		line = bytes.TrimSpace(line)
		escape := bytes.IndexByte(line, '#')
		if escape != -1 {
			line = line[:escape]
		}

		bs := bytes.Split(line, []byte(" "))

		if len(bs) < 1 {
			continue
		}
		bs = bytes.Split(bs[0], []byte("="))

		if len(bs) < 2 {
			continue
		}

		c.temp[string(bs[0])] = bs[1]

		if c.do.PersistToOSEnv {
			err := os.Setenv(string(bs[0]), string(bs[1]))
			if err != nil {
				var zeroValue T
				return zeroValue, nil
			}
		}
	}

	err := c.scanWithNestedPrefix(&value)

	return value, err
}

func (c *Codec[T]) scanWithNestedPrefix(v *T) error {
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		panic("value is not pointer")
	}

	vt := reflect.ValueOf(v).Elem()
	parent := reflect.TypeOf(v)
	c.scanNestedWithNestedPrefix(parent, vt, "")

	return nil
}

func (c *Codec[T]) scanNestedWithNestedPrefix(
	parent reflect.Type, v reflect.Value, nestedPrefix string,
) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := v.Type().Field(i)

		if structField.Type.Kind() == reflect.Struct && structField.Type != parent {
			nestedName := structField.Tag.Get(string(shared.GetTagNestedName()))
			if nestedName == "-" {
				continue
			}
			if nestedName == "" {
				nestedName = utility.PascalToUpperSnakeCase(structField.Name)
			}
			if nestedPrefix != "" {
				nestedName = nestedPrefix + "_" + nestedName
			}
			c.scanNestedWithNestedPrefix(parent, field, nestedName)
			continue
		}

		var name string
		name = structField.Tag.Get(string(shared.GetTagName()))
		if name == "-" {
			continue
		}
		if name == "" {
			name = utility.PascalToUpperSnakeCase(structField.Name)
		}

		if nestedPrefix != "" {
			sub := nestedPrefix + "_"
			name = sub + name
		}
		name = strings.ToUpper(name)

		var (
			val []byte
			ok  bool
		)

		if c.do.AutomaticEnv {
			if c.do.PreferFileOverEnv {

				val, ok = c.temp[name]
				if !ok {
					r := os.Getenv(name)
					if r != "" {
						val, ok = []byte(r), true
					}
				}
			} else {
				r := os.Getenv(name)
				if r == "" {
					val, ok = c.temp[name]
				} else {
					val, ok = []byte(r), true
				}

			}
		} else {
			val, ok = c.temp[name]
		}

		if !ok || !field.CanSet() {
			continue
		}

		setValue(field, string(val))
	}
}

func (c *Codec[T]) flattenWithNestedPrefix(v T) error {
	vt := reflect.ValueOf(v)
	if vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
	}
	parent := reflect.TypeOf(v)
	c.flattenNestedWithNestedPrefix(parent, vt, "")

	return nil
}

func (c *Codec[T]) flattenNestedWithNestedPrefix(
	parent reflect.Type, v reflect.Value, nestedPrefix string,
) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := v.Type().Field(i)

		if structField.Type.Kind() == reflect.Struct && structField.Type != parent {
			nestedName := structField.Tag.Get(string(shared.GetTagNestedName()))
			if nestedName == "-" {
				continue
			}
			if nestedName == "" {
				nestedName = utility.PascalToUpperSnakeCase(structField.Name)
			}
			if nestedPrefix != "" {
				nestedName = nestedPrefix + "_" + nestedName
			}
			c.scanNestedWithNestedPrefix(parent, field, nestedName)
			continue
		}

		var name string
		name = structField.Tag.Get(string(shared.GetTagName()))
		if name == "-" {
			continue
		}
		if name == "" {
			name = utility.PascalToUpperSnakeCase(structField.Name)
		}

		if nestedPrefix != "" {
			sub := nestedPrefix + "_"
			name = sub + name
		}
		name = strings.ToUpper(name)

		c.temp[name] = parseToBytes(field)
	}
}

func setValue(field reflect.Value, val string) {
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Elem()
	}

	// Basic kinds
	switch field.Kind() {
	case reflect.String:
		field.SetString(val)
	case reflect.Int, reflect.Int64:
		i64, err := strconv.ParseInt(val, 0, 64)
		if err != nil {
			log.Fatalf("convert string to int error: %+v", err)
		}
		field.SetInt(i64)
	case reflect.Float64:
		f64, err := strconv.ParseFloat(val, 64)
		if err != nil {
			log.Fatalf("convert string to float error: %+v", err)
		}
		field.SetFloat(f64)
	case reflect.Bool:
		bVal, err := strconv.ParseBool(val)
		if err != nil {
			log.Fatalf("convert string to bool error: %+v", err)
		}
		field.SetBool(bVal)
	}
}

func setValueAny(field reflect.Value, val any) {
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Elem()
	}

	// Basic kinds
	switch field.Kind() {
	case reflect.String:
		s, ok := val.(string)
		if !ok {
			return
		}
		field.SetString(s)

	case reflect.Int, reflect.Int64:
		i64, ok := val.(int)
		if !ok {
			log.Fatalf("convert string to int error: %+v", ok)
		}
		field.SetInt(int64(i64))
	case reflect.Float64:
		f64, ok := val.(float64)
		if !ok {
			log.Fatalf("convert string to float error: %+v", ok)
		}
		field.SetFloat(f64)
	case reflect.Bool:
		bVal, ok := val.(bool)
		if !ok {
			log.Fatalf("convert string to bool error: %+v", ok)
		}
		field.SetBool(bVal)
	}
}

func parseToBytes(field reflect.Value) []byte {
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Elem()
	}

	// Basic kinds
	switch field.Kind() {
	case reflect.String:
		return []byte(field.String())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(strconv.FormatInt(field.Int(), 10))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte(strconv.FormatUint(field.Uint(), 10))

	case reflect.Float32, reflect.Float64:
		return []byte(strconv.FormatFloat(field.Float(), 'f', -1, 64))

	case reflect.Bool:
		return []byte(strconv.FormatBool(field.Bool()))
	}
	return nil
}
