package handlers

import (
	"project/model"
	"project/service"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io"
	"mime"
	"net/http"
	"strconv"
)

type ConfigHandler struct {
	service service.ConfigService
}

func NewConfigHandler(service service.ConfigService) ConfigHandler {
	return ConfigHandler{
		service: service,
	}
}

func decodeBody(reader io.Reader) (*model.Config, error) {
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	var rt model.Config
	if err := decoder.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func renderJSON(writer http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

func (c ConfigHandler) Get(writer http.ResponseWriter, req *http.Request) {
	name := mux.Vars(req)["name"]
	version := mux.Vars(req)["version"]

	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	config, err := c.service.Get(name, versionInt)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(config)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content−Type", "application/json")
	writer.Write(resp)
}

func (c ConfigHandler) Add(writer http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != " application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(writer, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeBody(req.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	con := model.Config{
		Name:    rt.Name,
		Version: rt.Version,
		Params:  rt.Params,
	}
	c.service.Add(con)

	renderJSON(writer, rt)
}

func (c ConfigHandler) Delete(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	name := vars["name"]
	versionStr := vars["version"]

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(writer, "Invalid version", http.StatusBadRequest)
		return
	}

	err = c.service.Delete(name, version)
	if err != nil {
		http.Error(writer, "Failed to delete config", http.StatusInternalServerError)
		return
	}

	renderJSON(writer, "Deleted")
}