package hashicups

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

type dataSourceOrderType struct{}

func (r dataSourceOrderType) GetSchema(_ context.Context) (schema.Schema, []*tfprotov6.Diagnostic) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": {
				Type:     types.NumberType,
				Computed: true,
			},
			"items": {
				//tf will throw error if user doesn't specify value - optional - can or choose not to supply a value
				Required: false,
				Attributes: schema.SingleNestedAttributes(map[string]schema.Attribute{
					"id": {
						Type:     types.NumberType,
						Required: true,
					},
					"name": {
						Required: true,
						Type:     types.StringType,
					},
					"teaser": {
						Type:     types.StringType,
						Required: true,
					},
					"description": {
						Type:     types.StringType,
						Required: true,
					},
					"price": {
						Type:     types.NumberType,
						Required: true,
					},
					"image": {
						Type:     types.StringType,
						Required: true,
					},
					"coffee_quantity": {
						Type:     types.NumberType,
						Required: true,
					},
				}),
			},
		},
	}, nil
}

func (r dataSourceOrderType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, []*tfprotov6.Diagnostic) {
	return dataSourceOrder{
		p: *(p.(*provider)),
	}, nil
}

type dataSourceOrder struct {
	p provider
}

func (r dataSourceOrder) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	fmt.Fprintln(stderr, "[DEBUG]-read-error3:", req.Config.Raw)

	var state Order
	err := req.Config.Get(ctx, &state)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error reading coffee5",
			Detail:   "An unexpected error was encountered while reading the datasource_coffee: " + err.Error(),
		})
		return
	}

	t := strconv.Itoa(state.ID)

	orderID := t

	r.p.client.GetCoffeeIngredients(orderID)
}
