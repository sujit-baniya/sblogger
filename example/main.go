package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/log"
	"github.com/sujit-baniya/sblogger"
	"os"
)

func main() {

	log.DefaultLogger = log.Logger{
		TimeField:  "timestamp",
		TimeFormat: "2006-01-02 15:04:05",
		Writer: &log.MultiWriter{
			InfoWriter:    &log.FileWriter{Filename: "storage/logs/INFO.log", EnsureFolder: true, TimeFormat: "2006-01-02"},
			WarnWriter:    &log.FileWriter{Filename: "storage/logs/WARN.log", EnsureFolder: true, TimeFormat: "2006-01-02"},
			ErrorWriter:   &log.FileWriter{Filename: "storage/logs/ERROR.log", EnsureFolder: true, TimeFormat: "2006-01-02"},
			ConsoleWriter: &log.IOWriter{os.Stderr},
			ConsoleLevel:  log.InfoLevel,
		},
	}
	app := fiber.New()
	app.Use(sblogger.New(sblogger.Config{
		Logger:    &log.DefaultLogger,
		LogWriter: log.DefaultLogger.Writer,
	}),
	)
	app.Get("test", func(c *fiber.Ctx) error {
		return c.JSON("Hello")
	})
	app.Listen(":8081")

}
