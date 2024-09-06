package drw6

import (
	"fmt"
	"os"
	"path/filepath"

	"drw6/pkg/fileutils"
	"drw6/pkg/shell"

	"golang.org/x/sync/errgroup"
)

const (
	cmd = `.\drwupdater.exe /DBG /QU /GO /ST /UA /DIR:bases`
)

var (
	bases      = filepath.Join("bases")
	repository = filepath.Join("repository")
	source     = filepath.Join(repository, "10-drwbases", "common")
	timestamp  = filepath.Join(bases, "timestamp")
	idBackup   = filepath.Join(repository, "10-drwbases", "id_backup")
)

type loader struct {
	loadcmd string
}

func New() *loader {
	return &loader{
		loadcmd: cmd,
	}
}

func (l *loader) Create() error {
	if err := l.download(); err != nil {
		return fmt.Errorf("failed download: %w", err)
	}
	if err := l.copybases(); err != nil {
		return fmt.Errorf("failed copy vdb files: %w", err)
	}
	if err := l.timestemp(); err != nil {
		return fmt.Errorf("failed write timestemp: %w", err)
	}
	if err := l.zip(); err != nil {
		return fmt.Errorf("failed zip: %w", err)
	}
	return nil
}

func (l *loader) download() error {
	if _, err := shell.Command(l.loadcmd); err != nil {
		return fmt.Errorf("failed to load bases: %w", err)
	}
	return nil
}

func (l *loader) copybases() error {
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

func (l *loader) timestemp() error {
	time, err := os.ReadFile(timestamp)
	if err != nil {
		return fmt.Errorf("failed to read timestamp from %s: %w", timestamp, err)
	}

	if err := fileutils.WriteString(idBackup, fmt.Sprintf("%s", time)); err != nil {
		return fmt.Errorf("failed to write timestemp %s: %w", idBackup, err)
	}
	return nil
}

func (l *loader) zip() error {
	if err := fileutils.ZipDir("DRW_ESS6.zip", repository); err != nil {
		return fmt.Errorf("failed to create bases zip: %w", err)
	}
	return nil
}
