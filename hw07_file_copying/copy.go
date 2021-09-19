package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	file, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileReader := bufio.NewReader(file)
	discarded, err := fileReader.Discard(int(offset))
	log.Printf("Discarded bytes %d", discarded)
	if err != nil {
		return ErrOffsetExceedsFileSize
	}
	// start new bar
	bar := pb.Full.Start64(limit)
	defer bar.Finish()

	barReader := bar.NewProxyReader(fileReader)

	dst, err := os.Create(toPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer dst.Close()

	var written int64
	if limit == 0 {
		written, err = io.Copy(dst, barReader)
	} else {
		written, err = io.CopyN(dst, barReader, limit)
	}

	if err != nil && err != io.EOF {
		log.Printf("Failed during writing %d bytes", written)
		return err
	}
	log.Printf("Written bytes %d", written)

	return nil
}
