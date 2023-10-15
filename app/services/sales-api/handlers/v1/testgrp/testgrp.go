package testgrp

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"

	"github.com/neogan74/svc3/app/fondation/web"
	"github.com/neogan74/svc3/business/sys/validate"
	"go.uber.org/zap"
)

type Handlers struct {
	Log *zap.SugaredLogger
}

func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	if n := rand.Intn(100); n%2 == 0 {
		return errors.New("Untrusted error")
	}

	if n := rand.Intn(200); n%2 == 0 {
		return validate.NewRequestError(errors.New("Trusted error"), http.StatusBadRequest)
	}

	if n := rand.Intn(200); n%2 == 0 {
		return web.NewShutdownError("going down now")
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
