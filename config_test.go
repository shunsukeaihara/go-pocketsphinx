package pocketsphinx

import (
	"testing"
)

func TestNewNullConfig(t *testing.T) {
	conf := Config{}
	if conf.Debug.Valid {
		t.Error()
	}
}

func TestNewConfig(t *testing.T) {
	conf := Config{Hmm: ".", Kws_threshold: FloatParam(1e-20)}
	if conf.Hmm != "." {
		t.Error()
	}
}
