terraform {
  required_providers {
    opnsense = {
      version = "~> x.0"
      source  = "browningluke/opnsense"
    }
  }
}

provider "opnsense" {
  uri = "https://opnsense.example.com"
  api_key = "..."
  api_secret = "..."
}
