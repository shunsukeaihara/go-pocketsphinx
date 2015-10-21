package pocketsphinx

/*
#cgo CFLAGS: -I/usr/local/include/pocketsphinx -I/usr/local/include/sphinxbase/
#cgo LDFLAGS: -L/usr/local/lib -lpocketsphinx -lsphinxbase
#include <pocketsphinx.h>
*/
import "C"

type Segments struct {
	ps *C.ps_decoder_t
	nb *C.ps_nbest_t
}

type Segment struct {
	Word       string
	StartFrame int64
	EndFrame   int64
	Prob       int64
	AcProb     int64
	LmProb     int64
	LbackProb  int64
}

func NewSegments(ps *C.ps_decoder_t) Segments {
	return Segments{ps: ps}
}

func NewSegmentsForNbest(nb *C.ps_nbest_t) Segments {
	return Segments{nb: nb}
}

func (s Segments) GetBesyHypSegments() []Segment {
	var score C.int32
	segIt := C.ps_seg_iter(s.ps, &score)
	return s.getSegmentsFromIter(segIt)
}

func (s Segments) GetNbestHypSegments() []Segment {
	var score C.int32
	segIt := C.ps_nbest_seg(s.nb, &score)
	return s.getSegmentsFromIter(segIt)
}

func (s Segments) getSegmentsFromIter(segIt *C.ps_seg_t) []Segment {
	ret := make([]Segment, 0, 10)
	for {
		if segIt == nil {
			break
		}
		seg := s.getCurrentSegment(segIt)
		ret = append(ret, seg)
		segIt = C.ps_seg_next(segIt)
	}
	return ret
}

func (s Segments) getCurrentSegment(segIt *C.ps_seg_t) Segment {
	var start, end C.int
	word := C.GoString(C.ps_seg_word(segIt))
	C.ps_seg_frames(segIt, &start, &end)

	var acousticProb, lmProb, lbackProb C.int32
	segProb := C.ps_seg_prob(segIt, &acousticProb, &lmProb, &lbackProb)

	seg := Segment{word, int64(start), int64(end),
		int64(segProb), int64(acousticProb), int64(lmProb), int64(lbackProb),
	}
	return seg

}
