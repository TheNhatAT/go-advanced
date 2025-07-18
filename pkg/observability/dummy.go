// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package observability

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

const XTimes = 100

func Prepare() { fmt.Println("initializing operation!") }

func DoOperation() error {
	// Do some dummy, randomized heavy work (both in terms of latency, CPU and memory usage).
	alloc := make([]byte, 1e6)
	for i := 0; i < int(rand.Float64()*100); i++ {
		_ = fmt.Sprintf("doing stuff! %+v", alloc)
	}

	runtime.GC() // To have more interesting GC metrics.

	switch rand.Intn(3) {
	case 0:
		return nil
	case 1:
		return errors.New("error first")
	case 2:
		return errors.New("error other")
	}
	return nil
}

func DoOperationWithCtx(ctx context.Context) error {
	tracer := otel.Tracer("observability")

	ctx, span := tracer.Start(ctx, "first operation")
	defer span.End()

	// Do some dummy, randomized heavy work (both in terms of latency, CPU and memory usage).
	alloc := make([]byte, 1e6)
	for i := 0; i < int(rand.Float64()*125); i++ {
		_ = fmt.Sprintf("doing stuff! %+v", alloc) // Hope for this to not get cleared by compiler.
	}

	runtime.GC() // To have more interesting GC metrics.

	// ignore handling error
	_ = doInSpan(ctx, "sub operation2", func(ctx context.Context) error {
		return nil
	})
	_ = doInSpan(ctx, "sub operation3", func(ctx context.Context) error {
		return nil
	})

	return doInSpan(
		ctx,
		"choosing error",
		func(ctx context.Context) error {
			switch rand.Intn(3) {
			default:
				return nil
			case 1:
				time.Sleep(300 * time.Millisecond) // For more interesting results.
				return errors.New("error first")
			case 2:
				return errors.New("error other")
			}
		})
}

func doInSpan(ctx context.Context, name string, fn func(context.Context) error) error {
	tracer := otel.Tracer("observability")
	ctx, span := tracer.Start(ctx, name)
	defer span.End()

	err := fn(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}

func TearDown() { fmt.Println("closing operation!") }

func ErrorType(err error) string {
	if err != nil {
		if err.Error() == "error first" {
			return "error1"
		}
		return "other_error"
	}
	return ""
}

func PrintPrometheusMetrics(reg prometheus.Gatherer) {
	rec := httptest.NewRecorder()
	promhttp.HandlerFor(reg, promhttp.HandlerOpts{DisableCompression: true, EnableOpenMetrics: true}).ServeHTTP(rec, &http.Request{})
	if rec.Code != 200 {
		panic("unexpected error")
	}

	fmt.Println(rec.Body.String())
}
