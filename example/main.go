package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
	"os"
	"path/filepath"
	"sblogger"
	"time"
)

func logFile(fileName string) *os.File {
	// ext := filepath.Ext(fileName)
	date := time.Now().Format("2006-01-02")
	// name := strings.TrimSuffix(fileName, ext)
	// fileName = fmt.Sprintf("%s-%s%s", name, date, ext)
	path := filepath.Join("storage/logs", date)
	os.MkdirAll(path, os.ModePerm)
	path = filepath.Join(path, fileName)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(fmt.Sprintf("open file error: %+v", err))
	}
	return file
}

func main() {

	log.DefaultLogger = log.Logger{
		TimeField:  "timestamp",
		TimeFormat: "2006-01-02 15:04:05",
	}

	log.DefaultLogger.Writer = &log.MultiWriter{
		InfoWriter:    &log.IOWriter{logFile("INFO.log")},
		WarnWriter:    &log.IOWriter{logFile("WARN.log")},
		ErrorWriter:    &log.IOWriter{logFile("ERROR.log")},
		ConsoleWriter: &log.IOWriter{os.Stderr},
		ConsoleLevel:  log.InfoLevel,
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
	app.Listen(":8080")

}