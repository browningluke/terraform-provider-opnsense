# The frontend itself is managed outside Terraform; only its rule list is set here.
# The web rule runs before the ACME client plugin's challenge rule. That rule stays
# in the list, otherwise the plugin appends it again and the next plan shows a diff.
resource "opnsense_haproxy_frontend_actions" "http" {
  frontend_id = "00000000-0000-4000-8000-000000000001" # UUID of the port-80 frontend
  action_ids = [
    opnsense_haproxy_action.web_acme.id,
    "00000000-0000-4000-8000-000000000002", # UUID of redirect_acme_challenges
  ]
}
