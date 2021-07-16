package hashicups

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

type dataSourceOrderType struct{
	p provider
}

func (r dataSourceOrderType) GetSchema(_ context.Context) (schema.Schema, []*tfprotov6.Diagnostic) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": {
				Type:     types.NumberType,
				Computed: true,
			},
			"items": {
				//tf will throw error if user doesn't specify palue - optional - can or choose not to supply a value
				Required: false,
				Attributes: schema.ListNestedAttributes(map[string]schema.Attribute{
					"coffee_id": {
						Type:     types.NumberType,
						Required: true,
					},
					"coffee_name": {
						Required: true,
						Type: types.StringType,
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

func (r dataSourceOrderType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSourceType, []*tfprotov6.Diagnostic) {
	return dataSourceOrderType{
		p: *(p.(*provider)),
}

var dataSourceOrdersSchema = &tfprotov6.Schema{
	Block: &tfprotov6.SchemaBlock{
		Attributes: []*tfprotov6.SchemaAttribute{

			{
				Name:     "id",
				Type:     tftypes.String,
				Computed: true,
			},
			{
				Name:     "items",
				Computed: true,
				Type:     tftypes.String,
			},
			{
				Name:     "coffee_id",
				Computed: true,
				Type:     tftypes.String,
			},
			{
				Name:     "coffee_name",
				Computed: true,
				Type:     tftypes.String,
			},
			{
				Name:     "coffee_price",
				Computed: true,
				Type:     tftypes.String,
			},
			{
				Name:     "coffee_image",
				Computed: true,
				Type:     tftypes.String,
			},
		},
	},
}

var dataSourceOrdersTypeCoffeesType = tftypes.Object{
	AttributeTypes: map[string]tftypes.Type{
		"name":        tftypes.String,
		"orderid":     tftypes.String,
		"teaser":      tftypes.String,
		"description": tftypes.String,
		"image":       tftypes.String,
		"price":       tftypes.String,
	},
}

type dataSourceServeOrder struct {
	provider *dataSourceOrdersTypeCoffeesType

}


func (r dataSourceOrderType) Read(ctx context.Context, p tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	r.p.client.GetOrder(orderID)
}