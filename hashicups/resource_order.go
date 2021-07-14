package hashicups

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

type resourceOrderType struct{}

//Create the schema for the resource - what attributes are expected of a resource & what does it look like? Nested attributes feel weird
func (r resourceOrderType) GetSchema(_ context.Context) (schema.Schema, []*tfprotov6.Diagnostic) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"orderID": {
				Type:     types.StringType,
				Computed: true,
			},
			"last_updated": {
				Type: types.StringType,
				// provider will set value, user cannot specify
				Computed: true,
			},
			"items": {
				//tf will throw error if user doesn't specify palue - optional - can or choose not to supply a value
				Required: false,
				Attributes: schema.ListNestedAttributes(map[string]schema.Attribute{
					"coffee": {
						Required: true,
						Attributes: schema.ListNestedAttributes(map[string]schema.Attribute{
							"id": {
								Type:     types.NumberType,
								Required: true,
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
					"quantity": {
						Type:     types.NumberType,
						Required: true,
					},
				}, schema.ListNestedAttributesOptions{}),
			},
		},
	}, nil
}

// new resource instance
func (r resourceOrderType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, []*tfprotov6.Diagnostic) {
	return resourceOrder{
		p: *(p.(*provider)),
	}, nil
}

type resourceOrder struct {
	p provider
}

type resourceCoffeeData struct {
	ID          int          `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Teaser      types.String `tfsdk:"teaser"`
	Description types.String `tfsdk:"description"`
	Price       float64      `tfsdk:"price"`
	Image       types.String `tfsdk:"image"`
}

type resourceItemData struct {
	coffee   resourceCoffeeData
	quantity int `tfsdk:"quantity"`
}

type resourceOrderData struct {
	items        []resourceItemData `tfsdk:"items"`
	last_updated types.String       `tfsdk:"last_updated"`
	orderID      int                `tfsdk:"orderID"`
}

//create a new resource
func (r resourceOrder) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Provider not configured",
			Detail:   "The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		})
		return
	}
	var ticket resourceOrderData
	err := req.Plan.Get(ctx, &ticket)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error reading plan",
			Detail:   "An unexpected error was encountered while reading the plan: " + err.Error(),
		})
		return
	}

	var items []hashicups.OrderItem

	for _, item := range ticket.items {
		items = append(items, hashicups.OrderItem{
			Coffee: hashicups.Coffee{
				ID: item.coffee.ID,
			},
			Quantity: item.quantity,
		})
	}
	order, err := r.p.client.CreateOrder(items)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error creating order",
			Detail:   "Could not create order, unexpected error: " + err.Error(),
		})
		return
	}
	ticket.orderID = order.ID
	now := time.Now().Format(time.RFC850)
	ticket.last_updated = types.String{Value: string(now)}
	for _, planItem := range ticket.items {
		for _, item := range order.Items {
			if item.Coffee.ID == planItem.coffee.ID {
				planItem.coffee.Name = types.String{Value: item.Coffee.Name}
				planItem.coffee.Teaser = types.String{Value: item.Coffee.Teaser}
				planItem.coffee.Description = types.String{Value: item.Coffee.Description}
				planItem.coffee.Price = item.Coffee.Price
				planItem.coffee.Image = types.String{Value: item.Coffee.Image}
			}
		}
	}
	err = resp.State.Set(ctx, order)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error setting state",
			Detail:   "Could not set state, unexpected error: " + err.Error(),
		})
		return
	}
}

//Read
func (r resourceOrder) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	fmt.Fprintln(stderr, "[DEBUG] Got state in provider:", req.State.Raw)
	var state resourceOrderData
	err := req.State.Get(ctx, &state)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error reading state",
			Detail:   "An unexpected error was encountered while reading the state: " + err.Error(),
		})
		return
	}
	// get order from API and then update what is in state from what the API returns

	//Set on state var state resourceOrderData will hold what the API returns
	order, err := r.p.client.GetOrder(strconv.Itoa(state.orderID))
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error creating order",
			Detail:   "Could not create order, unexpected error: " + err.Error(),
		})
		return
	}

	state.items = []resourceItemData{}
	for _, item := range order.Items {
		state.items = append(state.items, resourceItemData{
			coffee: resourceCoffeeData{
				Name:        types.String{Value: item.Coffee.Name},
				Teaser:      types.String{Value: item.Coffee.Teaser},
				Description: types.String{Value: item.Coffee.Description},
				Price:       item.Coffee.Price,
				Image:       types.String{Value: item.Coffee.Image},
				ID:          item.Coffee.ID,
			},
			quantity: item.Quantity,
		})
	}

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error setting state",
			Detail:   "Unexpected error encountered trying to set new state: " + err.Error(),
		})
		return
	}
}

//update
func (r resourceOrder) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var plan resourceOrderData
	err := req.Plan.Get(ctx, &plan)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error reading plan",
			Detail:   "An unexpected error was encountered while reading the plan: " + err.Error(),
		})
		return
	}
	var state resourceOrderData
	err = req.State.Get(ctx, &state)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error reading prior state",
			Detail:   "An unexpected error was encountered while reading the prior state: " + err.Error(),
		})
		return
	}

	var items []hashicups.OrderItem

	for _, item := range plan.items {
		items = append(items, hashicups.OrderItem{
			Coffee: hashicups.Coffee{
				ID: item.coffee.ID,
			},
			Quantity: item.quantity,
		})
	}
	order, err := r.p.client.UpdateOrder(strconv.Itoa(state.orderID), []hashicups.OrderItem{})
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error updating order",
			Detail:   "Could not update order, unexpected error: " + err.Error(),
		})
		return
	}
	state.items = []resourceItemData{}
	for _, item := range order.Items {
		state.items = append(state.items, resourceItemData{
			coffee: resourceCoffeeData{
				Name:        types.String{Value: item.Coffee.Name},
				Teaser:      types.String{Value: item.Coffee.Teaser},
				Description: types.String{Value: item.Coffee.Description},
				Price:       item.Coffee.Price,
				Image:       types.String{Value: item.Coffee.Image},
				ID:          item.Coffee.ID,
			},
			quantity: item.Quantity,
		})
	}
	err = resp.State.Set(ctx, order)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error setting state",
			Detail:   "Could not set state, unexpected error: " + err.Error(),
		})
		return
	}
}

//Delete

func (r resourceOrder) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state resourceOrderData
	err := req.State.Get(ctx, &state)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error reading configuration",
			Detail:   "An unexpected error was encountered while reading the configuration: " + err.Error(),
		})
		return
	}
	// original framework test provider created a file on the file system and needed to destroy an on disk
	// Would delete in hashicups be removing the item from the state and API?
	//call hashicups API for DeleteOrder
	err = r.p.client.DeleteOrder(strconv.Itoa(state.orderID))
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error deleting order",
			Detail:   "Could not delete orderID " + strconv.Itoa(state.orderID) + ": " + err.Error(),
		})
		return
	}
	resp.State.RemoveResource(ctx)
}
