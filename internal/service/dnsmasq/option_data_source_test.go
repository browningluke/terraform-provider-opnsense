package dnsmasq_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/dnsmasq"
	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOptionDataSource(t *testing.T) {
	acctest.AccPreCheck(t)
	client := acctest.Client(t)

	const (
		testType        = "set"
		testOption      = "150"
		testValue       = "192.168.69.10"
		testForce       = "1"
		testDescription = "terraform-provider-opnsense acceptance test"
	)

	option := &dnsmasq.Option{
		Type:        api.SelectedMap(testType),
		OptionV4:    api.SelectedMap(testOption),
		Value:       testValue,
		Force:       testForce,
		Description: testDescription,
	}

	id := acctest.CreateDataSourceTestResource(t, acctest.ResourceLifecycle[*dnsmasq.Option]{
		Create: client.Dnsmasq().AddOption,
		Delete: client.Dnsmasq().DeleteOption,
	}, option)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOptionDataSourceConfig(id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_option.test", "id", id),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_option.test", "type", testType),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_option.test", "option", testOption),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_option.test", "value", testValue),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_option.test", "force", "true"),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_option.test", "description", testDescription),
				),
			},
		},
	})

}

func TestAccOption6DataSource(t *testing.T) {
	acctest.AccPreCheck(t)
	client := acctest.Client(t)

	const (
		testType        = "set"
		testOption6     = "23"
		testValue       = "[2001:db8::10]"
		testForce       = "1"
		testDescription = "terraform-provider-opnsense IPv6 acceptance test"
	)

	option := &dnsmasq.Option{
		Type:        api.SelectedMap(testType),
		OptionV6:    api.SelectedMap(testOption6),
		Value:       testValue,
		Force:       testForce,
		Description: testDescription,
	}

	id := acctest.CreateDataSourceTestResource(t, acctest.ResourceLifecycle[*dnsmasq.Option]{
		Create: client.Dnsmasq().AddOption,
		Delete: client.Dnsmasq().DeleteOption,
	}, option)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOptionDataSourceConfig(id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_option.test", "id", id),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_option.test", "type", testType),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_option.test", "option6", testOption6),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_option.test", "value", testValue),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_option.test", "force", "true"),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_option.test", "description", testDescription),
				),
			},
		},
	})

}

func testAccOptionDataSourceConfig(id string) string {
	return fmt.Sprintf(`data "opnsense_dnsmasq_option" "test" {
  id = "%s"
}`, id)
}
