package goapi

import (
	"context"

	"cloud.google.com/go/spanner"
)

func dbString(config *Config) string {
	return "projects/" + config.Spanner.Project + "/instances/" + config.Spanner.Instance + "/databases/" + config.Spanner.Database
}

var SpannerClient *spanner.Client

// SpannerClient ...
func (config *Config) SpannerClient() {
	var err error
	ctx := context.Background()

	SpannerClient, err = spanner.NewClient(ctx, dbString(config))
	if err != nil {
		panic("Error! spanner.NewClient: " + err.Error())
	}
}
