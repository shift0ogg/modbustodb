
package utils

import (
	"encoding/binary"
)


func UShortToBytes(i uint16) []byte {
    var buf = make([]byte, 8)
    binary.BigEndian.PutUint16(buf, uint16(i))
    return buf
}

func BytesToUShort(buf []byte) uint16 {
    return uint16(binary.BigEndian.Uint16(buf))
}

func ShortToBytes(i int16) []byte {
    var buf = make([]byte, 8)
    binary.BigEndian.PutUint16(buf, uint16(i))
    return buf
}

func BytesToShort(buf []byte) int16 {
    return int16(binary.BigEndian.Uint16(buf))
}




func Int64ToBytes(i int64) []byte {
    var buf = make([]byte, 8)
    binary.BigEndian.PutUint64(buf, uint64(i))
    return buf
}

func BytesToInt64(buf []byte) int64 {
    return int64(binary.BigEndian.Uint64(buf))
}
