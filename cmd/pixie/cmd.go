package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/alexflint/go-arg"
	"github.com/p3lim/pixie/pkg/http"
	"github.com/p3lim/pixie/pkg/tftp"
)

var version = "dev"

type args struct {
	TFTP    string `arg:"-t" default:"0.0.0.0:69" placeholder:"ADDR" help:"tftp server address and port"`
	HTTP    string `arg:"-u" default:"0.0.0.0:80" placeholder:"ADDR" help:"http server address and port"`
	Scripts string `arg:"-d,required" placeholder:"DIR" help:"path to iPXE scripts"`
}

func (args) Version() string {
	return fmt.Sprintf("pixie %s", version)
}

func (args) Validate() error {
	// TODO: validate address
	// TODO: validate port as valid and not in use (>0<65535)
	// TODO: validate scripts dir and permissions
	return nil
}

func main() {
	var cfg args
	arg.MustParse(&cfg)

	if err := cfg.Validate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		tftpServer := tftp.NewServer(cfg.TFTP, cfg.HTTP)
		log.Printf("tftp server listening on %s", cfg.TFTP)
		log.Fatal(tftpServer.Serve())
		wg.Done()
	}()

	go func() {
		httpServer := http.NewServer(cfg.HTTP, cfg.Scripts)
		log.Printf("http server listening on %s", cfg.HTTP)
		log.Fatal(httpServer.Serve())
		wg.Done()
	}()

	wg.Wait()
}
