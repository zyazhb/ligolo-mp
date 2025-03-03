package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/ttpreport/ligolo-mp/cmd/client/tui"
	"github.com/ttpreport/ligolo-mp/internal/certificate"
	"github.com/ttpreport/ligolo-mp/internal/config"
	"github.com/ttpreport/ligolo-mp/internal/crl"
	"github.com/ttpreport/ligolo-mp/internal/operator"
	"github.com/ttpreport/ligolo-mp/internal/storage"
	"github.com/ttpreport/ligolo-mp/pkg/logger"
)

func main() {
	var verbose = flag.Bool("v", false, "enable verbose mode")

	flag.Parse()

	loggingOpts := &slog.HandlerOptions{}
	if *verbose {
		lvl := new(slog.LevelVar)
		lvl.Set(slog.LevelDebug)
		loggingOpts = &slog.HandlerOptions{
			Level: lvl,
		}
	} else {
		lvl := new(slog.LevelVar)
		lvl.Set(slog.LevelInfo)
		loggingOpts = &slog.HandlerOptions{
			Level: lvl,
		}
	}
	logHandler := slog.New(slog.NewTextHandler(os.Stdout, loggingOpts))
	slog.SetDefault(logHandler)

	cfg := &config.Config{
		Environment: "client",
	}

	storage, err := storage.New(cfg.GetRootAppDir())
	if err != nil {
		panic(fmt.Sprintf("could not connect to storage: %v", err))
	}

	operRepo, err := operator.NewOperatorRepository(storage)
	if err != nil {
		panic(err)
	}

	certRepo, err := certificate.NewCertificateRepository(storage)
	if err != nil {
		panic(err)
	}

	crlRepo, err := crl.NewCRLRepository(storage)
	if err != nil {
		panic(err)
	}

	crlService := crl.NewCRLService(crlRepo)
	certService := certificate.NewCertificateService(certRepo, crlService)
	operService := operator.NewOperatorService(cfg, operRepo, certService)

	app := tui.NewApp(operService)
	logHandler = slog.New(logger.NewLogHandler(app.Logs, loggingOpts))
	slog.SetDefault(logHandler)
	app.Run()
}
