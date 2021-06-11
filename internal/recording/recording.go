package recording

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"
	"sync"
)

type Recording struct {
	Name   string
	File   *os.File
	Reader *bufio.Reader

	sync.Mutex
}

type Frame struct {
	Timestamp int64
	Length    int16
	Data      []byte
}

func Open(name string) (*Recording, error) {
	r := &Recording{}

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
