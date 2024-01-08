package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFileInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	if !fromFileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	fromFileSize := fromFileInfo.Size()
	if fromFileSize < offset {
		return ErrOffsetExceedsFileSize
	}

	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	_, err = fromFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	if limit <= 0 || offset+limit > fromFileSize {
		limit = fromFileSize - offset
	}

	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(fromFile)
	_, err = io.CopyN(toFile, barReader, limit)
	if err != nil {
		return err
	}
	bar.Finish()
	return nil
}
