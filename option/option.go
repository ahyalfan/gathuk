// Package option
package option

type DecodeOption struct {
	AutomaticEnv      bool // jika true, baca OS environment otomatis
	PersistToOSEnv    bool // jika true, hasil decode disimpan di OS env juga
	PreferFileOverEnv bool // jika true, config file diutamakan dibanding OS env / string
}

type EncodeOption struct {
	AutomaticEnv      bool // jika true, baca OS environment otomatis
	PreferFileOverEnv bool // jika true, config file diutamakan dibanding OS env / string
}

type DecodeOptionApplier interface {
	ApplyDecodeOption(*DecodeOption)
	CheckDecodeOption() bool
}
type EncodeOptionApplier interface {
	ApplyEncodeOption(*EncodeOption)
	CheckEncodeOption() bool
}
