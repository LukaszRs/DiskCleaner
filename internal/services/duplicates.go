package services

import (
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type FileDef struct {
	Path       string
	Filename   string
	Size       int64
	CreatedAt  string
	ModifiedAt string
	Checksum   string
	Hash       string
}

func GoOverFiles(path string, pool *Pool) {
	fmt.Println("Searching duplicates in :", path)

	items, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory", err)
		return
	}

	for _, item := range items {
		if item.Name() == "." || item.Name() == ".." {
			continue
		}
		if item.IsDir() {
			fmt.Println("Found directory " + item.Name() + " spinning new thread")
			pool.AddTask(DirectoryToProcess{Path: filepath.Join(path, item.Name())})
		} else {
			fmt.Println("Found file " + item.Name())
			stat := item.Sys().(*syscall.Stat_t)
			checksum, err := calculateChecksum(filepath.Join(path, item.Name()))
			if err != nil {
				fmt.Println("Error calculating checksum for ", path, err)
			}
			pool.outputChan <- FileDef{
				Path:       path,
				Filename:   item.Name(),
				Size:       item.Size(),
				CreatedAt:  time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec)).String(),
				ModifiedAt: item.ModTime().String(),
				Checksum:   checksum,
				Hash:       "",
			}
		}
	}
}

func calculateChecksum(path string) (string, error) {
	readBytes := func(file *os.File, pos int64, where int) ([]byte, error) {
		_, err := file.Seek(pos, where)
		if err != nil {
			return nil, err
		}
		b := make([]byte, 2)
		_, err = file.Read(b)
		if err != nil {
			return nil, err
		}
		return b, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return "", nil
	}
	defer f.Close()
	f.SetDeadline(time.Now().Add(200 * time.Millisecond))
	f.SetReadDeadline(time.Now().Add(200 * time.Millisecond))

	result := make([]byte, 6)
	b, err := readBytes(f, 0, io.SeekStart)
	if err != nil {
		return "", nil
	}
	copy(result[0:2], b)
	b, err = readBytes(f, -2, io.SeekEnd)
	if err != nil {
		return hex.EncodeToString(result), nil
	}
	copy(result[4:6], b)
	fileSize, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return hex.EncodeToString(result), nil
	}
	b, err = readBytes(f, fileSize/2-2, io.SeekStart)
	if err != nil {
		return hex.EncodeToString(result), nil
	}
	copy(result[2:4], b)
	return hex.EncodeToString(result), nil
}
