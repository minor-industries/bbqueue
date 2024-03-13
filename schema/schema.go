package schema

type ThermocoupleData struct {
	Temperature float32 // celsius
	Description [16]byte
}

type PowerData struct {
	Voltage     float32
	Current     float32
	Power       float32
	Addr        uint8
	Description [16]byte
}
