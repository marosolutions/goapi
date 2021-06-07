package goapi

import (
	"context"

	"cloud.google.com/go/spanner"
)

// spannerClient ...
func (config *Config) newSpannerClient() *spanner.Client {
	dbString := "projects/" + config.Spanner.Project +
		"/instances/" + config.Spanner.Instance +
		"/databases/" + config.Spanner.Database

	client, err := spanner.NewClient(context.Background(), dbString)
	if err != nil {
		panic("Error! spanner.NewClient: " + err.Error())
	}
	return client
}
