package plugin

import (
	"github.com/cloudquery/plugin-sdk/plugins/source"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/guardian/cq-source-cloudquery-riffraff/client"
	"github.com/guardian/cq-source-cloudquery-riffraff/resources"
)

var (
	Version = "development"
)

func Plugin() *source.Plugin {
	return source.NewPlugin(
		"guardian-cloudquery-riffraff",
		Version,
		schema.Tables{
			resources.HistoryTable(),
		},
		client.New,
	)
}
