package hashicups

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

type dataSourceCoffeesType struct{}

func (r dataSourceCoffeesType) GetSchema(_ context.Context) (schema.Schema, []*tfprotov6.Diagnostic) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"coffees": {
				Computed: true,
				Attributes: schema.ListNestedAttributes(map[string]schema.Attribute{
					"orderid": {
						Type:     types.NumberType,
						Computed: true,
					},
					"name": {
						Type:     types.StringType,
						Computed: true,
					},
					"teaser": {
						Type:     types.StringType,
						Computed: true,
					},
					"description": {
						Type:     types.StringType,
						Computed: true,
					},
					"price": {
						Type:     types.NumberType,
						Computed: true,
					},
					"image": {
						Type:     types.StringType,
						Computed: true,
					},
				}, schema.ListNestedAttributesOptions{}),
			},
		},
	}, nil
}

func (r dataSourceCoffeesType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, []*tfprotov6.Diagnostic) {
	return dataSourceCoffees{
		p: *(p.(*provider)),
	}, nil
}

type dataSourceCoffees struct {
	p provider
}

func (r dataSourceCoffees) Read(ctx context.Context, p tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	fmt.Fprintln(stderr, "[DEBUG]-read-error1:", p.Config.Raw)

	var state struct {
		Coffee []Coffee `tfsdk:"coffees"`
	}
	err := p.Config.Get(ctx, &state)

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error reading coffee",
			Detail:   "An unexpected error was encountered while reading the datasource_coffee: " + err.Error(),
		})
	}
	r.p.client.GetCoffees()
}
