package hashicups

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
