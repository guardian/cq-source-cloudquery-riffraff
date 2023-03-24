package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudquery/plugin-sdk/plugins/source"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/guardian/cq-source-cloudquery-riffraff/riffraff"
	"github.com/rs/zerolog"
)

type Spec struct {
	APIKey      string `json:"api_key"`
	RiffraffURL string `json:"riffraff_url"`
}

type Client struct {
	Logger   zerolog.Logger
	Riffraff riffraff.Client
}

func (c *Client) ID() string {
	return "riffraff"
}

func New(ctx context.Context, logger zerolog.Logger, s specs.Source, opts source.Options) (schema.ClientMeta, error) {
	var pluginSpec Spec

	if err := s.UnmarshalSpec(&pluginSpec); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plugin spec: %w", err)
	}

	if pluginSpec.APIKey == "" {
		return nil, errors.New("api_key cannot be empty")
	}

	if pluginSpec.RiffraffURL == "" {
		return nil, errors.New("riffraff_url cannot be empty")
	}

	return &Client{
		Logger:   logger,
		Riffraff: riffraff.New(pluginSpec.RiffraffURL, pluginSpec.APIKey),
	}, nil
}
