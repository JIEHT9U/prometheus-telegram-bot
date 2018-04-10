package msg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

//Alerts struct
type Alerts struct {
	Alerts            []Alert           `json:"Alerts"`
	CommonAnnotations map[string]string `json:"CommonAnnotations"`
	CommonLabels      map[string]string `json:"CommonLabels"`
	ExternalURL       string            `json:"ExternalURL"`
	GroupKey          string            `json:"GroupKey"`
	GroupLabels       map[string]string `json:"GroupLabels"`
	Receiver          string            `json:"Receiver"`
	Status            string            `json:"Status"`
	Version           string            `json:"Version"`
}

//Alert struct
type Alert struct {
	Annotations  map[string]string `json:"Annotations"`
	GeneratorURL string            `json:"GeneratorURL"`
	Labels       map[string]string `json:"Labels"`
	EndsAt       string            `json:"SendsAt"`
	StartsAt     string            `json:"StartsAt"`
	// Labels       Label             `json:"labels"`
}

// //Label struct
// type Label struct {
// 	Alertname string `json:"Alertname,omitempty"`
// 	Env       string `json:"Env,omitempty"`
// 	Instance  string `json:"Instance,omitempty"`
// 	Job       string `json:"Job,omitempty"`
// 	Monitor   string `json:"Monitor,omitempty"`
// 	Severity  string `json:"Severity,omitempty"`
// }

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
