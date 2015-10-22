package pocketsphinx

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestNewPocketSphinx(t *testing.T) {
	conf := Config{Hmm: "./models/en-us/en-us",
		Dict:        "./models/en-us/cmudict-en-us.dict",
		Lm:          "./models/en-us/en-us.lm.bin",
		DisableInfo: true}

	ps := NewPocketSphinx(conf)
	ps.Free()
}

func TestProcessUTT(t *testing.T) {
	conf := Config{Hmm: "./models/en-us/en-us",
		Dict:        "./models/en-us/cmudict-en-us.dict",
		Lm:          "./models/en-us/en-us.lm.bin",
		DisableInfo: false}

	ps := NewPocketSphinx(conf)
	defer ps.Free()
	dat, err := ioutil.ReadFile("./data/goforward.raw")
	results, err := ps.ProcessUtt(dat, 2, false)
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
	defer ps.Free()
	dat, _ := ioutil.ReadFile("./data/goforward.raw")
	ps.StartUtt()
	ps.ProcessRaw(dat, false, false)
	ps.EndUtt()
	r := ps.GetHyp(true)
	if r.Text == "" {
		t.Error("could not recognize")
	}

	fmt.Println(ps.GetNbest(10, true))
}
