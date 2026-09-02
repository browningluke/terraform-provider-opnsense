package dnsmasq_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/dnsmasq"
	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDomainDataSource(t *testing.T) {
	acctest.AccPreCheck(t)
	client := acctest.Client(t)

	const (
		testSequence    = "123"
		testDomain      = "test.domain.tld"
		testSourceIp    = "192.168.69.2"
		testPort        = "21781"
		testIp          = "192.168.69.3"
		testDescription = "terraform-provider-opnsense acceptance test"
	)

	domain := &dnsmasq.Domain{
		Sequence:    testSequence,
		Domain:      testDomain,
		SourceIp:    testSourceIp,
		Port:        testPort,
		Ip:          testIp,
		Description: testDescription,
	}

	id := acctest.CreateDataSourceTestResource(t, acctest.ResourceLifecycle[*dnsmasq.Domain]{
		Create: client.Dnsmasq().AddDomain,
		Delete: client.Dnsmasq().DeleteDomain,
	}, domain)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainDataSourceConfig(id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_domain.test", "id", id),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_domain.test", "sequence", testSequence),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_domain.test", "domain", testDomain),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_domain.test", "srcip", testSourceIp),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_domain.test", "port", testPort),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_domain.test", "ip", testIp),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_domain.test", "description", testDescription),
				),
			},
		},
	})
}

func testAccDomainDataSourceConfig(id string) string {
	return fmt.Sprintf(`
data "opnsense_dnsmasq_domain" "test" {
  id = "%s"
}
`, id)
}
