# This is a singleton resource. It must be imported before use:
# terraform import opnsense_firewall_geoip_settings.settings firewall_geoip_settings

resource "opnsense_firewall_geoip_settings" "settings" {
  url = "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country-CSV&license_key=<key>&suffix=zip"
}
