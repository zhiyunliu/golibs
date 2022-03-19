package xlog

import (
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/zhiyunliu/golibs/bytesconv"
)

//CreateSession create logger session
func CreateSession() string {
	var buf [32]byte
	uval := uuid.New()
	hex.Encode(buf[:], uval[:])
	return bytesconv.BytesToString(buf[:])
}
