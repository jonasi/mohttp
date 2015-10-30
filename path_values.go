package mohttp

import (
	"github.com/julienschmidt/httprouter"
	"net/url"
	"strconv"
	"strings"
)

type PathValues struct {
	Params StringSources
	Query  StringSources
}

type StringSources interface {
	Get(string) StringSource
	Has(string) bool
}

type query struct {
	url.Values
}

func (q *query) Get(s string) StringSource {
	return StringSource(q.Values.Get(s))
}

func (q *query) Has(s string) bool {
	return q.Values.Get(s) != ""
}

type params struct {
	httprouter.Params
}

func (p *params) Get(s string) StringSource {
	return StringSource(p.ByName(s))
}

func (p *params) Has(s string) bool {
	return p.ByName(s) != ""
}

type StringSource string

func (v StringSource) String() string   { return string(v) }
func (v StringSource) Bool() bool       { return parseBool(string(v)) }
func (v StringSource) Int() int         { return int(parseInt(string(v), 0)) }
func (v StringSource) Int8() int8       { return int8(parseInt(string(v), 8)) }
func (v StringSource) Int16() int16     { return int16(parseInt(string(v), 16)) }
func (v StringSource) Int32() int32     { return int32(parseInt(string(v), 32)) }
func (v StringSource) Int64() int64     { return int64(parseInt(string(v), 64)) }
func (v StringSource) Uint() uint       { return uint(parseUint(string(v), 0)) }
func (v StringSource) Uint8() uint8     { return uint8(parseUint(string(v), 8)) }
func (v StringSource) Uint16() uint16   { return uint16(parseUint(string(v), 16)) }
func (v StringSource) Uint32() uint32   { return uint32(parseUint(string(v), 32)) }
func (v StringSource) Uint64() uint64   { return uint64(parseUint(string(v), 64)) }
func (v StringSource) Float32() float32 { return float32(parseFloat(string(v), 32)) }
func (v StringSource) Float64() float64 { return float64(parseFloat(string(v), 64)) }

func (v StringSource) SliceString(sep string) []string {
	return strings.Split(string(v), sep)
}

func (v StringSource) SliceBool(sep string) []bool {
	var (
		strs = strings.Split(string(v), sep)
		vals = make([]bool, len(strs))
	)

	for i := range strs {
		vals[i] = parseBool(strs[i])
	}

	return vals
}

func (v StringSource) SliceInt(sep string) []int {
	var (
		strs = strings.Split(string(v), sep)
		vals = make([]int, len(strs))
	)

	for i := range strs {
		vals[i] = int(parseInt(strs[i], 0))
	}

	return vals
}

func (v StringSource) SliceInt8(sep string) []int8 {
	var (
		strs = strings.Split(string(v), sep)
		vals = make([]int8, len(strs))
	)

	for i := range strs {
		vals[i] = int8(parseInt(strs[i], 8))
	}

	return vals
}

func (v StringSource) SliceInt16(sep string) []int16 {
	var (
		strs = strings.Split(string(v), sep)
		vals = make([]int16, len(strs))
	)

	for i := range strs {
		vals[i] = int16(parseInt(strs[i], 16))
	}

	return vals
}

func (v StringSource) SliceInt32(sep string) []int32 {
	var (
		strs = strings.Split(string(v), sep)
		vals = make([]int32, len(strs))
	)

	for i := range strs {
		vals[i] = int32(parseInt(strs[i], 32))
	}

	return vals
}

func (v StringSource) SliceInt64(sep string) []int64 {
	var (
		strs = strings.Split(string(v), sep)
		vals = make([]int64, len(strs))
	)

	for i := range strs {
		vals[i] = int64(parseInt(strs[i], 64))
	}

	return vals
}

func (v StringSource) SliceUint(sep string) []uint {
	var (
		strs = strings.Split(string(v), sep)
		vals = make([]uint, len(strs))
	)

	for i := range strs {
		vals[i] = uint(parseUint(strs[i], 0))
	}

	return vals
}

func (v StringSource) SliceUint8(sep string) []uint8 {
	var (
		strs = strings.Split(string(v), sep)
		vals = make([]uint8, len(strs))
	)

	for i := range strs {
		vals[i] = uint8(parseUint(strs[i], 8))
	}

	return vals
}

func (v StringSource) SliceUint16(sep string) []uint16 {
	var (
		strs = strings.Split(string(v), sep)
		vals = make([]uint16, len(strs))
	)

	for i := range strs {
		vals[i] = uint16(parseUint(strs[i], 16))
	}

	return vals
}

func (v StringSource) SliceUint32(sep string) []uint32 {
	var (
		strs = strings.Split(string(v), sep)
		vals = make([]uint32, len(strs))
	)

	for i := range strs {
		vals[i] = uint32(parseUint(strs[i], 32))
	}

	return vals
}

func (v StringSource) SliceUint64(sep string) []uint64 {
	var (
		strs = strings.Split(string(v), sep)
		vals = make([]uint64, len(strs))
	)

	for i := range strs {
		vals[i] = uint64(parseUint(strs[i], 64))
	}

	return vals
}

func (v StringSource) SliceFloat32(sep string) []float32 {
	var (
		strs = strings.Split(string(v), sep)
		vals = make([]float32, len(strs))
	)

	for i := range strs {
		vals[i] = float32(parseFloat(strs[i], 32))
	}

	return vals
}

func (v StringSource) SliceFloat64(sep string) []float64 {
	var (
		strs = strings.Split(string(v), sep)
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
