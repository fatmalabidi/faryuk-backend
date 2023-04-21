package utils

import (
	"FaRyuk/internal/types"
	"encoding/json"
	"net/http"
)

func WriteForbidden(w *http.ResponseWriter, m string) {
  (*w).WriteHeader(http.StatusForbidden)
  writeResponse(w, types.JSONReturn{Status: "Fail", Body: m})
}

func WriteNotFound(w *http.ResponseWriter, m string) {
  (*w).WriteHeader(http.StatusNotFound)
  writeResponse(w, types.JSONReturn{Status: "Fail", Body: m})
}

func WriteInternalError(w *http.ResponseWriter, m string) {
  (*w).WriteHeader(http.StatusInternalServerError)
  writeResponse(w, types.JSONReturn{Status: "Fail", Body: m})
}

func WriteObject(w *http.ResponseWriter, m interface{}) {
  writeResponse(w, types.JSONReturn{Status: "Success", Body: m})
}

func writeResponse(w *http.ResponseWriter, m types.JSONReturn) {
  (*w).Header().Add("Content-Type", "application/json")
  err := json.NewEncoder(*w).Encode(m)
	if err != nil {
		return
	}
}