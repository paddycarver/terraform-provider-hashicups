package hashicups

import (
	"context"
	"fmt"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

type dataSourceCoffeesType struct{}

func (r dataSourceCoffeesType) GetSchema(_ context.Context) (schema.Schema, []*tfprotov6.Diagnostic) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"coffee": {
				Computed: true,
				Attributes: schema.SingleNestedAttributes(map[string]schema.Attribute{
					"id": {
						Computed: true,
						Type:     types.NumberType,
					},
				}),
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

type dataSourceCoffeesData struct {
	ID          types.Number                `tfsdk:"id"`
	Name        types.String                `tfsdk:"name"`
	Teaser      types.String                `tfsdk:"teaser"`
	Description types.String                `tfsdk:"description"`
	Price       types.Number                `tfsdk:"price"`
	Image       types.String                `tfsdk:"image"`
	Ingredient  []dataSourceIngredientsData `tfsdk:"ingredients"`
}

func (r dataSourceCoffees) Read(ctx context.Context, p tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	fmt.Fprintln(stderr, "[DEBUG] Got state in provider:", p.Config.Raw)
	var order dataSourceCoffeesData
	err := p.Config.Get(ctx, order)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error reading state",
			Detail:   "An unexpected error was encountered while reading the state: " + err.Error(),
		})
		return
	}

	coffees, err := r.p.client.GetCoffees()
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error deleting order",
			Detail:   "Could not + " + err.Error(),
		})
	}

	var items []hashicups.Coffee

	for _, item := range coffees {
		coffees = append(items, hashicups.Coffee{
			Name:        item.Name,
			Teaser:      item.Teaser,
			Description: item.Description,
			ID:          item.ID,
			Price:       item.Price,
			Image:       item.Image,
		})
		return
	}
}
