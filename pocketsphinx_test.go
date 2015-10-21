package pocketsphinx

import (
	"testing"
)

func TestNewPocketSphinx(t *testing.T) {
	conf := Config{Hmm: "./models/en-us/en-us",
		Dict:        "./models/en-us/cmudict-en-us.dict",
		Lm:          "./models/en-us/en-us.lm.bin",
		DisableInfo: true}

	NewPocketSphinx(conf)
}
