package mohttp

import (
	"github.com/julienschmidt/httprouter"
	"net/url"
	"strconv"
	"strings"
)

type PathValues struct {
	Params StringValues
	Query  StringValues
}

type StringValues interface {
	String(string) string
	Bool(string) bool
	Int(string) int
	Int8(string) int8
	Int16(string) int16
	Int32(string) int32
	Int64(string) int64
	Uint(string) uint
	Uint8(string) uint8
	Uint16(string) uint16
	Uint32(string) uint32
	Uint64(string) uint64
	Float32(string) float32
	Float64(string) float64
	SliceInt(string, string) []int
	SliceInt8(string, string) []int8
	SliceInt16(string, string) []int16
	SliceInt32(string, string) []int32
	SliceInt64(string, string) []int64
	SliceUint(string, string) []uint
	SliceUint8(string, string) []uint8
	SliceUint16(string, string) []uint16
	SliceUint32(string, string) []uint32
	SliceUint64(string, string) []uint64
	SliceFloat32(string, string) []float32
	SliceFloat64(string, string) []float64
	Has(string) bool
}

func query(v url.Values) StringValues {
	return stringValues(func(k string) string { return v.Get(k) })
}

func params(p httprouter.Params) StringValues {
	return stringValues(func(k string) string { return p.ByName(k) })
}

type stringValues func(string) string

func (v stringValues) Has(k string) bool        { return v(k) != "" }
func (v stringValues) String(k string) string   { return v(k) }
func (v stringValues) Bool(k string) bool       { return parseBool(v(k)) }
func (v stringValues) Int(k string) int         { return int(parseInt(v(k), 0)) }
func (v stringValues) Int8(k string) int8       { return int8(parseInt(v(k), 8)) }
func (v stringValues) Int16(k string) int16     { return int16(parseInt(v(k), 16)) }
func (v stringValues) Int32(k string) int32     { return int32(parseInt(v(k), 32)) }
func (v stringValues) Int64(k string) int64     { return int64(parseInt(v(k), 64)) }
func (v stringValues) Uint(k string) uint       { return uint(parseUint(v(k), 0)) }
func (v stringValues) Uint8(k string) uint8     { return uint8(parseUint(v(k), 8)) }
func (v stringValues) Uint16(k string) uint16   { return uint16(parseUint(v(k), 16)) }
func (v stringValues) Uint32(k string) uint32   { return uint32(parseUint(v(k), 32)) }
func (v stringValues) Uint64(k string) uint64   { return uint64(parseUint(v(k), 64)) }
func (v stringValues) Float32(k string) float32 { return float32(parseFloat(v(k), 32)) }
func (v stringValues) Float64(k string) float64 { return float64(parseFloat(v(k), 64)) }

func (v stringValues) SliceString(k, sep string) []string {
	return strings.Split(v(k), sep)
}

func (v stringValues) SliceBool(k, sep string) []bool {
	var (
		strs = strings.Split(v(k), sep)
		vals = make([]bool, len(strs))
	)

	for i := range strs {
		vals[i] = parseBool(strs[i])
	}

	return vals
}

func (v stringValues) SliceInt(k, sep string) []int {
	var (
		strs = strings.Split(v(k), sep)
		vals = make([]int, len(strs))
	)

	for i := range strs {
		vals[i] = int(parseInt(strs[i], 0))
	}

	return vals
}

func (v stringValues) SliceInt8(k, sep string) []int8 {
	var (
		strs = strings.Split(v(k), sep)
		vals = make([]int8, len(strs))
	)

	for i := range strs {
		vals[i] = int8(parseInt(strs[i], 8))
	}

	return vals
}

func (v stringValues) SliceInt16(k, sep string) []int16 {
	var (
		strs = strings.Split(v(k), sep)
		vals = make([]int16, len(strs))
	)

	for i := range strs {
		vals[i] = int16(parseInt(strs[i], 16))
	}

	return vals
}

func (v stringValues) SliceInt32(k, sep string) []int32 {
	var (
		strs = strings.Split(v(k), sep)
		vals = make([]int32, len(strs))
	)

	for i := range strs {
		vals[i] = int32(parseInt(strs[i], 32))
	}

	return vals
}

func (v stringValues) SliceInt64(k, sep string) []int64 {
	var (
		strs = strings.Split(v(k), sep)
		vals = make([]int64, len(strs))
	)

	for i := range strs {
		vals[i] = int64(parseInt(strs[i], 64))
	}

	return vals
}

func (v stringValues) SliceUint(k, sep string) []uint {
	var (
		strs = strings.Split(v(k), sep)
		vals = make([]uint, len(strs))
	)

	for i := range strs {
		vals[i] = uint(parseUint(strs[i], 0))
	}

	return vals
}

func (v stringValues) SliceUint8(k, sep string) []uint8 {
	var (
		strs = strings.Split(v(k), sep)
		vals = make([]uint8, len(strs))
	)

	for i := range strs {
		vals[i] = uint8(parseUint(strs[i], 8))
	}

	return vals
}

func (v stringValues) SliceUint16(k, sep string) []uint16 {
	var (
		strs = strings.Split(v(k), sep)
		vals = make([]uint16, len(strs))
	)

	for i := range strs {
		vals[i] = uint16(parseUint(strs[i], 16))
	}

	return vals
}

func (v stringValues) SliceUint32(k, sep string) []uint32 {
	var (
		strs = strings.Split(v(k), sep)
		vals = make([]uint32, len(strs))
	)

	for i := range strs {
		vals[i] = uint32(parseUint(strs[i], 32))
	}

	return vals
}

func (v stringValues) SliceUint64(k, sep string) []uint64 {
	var (
		strs = strings.Split(v(k), sep)
		vals = make([]uint64, len(strs))
	)

	for i := range strs {
		vals[i] = uint64(parseUint(strs[i], 64))
	}

	return vals
}

func (v stringValues) SliceFloat32(k, sep string) []float32 {
	var (
		strs = strings.Split(v(k), sep)
		vals = make([]float32, len(strs))
	)

	for i := range strs {
		vals[i] = float32(parseFloat(strs[i], 32))
	}

	return vals
}

func (v stringValues) SliceFloat64(k, sep string) []float64 {
	var (
		strs = strings.Split(v(k), sep)
		vals = make([]float64, len(strs))
	)

	for i := range strs {
		vals[i] = float64(parseUint(strs[i], 64))
	}

	return vals
}

func parseInt(v string, bits int) int64 {
	iv, _ := strconv.ParseInt(v, 10, bits)
	return iv
}

func parseUint(v string, bits int) uint64 {
	iv, _ := strconv.ParseUint(v, 10, bits)
	return iv
}

func parseFloat(v string, bits int) float64 {
	iv, _ := strconv.ParseFloat(v, bits)
	return iv
}

func parseBool(v string) bool {
	iv, _ := strconv.ParseBool(v)
	return iv
}
