package pocketsphinx

import (
	"strings"
	"testing"
)

const (
	yamlStr = `hmm: "aaa"
dict: "bbb"
lm: "lm"
kws_threshold: "1e-20"
debug: "2"
`
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

func TestUnmarshallYaml(t *testing.T) {
	r := strings.NewReader(yamlStr)
	_, err := Load(r)
	if err != nil {
		t.Error(err)
	}
}
