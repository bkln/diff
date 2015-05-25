package puzz

import "unsafe"
import "fmt"
import "crypto/rand"

// #cgo LDFLAGS: -lpuzzle
// #include <stdio.h>
// #include "/usr/include/puzzle.h"
import "C"

var (
	K = 10
	N = 100
)

func GetCVec(filename string) ([]byte, error) {
	var c C.PuzzleContext
	var cvec C.PuzzleCvec
	// Initialization is very cheap.
	C.puzzle_init_context(&c)
	C.puzzle_init_cvec(&c, &cvec)

	err_num := C.puzzle_fill_cvec_from_file(&c, &cvec, C.CString(filename))
	if err_num != 0 {
		return nil, fmt.Errorf("error reading file %s", filename)
	}

	bb := C.GoBytes(unsafe.Pointer(cvec.vec), C.int(cvec.sizeof_vec))
	b := make([]byte, len(bb), 600)
	copy(b, bb)

	C.puzzle_free_cvec(&c, &cvec)
	C.puzzle_free_context(&c)
	if got, want := len(b), 544; got != want {
		return nil, fmt.Errorf("vector has the wrong length: got %d; want %d", got, want)
	}
	return b, nil
}

func GetRandCVec() []byte {
	rb := make([]byte, 544)
	_, err := rand.Read(rb)

	if err != nil {
		fmt.Println(err)
	}
	return rb

}

func GetSubCVec(cvec []byte) [][]byte {
	result := [][]byte{}
	// Make sure N < 544 - K + 1
	for i := 0; i < N; i++ {
		sub := make([]byte, K)
		copy(sub, cvec[i:i+K])
		result = append(result, sub)
	}
	return result
}

func CompareCVecs(cbvec1, cbvec2 []byte) float64 {
	var c C.PuzzleContext
	var cvec1, cvec2 C.PuzzleCvec
	// Initialization is very cheap.
	C.puzzle_init_context(&c)

	cvec1.sizeof_vec = C.size_t(len(cbvec1))
	cvec1.vec = (*C.schar)(unsafe.Pointer(&cbvec1[0]))

	cvec2.sizeof_vec = C.size_t(len(cbvec2))
	cvec2.vec = (*C.schar)(unsafe.Pointer(&cbvec2[0]))

	// Diff
	diff := C.puzzle_vector_normalized_distance(&c, &cvec1, &cvec2, 0)

	C.puzzle_free_context(&c)
	return float64(diff)
}
