package pocketsphinx

/*
#cgo CFLAGS: -I/usr/local/include/pocketsphinx -I/usr/local/include/sphinxbase/
#cgo LDFLAGS: -L/usr/local/lib -lpocketsphinx -lsphinxbase
#include <pocketsphinx.h>
*/
import "C"

import (
	"unsafe"

	"gopkg.in/guregu/null.v3"
)

func IntParam(i int64) null.Int {
	return null.IntFrom(i)
}

func FloatParam(f float64) null.Float {
	return null.FloatFrom(f)
}

type Config struct {
	Hmm           string
	Dict          string
	Lm            string
	Jsgf          string
	Keyphrase     string
	Kws           string
	Kws_threshold null.Float
	Kws_plp       null.Float
	Debug         null.Int
	SamplingRate  null.Int
	DisableInfo   bool
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
