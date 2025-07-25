// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/aws-gopher/unstructured-sdk-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure Provider satisfies various provider interfaces.
var _ provider.Provider = &Provider{}

// Provider defines the provider implementation.
type Provider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
	client  *unstructured.Client
}

// ProviderModel describes the provider data model.
type ProviderModel struct {
	APIKey   types.String `tfsdk:"api_key"`
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "unstructured"
	resp.Version = p.version
}

func (p *Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Optional:    true,
				Description: "The API key for the Unstructured API",
				Sensitive:   true,
			},
			"endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "The endpoint of the API",
			},
		},
	}
}

func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ProviderModel

	opts := []unstructured.Option{}

	// Check environment variables
	apiKey := os.Getenv("UNSTRUCTURED_API_KEY")
	endpoint := os.Getenv("UNSTRUCTURED_API_URL")

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check configuration data, which should take precedence over
	// environment variable data, if found.
	if data.APIKey.ValueString() != "" {
		apiKey = data.APIKey.ValueString()
	}

	if data.Endpoint.ValueString() != "" {
		endpoint = data.Endpoint.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing API Key Configuration",
			"While configuring the provider, the API key was not found in "+
				"the UNSTRUCTURED_API_KEY environment variable or provider "+
				"configuration block api_key attribute.",
		)
		// Not returning early allows the logic to collect all errors.
	}

	// Create data/clients and persist to resp.DataSourceData, resp.ResourceData,
	// and resp.EphemeralResourceData as appropriate.

	if apiKey != "" {
		opts = append(opts, unstructured.WithKey(apiKey))
	}

	if endpoint != "" {
		opts = append(opts, unstructured.WithEndpoint(endpoint))
	}

	client, err := unstructured.New(opts...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Unstructured client",
			err.Error(),
		)
	}
	p.client = client
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *Provider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewWorkflowResource,
		NewSourceResource,
		NewDestinationResource,
	}
}

func (p *Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewWorkflowDataSource,
		NewSourceDataSource,
		NewDestinationDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &Provider{
			version: version,
		}
	}
}
