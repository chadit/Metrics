package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/chadit/Metrics/Internal/cors"
	"github.com/chadit/Metrics/Internal/logger"
	"github.com/chadit/Metrics/Internal/metric"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type service struct {
	collection *metric.Collection
	logger     *zap.SugaredLogger
}

func newService() *service {
	return &service{
		collection: metric.New(),
		logger:     logger.New(),
	}
}

func (s service) create(ctx *fasthttp.RequestCtx) {
	key := ctx.UserValue("key").(string)

	var body metric.Request

	if err := json.Unmarshal(ctx.PostBody(), &body); err != nil {
		ctx.Error("Invalid data", fasthttp.StatusBadRequest)
		return
	}

	if err := s.collection.Add(key, body.Value); err != nil {
		ctx.Error("Saving metric", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (s service) sum(ctx *fasthttp.RequestCtx) {
	key := ctx.UserValue("key").(string)

	defaultDuration := time.Duration(90 * time.Minute)

	sum, err := s.collection.Sum(key, defaultDuration)
	if err != nil {
		if err == metric.ErrNotFound {
			ctx.NotFound()
			return
		}

		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	body := metric.Response{
		Value: sum,
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")

	if err := json.NewEncoder(ctx).Encode(body); err != nil {
		ctx.Error("returning total", fasthttp.StatusInternalServerError)
		return
	}

}

func (s service) cleanup() {
	time.Sleep(250 * time.Millisecond)
	if err := s.collection.Purge(time.Duration(90 * time.Minute)); err != nil {
		log.Printf("error purging %v", err)
	}
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	s := newService()
	r := router.New()

	r.POST("/metric/{key}", s.create)
	r.GET("/metric/{key}/sum", s.sum)

	go s.cleanup()

	s.logger.Info("Service started on port: 3000")
	s.logger.Fatal(fasthttp.ListenAndServe(":3000", cors.New(r.Handler)))
}
