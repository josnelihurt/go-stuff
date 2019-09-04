package service

import (
	"context"

	log "github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   DataHarvestService //? ptr?
}

// LoggingMiddleware takes a logger as a dependency
// and returns a pingService Middleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next DataHarvestService) DataHarvestService {
		return &loggingMiddleware{logger, next}
	}
}

//Collect implemented from Service
func (context *loggingMiddleware) Collect(ctx context.Context, param DataHarvestServiceParam) (result DataHarvestServiceResult, err error) {
	defer func() {
		context.logger.Log("method", "Collect", "param", param, "result", result, "err", err)
	}()
	return context.next.Collect(ctx, param)
}
func (context *loggingMiddleware) Status(ctx context.Context) (result bool, err error) {
	defer func() {
		context.logger.Log("method", "Status", "result", result, "err", err)
	}()
	return false, nil
}
