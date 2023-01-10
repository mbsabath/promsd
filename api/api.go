package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type sdHandler struct {
	Targets map[string][]map[string]string
}

type staticConfig struct {
	Target string              `json:"target"`
	Labels []map[string]string `json:"labels"`
}

func NewSdHandler() *sdHandler {
	out := sdHandler{
		Targets: make(map[string][]map[string]string),
	}

	return &out
}

func (h *sdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	switch r.Method {
	case http.MethodPut:
		h.addTarget(w, r)
	case http.MethodDelete:
		h.deleteTarget(w, r)
	case http.MethodGet:
		h.returnTargets(w)
	}
}

func (h *sdHandler) addTarget(w http.ResponseWriter, r *http.Request) {
	request_data := staticConfig{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&request_data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Data Passed", http.StatusBadRequest)
	} else {
		h.Targets[request_data.Target] = request_data.Labels
		w.Write([]byte(fmt.Sprintf("Target %s added\n", request_data.Target)))
	}

}

func (h *sdHandler) deleteTarget(w http.ResponseWriter, r *http.Request) {
	target, ok := r.Form["target"]
	if !ok {
		http.Error(w, "No deletion target specified\n", http.StatusBadRequest)
		return
	}

	_, ok = h.Targets[target[0]]
	if !ok {
		http.Error(w, "Invalid target specified\n", http.StatusBadRequest)
	} else {
		delete(h.Targets, target[0])
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Target %s deleted\n", target[0])))
		fmt.Printf("Target %s deleted\n", target[0])
		fmt.Println(h.Targets)
	}
}

func (h *sdHandler) returnTargets(w http.ResponseWriter) {
	config_array := make([]staticConfig, len(h.Targets))
	i := 0
	for target, labels := range h.Targets {
		config_array[i] = staticConfig{Target: target, Labels: labels}
		i++
	}
	out, err := json.Marshal(config_array)
	if err != nil {
		http.Error(w, "Error creating response", http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	}

}
