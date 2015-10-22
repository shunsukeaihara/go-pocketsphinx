package pocketsphinx

/*
#cgo CFLAGS: -I/usr/local/include/pocketsphinx -I/usr/local/include/sphinxbase/
#cgo LDFLAGS: -L/usr/local/lib -lpocketsphinx -lsphinxbase
#include <pocketsphinx.h>
#include <err.h>

cmd_ln_t *default_config(){
    return cmd_ln_parse_r(NULL, ps_args(), 0, NULL, FALSE);
}

int process_raw(ps_decoder_t *ps, char const *data, size_t n_samples, int no_search, int full_utt){
    n_samples /= sizeof(int16);
    return ps_process_raw(ps, (int16 *)data, n_samples, no_search, full_utt);
}


//nbest


*/
import "C"

import (
	"fmt"
	"unsafe"
)

//Result is a speech recognition result
type Result struct {
	Text     string
	Score    int64
	Segments []Segment
}

//PocketSphinx is a speech recognition decoder object
type PocketSphinx struct {
	ps     *C.ps_decoder_t
	Config Config
	conf   *C.cmd_ln_t
}

//NewPocketSphinx creates PocketSphinx instance with Config.
func NewPocketSphinx(config Config) *PocketSphinx {
	var psConfig *C.cmd_ln_t
	psConfig = C.default_config()
	config.SetParams(psConfig)

	if config.DisableInfo {
		path := C.CString("/dev/null")
		defer C.free(unsafe.Pointer(path))
		C.err_set_logfile(path)
	}

	var ps *C.ps_decoder_t
	ps = C.ps_init(psConfig)

	return &PocketSphinx{ps: ps, Config: config, conf: psConfig}
}

//Free releases all resources associated with the PocketSphinx.
func (p *PocketSphinx) Free() {
	C.ps_free(p.ps)
	C.cmd_ln_free_r(p.conf)
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
	raw_byte := (*C.char)(unsafe.Pointer(&raw))
	numByte := len(raw)
	errorcode := C.process_raw(p.ps, raw_byte, C.size_t(numByte), C.int(bool2int(noSearch)), C.int(bool2int(fullUtt)))
	if errorcode < 0 {
		return fmt.Errorf("process_raw error:%d", errorcode)
	}
	return nil
}

//GetHyp gets speech recognition result for best hypothesis. If segment is true, Result contains word segments in recognized text.
func (p *PocketSphinx) GetHyp(segment bool) Result {
	var score C.int32
	text := C.GoString(C.ps_get_hyp(p.ps, &score))
	ret := Result{Text: text, Score: int64(score)}
	if segment {
		ret.Segments = GetSegments(p.ps)
	}
	return ret
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

	nbestIt := C.ps_nbest(p.ps, 0, -1, nil, nil)
	for {
		if nbestIt == nil {
			break
		}

		hyp := p.getNbestHyp(nbestIt, segment)
		if hyp.Text == "" {
			C.ps_nbest_free(nbestIt)
			break
		}
		ret = append(ret, hyp)
		if len(ret) == numNbest {
			C.ps_nbest_free(nbestIt)
			break
		}
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
	err = p.ProcessRaw(raw, false, false)
	if err != nil {
		return ret, err
	}
	err = p.EndUtt()
	if err != nil {
		return ret, err
	}
	r := p.GetHyp(segment)
	if r.Text != "" {
		ret = append(ret, r)
	}

	ret = append(ret, p.GetNbest(numNbest-1, segment)...)
	return ret, nil
}
