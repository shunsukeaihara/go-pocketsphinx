package pocketsphinx

/*

#cgo CFLAGS: -I/usr/local/include/pocketsphinx -I/usr/local/include/sphinxbase/
#cgo LDFLAGS: -L/usr/local/lib -lpocketsphinx -lsphinxbase
#include <pocketsphinx.h>

cmd_ln_t *default_config(){
    return cmd_ln_parse_r(NULL, ps_args(), 0, NULL, FALSE);
}


*/
import "C"

type Pocketsphinx struct {
	ps     *C.ps_decoder_t
	Config Config
}

func NewPocketSphinx(config Config) *Pocketsphinx {
	var psConfig *C.cmd_ln_t
	psConfig = C.default_config()
	config.SetParams(psConfig)
	return nil
}
