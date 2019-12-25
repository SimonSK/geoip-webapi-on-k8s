package webapi

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	IPv6HexMaxLen      = 39 // Maximum length of a valid IPv6 address in hex, including colons
	routeVariableKeyIP = "ipAddress"
)

var (
	rootPath                     = "/"
	apiRootPath                  = "/api"
	apiCityDetailsPath           = apiRootPath + fmt.Sprintf("/{%s}", routeVariableKeyIP)
	apiCityLocationPath          = apiCityDetailsPath + "/location"
	apiCityLocationGPSCoordsPath = apiCityLocationPath + "/coords"
)

func respond(w http.ResponseWriter, status int, contentType string, data interface{}) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", contentType)
	switch contentType {
	case "application/json":
		return json.NewEncoder(w).Encode(data)
	default:
		_, err := fmt.Fprintf(w, "%v", data)
		return err
	}
}

func (s *Server) getDefault(w http.ResponseWriter, r *http.Request) {
	route := "default"
	contentType := "text/plain"
	var (
		status int
		result interface{}
	)
	s.Log.Infof("[method=%s route=%s remoteAddr=%s] req", r.Method, route, r.RemoteAddr)
	switch r.Method {
	case "GET":
		status = http.StatusOK
		result = s.info()
	default:
		status = http.StatusNotImplemented
	}
	s.Log.Infof("[status=%d route=%s remoteAddr=%s] resp", status, route, r.RemoteAddr)
	if err := respond(w, status, contentType, result); err != nil {
		s.Log.Error(err)
	}
	return
}

func (s *Server) getCityDetails(w http.ResponseWriter, r *http.Request) {
	route := "details"
	contentType := "application/json"
	var (
		status int
		result interface{}
	)
	ipAddrStr := mux.Vars(r)[routeVariableKeyIP]
	s.Log.Infof("[method=%s route=%s remoteAddr=%s query=%.*s] req", r.Method, route, r.RemoteAddr, IPv6HexMaxLen, ipAddrStr)
	ipAddr := net.ParseIP(ipAddrStr)
	if ipAddr == nil {
		s.Log.Debugf("unable to parse %s as an IP address", ipAddrStr)
		status = http.StatusBadRequest
	} else {
		switch r.Method {
		case "GET":
			if record, err := s.db.GetRecord(ipAddr); err == nil {
				status = http.StatusOK
				result = record
			} else {
				status = http.StatusInternalServerError
			}
		default:
			status = http.StatusNotImplemented
		}
	}
	s.Log.Infof("[status=%d route=%s remoteAddr=%s] resp", status, route, r.RemoteAddr)
	if err := respond(w, status, contentType, result); err != nil {
		s.Log.Error(err)
	}
	return
}

func (s *Server) getCityLocation(w http.ResponseWriter, r *http.Request) {
	route := "location"
	contentType := "application/json"
	var (
		status int
		result interface{}
	)
	ipAddrStr := mux.Vars(r)[routeVariableKeyIP]
	s.Log.Infof("[method=%s route=%s remoteAddr=%s query=%.*s] req", r.Method, route, r.RemoteAddr, IPv6HexMaxLen, ipAddrStr)
	ipAddr := net.ParseIP(ipAddrStr)
	if ipAddr == nil {
		s.Log.Debugf("unable to parse %s as an IP address", ipAddrStr)
		status = http.StatusBadRequest
	} else {
		status = http.StatusInternalServerError
		switch r.Method {
		case "GET":
			if record, err := s.db.GetRecord(ipAddr); err == nil {
				if loc, ok := record.(map[string]interface{})["location"]; ok {
					status = http.StatusOK
					result = loc
				}
			}
		default:
			status = http.StatusNotImplemented
		}
	}
	s.Log.Infof("[status=%d route=%s remoteAddr=%s] resp", status, route, r.RemoteAddr)
	if err := respond(w, status, contentType, result); err != nil {
		s.Log.Error(err)
	}
	return
}

func (s *Server) getCityLocationGPSCoords(w http.ResponseWriter, r *http.Request) {
	route := "coords"
	contentType := "application/json"
	var (
		status int
		result interface{}
	)
	ipAddrStr := mux.Vars(r)[routeVariableKeyIP]
	s.Log.Infof("[method=%s route=%s remoteAddr=%s query=%.*s] req", r.Method, route, r.RemoteAddr, IPv6HexMaxLen, ipAddrStr)
	ipAddr := net.ParseIP(ipAddrStr)
	if ipAddr == nil {
		s.Log.Debugf("unable to parse %s as an IP address", ipAddrStr)
		status = http.StatusBadRequest
	} else {
		status = http.StatusInternalServerError
		switch r.Method {
		case "GET":
			if record, err := s.db.GetRecord(ipAddr); err == nil {
				if loc, ok := record.(map[string]interface{})["location"]; ok {
					lat, latOk := loc.(map[string]interface{})["latitude"].(float64)
					lon, lonOk := loc.(map[string]interface{})["longitude"].(float64)
					if latOk && lonOk {
						status = http.StatusOK
						result = map[string]float64{
							"lat": lat,
							"lon": lon,
						}
					}
				}
			}
		default:
			status = http.StatusNotImplemented
		}
	}
	s.Log.Infof("[status=%d route=%s remoteAddr=%s] resp", status, route, r.RemoteAddr)
	if err := respond(w, status, contentType, result); err != nil {
		s.Log.Error(err)
	}
	return
}
