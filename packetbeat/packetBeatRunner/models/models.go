package models

// PacketBeatStatus is a response for start and top packetbeat
// swagger:model resPBStatus
type PacketBeatStatus struct {
	//Msg status
	//
	// example: Packetbeat started
	Msg string `json:"msg"`
	//PID process ID of packetbeat
	//
	// example: 1234
	PID string `json:"pid"`
}
