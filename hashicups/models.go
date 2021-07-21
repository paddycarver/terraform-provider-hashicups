package hashicups

import "github.com/hashicorp/terraform-plugin-framework/types"

// Order -
type Order struct {
	ID    int         `tfsdk:"id"`
	Items []OrderItem `tfsdk:"items"`
}

// OrderItem -
type OrderItem struct {
	Coffee   Coffee `tfsdk:"coffee"`
	Quantity int    `tfsdk:"quantity"`
}

// Coffee -
type Coffee struct {
	ID            int          `tfsdk:"orderid"`
	Name          string       `tfsdk:"name"`
	Teaser        string       `tfsdk:"teaser"`
	Description   string       `tfsdk:"description"`
	Price         float64      `tfsdk:"price"`
	Image         string       `tfsdk:"image"`
	IngredientIDs types.String `tfsdk:"ingredient_id"`
}

// Ingredient -
type Ingredient struct {
	ID       int    `tfsdk:"ingredient_id"`
	Name     string `tfsdk:"name"`
	Quantity int    `tfsdk:"quantity"`
	Unit     string `tfsdk:"unit"`
}
