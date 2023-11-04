package testgrp

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/neogan74/svc3/fondation/web"
	"go.uber.org/zap"
)

type Handlers struct {
	Log *zap.SugaredLogger
}

func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	// if n := rand.Intn(100); n%2 == 0 {
	// 	return errors.New("Untrusted error")
	// }

	// if n := rand.Intn(200); n%2 == 0 {
	// 	return validate.NewRequestError(errors.New("Trusted error"), http.StatusBadRequest)
	// }

	// if n := rand.Intn(200); n%2 == 0 {
	// 	return web.NewShutdownError("going down now")
	// }

	if n := rand.Intn(1000); n%2 == 0 {
		panic("testing panic")
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}
	json.NewEncoder(w).Encode(status)

	statusCode := http.StatusOK

	h.Log.Infow("v1.test", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

	return web.Respond(ctx, w, status, statusCode)
}
