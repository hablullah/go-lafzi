package lafzi

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
)

func peek(runes []rune, idx int) (rune, bool) {
	if idx < 0 || idx >= len(runes) {
		return 0, false
	}

	return runes[idx], true
}

func int64ToBytes(i int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

func bytesToInt64(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(b))
}

func arrayIntToBytes(arr []int) []byte {
	buffer := bytes.NewBuffer(nil)
	gob.NewEncoder(buffer).Encode(&arr)
	return buffer.Bytes()
}

func bytesToArrayInt(b []byte) []int {
	var arr []int
	reader := bytes.NewReader(b)
	gob.NewDecoder(reader).Decode(&arr)
	return arr
}
