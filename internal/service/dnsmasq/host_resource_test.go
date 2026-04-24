package dnsmasq_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDnsmasqHostResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDnsmasqHostResourceConfig("testhost", "192.168.1.100", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "hostname", "testhost"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "ip_addresses.#", "1"),
					resource.TestCheckTypeSetElemAttr("opnsense_dnsmasq_host.test", "ip_addresses.*", "192.168.1.100"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "is_local_domain", "false"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "is_ignored", "false"),
					resource.TestCheckResourceAttrSet("opnsense_dnsmasq_host.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_dnsmasq_host.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDnsmasqHostResourceConfig("testhostupdated", "192.168.1.200", "Updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "hostname", "testhostupdated"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "ip_addresses.#", "1"),
					resource.TestCheckTypeSetElemAttr("opnsense_dnsmasq_host.test", "ip_addresses.*", "192.168.1.200"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "description", "Updated description"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDnsmasqHostResource_MultipleIPs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsmasqHostResourceConfigMultipleIPs("multiiphost"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "hostname", "multiiphost"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "ip_addresses.#", "3"),
					resource.TestCheckTypeSetElemAttr("opnsense_dnsmasq_host.test", "ip_addresses.*", "192.168.1.100"),
					resource.TestCheckTypeSetElemAttr("opnsense_dnsmasq_host.test", "ip_addresses.*", "192.168.1.101"),
					resource.TestCheckTypeSetElemAttr("opnsense_dnsmasq_host.test", "ip_addresses.*", "192.168.1.102"),
				),
			},
		},
	})
}

func TestAccDnsmasqHostResource_WithOptionalFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read with all optional fields
			{
				Config: testAccDnsmasqHostResourceConfigOptional(
					"optfieldhost",
					"test.local",
					true,
					false,
					"192.168.1.100",
					[]string{"alias.test.local"},
					[]string{"alternate.test.local"},
					[]string{"aa:bb:cc:dd:ee:ff"},
					"00:03:00:01:aa:bb:cc:dd:ee:ff",
					"Test description",
					"Test comment",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "hostname", "optfieldhost"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "domain", "test.local"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "is_local_domain", "true"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "is_ignored", "false"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "ip_addresses.#", "1"),
					resource.TestCheckTypeSetElemAttr("opnsense_dnsmasq_host.test", "ip_addresses.*", "192.168.1.100"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "alias_records.#", "1"),
					resource.TestCheckTypeSetElemAttr("opnsense_dnsmasq_host.test", "alias_records.*", "alias.test.local"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "cname_records.#", "1"),
					resource.TestCheckTypeSetElemAttr("opnsense_dnsmasq_host.test", "cname_records.*", "alternate.test.local"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "hardware_addresses.#", "1"),
					resource.TestCheckTypeSetElemAttr("opnsense_dnsmasq_host.test", "hardware_addresses.*", "aa:bb:cc:dd:ee:ff"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "client_id", "00:03:00:01:aa:bb:cc:dd:ee:ff"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "description", "Test description"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "comment", "Test comment"),
					resource.TestCheckResourceAttrSet("opnsense_dnsmasq_host.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "opnsense_dnsmasq_host.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDnsmasqHostResourceConfigOptional(
					"optfieldhost",
					"test.local",
					false,
					true,
					"192.168.1.100",
					[]string{"alias.test.local"},
					[]string{"alternate.test.local", "alternate2.test.local"},
					[]string{"aa:bb:cc:dd:ee:ff", "11:22:33:44:55:66"},
					"00:03:00:01:aa:bb:cc:dd:ee:ff",
					"Updated description",
					"Updated comment",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "is_local_domain", "false"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "is_ignored", "true"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "cname_records.#", "2"),
					resource.TestCheckTypeSetElemAttr("opnsense_dnsmasq_host.test", "cname_records.*", "alternate.test.local"),
					resource.TestCheckTypeSetElemAttr("opnsense_dnsmasq_host.test", "cname_records.*", "alternate2.test.local"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "hardware_addresses.#", "2"),
					resource.TestCheckTypeSetElemAttr("opnsense_dnsmasq_host.test", "hardware_addresses.*", "aa:bb:cc:dd:ee:ff"),
					resource.TestCheckTypeSetElemAttr("opnsense_dnsmasq_host.test", "hardware_addresses.*", "11:22:33:44:55:66"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("opnsense_dnsmasq_host.test", "comment", "Updated comment"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDnsmasqHostResourceConfig(hostname, ip, description string) string {
	descriptionLine := ""
	if description != "" {
		descriptionLine = fmt.Sprintf("  description  = %q\n", description)
	}
	return fmt.Sprintf(`
resource "opnsense_dnsmasq_host" "test" {
  hostname     = %[1]q
  ip_addresses = [%[2]q]
%[3]s}
`, hostname, ip, descriptionLine)
}

func testAccDnsmasqHostResourceConfigMultipleIPs(hostname string) string {
	return fmt.Sprintf(`
resource "opnsense_dnsmasq_host" "test" {
  hostname     = %[1]q
  ip_addresses = ["192.168.1.100", "192.168.1.101", "192.168.1.102"]
}
`, hostname)
}

func testAccDnsmasqHostResourceConfigOptional(
	hostname, domain string,
	isLocalDomain, isIgnored bool,
	ip string,
	aliasRecords, cnameRecords, hwAddresses []string,
	clientID, description, comment string,
) string {
	aliasRecordsLine := ""
	if len(aliasRecords) > 0 {
		aliasRecordsLine = fmt.Sprintf("  alias_records       = [\"%s\"]\n", strings.Join(aliasRecords, `", "`))
	}
	cnameRecordsLine := ""
	if len(cnameRecords) > 0 {
		cnameRecordsLine = fmt.Sprintf("  cname_records       = [\"%s\"]\n", strings.Join(cnameRecords, `", "`))
	}
	hwAddressesLine := ""
	if len(hwAddresses) > 0 {
		hwAddressesLine = fmt.Sprintf("  hardware_addresses  = [\"%s\"]\n", strings.Join(hwAddresses, `", "`))
	}
	return fmt.Sprintf(`
resource "opnsense_dnsmasq_host" "test" {
  hostname        = %[1]q
  domain          = %[2]q
  is_local_domain = %[3]t
  is_ignored      = %[4]t
  ip_addresses    = [%[5]q]
%[6]s%[7]s%[8]s  client_id       = %[9]q
  description     = %[10]q
  comment         = %[11]q
}
`, hostname, domain, isLocalDomain, isIgnored, ip, aliasRecordsLine, cnameRecordsLine, hwAddressesLine, clientID, description, comment)
}
