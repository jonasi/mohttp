package mohttp

import (
	"github.com/julienschmidt/httprouter"
	"net/url"
	"strconv"
)

type PathValues struct {
	Params PathValueStore
	Query  PathValueStore
}

type PathValueStore interface {
	Get(string) PathValue
	Has(string) bool
}

type queryStore struct {
	url.Values
}

func (q *queryStore) Get(s string) PathValue {
	return PathValue(q.Values.Get(s))
}

func (q *queryStore) Has(s string) bool {
	return q.Values.Get(s) != ""
}

type paramsStore struct {
	httprouter.Params
}

func (p *paramsStore) Get(s string) PathValue {
	return PathValue(p.ByName(s))
}

func (p *paramsStore) Has(s string) bool {
	return p.ByName(s) != ""
}

type PathValue string

func (v PathValue) String() string { return string(v) }

func (v PathValue) Bool() bool {
	iv, _ := strconv.ParseBool(string(v))
	return iv
}

func (v PathValue) Int() int {
	iv, _ := strconv.ParseInt(string(v), 10, 0)
	return int(iv)
}

func (v PathValue) Int8() int8 {
	iv, _ := strconv.ParseInt(string(v), 10, 8)
	return int8(iv)
}

func (v PathValue) Int16() int16 {
	iv, _ := strconv.ParseInt(string(v), 10, 16)
	return int16(iv)
}

func (v PathValue) Int32() int32 {
	iv, _ := strconv.ParseInt(string(v), 10, 32)
	return int32(iv)
}

func (v PathValue) Int64() int64 {
	iv, _ := strconv.ParseInt(string(v), 10, 64)
	return iv
}

func (v PathValue) Uint() uint {
	iv, _ := strconv.ParseUint(string(v), 10, 0)
	return uint(iv)
}

func (v PathValue) Uint8() uint8 {
	iv, _ := strconv.ParseUint(string(v), 10, 8)
	return uint8(iv)
}

func (v PathValue) Uint16() uint16 {
	iv, _ := strconv.ParseUint(string(v), 10, 16)
	return uint16(iv)
}

func (v PathValue) Uint32() uint32 {
	iv, _ := strconv.ParseUint(string(v), 10, 32)
	return uint32(iv)
}

func (v PathValue) Uint64() uint64 {
	iv, _ := strconv.ParseUint(string(v), 10, 64)
	return iv
}

func (v PathValue) Float32() float32 {
	iv, _ := strconv.ParseFloat(string(v), 32)
	return float32(iv)
}

func (v PathValue) Float64() float64 {
	iv, _ := strconv.ParseFloat(string(v), 64)
	return iv
}
