package pocketsphinx

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestProcessUTT(t *testing.T) {
	conf := Config{Hmm: "./models/en-us/en-us",
		Dict:        "./models/en-us/cmudict-en-us.dict",
		Lm:          "./models/en-us/en-us.lm.bin",
		DisableInfo: true}

	ps := NewPocketSphinx(conf)
	defer ps.Free()
	ps.StartStream()
	dat, err := ioutil.ReadFile("./data/goforward.raw")
	results, err := ps.ProcessUtt(dat, 2, true)
	if err != nil {
		t.Error(err)
		return
	}
	if len(results) == 0 {
		t.Error(results)
		return
	}
	if results[0].Text != "go forward ten meters" {
		t.Error("could not recognize", results[0].Text)
	}
}

func TestProcessRaw(t *testing.T) {
	conf := Config{Hmm: "./models/en-us/en-us",
		Dict: "./models/en-us/cmudict-en-us.dict",
		Lm:   "./models/en-us/en-us.lm.bin",
		//Beam:        FloatParam(1e-400),
		//Wbeam:       FloatParam(1e-400),
		//Bestpath:    "no",
		DisableInfo: true}

	ps := NewPocketSphinx(conf)
	defer ps.Free()
	ps.StartStream()
	dat, _ := ioutil.ReadFile("./data/something.raw")
	err := ps.StartUtt()
	if err != nil {
		t.Error(err)
	}
	err = ps.ProcessRaw(dat, false, false)
	if err != nil {
		t.Error(err)
	}
	err = ps.EndUtt()
	if err != nil {
		t.Error(err)
	}
	r, err := ps.GetHyp(false)
	if err != nil {
		t.Error(err)
	}
	if r.Text != "go somewhere and do something" {
		t.Error("could not recognize", r.Text, r.Score)
	}
	dat, _ = ioutil.ReadFile("./data/goforward.raw")
	err = ps.StartUtt()
	if err != nil {
		t.Error(err)
	}
	err = ps.ProcessRaw(dat, false, true)
	if err != nil {
		t.Error(err)
	}
	err = ps.EndUtt()
	if err != nil {
		t.Error(err)
	}
	r, _ = ps.GetHyp(true)

	if r.Text != "go forward ten meters" {
		t.Error("could not recognize", r.Text, r.Score)
	}

}

func TestProcessRawIncremental(t *testing.T) {
	conf := Config{Hmm: "./models/en-us/en-us",
		Dict:        "./models/en-us/cmudict-en-us.dict",
		Lm:          "./models/en-us/en-us.lm.bin",
		DisableInfo: true}
	ps := NewPocketSphinx(conf)
	defer ps.Free()
	file, _ := os.Open("./data/goforward.raw")
	defer file.Close()
	buf := make([]byte, 1024)
	var lastr Result
	ps.StartUtt()
	for {
		size, err := file.Read(buf)
		if err != nil {
			break
		}
		ps.ProcessRaw(buf[:size], false, false)
		r, err := ps.GetHyp(false)
		if err == nil {
			lastr = r
		}
	}
	ps.EndUtt()
	if lastr.Text != "go forward ten meters" {
		t.Error("recognition failed")
	}
}

func TestWordSpottingUtt(t *testing.T) {
	conf := Config{Hmm: "./models/en-us/en-us",
		Dict:        "./models/en-us/cmudict-en-us.dict",
		Keyphrase:   "forward",
		DisableInfo: true}

	ps := NewPocketSphinx(conf)
	defer ps.Free()
	dat, _ := ioutil.ReadFile("./data/goforward.raw")
	ps.StartUtt()
	ps.ProcessRaw(dat, false, true)
	ps.EndUtt()
	r, _ := ps.GetHyp(false)
	if r.Text != "forward" {
		t.Error("could not recognize", r.Text)
	}
}
func TestWordSpotting(t *testing.T) {
	conf := Config{Hmm: "./models/en-us/en-us",
		Dict:        "./models/en-us/cmudict-en-us.dict",
		Keyphrase:   "forward",
		DisableInfo: true}
	ps := NewPocketSphinx(conf)
	defer ps.Free()
	file, _ := os.Open("./data/goforward.raw")
	defer file.Close()
	c := 0
	buf := make([]byte, 1024)
	ps.StartUtt()
	for {
		size, err := file.Read(buf)
		if err != nil {
			break
		}
		ps.ProcessRaw(buf[:size], false, false)
		r, err := ps.GetHyp(false)
		if err == nil && r.Text == "forward" {
			c += 1
			ps.EndUtt()
			ps.StartUtt()
		}
	}
	ps.EndUtt()
	if c == 0 {
		t.Error("keyphrase not found")
	}
}

func TestEndUttErr(t *testing.T) {
	conf := Config{Hmm: "./models/en-us/en-us",
		Dict:        "./models/en-us/cmudict-en-us.dict",
		Keyphrase:   "forward",
		DisableInfo: true}
	ps := NewPocketSphinx(conf)
	defer ps.Free()
	if ps.EndUtt() == nil {
		t.Error("call EndUtt befoer StartUtt must raise error")
	}
}

func TestStartUttErr(t *testing.T) {
	conf := Config{Hmm: "./models/en-us/en-us",
		Dict:        "./models/en-us/cmudict-en-us.dict",
		Keyphrase:   "forward",
		DisableInfo: true}
	ps := NewPocketSphinx(conf)
	defer ps.Free()
	if err := ps.StartUtt(); err != nil {
		t.Error(err)
		return
	}
	if ps.StartUtt() == nil {
		t.Error("call StartUtt after calling StartUtt without ErrUtt, must raise error")
	}
}
