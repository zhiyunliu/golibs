package aes

type Option func(*options)

type options struct {
	IV        []byte
	BlockSize int
}
