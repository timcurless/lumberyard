package lumberyard

import (
  "context"
  "time"

  "github.com/go-kit/kit/log"
)

type Middleware func(Store) Store

func LoggingMiddleware(logger log.Logger) Middleware {
  return func(next Store) Store {
    return &loggingMiddleware{
      next: next,
      logger: logger,
    }
  }
}

type loggingMiddleware struct {
  next Store
  logger log.Logger
}

func (mw loggingMiddleware) PostDeployment(ctx context.Context, d Deployment) (err error) {
  defer func(begin time.Time) {
    mw.logger.Log("method", "PostDeployment", "name", d.Name, "took", time.Since(begin), "err", err)
  }(time.Now())
  return mw.next.PostDeployment(ctx, d)
}

func (mw loggingMiddleware) GetDeployments(ctx context.Context) (deployments []Deployment, err error) {
  defer func(begin time.Time) {
    mw.logger.Log("method", "GetDeployments", "took", time.Since(begin), "err", err)
  }(time.Now())
  return mw.next.GetDeployments(ctx)
}
