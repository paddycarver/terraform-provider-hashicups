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
					"coffee_id": {
						Type:     types.NumberType,
						Required: true,
					},
					"coffee_name": {
						Required: true,
						Type:     types.StringType,
					},
					"coffee_teaser": {
						Type:     types.StringType,
						Required: true,
					},
					"coffee_description": {
						Type:     types.StringType,
						Required: true,
					},
					"coffee_price": {
						Type:     types.NumberType,
						Required: true,
					},
					"coffee_image": {
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
type Order struct {
	ID    types.Number `tfsdk:"id"`
	Items []OrderItem  `tfsdk:"items"`
}

type OrderItem struct {
	Coffee   Coffee `tfsdk:"coffee"`
	Quantity int    `tfsdk:"quantity"`
}

type Coffee struct {
	ID          int          `tfsdk:"id"`
	Name        string       `tfsdk:"name"`
	Teaser      string       `tfsdk:"teaser"`
	Description string       `tfsdk:"description"`
	Price       float64      `tfsdk:"price"`
	Image       string       `tfsdk:"image"`
	Ingredient  []Ingredient `tfsdk:"ingredients"`
}

// Ingredient -
type Ingredient struct {
	ID       int    `tfsdk:"id"`
	Name     string `tfsdk:"name"`
	Quantity int    `tfsdk:"quantity"`
	Unit     string `tfsdk:"unit"`
}

// func (r dataSourceOrder) Read(ctx context.Context, p tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
// 	fmt.Fprintln(stderr, "[DEBUG] Got state in provider:", p.Config.Raw)
// 	var order Order
// 	err := p.Config.Get(ctx, order)
// 	if err != nil {
// 		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
// 			Severity: tfprotov6.DiagnosticSeverityError,
// 			Summary:  "Error reading state",
// 			Detail:   "An unexpected error was encountered while reading the state: " + err.Error(),
// 		})
// 		return
// 	}
// 	orderID, acc := order.ID.Value.Int64()
// 	if acc != big.Exact {
// 		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
// 			Severity: tfprotov6.DiagnosticSeverityError,
// 			Summary:  "Invalid Order ID",
// 			Detail:   "OrderID must be an integer, cannot be a float.",
// 		})
// 		return
// 	}
// 	order, err := r.p.client.GetOrder(*orderID)
// 	config.Items = []Order{}

// 	for _, item := range order.Items {
// 		order.Items = append(order.ID, Order{
// 			Coffee: Coffee{
// 				Name:        types.String{Value: item.Coffee.Name},
// 				Teaser:      types.String{Value: item.Coffee.Teaser},
// 				Description: types.String{Value: item.Coffee.Description},
// 				Price:       item.Coffee.Price,
// 				Image:       types.String{Value: item.Coffee.Image},
// 				ID:          item.Coffee.ID,
// 			},
// 			Quantity: item.Quantity,
// 		})
// 	}

// 	if err != nil {
// 		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
// 			Severity: tfprotov6.DiagnosticSeverityError,
// 			Summary:  "Error setting state",
// 			Detail:   "Unexpected error encountered trying to set new state: " + err.Error(),
// 		})
// 		return
// 	}
// }

func (r dataSourceOrder) Read(ctx context.Context, p tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	fmt.Fprintln(stderr, "[DEBUG] Got state in provider:", p.Config.Raw)
	var order Order
	err := p.Config.Get(ctx, order)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error reading state",
			Detail:   "An unexpected error was encountered while reading the state: " + err.Error(),
		})
		return
	}
	orderID, acc := order.ID.Value.Int64()
	if acc != big.Exact {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Invalid Order ID",
			Detail:   "OrderID must be an integer, cannot be a float.",
		})
		return
	}

	r.p.client.GetOrder(strconv.FormatInt(orderID, 10))
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error deleting order",
			Detail:   "Could not delete orderID " + strconv.FormatInt(orderID, 10) + ": " + err.Error(),
		})
	}
}
