package web

import (
	"context"
	"encoding/json"
	"net/http"
)

func Response(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {

	SetStatusCode(ctx, statusCode)

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// Set content type and headers once we know marshaling has suceeded
	w.Header().Set("Content-type", "application/json")

	// Write status code to the responce
	w.WriteHeader(statusCode)

	// Send the resule back to the client.
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
