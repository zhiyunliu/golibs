package rsa

type Option func(*options)

type options struct {
	PkcsType string
}
