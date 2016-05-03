package pocketsphinx

/*
#cgo pkg-config: pocketsphinx sphinxbase
#include <pocketsphinx.h>
*/
import "C"

type segments struct {
	ps *C.ps_decoder_t
	nb *C.ps_nbest_t
}

//Segment represents a word segment contains frame infomations and probabirity
type Segment struct {
	Word       string `json:"word"`
	StartFrame int64  `json:"start_frame"`
	EndFrame   int64  `json:"end_frame"`
	Prob       int64  `json:"log_posterior_probability. "`
	AcScore    int64  `json:"acoustic_score"`
	LmScore    int64  `json:"languagemodel_score"`
	Lmbackoff  int64  `json:"languagemodel_backoff"`
}

// GetSegments returns word segment list for best hypotesis
func GetSegments(ps *C.ps_decoder_t) []Segment {
	s := segments{ps: ps}
	return s.getBestHypSegments()
}

// GetSegments returns word segment list for nbest_t
func GetSegmentsForNbest(nb *C.ps_nbest_t) []Segment {
	s := segments{nb: nb}
	return s.getNbestHypSegments()
}

func (s segments) getBestHypSegments() []Segment {
	segIt := C.ps_seg_iter(s.ps)
	return s.getSegmentsFromIter(segIt)
}

func (s segments) getNbestHypSegments() []Segment {
	segIt := C.ps_nbest_seg(s.nb)
	return s.getSegmentsFromIter(segIt)
}

func (s segments) getSegmentsFromIter(segIt *C.ps_seg_t) []Segment {
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

func (s segments) getCurrentSegment(segIt *C.ps_seg_t) Segment {
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
