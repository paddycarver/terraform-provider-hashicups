package hashicups

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

type dataSourceIngredientsType struct{}

func (r dataSourceIngredientsType) GetSchema(_ context.Context) (schema.Schema, []*tfprotov6.Diagnostic) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"coffee_id": {
				Type:     types.NumberType,
				Required: true,
			},
			"ingredients": {
				Required: true,
				Attributes: schema.SingleNestedAttributes(map[string]schema.Attribute{
					"id": {
						Type:     types.NumberType,
						Computed: true,
					},
					"name": {
						Type:     types.StringType,
						Computed: true,
					},
					"quantity": {
						Type:     types.StringType,
						Computed: true,
					},
					"unit": {
						Type:     types.StringType,
						Computed: true,
					},
				}),
			},
		},
	}, nil
}

func (r dataSourceIngredientsType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, []*tfprotov6.Diagnostic) {
	return dataSourceIngredients{
		p: *(p.(*provider)),
	}, nil
}

type dataSourceIngredients struct {
	p provider
}

type dataSourceItemData struct {
	ID       int          `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Quantity int          `tfsdk:"quantity"`
	Unit     types.String `tfsdk:"unit"`
}

type dataSourceIngredientsData struct {
	Ingredients []dataSourceItemData `tfsdk:"ingredients"`
	CoffeeID    types.Number         `tfsdk:"coffee_id"`
}

func (r dataSourceIngredients) Read(ctx context.Context, p tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	fmt.Fprintln(stderr, "[DEBUG] Got state in provider:", p.Config.Raw)
	var order dataSourceIngredientsData
	err := p.Config.Get(ctx, order)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error reading state",
			Detail:   "An unexpected error was encountered while reading the state: " + err.Error(),
		})
		return
	}
	orderID, acc := order.CoffeeID.Value.Int64()
	if acc != big.Exact {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Invalid Order ID",
			Detail:   "OrderID must be an integer, cannot be a float.",
		})
		return
	}

	ings, err := r.p.client.GetCoffeeIngredients(strconv.FormatInt(orderID, 10))
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error deleting order",
			Detail:   "Could not delete orderID " + strconv.FormatInt(orderID, 10) + ": " + err.Error(),
		})
	}

	order.Ingredients = []dataSourceItemData{}
	for _, items := range ings {
		order.Ingredients = append(order.Ingredients, dataSourceItemData{
			Name:     types.String{Value: items.Name},
			Quantity: items.Quantity,
			Unit:     types.String{Value: items.Unit},
			ID:       items.ID,
		})
		return
	}
}
