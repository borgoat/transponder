terraform {
  backend "http" {
    address = "http://localhost:1492/terraform/example"
    lock_address = "http://localhost:1492/terraform/example"
    unlock_address = "http://localhost:1492/terraform/example"
  }

  required_version = "~> 0.12.0"
}

resource "random_pet" "pet" {
  count = 100
}