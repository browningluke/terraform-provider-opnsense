package dnsmasq_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/dnsmasq"
	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRangeDataSource(t *testing.T) {
	acctest.AccPreCheck(t)
	client := acctest.Client(t)

	const (
		testStartAddress = "192.168.100.100"
		testDescription  = "terraform-provider-opnsense range data source acceptance test"
	)

	mode := []string{
		"static",
	}

	rng := &dnsmasq.Range{
		StartAddress: testStartAddress,
		DomainType:   api.SelectedMap("range"),
		Mode:         api.SelectedMapList(mode),
		Description:  testDescription,
	}

	id := acctest.CreateDataSourceTestResource(t, acctest.ResourceLifecycle[*dnsmasq.Range]{
		Create: client.Dnsmasq().AddRange,
		Delete: client.Dnsmasq().DeleteRange,
	}, rng)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccRangeDataSourceConfig(id),

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_range.test", "id", id),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_range.test", "start_address", testStartAddress),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_range.test", "end_address", ""),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_range.test", "subnet_mask", ""),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_range.test", "domain_type", "range"),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_range.test", "description", testDescription),

					// SelectedMapList -> Terraform Set
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_range.test", "mode.#", "1"),
					resource.TestCheckTypeSetElemAttr("data.opnsense_dnsmasq_range.test", "mode.*", "static"),
				),
			},
		},
	})
}

func TestAccRange6DataSource(t *testing.T) {
	acctest.AccPreCheck(t)
	client := acctest.Client(t)

	const (
		testInterface    = "wan"
		testStartAddress = "::"
		testDescription  = "terraform-provider-opnsense IPv6 range data source acceptance test"
	)

	raMode := []string{
		"ra-stateless",
		"ra-names",
	}

	rng := &dnsmasq.Range{
		Interface:    api.SelectedMap(testInterface),
		StartAddress: testStartAddress,
		Constructor:  api.SelectedMap(testInterface),
		DomainType:   api.SelectedMap("interface"),
		RaMode:       api.SelectedMapList(raMode),
		Description:  testDescription,
	}

	id := acctest.CreateDataSourceTestResource(t, acctest.ResourceLifecycle[*dnsmasq.Range]{
		Create: client.Dnsmasq().AddRange,
		Delete: client.Dnsmasq().DeleteRange,
	}, rng)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccRangeDataSourceConfig(id),

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_range.test", "id", id),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_range.test", "start_address", testStartAddress),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_range.test", "interface", testInterface),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_range.test", "domain_type", "interface"),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_range.test", "description", testDescription),

					// SelectedMapList -> Terraform Set
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_range.test", "ra_mode.#", "2"),
					resource.TestCheckTypeSetElemAttr("data.opnsense_dnsmasq_range.test", "ra_mode.*", "ra-stateless"),
					resource.TestCheckTypeSetElemAttr("data.opnsense_dnsmasq_range.test", "ra_mode.*", "ra-names"),
				),
			},
		},
	})
}

func testAccRangeDataSourceConfig(id string) string {
	return fmt.Sprintf(`data "opnsense_dnsmasq_range" "test" {
  id = "%s"
}`, id)
}
