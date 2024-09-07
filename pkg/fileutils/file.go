package fileutils

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

type File struct {
	path    string
	size    int64
	modTime time.Time
	Href    string `json:"href"`
	Name    string `json:"name"`
	IsDir   bool   `json:"dir"`
	Size    string `json:"size"`
	ModTime string `json:"modify"`
}

const (
	_ = 1 << (10 * iota)
	KB
	MB
	GB
)

func formatSize(size int64) string {
	switch {
	case size < KB:
		return fmt.Sprintf("%d B", size)
	case size < MB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	case size < GB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	default:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	}
}

func DirContent(path string) ([]File, error) {
	var c []File
	pathes, err := filepath.Glob(path + "/*")
	if err != nil {
		return nil, err
	}
	for _, p := range pathes {
		file, err := New(p)
		if err != nil {
			return nil, err
		}
		c = append(c, *file)
	}
	return c, nil
}

func (f *File) stat() error {
	stat, err := os.Stat(f.path)
	if err != nil {
		return fmt.Errorf("get file info: %w for file %s", err, f.path)
	}
	f.IsDir = stat.IsDir()
	f.modTime = stat.ModTime()
	f.size = stat.Size()
	return nil
}

func pathToUrl(path string) string {
	path = strings.ReplaceAll(path, string(os.PathSeparator), "/")
	segments := strings.Split(path, "/")
	for i, segment := range segments {
		segments[i] = url.PathEscape(segment)
	}
	return strings.Join(segments, "/")
}

func (f *File) href() error {
	base, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get current folder: %w", err)
	}
	path := strings.TrimPrefix(f.path, base)
	f.Href = pathToUrl(path)
	return nil
}

func New(_path string) (*File, error) {
	f := File{
		path: _path,
		Name: filepath.Base(_path),
	}
	if err := f.stat(); err != nil {
		return nil, err
	}
	if err := f.href(); err != nil {
		return nil, err
	}
	f.Size = formatSize(f.size)
	f.ModTime = f.modTime.Format(timeFormat)
	return &f, nil
}
