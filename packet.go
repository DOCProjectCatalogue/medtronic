package medtronic

import (
	"fmt"

	"github.com/ecc1/radio"
)

func (pump *Pump) DecodePacket(packet radio.Packet) ([]byte, error) {
	data, err := Decode6b4b(packet.Data)
	if err != nil {
		pump.DecodingErrors++
		return nil, err
	}
	crc := Crc8(data[:len(data)-1])
	if data[len(data)-1] != crc {
		pump.CrcErrors++
		return data, fmt.Errorf("CRC should be %X, not %X", crc, data[len(data)-1])
	}
	return data, nil
}

func EncodePacket(data []byte) radio.Packet {
	// Don't use append() to add the CRC, because append
	// may write into the array underlying the caller's slice.
	buf := make([]byte, len(data)+1)
	copy(buf, data)
	buf[len(data)] = Crc8(data)
	return radio.Packet{Data: Encode4b6b(buf)}
}

func (pump *Pump) PrintStats() {
	stats := pump.Radio.Statistics()
	good := stats.Packets.Received - pump.DecodingErrors - pump.CrcErrors
	fmt.Printf("\nTX: %6d    RX: %6d    decode errs: %6d    CRC errs: %6d\n", stats.Packets.Sent, good, pump.DecodingErrors, pump.CrcErrors)
	fmt.Printf("State: %s\n", pump.Radio.State())
}
