// Package hotspot provides an interface to HotSpot, a thermal model and
// simulator used in computer-architecture studies.
//
// http://lava.cs.virginia.edu/HotSpot/
package hotspot

// #cgo CFLAGS: -Ibuild
// #cgo LDFLAGS: -Wl,-no_compact_unwind -Lbuild -lhotspot
// #include <stdlib.h>
// #include "hotspot.h"
import "C"

import "unsafe"

// Model represents the block variant of the HotSpot model. The thermal system
// under consideration is as follows:
//
//     diag(C) * dT/dt + G * (T - Tamb) = P
//
// where
//
//     Cores is the number of cores (active thermal nodes),
//     Nodes is the number of thermal nodes (4 * Cores + 12),
//
//     T is a Nodes-element vector of temperature,
//     P is a Cores-element vector of power,
//     C is a Nodes-element vector of capacitance, and
//     G is a (Nodes x Nodes) matrix of conductance.
type Model struct {
	Cores uint16
	Nodes uint16

	C []float64
	G []float64
}

// Load returns the block HotSpot model corresponding to the given floorplan
// file, configuration file, and parameter line. The parameter line bears the
// same meaning as the command-line arguments of the HotSpot tool. The names of
// parameters should not include dashes in front of them; for instance, params
// can be "t_chip 0.00015 k_chip 100.0".
func Load(floorplan string, config string, params string) *Model {
	cfloorplan := C.CString(floorplan)
	defer C.free(unsafe.Pointer(cfloorplan))

	cconfig := C.CString(config)
	defer C.free(unsafe.Pointer(cconfig))

	cparams := C.CString(params)
	defer C.free(unsafe.Pointer(cparams))

	h := C.newHotSpot(cfloorplan, cconfig, cparams)
	defer C.freeHotSpot(h)

	cc := uint16(h.cores)
	nc := uint32(h.nodes)

	m := &Model{
		Cores: cc,
		Nodes: uint16(nc),

		C: make([]float64, nc),
		G: make([]float64, nc*nc),
	}

	C.copyC(h, (*C.double)(unsafe.Pointer(&m.C[0])))
	C.copyG(h, (*C.double)(unsafe.Pointer(&m.G[0])))

	return m
}
