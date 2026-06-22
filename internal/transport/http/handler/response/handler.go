package response

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Zerferrous/Task-Manager-API/internal/logger"
)

type Handler struct {
	logger *logger.Logger
	w      http.ResponseWriter
}

func NewHandler(logger *logger.Logger, w http.ResponseWriter) *Handler {
	return &Handler{
		logger: logger,
		w:      w,
	}
}

func (h *Handler) Panic(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.logger.Error(msg, slog.Any("error", err))
	h.w.WriteHeader(statusCode)

	response := map[string]string{
		"msg":   msg,
		"error": err.Error(),
	}

	if err := json.NewEncoder(h.w).Encode(response); err != nil {
		h.logger.Error("write http response error", slog.Any("error", err))
	}
}
