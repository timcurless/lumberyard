package service

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// Middleware is a middleware wrapped around our service
type Middleware func(Service) Service

// LoggingMiddleware is a middleware for logging
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) PostProject(ctx context.Context, p Project) (id string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostProject", "id", p.ID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostProject(ctx, p)
}

func (mw loggingMiddleware) GetProject(ctx context.Context, id string) (p Project, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetProject", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetProject(ctx, id)
}

func (mw loggingMiddleware) PostStack(ctx context.Context, projectID string, st Stack) (id string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostStack", "projectID", projectID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostStack(ctx, projectID, st)
}
