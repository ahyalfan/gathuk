// Package gathuk
package gathuk

type Tag string

func (t *Tag) Set(v Tag) {
	*t = v
}

func (t Tag) Get() Tag {
	return t
}
