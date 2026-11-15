# The web rule runs before the ACME client plugin's challenge rule. That rule stays
# in the list, otherwise the plugin appends it again and the next plan shows a diff.
# frontend_id can also be the UUID of a frontend managed outside Terraform.
resource "opnsense_haproxy_frontend_actions" "http" {
  frontend_id = opnsense_haproxy_frontend.http.id
  action_ids = [
    opnsense_haproxy_action.web_acme.id,
    "00000000-0000-4000-8000-000000000002", # UUID of redirect_acme_challenges
  ]
}
