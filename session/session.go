package session

import (
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/zhiyunliu/golibs/bytesconv"
)

//Create create logger session
func Create() string {
	var buf [32]byte
	uval := uuid.New()
	hex.Encode(buf[:], uval[:])
	return bytesconv.BytesToString(buf[:16])
}
