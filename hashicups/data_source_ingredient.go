package hashicups

import (
	"context"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type dataSourceIngredientsType struct {
	p provider
}

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
	return dataSourceCoffeesType{
		p: *(p.(*provider)),
	}, nil
}

var dataSourceIngredientsSchema = &tfprotov6.Schema{
	Block: &tfprotov6.SchemaBlock{
		Attributes: []*tfprotov6.SchemaAttribute{
			{
				Name:     "id",
				Type:     tftypes.String,
				Computed: true,
			},
			{
				Name:     "name",
				Computed: true,
				Type:     tftypes.String,
			},
			{
				Name:     "quantity",
				Computed: true,
				Type:     tftypes.String,
			},
			{
				Name:     "unit",
				Computed: true,
				Type:     tftypes.String,
			},
		},
	},
}

var dataSourceIngredientsTypeCoffeesType = tftypes.Object{
	AttributeTypes: map[string]tftypes.Type{
		"id":       tftypes.String,
		"name":     tftypes.String,
		"quantity": tftypes.String,
		"unit":     tftypes.String,
	},
}

type dataSourceServeIngredients struct {
	provider *resourceCoffeeData
}

func (r dataSourceIngredientsType) Read(ctx context.Context, p tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	r.p.client.GetCoffeeIngredients("1", []hashicups.Ingredient())
}
