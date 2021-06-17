package recording

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"
	"sync"
)

type Recording struct {
	Name       string
	File       *os.File
	Reader     *bufio.Reader
	InMemory   []*Frame
	NbInMemory uint32
	ReadIdx    uint32

	sync.Mutex
}

type Frame struct {
	Timestamp int64
	Length    int16
	Data      []byte
}

func Open(name string) (*Recording, error) {
	r := &Recording{ReadIdx: 0}

	r.Name = name

	file, err := os.Open(name)

	if err != nil {
		return nil, err
	}

	r.File = file

	return r, nil
}

func (recording *Recording) Close() {
	if recording.File != nil {
		recording.File.Close()
	}
}

func (recording *Recording) Reset() {
	recording.Lock()
	recording.File.Seek(0, io.SeekStart)
	recording.Reader.Reset(recording.File)
	recording.Unlock()
}

func (recording *Recording) ReadFrame() (*Frame, error) {
	recording.Lock()
	defer recording.Unlock()

	if recording.Reader == nil {
		reader := bufio.NewReader(recording.File)
		recording.Reader = reader
	}

	timestamp := make([]byte, 8)
	length := make([]byte, 2)

	_, err := io.ReadFull(recording.Reader, timestamp)

	if err != nil {
		return nil, err
	}

	_, err = io.ReadFull(recording.Reader, length)
	if err != nil {
		return nil, err
	}

	frame := &Frame{
		Timestamp: int64(binary.BigEndian.Uint64(timestamp)),
		Length:    int16(binary.BigEndian.Uint16(length)),
	}
	frame.Data = make([]byte, frame.Length)

	_, err = io.ReadFull(recording.Reader, frame.Data)

	if err != nil {
		return nil, err
	}

	return frame, nil
}

func (rec *Recording) LoadInMemory(nbFrame int) error {
	rec.NbInMemory = uint32(nbFrame)
	rec.InMemory = make([]*Frame, nbFrame)

	for i := 0; i < nbFrame; i++ {
		f, err := rec.ReadFrame()
		if err != nil {
			return err
		}
		rec.InMemory[i] = f
	}

	return nil
}

func (rec *Recording) GetInMemoryFrame() *Frame {
	rec.Lock()
	defer rec.Unlock()

	if rec.ReadIdx >= rec.NbInMemory {
		rec.ReadIdx = 0
	}

	frame := rec.InMemory[rec.ReadIdx]
	rec.ReadIdx += 1

	return frame
}
