package sblogger

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
	"github.com/sujit-baniya/xid"
	"time"
)

type Config struct {
	Logger    *log.Logger
	LogWriter log.Writer
}

//Middleware requestid + logger + recover for request traceability
func New(config Config) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		if config.Logger == nil {
			config.Logger = &log.Logger{
				TimeField:  "timestamp",
				TimeFormat: "2006-01-02 15:04:05",
			}
		}
		if config.LogWriter != nil {
			config.Logger.Writer = config.LogWriter
		}
		rid := c.Get(fiber.HeaderXRequestID)
		if rid == "" {
			rid = xid.New().String()
			c.Set(fiber.HeaderXRequestID, rid)
		}
		ctx := log.NewContext(nil).
			Str("request_id", rid).
			Str("remote_ip", c.IP()).
			Str("method", c.Method()).
			Str("host", c.Hostname()).
			Str("path", c.Path()).
			Str("protocol", c.Protocol()).
			Int("status", c.Response().StatusCode()).
			Str("latency", fmt.Sprintf("%s", time.Since(start))).
			Str("ua", c.Get(fiber.HeaderUserAgent)).
			Value()

		switch {
		case c.Response().StatusCode() >= 500:
			config.Logger.Error().Context(ctx).Msg("server error")
		case c.Response().StatusCode() >= 400:
			config.Logger.Error().Context(ctx).Msg("client error")
		case c.Response().StatusCode() >= 300:
			config.Logger.Warn().Context(ctx).Msg("redirect")
		case c.Response().StatusCode() >= 200:
			config.Logger.Info().Context(ctx).Msg("success")
		case c.Response().StatusCode() >= 100:
			config.Logger.Info().Context(ctx).Msg("informative")
		default:
			config.Logger.Warn().Context(ctx).Msg("unknown status")
		}
		return c.Next()
	}
}
