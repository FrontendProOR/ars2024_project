package handler

import (
	"project/model"
	"project/service"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"mime"
	"net/http"
	"strconv"
)

type ConfigGroupHandler struct {
	service service.ConfigGroupService
}

func NewConfigGruopHandler(service service.ConfigGroupService) ConfigGroupHandler {
	return ConfigGroupHandler{
		service: service,
	}
}

func decodeBodyCG(reader io.Reader) (*model.ConfigGroup, error) {
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	var rt model.ConfigGroup
	if err := decoder.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func (cgh ConfigGroupHandler) Get(writer http.ResponseWriter, req *http.Request) {
	name := mux.Vars(req)["name"]
	version := mux.Vars(req)["version"]

	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	config, err := cgh.service.Get(name, versionInt)

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

func (cgh ConfigGroupHandler) Add(writer http.ResponseWriter, req *http.Request) {
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

	rt, err := decodeBodyCG(req.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	configGroup := model.ConfigGroup{
		Name:    rt.Name,
		Version: rt.Version,
		Configs: rt.Configs,
	}
	cgh.service.Add(configGroup)

	renderJSON(writer, rt)
}

func (cgh ConfigGroupHandler) Delete(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	name := vars["name"]
	versionStr := vars["version"]

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(writer, "Invalid version", http.StatusBadRequest)
		return
	}

	err = cgh.service.Delete(name, version)
	if err != nil {
		http.Error(writer, "Failed to delete config group", http.StatusInternalServerError)
		return
	}

	renderJSON(writer, "Deleted")
}

func (cgh ConfigGroupHandler) AddConfToGroup(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	nameG := vars["nameG"]
	versionGStr := vars["versionG"]
	nameC := vars["nameC"]
	versionCStr := vars["versionC"]

	versionG, err := strconv.Atoi(versionGStr)
	if err != nil {
		http.Error(writer, "Invalid version", http.StatusBadRequest)
		return
	}

	versionC, err := strconv.Atoi(versionCStr)
	if err != nil {
		http.Error(writer, "Invalid version", http.StatusBadRequest)
		return
	}

	key := fmt.Sprintf("%s/%d", nameC, versionC)

	group, _ := cgh.service.Get(nameG, versionG)
	conf, _ := service.ConfigService{}.Get(nameC, versionC)
	group.Configs[key] = &conf

	renderJSON(writer, "success Put")
}

func (cgh ConfigGroupHandler) RemoveConfFromGroup(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	nameG := vars["nameG"]
	versionGStr := vars["versionG"]
	nameC := vars["nameC"]
	versionCStr := vars["versionC"]

	versionG, err := strconv.Atoi(versionGStr)
	if err != nil {
		http.Error(writer, "Invalid version", http.StatusBadRequest)
		return
	}

	versionC, err := strconv.Atoi(versionCStr)
	if err != nil {
		http.Error(writer, "Invalid version", http.StatusBadRequest)
		return
	}

	key := fmt.Sprintf("%s/%d", nameC, versionC)

	group, _ := cgh.service.Get(nameG, versionG)
	delete(group.Configs, key)

	renderJSON(writer, "success Put")
}