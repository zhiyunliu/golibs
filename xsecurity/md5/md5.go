package md5

import (
	"crypto/md5"
	"fmt"

	"github.com/zhiyunliu/golibs/bytesconv"
)

func Str(val string) (r string) {
	return Bytes(bytesconv.StringToBytes(val))
}

func Bytes(buffer []byte) (r string) {
	md5Ctx := md5.New()
	md5Ctx.Write(buffer)
	return fmt.Sprintf("%x", md5Ctx.Sum(nil))
}
