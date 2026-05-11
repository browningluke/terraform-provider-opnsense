package interfaces_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInterfacesOverviewInterfaceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOverviewInterfaceDataSourceConfig("vtnet0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.opnsense_interfaces_overview.test", "device", "vtnet0"),
					resource.TestCheckResourceAttrSet("data.opnsense_interfaces_overview.test", "identifier"),
					resource.TestCheckResourceAttrSet("data.opnsense_interfaces_overview.test", "macaddr"),
					resource.TestCheckResourceAttrSet("data.opnsense_interfaces_overview.test", "mtu"),
					resource.TestCheckResourceAttr("data.opnsense_interfaces_overview.test", "is_physical", "true"),
					resource.TestCheckResourceAttr("data.opnsense_interfaces_overview.test", "enabled", "true"),
				),
			},
		},
	})
}

func TestAccInterfacesOverviewInterfaceAllDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOverviewInterfaceAllDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// At least one interface must be present
					resource.TestCheckResourceAttrSet("data.opnsense_interfaces_overview_all.test", "interfaces.#"),
					// device and mtu are guaranteed non-empty for every interface type,
					resource.TestCheckResourceAttrSet("data.opnsense_interfaces_overview_all.test", "interfaces.0.device"),
					resource.TestCheckResourceAttrSet("data.opnsense_interfaces_overview_all.test", "interfaces.0.mtu"),
				),
			},
		},
	})
}

func testAccOverviewInterfaceDataSourceConfig(device string) string {
	return fmt.Sprintf(`
data "opnsense_interfaces_overview" "test" {
  device = %[1]q
}
`, device)
}

func testAccOverviewInterfaceAllDataSourceConfig() string {
	return `
data "opnsense_interfaces_overview_all" "test" {}
`
}
