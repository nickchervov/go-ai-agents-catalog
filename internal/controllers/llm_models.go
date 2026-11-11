package controllers

import "net/http"

func (h *AgentsHandler) Models(w http.ResponseWriter, r *http.Request) {
	output, err := h.svc.Models(r.Context())
	if err != nil {
		JSONError(r.Context(), w, err)
		return
	}

	JSON(w, http.StatusOK, output)
}
