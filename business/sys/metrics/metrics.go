package metrics

import (
	"context"
	"expvar"
)

// This holds the single instance of the metrics value needed for
// collecting metrics. The expvar package is already based on a singleton
// for the different metrics that are registered with the package so there
// isn't much choice here.
var m *metrics

// Metrics represents the sert of metrics we gather. These fields are
// safe to be accessed concurectly thant to expvar. No extra abstraction is required.
type metrics struct {
	goroutines *expvar.Int
	reuests    *expvar.Int
	errors     *expvar.Int
	panics     *expvar.Int
}

func init() {
	m = &metrics{
		goroutines: expvar.NewInt("goroutines"),
		reuests:    expvar.NewInt("requests"),
		errors:     expvar.NewInt("errors"),
		panics:     expvar.NewInt("panics"),
	}
}

type ctxKey int

// key is now metrics
const key ctxKey = 1

func Set(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, m)
}

func AddGoroutines(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		if v.reuests.Value()%100 == 0 {
			v.goroutines.Add(1)
		}
	}
}

func AddRequests(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.reuests.Add(1)
	}
}

func AddErrors(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.errors.Add(1)
	}
}

func AddPanics(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.panics.Add(1)
	}
}
