package handlers

import (
	"drw6/pkg/fileutils"
	"net/url"
	"os"
	"path/filepath"
)

func URLToFilepath(u string) (string, error) {
	decoded, err := url.PathUnescape(u)
	if err != nil {
		return "", err
	}
	base, err := os.Getwd()
	if err != nil {
		return "", err
	}
	abs, err := filepath.Abs(filepath.Join(base, decoded))
	if err != nil {
		return "", err
	}
	return abs, nil
}

func FileList(h *Handler) error {
	path, err := URLToFilepath(h.ctx.Path())
	if err != nil {
		return h.ctx.Next()
	}
	if stat, err := os.Stat(path); err != nil || !stat.IsDir() {
		return h.ctx.Next()
	}
	files, err := fileutils.DirContent(path)
	if err != nil {
		return h.ctx.Next()
	}
	return h.ctx.JSON(
		files,
	)
}
