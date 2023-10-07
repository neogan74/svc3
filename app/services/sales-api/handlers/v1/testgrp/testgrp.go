package testgrp

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/neogan74/svc3/app/fondation/web"
	"go.uber.org/zap"
)

type Handlers struct {
	Log *zap.SugaredLogger
}

func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "OK",
	}
	json.NewEncoder(w).Encode(status)

	statusCode := http.StatusOK

	h.Log.Infow("v1.test", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

	return web.Response(ctx, w, status, statusCode)
}
