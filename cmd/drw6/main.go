package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "embed"

	"drw6/pkg/fileutils"
	"drw6/pkg/shell"

	"golang.org/x/sync/errgroup"
)

const (
	DrwUpdater = "drwupdater.exe"
	UpdateDrl  = "update.drl"
	Drweb32    = "drweb32.ini"
)

//go:embed drwupdater.exe
var DrwUpdaterBin []byte

//go:embed update.drl
var UpdateDrlBin []byte

//go:embed drweb32.ini
var Drweb32Bin []byte

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

func New() (*Loader, error) {
	if err := fileutils.WriteBytes(DrwUpdater, DrwUpdaterBin); err != nil {
		return nil, err
	}
	if err := fileutils.WriteBytes(UpdateDrl, UpdateDrlBin); err != nil {
		return nil, err
	}
	if err := fileutils.WriteBytes(Drweb32, Drweb32Bin); err != nil {
		return nil, err
	}
	return &Loader{
		LoadCmd: cmd,
	}, nil
}

func (l *Loader) Load() error {
	if _, err := shell.Command(l.LoadCmd); err != nil {
		return fmt.Errorf("failed to load bases: %w", err)
	}
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

	if err := fileutils.WriteString(idBackup, fmt.Sprintf("%s", time)); err != nil {
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
	loader, err := New()
	if err != nil {
		log.Fatal(err)
	}
	if err := loader.Load(); err != nil {
		log.Fatal(err)
	}
	if err := loader.CopyBasesFiles(); err != nil {
		log.Fatal(err)
	}
	if err := loader.WriteTimestemp(); err != nil {
		log.Fatal(err)
	}
	if err := loader.CreateZip(); err != nil {
		log.Fatal(err)
	}
}
