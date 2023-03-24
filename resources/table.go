package resources

import (
	"context"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"github.com/guardian/cq-source-cloudquery-riffraff/client"
	"github.com/guardian/cq-source-cloudquery-riffraff/riffraff"
)

func HistoryTable() *schema.Table {
	return &schema.Table{
		Name:      "riffraff_history",
		Resolver:  fetchHistoryTable,
		Transform: transformers.TransformWithStruct(&riffraff.HistoryItem{}),
	}
}

// It doesn't look like cursor-based pagination is supported in RR so to get
// started we simply read deploys from the last 30 days. It would be great to
// make this incremental though and ingest the entire history of data.
func fetchHistoryTable(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	c := meta.(*client.Client)
	pageNum := 1
	params := riffraff.HistoryParams{Page: pageNum, PageSize: 100, MaxDaysAgo: 7}
	paginator := riffraff.NewHistoryPaginator(c.Riffraff, params)

	for paginator.HasMorePages() {
		items, err := paginator.NextPage()
		if err != nil {
			return err
		}

		for _, item := range items {
			res <- item
		}
	}

	return nil
}
