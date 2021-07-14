package main

import (
	"context"

	"github.com/tr0njavolta/terraform-provider-hashicups/hashicups"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func main() {
	tfsdk.Serve(context.Background(), hashicups.New, tfsdk.ServeOpts{
		Name: "hashicups",
	})
}
