package msg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

//Alerts struct
type Alerts struct {
	Alerts            []Alert           `json:"alerts"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	CommonLabels      map[string]string `json:"commonLabels"`
	ExternalURL       string            `json:"externalURL"`
	GroupKey          string            `json:"groupKey"`
	GroupLabels       map[string]string `json:"groupLabels"`
	Receiver          string            `json:"receiver"`
	Status            string            `json:"status"`
	Version           string            `json:"version"`
}

//Alert struct
type Alert struct {
	Annotations  map[string]string `json:"annotations"`
	GeneratorURL string            `json:"generatorURL"`
	// Labels       Label             `json:"labels"`
	Labels   map[string]string `json:"labels"`
	EndsAt   string            `json:"sendsAt"`
	StartsAt string            `json:"startsAt"`
}

//Label struct
type Label struct {
	Alertname string `json:"alertname,omitempty"`
	Env       string `json:"env,omitempty"`
	Instance  string `json:"instance,omitempty"`
	Job       string `json:"job,omitempty"`
	Monitor   string `json:"monitor,omitempty"`
	Severity  string `json:"severity,omitempty"`
}

//Parser use for parsing input io msg end return Alerts
func Parser(r io.Reader) (Alerts, string, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)

	var alerts Alerts
	if err := json.NewDecoder(tee).Decode(&alerts); err != nil {
		return alerts, "", fmt.Errorf("Error decode alert message %v", err)
	}
	return alerts, buf.String(), nil
}
