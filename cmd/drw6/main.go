package main

import (
	"log"

	_ "embed"

	"drw6/internal/drw6"

	"drw6/pkg/fileutils"
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

func main() {
	if err := fileutils.WriteBytes(DrwUpdater, DrwUpdaterBin); err != nil {
		log.Fatal(err)
	}
	if err := fileutils.WriteBytes(UpdateDrl, UpdateDrlBin); err != nil {
		log.Fatal(err)
	}
	if err := fileutils.WriteBytes(Drweb32, Drweb32Bin); err != nil {
		log.Fatal(err)
	}
	_drw6 := drw6.New()
	if err := _drw6.Create(); err != nil {
		log.Fatal(err)
	}
}
