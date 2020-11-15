package sblogger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
	"github.com/sujit-baniya/xid"
	"time"
)

type Config struct {
	Logger *log.Logger
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

		fields := map[string]interface{} {
			"request_id":       rid,
			"remote_ip": c.IP(),
			"method":   c.Method(),
			"host":     c.Hostname(),
			"path":     c.Path(),
			"protocol": c.Protocol(),
			"status": c.Response().StatusCode(),
			"latency" : time.Since(start).Seconds(),
			"ua": c.Get(fiber.HeaderUserAgent),
		}

		switch {
		case c.Response().StatusCode() >= 500:
			config.Logger.Error().Fields(fields).Msg("server error")
		case c.Response().StatusCode() >= 400:
			config.Logger.Error().Fields(fields).Msg("client error")
		case c.Response().StatusCode() >= 300:
			config.Logger.Warn().Fields(fields).Msg("redirect")
		case c.Response().StatusCode() >= 200:
			config.Logger.Info().Fields(fields).Msg("success")
		case c.Response().StatusCode() >= 100:
			config.Logger.Info().Fields(fields).Msg("informative")
		default:
			config.Logger.Warn().Fields(fields).Msg("unknown status")
		}
		return c.Next()
	}
}