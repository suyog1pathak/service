package healthcheck

import "time"

type Components struct {
	Datastore bool `json:"datastore" yaml:"datastore"`
} //@name Components

// Response is a response body for healthcheck endpoint.
type Response struct {
	Status     string     `json:"status" yaml:"status"`
	StatusCode int        `json:"statusCode" yaml:"statusCode"`
	Components Components `json:"components" yaml:"components"`
	Timestamp  time.Time  `json:"timestamp" yaml:"timestamp"`
} //@name HealthcheckResponse
