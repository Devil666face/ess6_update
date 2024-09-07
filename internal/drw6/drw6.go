package drw6

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"drw6/pkg/fileutils"
	"drw6/pkg/shell"

	"github.com/reugn/go-quartz/job"
	"github.com/reugn/go-quartz/quartz"
	"golang.org/x/sync/errgroup"
)

const (
	cmd = `.\drwupdater.exe /DBG /QU /GO /ST /UA /DIR:bases`
)

const (
	loadmess      = "download bases"
	copymess      = "copy bases"
	timestempmess = "set timestemp"
	zipmess       = "zip bases"
	successmess   = "successful loaded"
)

var (
	bases      = filepath.Join("bases")
	repository = filepath.Join("repository")
	source     = filepath.Join(repository, "10-drwbases", "common")
	timestamp  = filepath.Join(bases, "timestamp")
	idBackup   = filepath.Join(repository, "10-drwbases", "id_backup")
)

type Drw6 struct {
	loadcmd string
	State   *LoadState
	sched   quartz.Scheduler
}

func New(trigger string) (*Drw6, error) {
	cron, err := quartz.NewCronTrigger(trigger)
	if err != nil {
		return nil, fmt.Errorf("error parse cron trigger: %w", err)
	}

	_drw6 := Drw6{
		loadcmd: cmd,
		State:   &LoadState{},
		sched:   quartz.NewStdScheduler(),
	}

	_drw6.sched.Start(context.Background())
	function := job.NewFunctionJob(
		func(_ context.Context) (bool, error) {
			if _drw6.State.IsLoad() {
				return false, fmt.Errorf("already in loading")
			}
			if err := _drw6.Update(); err != nil {
				return false, err
			}
			return true, nil
		},
	)
	_drw6.sched.ScheduleJob(
		quartz.NewJobDetail(
			function,
			quartz.NewJobKey("update"),
		),
		cron,
	)
	return &_drw6, nil
}

func (d *Drw6) UpdateMust() {
	if err := d.Update(); err != nil {
		log.Print(err)
	}
}

func (d *Drw6) Update() error {
	d.State.Start()
	defer d.State.Stop()

	if err := func() error {
		if err := d.download(); err != nil {
			return fmt.Errorf("failed download: %w", err)
		}
		if err := d.copybases(); err != nil {
			return fmt.Errorf("failed copy vdb files: %w", err)
		}
		if err := d.timestemp(); err != nil {
			return fmt.Errorf("failed write timestemp: %w", err)
		}
		if err := d.zip(); err != nil {
			return fmt.Errorf("failed zip: %w", err)
		}
		return nil
	}(); err != nil {
		d.State.SetError(err)
		d.State.SetMessage("")
		return err
	}
	d.State.SetMessage(successmess)
	return nil
}

func (d *Drw6) download() error {
	d.State.SetMessage(loadmess)
	if _, err := shell.Command(d.loadcmd); err != nil {
		return fmt.Errorf("failed to load bases: %w", err)
	}
	return nil
}

func (d *Drw6) copybases() error {
	d.State.SetMessage(copymess)

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

func (d *Drw6) timestemp() error {
	d.State.SetMessage(timestempmess)

	time, err := os.ReadFile(timestamp)
	if err != nil {
		return fmt.Errorf("failed to read timestamp from %s: %w", timestamp, err)
	}

	if err := fileutils.WriteString(idBackup, fmt.Sprintf("%s", time)); err != nil {
		return fmt.Errorf("failed to write timestemp %s: %w", idBackup, err)
	}
	return nil
}

func (d *Drw6) zip() error {
	d.State.SetMessage(zipmess)

	if err := fileutils.ZipDir("media/DRW_ESS6.zip", repository); err != nil {
		return fmt.Errorf("failed to create bases zip: %w", err)
	}
	return nil
}
