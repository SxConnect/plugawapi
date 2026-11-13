package call

import (
	"encoding/binary"
	"os"
)

func writeWAV(path string, pcm []byte, rate int) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	dataSize := uint32(len(pcm))
	byteRate := uint32(rate * 2)
	chunkSize := 36 + dataSize
	f.Write([]byte("RIFF"))
	binary.Write(f, binary.LittleEndian, chunkSize)
	f.Write([]byte("WAVE"))
	f.Write([]byte("fmt "))
	binary.Write(f, binary.LittleEndian, uint32(16))
	binary.Write(f, binary.LittleEndian, uint16(1))
	binary.Write(f, binary.LittleEndian, uint16(1))
	binary.Write(f, binary.LittleEndian, uint32(rate))
	binary.Write(f, binary.LittleEndian, byteRate)
	binary.Write(f, binary.LittleEndian, uint16(2))
	binary.Write(f, binary.LittleEndian, uint16(16))
	f.Write([]byte("data"))
	binary.Write(f, binary.LittleEndian, dataSize)
	f.Write(pcm)
}
