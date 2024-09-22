package logger

import (
	"IM/pkg/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	luberjack "gopkg.in/natefinch/lumberjack.v2"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

func Init(cfg *config.LogConfig) error {
	encoder := GetEncoder()
	writeSyncer := GetWriteSyncer(cfg)
	var l = new(zapcore.Level)
	err := l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return err
	}
	core := zapcore.NewCore(encoder, writeSyncer, l)

	Logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(Logger) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	return nil
}

func GetWriteSyncer(cfg *config.LogConfig) zapcore.WriteSyncer {
	//WriteSyncer:指定将日志写道哪去，使用zapcore.AddSync()函数将打开的文件句柄传进去。
	//file, _ := os.OpenFile("./test.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0744)
	//利用io.MultiWriter支持文件和终端两个输出目标
	//ws := io.MultiWriter(file, os.Stdout)
	//writerSync := zapcore.AddSync(ws)

	//return writerSync

	// 使用lumberJackLogger配置，可以使得logger可以对日志切割分档
	lumberJackLogger := &luberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,    // M为单位
		MaxBackups: cfg.MaxBackups, // 最大备份数量
		MaxAge:     cfg.MaxAge,     // 最大过期时间，天为单位
		Compress:   cfg.Compress,   // 是否压缩
	}
	return zapcore.AddSync(lumberJackLogger)

}

func GetEncoder() zapcore.Encoder {
	//Encoder：编码器（如何写入日志），使用NewJSONEncoder（）并使用预先设置的NewProductionEncoderConfig()
	//encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	//使用zapcore.NewConsoleEncoder可以改为普通的Encoder
	//encoder := zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
	//return encoder

	//使用自定义配置
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder   //修改时间编码器
	config.EncodeLevel = zapcore.CapitalLevelEncoder //使用大写字母日志级别

	return zapcore.NewConsoleEncoder(config)
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		zap.L().Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
