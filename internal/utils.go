package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SendJson(data any, w http.ResponseWriter) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("sendJson: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)

	return nil
}
