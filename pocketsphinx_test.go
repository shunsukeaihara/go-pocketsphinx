package pocketsphinx

import (
	"io/ioutil"
	"testing"
)

func TestNewPocketSphinx(t *testing.T) {
	conf := Config{Hmm: "./models/en-us/en-us",
		Dict:        "./models/en-us/cmudict-en-us.dict",
		Lm:          "./models/en-us/en-us.lm.bin",
		DisableInfo: true}

	NewPocketSphinx(conf)
}

func TestProcessUTT(t *testing.T) {
	conf := Config{Hmm: "./models/en-us/en-us",
		Dict:        "./models/en-us/cmudict-en-us.dict",
		Lm:          "./models/en-us/en-us.lm.bin",
		DisableInfo: false}

	ps := NewPocketSphinx(conf)

	dat, err := ioutil.ReadFile("./data/goforward.raw")
	results, err := ps.ProcessUtt(dat, 2)
	if err != nil {
		t.Error(err)
	}
	if len(results) == 0 {
		t.Error(results)
	}
}

func TestProcessRaw(t *testing.T) {
	conf := Config{Hmm: "./models/en-us/en-us",
		Dict:        "./models/en-us/cmudict-en-us.dict",
		Lm:          "./models/en-us/en-us.lm.bin",
		DisableInfo: true}

	ps := NewPocketSphinx(conf)

	dat, _ := ioutil.ReadFile("./data/goforward.raw")
	ps.StartUtt()
	ps.ProcessRaw(dat, false, true)
	ps.EndUtt()
	r := ps.GetHyp()
	if r.Text == "" {
		t.Error("could not recognize")
	}
}
