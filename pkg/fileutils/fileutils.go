package fileutils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const BUFFERSIZE = 1024 * 1024

func getBufferSize(buffersize []int) int {
	if len(buffersize) > 0 {
		return buffersize[0]
	}
	return BUFFERSIZE
}

func GetFilesInDir(dir string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read source bases dir %s: %w", dir, err)
	}
	for _, e := range entries {
		files = append(files, filepath.Join(dir, e.Name()))
	}
	return files, nil
}

func Copy(src string, dst string, buffersize ...int) error {
	var _buffersize = getBufferSize(buffersize)

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	buf := make([]byte, _buffersize)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}

func WriteString(src string, content string) error {
	file, err := os.OpenFile(src, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", src, err)
	}
	defer file.Close()

	if _, err = file.WriteString(content); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}

func WriteBytes(src string, content []byte) error {
	file, err := os.OpenFile(src, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", src, err)
	}
	defer file.Close()

	if _, err = file.Write(content); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}

func ZipDir(zipname string, src string) error {
	type fileMeta struct {
		Path  string
		IsDir bool
	}

	var files []fileMeta

	if err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		files = append(files, fileMeta{Path: path, IsDir: info.IsDir()})
		return nil
	}); err != nil {
		return err
	}

	z, err := os.Create(zipname)
	if err != nil {
		return err
	}
	defer z.Close()

	zw := zip.NewWriter(z)
	defer zw.Close()

	for _, f := range files {
		path := f.Path
		if f.IsDir {
			path = fmt.Sprintf("%s%c", path, os.PathSeparator)
		}

		w, err := zw.Create(path)
		if err != nil {
			return err
		}

		if !f.IsDir {
			file, err := os.Open(f.Path)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err = io.Copy(w, file); err != nil {
				return err
			}
		}
	}
	return err
}
