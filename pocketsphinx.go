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

type Result struct {
	Result string
	Score  int64
}

var noReult = make([]Result, 0)

type PocketSphinx struct {
	ps     *C.ps_decoder_t
	Config Config
}

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

	return &PocketSphinx{ps: ps, Config: config}
}

func (p *PocketSphinx) StartUtt() error {
	ret := C.ps_start_utt(p.ps)
	if ret != 0 {
		return fmt.Errorf("start_utt error:%d", ret)
	}
	return nil
}

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

func (p *PocketSphinx) ProcessRaw(raw []byte, noSearch, fullUtt bool) error {
	raw_byte := (*C.char)(unsafe.Pointer(&raw))
	numByte := len(raw)
	errorcode := C.process_raw(p.ps, raw_byte, C.size_t(numByte), C.int(bool2int(noSearch)), C.int(bool2int(fullUtt)))
	if errorcode != 0 {
		return fmt.Errorf("process_raw error:%d", errorcode)
	}
	return nil
}

func (p *PocketSphinx) getHyp() Result {
	var score C.int32
	ret := C.GoString(C.ps_get_hyp(p.ps, &score))
	return Result{ret, int64(score)}
}

func (p *PocketSphinx) getNbestHyp(nbest *C.ps_nbest_t) Result {
	var score C.int32
	ret := C.GoString(C.ps_nbest_hyp(nbest, &score))
	return Result{ret, int64(score)}
}

func (p *PocketSphinx) getNbest(numNbest int) []Result {
	ret := make([]Result, 0, numNbest)

	nbestIt := C.ps_nbest(p.ps, 0, -1, nil, nil)
	for {
		if nbestIt == nil {
			break
		}
		ret := append(ret, p.getNbestHyp(nbestIt))
		if len(ret) == numNbest {
			break
		}
	}
	return ret
}

func (p *PocketSphinx) ProcessUtt(raw []byte, numNbest int) ([]Result, error) {
	err := p.StartUtt()
	if err != nil {
		return noReult, err
	}
	err = p.ProcessRaw(raw, false, true)
	if err != nil {
		return noReult, err
	}
	err = p.EndUtt()
	if err != nil {
		return noReult, err
	}
	ret := p.getNbest(numNbest)
	return ret, nil
}
