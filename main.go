package main

import (
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/tr0njavolta/terraform-provider-hashicups/hashicups"
)

func main() {
	tfsdk.Serve(context.Background(), provider.New, tfsdk.ServeOpts{
		Name: "hashicups",
	})
}
