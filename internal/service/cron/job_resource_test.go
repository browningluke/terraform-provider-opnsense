package cron_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCronJobResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCronJobResourceConfig("firmware poll", "0", "4"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_cron_job.test", "command", "firmware poll"),
					resource.TestCheckResourceAttr("opnsense_cron_job.test", "minutes", "0"),
					resource.TestCheckResourceAttr("opnsense_cron_job.test", "hours", "4"),
					resource.TestCheckResourceAttr("opnsense_cron_job.test", "days", "*"),
					resource.TestCheckResourceAttr("opnsense_cron_job.test", "months", "*"),
					resource.TestCheckResourceAttr("opnsense_cron_job.test", "weekdays", "*"),
					resource.TestCheckResourceAttr("opnsense_cron_job.test", "who", "root"),
					resource.TestCheckResourceAttr("opnsense_cron_job.test", "enabled", "true"),
					resource.TestCheckResourceAttr("opnsense_cron_job.test", "description", "acc test job"),
					resource.TestCheckResourceAttrSet("opnsense_cron_job.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_cron_job.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCronJobResourceConfig("firmware poll", "0", "3"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_cron_job.test", "hours", "3"),
				),
			},
		},
	})
}

func TestAccCronJobResource_Disabled(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCronJobResourceConfigDisabled("firmware poll", "0", "4"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_cron_job.test", "enabled", "false"),
					resource.TestCheckResourceAttrSet("opnsense_cron_job.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_cron_job.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCronJobResource_WithDescription(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCronJobResourceConfigWithDescription("firmware poll", "0", "4", "acc test job"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_cron_job.test", "description", "acc test job"),
					resource.TestCheckResourceAttrSet("opnsense_cron_job.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_cron_job.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCronJobResourceConfig(command, minutes, hours string) string {
	return fmt.Sprintf(`
resource "opnsense_cron_job" "test" {
  command     = %[1]q
  minutes     = %[2]q
  hours       = %[3]q
  description = "acc test job"
}
`, command, minutes, hours)
}

func testAccCronJobResourceConfigDisabled(command, minutes, hours string) string {
	return fmt.Sprintf(`
resource "opnsense_cron_job" "test" {
  enabled     = false
  command     = %[1]q
  minutes     = %[2]q
  hours       = %[3]q
  description = "acc test job disabled"
}
`, command, minutes, hours)
}

func testAccCronJobResourceConfigWithDescription(command, minutes, hours, description string) string {
	return fmt.Sprintf(`
resource "opnsense_cron_job" "test" {
  command     = %[1]q
  minutes     = %[2]q
  hours       = %[3]q
  description = %[4]q
}
`, command, minutes, hours, description)
}
