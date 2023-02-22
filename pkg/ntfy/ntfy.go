package ntfy

import (
	"fmt"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Alert struct {
	Status       string
	Labels       map[string]string
	Annotations  map[string]string
	StartsAt     string
	EndsAt       string
	GeneratorURL string
}

type Payload struct {
	Receiver          string
	Status            string
	Alerts            []Alert
	GroupLabels       map[string]string
	CommonLabels      map[string]string
	CommonAnnotations map[string]string
	ExternalURL       string
	Version           string
	GroupKey          string
}

// https://prometheus.io/docs/alerting/latest/configuration/#webhook_config
// The Alertmanager will send HTTP POST requests in the following JSON format to the configured endpoint:
// {
//   "version": "4",
//   "groupKey": <string>,              // key identifying the group of alerts (e.g. to deduplicate)
//   "truncatedAlerts": <int>,          // how many alerts have been truncated due to "max_alerts"
//   "status": "<resolved|firing>",
//   "receiver": <string>,
//   "groupLabels": <object>,
//   "commonLabels": <object>,
//   "commonAnnotations": <object>,
//   "externalURL": <string>,           // backlink to the Alertmanager.
//   "alerts": [
//     {
//       "status": "<resolved|firing>",
//       "labels": <object>,
//       "annotations": <object>,
//       "startsAt": "<rfc3339>",
//       "endsAt": "<rfc3339>",
//       "generatorURL": <string>,      // identifies the entity that caused the alert
//       "fingerprint": <string>        // fingerprint to identify the alert
//     },
//     ...
//   ]
// }

func mapping(payload Payload) error {
	for _, alert := range payload.Alerts {
		req, err := http.NewRequest("POST", "https://ntfy.sh/dictum", strings.NewReader(alert.Labels["alertname"]))
		if err != nil {
			return err
		}
		// req.Header.Set("Attach", "https://nest.com/view/yAxkasd.jpg")
		// req.Header.Set("Actions", "http, Open door, https://api.nest.com/open/yAxkasd, clear=true")
		// req.Header.Set("Email", "phil@example.com")
		req.Header.Set("Title", fmt.Sprintf("%s: %s", alert.Status, alert.Annotations["description"]))
		// req.Header.Set("Click", alert.GeneratorURL)
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
	}

	return nil
}

func Send(c *gin.Context) {
	var payload Payload
	c.BindJSON(&payload)

	err := mapping(payload)
	if err != nil {
		util.Error(c, err)
	}

	util.Success(c, nil)
}
