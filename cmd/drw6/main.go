package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"drw6/pkg/fileutils"
	"drw6/pkg/shell"

	"golang.org/x/sync/errgroup"
)

var (
	bases      = filepath.Join("bases")
	repository = filepath.Join("repository")
	source     = filepath.Join(repository, "10-drwbases", "common")
	timestamp  = filepath.Join(bases, "timestamp")
	idBackup   = filepath.Join(repository, "10-drwbases", "id_backup")
	cmd        = `.\drwupdater.exe /DBG /QU /GO /ST /UA /DIR:bases`
)

type Loader struct {
	LoadCmd string
}

func New() *Loader {
	return &Loader{
		LoadCmd: cmd,
	}
}

func (l *Loader) Load() error {
	out, err := shell.Command(l.LoadCmd)
	if err != nil {
		return fmt.Errorf("failed to load bases: %w", err)
	}
	fmt.Println(out)
	return nil
}

func (l *Loader) CopyBasesFiles() error {
	var files []string

	entries, err := os.ReadDir(bases)
	if err != nil {
		return fmt.Errorf("failed to read source bases dir %s: %w", bases, err)
	}
	for _, e := range entries {
		files = append(files, filepath.Join(bases, e.Name()))
	}

	var g errgroup.Group

	for _, f := range files {
		_f := f
		g.Go(func() error {
			if err := fileutils.Copy(_f, filepath.Join(source, filepath.Base(_f))); err != nil {
				return fmt.Errorf("failed to copy %s: %w", _f, err)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func (l *Loader) WriteTimestemp() error {
	time, err := os.ReadFile(timestamp)
	if err != nil {
		return fmt.Errorf("failed to read timestamp from %s: %w", timestamp, err)
	}

	if err := fileutils.WriteFile(idBackup, fmt.Sprintf("%s", time)); err != nil {
		return fmt.Errorf("failed to write timestemp %s: %w", idBackup, err)
	}
	return nil
}

func (l *Loader) CreateZip() error {
	if err := fileutils.ZipDir("DRW_ESS6.zip", repository); err != nil {
		return fmt.Errorf("failed to create bases zip: %w", err)
	}
	return nil
}

func main() {
	_loader := New()

	if err := _loader.Load(); err != nil {
		log.Fatal(err)
	}
	if err := _loader.CopyBasesFiles(); err != nil {
		log.Fatal(err)
	}
	if err := _loader.WriteTimestemp(); err != nil {
		log.Fatal(err)
	}
	if err := _loader.CreateZip(); err != nil {
		log.Fatal(err)
	}
}
