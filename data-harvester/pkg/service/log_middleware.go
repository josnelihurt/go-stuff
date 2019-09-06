package service

import (
	"context"
	"encoding/json"

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
		val, _ := json.Marshal(result)
		context.logger.Log("method", "Collect", "param", "none", "result", val, "err", err)
	}()
	result, err = context.next.Collect(ctx, param)
	return result, err
}
func (context *loggingMiddleware) Status(ctx context.Context) (result bool, err error) {
	defer func() {
		context.logger.Log("method", "Status", "result", result, "err", err)
	}()
	return context.next.Status(ctx)
}
