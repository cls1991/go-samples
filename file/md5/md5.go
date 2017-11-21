package main

import (
	"crypto/md5"
	"io"
	"fmt"
	"os"
	"path/filepath"
	"encoding/json"
	"log"
	"strings"
	"io/ioutil"
)

type M map[string]string

const SplitChar = "/"

func main() {
	targetPath := "/home/cls1991/test"
	var fileList []string
	err := filepath.Walk(targetPath, func(path string, info os.FileInfo, err error) error {
		// ignore dir
		if !info.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	dat := make(map[string]M)
	for _, file := range fileList {
		if r, err := hashFile(file); err != nil {
			fmt.Printf("Err: %v\n", err)
		} else {
			prefix := targetPath[:]
			if fmt.Sprintf("%c", prefix[len(targetPath)-1]) != SplitChar {
				prefix += SplitChar
			}
			name := strings.TrimPrefix(file, prefix)
			dat[name] = M {"md5": fmt.Sprintf("%x", r)}
		}
	}
	if len(dat) > 0 {
		b, err := json.MarshalIndent(dat, "", "\t")
		if err != nil {
			fmt.Printf("Marshal err: %v\n", err)
		} else {
			err := ioutil.WriteFile("file/md5/project.manifest", b, 0644)
			if err != nil {
				log.Panic(err)
			}
		}
	}
}

func hashFile(filePath string) ([]byte, error) {
	var res []byte
	f, err := os.Open(filePath)
	if err != nil {
		return res, err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return res, err
	}
	return h.Sum(res), nil
}
