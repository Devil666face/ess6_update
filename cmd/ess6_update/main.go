package main

import (
	"log"

	_ "embed"

	"drw6/internal/config"
	"drw6/internal/drw6"
	"drw6/internal/web"
	"drw6/pkg/fileutils"
	"drw6/pkg/netutils"
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

//go:embed server.key
var TLSServerKey string

//go:embed server.crt
var TLSServerCert string

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

	_tlsconfig, err := netutils.LoadTlsCreds(
		TLSServerCert,
		TLSServerKey,
	)
	if err != nil {
		log.Fatal(err)
	}

	_config, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("parsed config: %v", _config)

	_drw6, err := drw6.New(
		_config.CronTrigger,
	)
	if err != nil {
		log.Fatal(err)
	}

	_web, err := web.New(
		_tlsconfig,
		_config,
		_drw6,
	)
	if err := _web.Listen(); err != nil {
		log.Fatal(err)
	}
}
