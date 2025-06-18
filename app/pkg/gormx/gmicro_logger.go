package gormx

import (
	"context"
	"fmt"
	"time"

	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"github.com/WeiXinao/daily_your_go/pkg/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type Config struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
	LogLevel                  logger.LogLevel
}

var _ logger.Interface = (*GMicroLogger)(nil)

type GMicroLogger struct {
	Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// New initialize logger
func New(config Config) logger.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	return &GMicroLogger{
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

// Error implements logger.Interface.
func (g *GMicroLogger) Error(ctx context.Context, s string, i ...interface{}) {
	if g.LogLevel >= logger.Error {
		log.ErrorfC(ctx, s, i...)
	}
}

// Info implements logger.Interface.
func (g *GMicroLogger) Info(ctx context.Context, s string, i ...interface{}) {
	if g.LogLevel >= logger.Info {
		log.InfofC(ctx, s, i...)
	}
}

// LogMode implements logger.Interface.
func (g *GMicroLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *g
	newLogger.LogLevel = level
	return &newLogger
}

// Trace implements logger.Interface.
func (g *GMicroLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if g.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && g.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !g.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			log.Infof(g.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Infof(g.traceErrStr,  utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > g.SlowThreshold && g.SlowThreshold != 0 && g.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", g.SlowThreshold)
		if rows == -1 {
			log.Infof(g.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Infof(g.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case g.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			log.Infof(g.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Infof(g.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

// Warn implements logger.Interface.
func (g *GMicroLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	if g.LogLevel >= logger.Warn {
		log.WarnfC(ctx, s, i...)
	}
}
