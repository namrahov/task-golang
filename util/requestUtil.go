package util

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"task-golang/model"
)

func DecodeBody(w http.ResponseWriter, r *http.Request, data interface{}) error {
	logger, ok := r.Context().Value(model.ContextLogger).(*log.Entry)
	if !ok {
		logger.Errorf("requestUtil.DecodeBody.error: cannot get ContextLogger")
		return errors.New(fmt.Sprintf("%s.can't-get-ContextLogger", model.Exception))
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Errorf("requestUtil.DecodeBody.error: cannot close Body %v", err)
			_ = errors.New(fmt.Sprintf("%s.can't-close-Body", model.Exception))
		}
	}(r.Body)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("Error reading request body", "error", err)
		http.Error(w, "Bad Request: Could not read body", http.StatusBadRequest) // Generic client error msg
		return err
	}

	if err := json.Unmarshal(body, data); err != nil {
		logger.Error("Error unmarshaling json body", "error", err)
		http.Error(w, "Bad Request: Invalid JSON body", http.StatusBadRequest)
		return err
	}

	return nil
}
