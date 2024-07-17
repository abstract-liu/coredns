package constant

type TunnelMode int

const (
	Global TunnelMode = iota
	Rule
	Direct
)
