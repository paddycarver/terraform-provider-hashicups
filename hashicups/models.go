package hashicups

import "github.com/hashicorp/terraform-plugin-framework/types"

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
	ID       types.Number `tfsdk:"id"`
	Name     string       `tfsdk:"name"`
	Quantity int          `tfsdk:"quantity"`
	Unit     string       `tfsdk:"unit"`
}
