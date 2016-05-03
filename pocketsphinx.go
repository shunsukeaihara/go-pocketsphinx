package pocketsphinx

/*
#cgo pkg-config: pocketsphinx sphinxbase
#include <pocketsphinx.h>
#include <err.h>
#include <stdio.h>

cmd_ln_t *default_config(){
    return cmd_ln_parse_r(NULL, ps_args(), 0, NULL, FALSE);
}

int process_raw(ps_decoder_t *ps, char const *data, size_t n_samples, int no_search, int full_utt){
    n_samples /= sizeof(int16);
    return ps_process_raw(ps, (int16 *)data, n_samples, no_search, full_utt);
}
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

//Result is a speech recognition result
type Result struct {
	Text     string    `json:"text"`
	Score    int64     `json:"score"`
	Prob     int64     `json:"prob"`
	Segments []Segment `json:"segments"`
}

//PocketSphinx is a speech recognition decoder object
type PocketSphinx struct {
	ps     *C.ps_decoder_t
	Config Config
}

//NewPocketSphinx creates PocketSphinx instance with Config.
func NewPocketSphinx(config Config) *PocketSphinx {
	psConfig := C.default_config()
	config.SetParams(psConfig)
	if config.DisableInfo {
		path := C.CString("/dev/null")
		defer C.free(unsafe.Pointer(path))
		C.err_set_logfile(path)
	}

	ps := C.ps_init(psConfig)
	C.cmd_ln_free_r(psConfig)

	return &PocketSphinx{ps: ps, Config: config}
}

func (p *PocketSphinx) GetInSpeech() bool {
	return C.ps_get_in_speech(p.ps) != 0
}

//Free releases all resources associated with the PocketSphinx.
func (p *PocketSphinx) Free() {
	C.ps_free(p.ps)
}

//StartUtt starts utterance processing.
func (p *PocketSphinx) StartStream() error {
	ret := C.ps_start_stream(p.ps)
	if ret != 0 {
		return fmt.Errorf("start_stream error:%d", ret)
	}
	return nil
}

//StartUtt starts utterance processing.
func (p *PocketSphinx) StartUtt() error {
	ret := C.ps_start_utt(p.ps)
	if ret != 0 {
		return fmt.Errorf("start_utt error:%d", ret)
	}
	return nil
}

//EndUtt ends utterance processing.
func (p *PocketSphinx) EndUtt() error {
	ret := C.ps_end_utt(p.ps)
	if ret != 0 {
		return fmt.Errorf("end_utt error:%d", ret)
	}
	return nil
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}

//ProcessRaw processes a single channel, 16-bit pcm signal. if noSearch is true, ProcessRaw performs only feature extraction but don't do any recognition yet. if fullUtt is true, this block of data is a full utterance worth of data.
func (p *PocketSphinx) ProcessRaw(raw []byte, noSearch, fullUtt bool) error {
	raw_byte := (*C.char)(unsafe.Pointer(&raw[0]))
	numByte := len(raw)
	processed := C.process_raw(p.ps, raw_byte, C.size_t(numByte), C.int(bool2int(noSearch)), C.int(bool2int(fullUtt)))
	if processed < 0 {
		return fmt.Errorf("process_raw error")
	}
	return nil
}

//GetHyp gets speech recognition result for best hypothesis. If segment is true, result contains word segments in recognized text.
func (p *PocketSphinx) GetHyp(segment bool) (Result, error) {
	var score C.int32
	charp := C.ps_get_hyp(p.ps, &score)
	if charp == nil {
		return Result{}, errors.New("no hypothesis")
	}
	text := C.GoString(charp)
	ret := Result{Text: text, Score: int64(score), Prob: int64(C.ps_get_prob(p.ps))}
	if segment {
		ret.Segments = GetSegments(p.ps)
	}
	return ret, nil
}

func (p *PocketSphinx) getNbestHyp(nbest *C.ps_nbest_t, segment bool) Result {
	var score C.int32
	text := C.GoString(C.ps_nbest_hyp(nbest, &score))
	ret := Result{Text: text, Score: int64(score)}
	if segment {
		ret.Segments = GetSegmentsForNbest(nbest)
	}
	return ret
}

func (p *PocketSphinx) GetNbest(numNbest int, segment bool) []Result {
	ret := make([]Result, 0, numNbest)

	nbestIt := C.ps_nbest(p.ps)
	for {
		if nbestIt == nil {
			break
		}
		if len(ret) == numNbest {
			C.ps_nbest_free(nbestIt)
			break
		}

		hyp := p.getNbestHyp(nbestIt, segment)
		if hyp.Text == "" {
			C.ps_nbest_free(nbestIt)
			break
		}
		ret = append(ret, hyp)
		nbestIt = C.ps_nbest_next(nbestIt)
	}

	return ret
}

func (p *PocketSphinx) ProcessUtt(raw []byte, numNbest int, segment bool) ([]Result, error) {
	ret := make([]Result, 0, numNbest)
	err := p.StartUtt()
	if err != nil {
		return ret, err
	}
	err = p.ProcessRaw(raw, false, true)
	if err != nil {
		return ret, err
	}
	err = p.EndUtt()
	if err != nil {
		return ret, err
	}
	r, err := p.GetHyp(segment)
	if err == nil {
		ret = append(ret, r)
	} else {
		return ret, err
	}

	ret = append(ret, p.GetNbest(numNbest-1, segment)...)
	return ret, nil
}
