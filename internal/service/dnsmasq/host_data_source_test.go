package dnsmasq_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/dnsmasq"
	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccHostDataSource(t *testing.T) {
	acctest.AccPreCheck(t)
	client := acctest.Client(t)

	const (
		testHostname          = "test-host"
		testDomain            = "test.domain.tld"
		testIsLocalDomain     = "0"
		testClientID          = "00:03:00:01:aa:bb:cc:dd:ee:ff"
		testTag               = "test-tag"
		testIsIgnored         = "0"
		testDescription       = "terraform-provider-opnsense acceptance test"
		testComment           = "terraform-provider-opnsense acceptance test comment"
		testIpAddress         = "192.168.69.3"
		testAliasRecord       = "alias.test.domain.tld"
		secondTestAliasRecord = "second.test.domain.tld"
		testCnameRecord       = "cname.test.domain.tld"
		secondTestCnameRecord = "second.cname.domain.tld"
		testHardwareAddress   = "00:11:22:33:44:55"
	)

	host := &dnsmasq.Host{
		Hostname:          testHostname,
		Domain:            testDomain,
		IsLocalDomain:     testIsLocalDomain,
		IpAddresses:       []string{testIpAddress},
		AliasRecords:      []string{testAliasRecord, secondTestAliasRecord},
		CnameRecords:      []string{testCnameRecord, secondTestCnameRecord},
		ClientId:          testClientID,
		HardwareAddresses: []string{testHardwareAddress},
		IsIgnored:         testIsIgnored,
		Description:       testDescription,
		Comments:          testComment,
	}

	id := acctest.CreateDataSourceTestResource(t, acctest.ResourceLifecycle[*dnsmasq.Host]{
		Create: client.Dnsmasq().AddHost,
		Delete: client.Dnsmasq().DeleteHost,
	}, host)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccHostDataSourceConfig(id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_host.test", "id", id),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_host.test", "hostname", testHostname),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_host.test", "domain", testDomain),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_host.test", "is_local_domain", "false"),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_host.test", "ip_addresses.#", "1"),
					resource.TestCheckTypeSetElemAttr("data.opnsense_dnsmasq_host.test", "ip_addresses.*", testIpAddress),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_host.test", "alias_records.#", "2"),
					resource.TestCheckTypeSetElemAttr("data.opnsense_dnsmasq_host.test", "alias_records.*", testAliasRecord),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_host.test", "cname_records.#", "2"),
					resource.TestCheckTypeSetElemAttr("data.opnsense_dnsmasq_host.test", "cname_records.*", testCnameRecord),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_host.test", "client_id", testClientID),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_host.test", "hardware_addresses.#", "1"),
					resource.TestCheckTypeSetElemAttr("data.opnsense_dnsmasq_host.test", "hardware_addresses.*", testHardwareAddress),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_host.test", "is_ignored", "false"),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_host.test", "description", testDescription),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_host.test", "comment", testComment),
				),
			},
		},
	})
}

func testAccHostDataSourceConfig(id string) string {
	return fmt.Sprintf(`
data "opnsense_dnsmasq_host" "test" {
  id = "%s"
}
`, id)
}
