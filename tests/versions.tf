terraform {
  required_version = ">= 1.6.0"
  required_providers {
    opnsense = {
      source  = "browningluke/opnsense"
      version = ">= 0.11.0"
    }
  }
}
