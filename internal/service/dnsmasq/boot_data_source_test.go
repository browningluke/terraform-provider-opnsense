package dnsmasq_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/dnsmasq"
	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBootDataSource(t *testing.T) {
	acctest.AccPreCheck(t)
	client := acctest.Client(t)

	const (
		testInterface     = "wan"
		testFileName      = "terraform-acc-test-pxelinux.0"
		testServerName    = "terraform-acc-test.example"
		testServerAddress = "192.0.2.10"
		testDescription   = "terraform-provider-opnsense acceptance test"
	)

	boot := &dnsmasq.Boot{
		Interface:     api.SelectedMap(testInterface),
		Filename:      testFileName,
		Servername:    testServerName,
		ServerAddress: testServerAddress,
		Description:   testDescription,
	}

	id := acctest.CreateDataSourceTestResource(t, acctest.ResourceLifecycle[*dnsmasq.Boot]{
		Create: client.Dnsmasq().AddBoot,
		Delete: client.Dnsmasq().DeleteBoot,
	}, boot)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBootDataSourceConfig(id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_boot.test", "id", id),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_boot.test", "interface", testInterface),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_boot.test", "file_name", testFileName),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_boot.test", "server_name", testServerName),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_boot.test", "server_address", testServerAddress),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_boot.test", "description", testDescription),
				),
			},
		},
	})
}

func testAccBootDataSourceConfig(id string) string {
	return fmt.Sprintf(`
data "opnsense_dnsmasq_boot" "test" {
  id = "%s"
}
`, id)
}
