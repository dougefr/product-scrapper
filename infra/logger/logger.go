package logger

import (
	"context"
	"time"

	logger2 "github.com/dougefr/product-scrapper/domain/contract/logger"
	context2 "github.com/dougefr/product-scrapper/infra/context"
	"go.uber.org/zap"
)

type logger struct {
	log            *zap.Logger
	tickerSync     *time.Ticker
	quitTickerSync chan struct{}
}

// NewLogger cria um novo logger
func NewLogger() logger2.Logger {
	log, _ := zap.NewProduction()
	log = log.WithOptions(zap.AddCallerSkip(1))

	// ticker para executar o sync
	tickerSync := time.NewTicker(1 * time.Second)
	quitTickerSync := make(chan struct{})
	go func() {
		for {
			select {
			case <-tickerSync.C:
				err := log.Sync()
				if err != nil {
					return
				}
			case <-quitTickerSync:
				tickerSync.Stop()
			}
		}
	}()

	return logger{
		log:            log,
		tickerSync:     tickerSync,
		quitTickerSync: quitTickerSync,
	}
}

func convertBodyToFields(ctx context.Context, valueMaps ...logger2.Body) []interface{} {
	fields := make([]interface{}, 0)
	for _, values := range valueMaps {
		for name, value := range values {
			fields = append(fields, name)
			fields = append(fields, value)
		}
	}

	requestID, ok := ctx.Value(context2.RequestID).(string)
	if ok {
		fields = append(fields, "request-id")
		fields = append(fields, requestID)
	}

	path, ok := ctx.Value(context2.Path).(string)
	if ok {
		fields = append(fields, "path")
		fields = append(fields, path)
	}

	method, ok := ctx.Value(context2.Method).(string)
	if ok {
		fields = append(fields, "method")
		fields = append(fields, method)
	}

	return fields
}

// Debug loga uma informação de depuração
func (l logger) Debug(ctx context.Context, message string, valueMaps ...logger2.Body) {
	l.log.Sugar().Debugw(message, convertBodyToFields(ctx, valueMaps...)...)
}

// Info loga uma informação
func (l logger) Info(ctx context.Context, message string, valueMaps ...logger2.Body) {
	l.log.Sugar().Infow(message, convertBodyToFields(ctx, valueMaps...)...)
}

// Warn loga uma advertência
func (l logger) Warn(ctx context.Context, message string, valueMaps ...logger2.Body) {
	l.log.Sugar().Warnw(message, convertBodyToFields(ctx, valueMaps...)...)
}

// Error loga um erro
func (l logger) Error(ctx context.Context, message string, valueMaps ...logger2.Body) {
	l.log.Sugar().Errorw(message, convertBodyToFields(ctx, valueMaps...)...)
	err := l.log.Sync()
	if err != nil {
		return
	}
}

// Fatal loga uma mensagem e encerra a aplicação
func (l logger) Fatal(ctx context.Context, message string, valueMaps ...logger2.Body) {
	l.log.Sugar().Fatalw(message, convertBodyToFields(ctx, valueMaps...)...)
	err := l.log.Sync()
	if err != nil {
		return
	}
}
