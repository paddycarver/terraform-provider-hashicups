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

type dataSourceIngredientsType struct{}

func (r dataSourceIngredientsType) GetSchema(_ context.Context) (schema.Schema, []*tfprotov6.Diagnostic) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"coffee_id": {
				Type:     types.NumberType,
				Required: true,
			},
			"ingredients": {
				Computed: true,
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

func (r dataSourceIngredients) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	fmt.Fprintln(stderr, "[DEBUG]-read-error2:", req.Config.Schema)
	var ing Ingredient
	err := req.Config.Get(ctx, &ing)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error reading ingredients",
			Detail:   "An unexpected error was encountered while reading the ingredients: " + err.Error(),
		})
		return
	}
	t := strconv.Itoa(ing.ID)

	ingID := t

	r.p.client.GetCoffeeIngredients(ingID)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error reading coffee2",
		})
	}
}
