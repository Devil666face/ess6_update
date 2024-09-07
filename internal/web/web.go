package web

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"

	"drw6/internal/config"
	"drw6/internal/drw6"
	"drw6/internal/web/handlers"
	"drw6/pkg/fileutils"

	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"

	"github.com/gofiber/fiber/v3"
)

const (
	media = "media"
)

type Web struct {
	fiber     *fiber.App
	drw6      *drw6.Drw6
	tlsconfig *tls.Config
	config    *config.Config
}

func New(
	_tlsconfig *tls.Config,
	_config *config.Config,
	_drw6 *drw6.Drw6,
) (*Web, error) {
	w := Web{
		fiber: fiber.New(
			fiber.Config{
				ErrorHandler: handlers.ErrorHandler,
			},
		),
		tlsconfig: _tlsconfig,
		config:    _config,
		drw6:      _drw6,
	}
	w.fiber.Use(logger.New())
	w.fiber.Use(recover.New())
	w.fiber.Use(w.wrap(handlers.AllowHost))
	w.fiber.Use(helmet.New(helmet.Config{
		XSSProtection:  "1",
		ReferrerPolicy: "same-origin",
	}))
	w.fiber.Use(helmet.New(helmet.Config{
		XSSProtection:  "1",
		ReferrerPolicy: "same-origin",
	}))
	path, err := fileutils.SetDir(media)
	if err != nil {
		return nil, fmt.Errorf("media directory not create of found: %w", err)
	}
	w.fiber.Get("/download", w.wrap(handlers.Download))
	w.fiber.Get("/status", w.wrap(handlers.Status))

	files := w.fiber.Group("/media")
	files.Use(w.wrap(handlers.FileList))
	files.Get("/*", static.New("", static.Config{
		FS:        os.DirFS(path),
		Browse:    false,
		ByteRange: true,
		Download:  true,
	}))
	return &w, nil
}

func (w *Web) wrap(handler func(*handlers.Handler) error) func(fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		return handler(
			handlers.New(
				c,
				w.drw6,
				w.config.AllowHost,
			),
		)
	}
}

func (w *Web) Listen() error {
	lis, err := net.Listen("tcp", w.config.ConnectHTTPS)
	if err != nil {
		return fmt.Errorf("error to listen: %w", err)
	}
	go w.RedirectServer()
	if err := w.fiber.Listener(tls.NewListener(lis, w.tlsconfig), fiber.ListenConfig{DisableStartupMessage: true}); err != nil {
		return fmt.Errorf("error static https share: %w", err)
	}
	return nil
}

func (w *Web) RedirectServer() {
	app := fiber.New()
	app.Use(func(c fiber.Ctx) error {
		return c.Redirect().To(w.config.HTTPSRedirect)
	})
	if err := app.Listen(w.config.ConnectHTTP, fiber.ListenConfig{DisableStartupMessage: true}); err != nil {
		log.Fatalf("failed to start redirect server: %s", err)
	}
}
