package hearthlands

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

const (
	offsetPopulation int64 = 0xa
	offsetGold       int64 = 0xe
)

type SaveFile struct {
	Header []byte
	Entry  []byte
}

func OpenSaveFile(filename string) (*SaveFile, error) {
	s := SaveFile{}

	if err := s.Open(filename); err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *SaveFile) Open(filename string) error {

	f, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}

	// search for beginning of embedded zip file
	i := bytes.Index(f, []byte{0x50, 0x4b, 0x03, 0x04})

	if i == -1 {
		return errors.New("invalid save game")
	}

	s.Header = f[0 : i-4]

	z, err := zip.NewReader(bytes.NewReader(f[i:]), int64(len(f)-i))

	if err != nil {
		return err
	}

	if len(z.File) != 1 || z.File[0].Name != "entry.txt" {
		return errors.New("invalid save game container")
	}

	entry, err := z.File[0].Open()

	if err != nil {
		return err
	}

	s.Entry, err = ioutil.ReadAll(entry)

	if err != nil {
		return err
	}

	return nil
}

func (s *SaveFile) readValue(offset int64) int {
	r := bytes.NewReader(s.Entry)
	r.Seek(offset, io.SeekStart)
	g := make([]byte, 4)
	r.Read(g)

	return int(binary.BigEndian.Uint32(g))
}

func (s *SaveFile) Gold() int {
	return s.readValue(offsetGold)
}

func (s *SaveFile) SetGold(v int) {
	binary.BigEndian.PutUint32(s.Entry[offsetGold:offsetGold+4], uint32(v))
}

func (s *SaveFile) Population() int {
	return s.readValue(offsetPopulation)
}

func (s *SaveFile) SetPopulation(v int) {
	binary.BigEndian.PutUint32(s.Entry[offsetPopulation:offsetPopulation+4], uint32(v))
}

func (s *SaveFile) Save(filename string) error {
	z := new(bytes.Buffer)
	zw := zip.NewWriter(z)
	zf, err := zw.Create("entry.txt")

	if err != nil {
		return err
	}

	_, err = io.Copy(zf, bytes.NewReader(s.Entry))

	if err != nil {
		return err
	}

	if err = zw.Close(); err != nil {
		return err
	}

	// output format:
	// - header
	// - zip file length (4 byte int)
	// - zip file

	f, err := os.Create(filename)

	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.Write(s.Header); err != nil {
		return err
	}

	if err = binary.Write(f, binary.BigEndian, uint32(z.Len())); err != nil {
		return err
	}

	if _, err = f.Write(z.Bytes()); err != nil {
		return err
	}

	return nil
}
