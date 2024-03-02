package xreflect

type MapScanner interface {
	MapScan(val any) error
}

type StructScanner interface {
	StructScan(vals ...any) error
}
