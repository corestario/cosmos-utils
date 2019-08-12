package storewrapper

import (
	"encoding/binary"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/types"
)

var _ types.KVStore = KVStore{}

const (
	DefaultChunkSize = 1024 * 5
)

type KVStore struct {
	types.KVStore
	ChunkSize int
}

// Create a wrapper around existing store
func NewKVStore(store types.KVStore, chunkSize int) *KVStore {
	if chunkSize < 16 || chunkSize > DefaultChunkSize {
		chunkSize = DefaultChunkSize
	}
	return &KVStore{
		KVStore:   store,
		ChunkSize: chunkSize,
	}
}

// Get all parts of stored value
func (k *KVStore) GetW(keySrc []byte) (res []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			nErr, nOk := r.(error)
			if nOk {
				err = nErr
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	parts := 0
	res = make([]byte, 0)
	b := make([]byte, 8)

	binary.LittleEndian.PutUint64(b, uint64(parts))
	key := append(keySrc, b...)
	for k.Has(key) {
		ba := k.Get(key)
		res = append(res, ba...)

		parts++
		binary.LittleEndian.PutUint64(b, uint64(parts))

		key = append(keySrc, b...)
	}

	return
}

// Check if store has value; return number of parts, hasValue, error
func (k *KVStore) HasW(keySrc []byte) (parts int, ok bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			nErr, nOk := r.(error)
			if nOk {
				err = nErr
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	parts = 0
	b := make([]byte, 8)

	binary.LittleEndian.PutUint64(b, uint64(parts))
	key := append(keySrc, b...)

	for k.Has(key) {
		parts++
		binary.LittleEndian.PutUint64(b, uint64(parts))

		key = append(keySrc, b...)
	}

	ok = parts > 0

	return
}

// Add or update value in store; if update - delete old parts; return parts stored, error
func (k *KVStore) SetW(keySrc, valueSrc []byte) (parts int, err error) {
	defer func() {
		if r := recover(); r != nil {
			nErr, nOk := r.(error)
			if nOk {
				err = nErr
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	chunkSize := k.ChunkSize
	lenVal := len(valueSrc)

	parts = 0
	b := make([]byte, 8)

	binary.LittleEndian.PutUint64(b, uint64(parts))
	key := append(keySrc, b...)

	// add new value
	for lenVal > 0 {
		if lenVal > chunkSize {
			k.Set(key, valueSrc[parts*chunkSize:(parts+1)*chunkSize])
		} else {
			k.Set(key, valueSrc[parts*chunkSize:])
		}

		parts++
		binary.LittleEndian.PutUint64(b, uint64(parts))
		key = append(keySrc, b...)

		lenVal -= chunkSize
	}

	i := parts

	// delete old values
	for k.Has(key) {
		k.Delete(key)
		i++
		binary.LittleEndian.PutUint64(b, uint64(i))
		key = append(keySrc, b...)
	}
	return
}

// Delete all parts of stored value
func (k *KVStore) DeleteW(keySrc []byte) (parts int, err error) {
	defer func() {
		if r := recover(); r != nil {
			nErr, nOk := r.(error)
			if nOk {
				err = nErr
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	parts = 0
	b := make([]byte, 8)

	binary.LittleEndian.PutUint64(b, uint64(parts))
	key := append(keySrc, b...)

	// delete old values
	for k.Has(key) {
		k.Delete(key)

		parts++
		binary.LittleEndian.PutUint64(b, uint64(parts))

		key = append(keySrc, b...)
	}

	return
}
