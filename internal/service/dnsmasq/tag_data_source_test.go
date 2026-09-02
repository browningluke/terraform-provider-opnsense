package dnsmasq_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/dnsmasq"
	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTagDataSource(t *testing.T) {
	acctest.AccPreCheck(t)
	client := acctest.Client(t)

	const (
		testTag = "terraform_test"
	)

	tag := &dnsmasq.Tag{
		Tag: testTag,
	}

	id := acctest.CreateDataSourceTestResource(t, acctest.ResourceLifecycle[*dnsmasq.Tag]{
		Create: client.Dnsmasq().AddTag,
		Delete: client.Dnsmasq().DeleteTag,
	}, tag)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccTagDataSourceConfig(id),

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_tag.test", "id", id),
					resource.TestCheckResourceAttr("data.opnsense_dnsmasq_tag.test", "tag", testTag),
				),
			},
		},
	})
}

func testAccTagDataSourceConfig(id string) string {
	return fmt.Sprintf(`data "opnsense_dnsmasq_tag" "test" {
  id = "%s"
}`, id)
}
