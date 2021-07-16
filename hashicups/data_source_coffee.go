package hashicups

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type dataSourceCoffeesType struct {
	p provider
}

func (r dataSourceCoffeesType) GetSchema(_ context.Context) (schema.Schema, []*tfprotov6.Diagnostic) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"coffee": {
				Required: true,
				Attributes: schema.SingleNestedAttributes(map[string]schema.Attribute{
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
				}),
			},
		},
	}, nil
}

func (r dataSourceCoffeesType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, []*tfprotov6.Diagnostic) {
	return dataSourceCoffeesType{
		p: *(p.(*provider)),
	}, nil
}

var dataSourceCoffeesSchema = &tfprotov6.Schema{
	Block: &tfprotov6.SchemaBlock{
		Attributes: []*tfprotov6.SchemaAttribute{

			{
				Name:     "orderid",
				Type:     tftypes.String,
				Computed: true,
			},
			{
				Name:     "name",
				Computed: true,
				Type:     tftypes.String,
			},
			{
				Name:     "teaser",
				Computed: true,
				Type:     tftypes.String,
			},
			{
				Name:     "description",
				Computed: true,
				Type:     tftypes.String,
			},
			{
				Name:     "price",
				Computed: true,
				Type:     tftypes.String,
			},
			{
				Name:     "image",
				Computed: true,
				Type:     tftypes.String,
			},
		},
	},
}

var dataSourceCoffeesTypeCoffeesType = tftypes.Object{
	AttributeTypes: map[string]tftypes.Type{
		"name":        tftypes.String,
		"orderid":     tftypes.String,
		"teaser":      tftypes.String,
		"description": tftypes.String,
		"image":       tftypes.String,
		"price":       tftypes.String,
	},
}

type dataSourceServeCoffee struct {
	provider *resourceCoffeeData
}

func (r dataSourceCoffeesType) Read(ctx context.Context, p tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	r.p.client.GetCoffees()
}
