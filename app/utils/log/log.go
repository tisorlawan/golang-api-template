package log

import (
	"fmt"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type cLogger struct {
	*zap.Logger
}

// Loggers
var middlewareLogger, _ = zap.NewProduction()
var logger, _ = zap.NewDevelopment()

// Logging methods
var Debug = logger.Debug
var Info = logger.Info
var Warn = logger.Warn
var Error = logger.Error

func Sync() {
	logger.Sync()
	middlewareLogger.Sync()
}

func GetMiddleware() echo.MiddlewareFunc {

	atom := zap.NewAtomicLevel()

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	middlewareLogger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))

	atom.SetLevel(zap.DebugLevel)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			fields := []zapcore.Field{
				zap.String("remote_ip", c.RealIP()),
				zap.String("latency", time.Since(start).String()),
				zap.String("host", req.Host),
				zap.String("request", fmt.Sprintf("%s %s", req.Method, req.RequestURI)),
				zap.Int("status", res.Status),
				zap.Int64("size", res.Size),
				zap.String("user_agent", req.UserAgent()),
			}

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
				fields = append(fields, zap.String("request_id", id))
			}

			n := res.Status
			switch {
			case n >= 500:
				middlewareLogger.Error("Server error", fields...)
			case n >= 400:
				middlewareLogger.Warn("Client error", fields...)
			case n >= 300:
				middlewareLogger.Info("Redirection", fields...)
			default:
				middlewareLogger.Info("Success", fields...)
			}

			return nil
		}
	}
}
