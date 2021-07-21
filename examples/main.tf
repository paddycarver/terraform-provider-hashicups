terraform {
  required_providers {
    hashicups = {
      version = "0.3"
      source = "hashicorp.com/edu/hashicups"
    }
  }
}

provider "hashicups" {
  username = "rachel"
  password = "test123"
  host = "http://localhost:19090"
}
data "hashicups_coffees" "all" {}

# Returns all coffees
output "all_coffees" {
  value = data.hashicups_coffees.all.coffees
}
