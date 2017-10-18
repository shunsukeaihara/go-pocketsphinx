package pocketsphinx

/*
#cgo pkg-config: pocketsphinx sphinxbase
#include <pocketsphinx.h>
*/
import "C"

import (
	"io"
	"io/ioutil"
	"strconv"
	"unsafe"

	"gopkg.in/guregu/null.v3"
	yaml "gopkg.in/yaml.v2"
)

type NullInt null.Int
type NullFloat null.Float

func IntParam(i int64) NullInt {
	return NullInt(null.IntFrom(i))
}

func FloatParam(f float64) NullFloat {
	return NullFloat(null.FloatFrom(f))
}

func (f *NullFloat) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		f.Valid = false
		return nil
	}
	var err error
	f.Float64, err = strconv.ParseFloat(str, 64)
	f.Valid = err == nil
	return err
}

func (f *NullInt) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		f.Valid = false
		return nil
	}
	var err error
	f.Int64, err = strconv.ParseInt(str, 10, 64)
	f.Valid = err == nil
	return err
}

func Load(r io.Reader) (Config, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

type Config struct {
	Hmm           string
	Dict          string
	Lm            string
	Jsgf          string
	Bestpath      string
	Beam          NullFloat
	Wbeam         NullFloat
	Keyphrase     string
	Kws           string
	Kws_threshold NullFloat
	Kws_plp       NullFloat
	Debug         NullInt
	SamplingRate  NullFloat
	DisableInfo   bool
	Language      string
}

func (c Config) SetParams(psConfig *C.cmd_ln_t) {
	if c.Hmm != "" {
		setStringParam(psConfig, "-hmm", c.Hmm)
	}
	if c.Dict != "" {
		setStringParam(psConfig, "-dict", c.Dict)
	}
	if c.Lm != "" {
		setStringParam(psConfig, "-lm", c.Lm)
	}
	if c.Jsgf != "" {
		setStringParam(psConfig, "-jsgf", c.Jsgf)
	}
	if c.Bestpath == "no" {
		setIntParam(psConfig, "-bestpath", 0)
	}
	if c.Keyphrase != "" {
		setStringParam(psConfig, "-keyphrase", c.Keyphrase)
	}
	if c.Kws != "" {
		setStringParam(psConfig, "-kws", c.Kws)
	}
	if c.Kws_threshold.Valid {
		setFloatParam(psConfig, "-kws_threshold", c.Kws_threshold.Float64)
	}
	if c.Kws_plp.Valid {
		setFloatParam(psConfig, "-kws_plp", c.Kws_plp.Float64)
	}
	if c.Beam.Valid {
		setFloatParam(psConfig, "-beam", c.Beam.Float64)
	}
	if c.Wbeam.Valid {
		setFloatParam(psConfig, "-wbeam", c.Wbeam.Float64)
	}
	if c.SamplingRate.Valid {
		setFloatParam(psConfig, "-samprate", c.SamplingRate.Float64)
	}

	if c.Debug.Valid {
		setIntParam(psConfig, "-debug", c.Debug.Int64)
	}
}

func setStringParam(psConfig *C.cmd_ln_t, key, val string) {
	keyPtr := C.CString(key)
	defer C.free(unsafe.Pointer(keyPtr))
	valPtr := C.CString(val)
	defer C.free(unsafe.Pointer(valPtr))
	C.cmd_ln_set_str_r(psConfig, keyPtr, valPtr)
}

func setFloatParam(psConfig *C.cmd_ln_t, key string, val float64) {
	keyPtr := C.CString(key)
	defer C.free(unsafe.Pointer(keyPtr))
	C.cmd_ln_set_float_r(psConfig, keyPtr, C.double(val))
}

func setIntParam(psConfig *C.cmd_ln_t, key string, val int64) {
	keyPtr := C.CString(key)
	defer C.free(unsafe.Pointer(keyPtr))
	C.cmd_ln_set_int_r(psConfig, keyPtr, C.long(val))
}

func (c Config) Mode() string {
	if c.Lm != "" {
		return "dictation"
	} else if c.Jsgf != "" {
		return "grammer"
	} else if c.Kws != "" || c.Keyphrase != "" {
		return "spotting"
	} else {
		return ""
	}
}
