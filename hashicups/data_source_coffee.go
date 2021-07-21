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
					// "ingredients": {
					// 	Computed: true,
					// 	Attributes: schema.ListNestedAttributes(map[string]schema.Attribute{
					// 		"ingredient_id": {
					// 			Computed: true,
					// 			Type:     types.NumberType,
					// 		},
					// 	}, schema.ListNestedAttributesOptions{}),
					// },
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

func (r dataSourceCoffees) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	fmt.Fprintf(stderr, "[DEBUG]-read-error1:%+v", req.Config.Raw)

	var state struct {
		OrderItem []Coffee `tfsdk:"coffees"`
		//	Coffee    []Ingredient `tfsdk:"ingredient_id"`
	}

	state.OrderItem = make([]Coffee, 0)
	//	state.Coffee = make([]Ingredient, 0)

	i, acc := r.p.client.GetCoffees()

	if acc != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error reading coffee2",
		})
	}

	for _, c := range i {
		state.OrderItem = append(state.OrderItem, Coffee{
			Name:        c.Name,
			Teaser:      c.Teaser,
			Description: c.Description,
			Price:       c.Price,
			Image:       c.Image,
			ID:          c.ID,
		})
	}
	// err := resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("coffees"), i)
	state.OrderItem = i
	err := resp.State.Set(ctx, state)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error reading coffee1",
			Detail:   "An unexpected error was encountered while reading the datasource_coffee: " + fmt.Sprint(i) + err.Error(),
		})
		return
	}
}
