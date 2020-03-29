package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/liatrio/terraform-provider-harbor/harbor"
)

func TestAccHarborRegistryBasic(t *testing.T) {
	registryName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckHarborRegistryDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testHarborRegistryBasic(registryName, "true"),
				Check:  testAccCheckHarborRegistryExists("harbor_registry.registry"),
			},
			//		{
			//			ResourceName:        "keycloak_group.group",
			//			ImportState:         true,
			//			ImportStateVerify:   true,
			//			ImportStateIdPrefix: realmName + "/",
			//		},
		},
	})
}

func TestAccHarborRegistryUpdate(t *testing.T) {
	registryName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckHarborRegistryDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testHarborRegistryBasic(registryName, "true"),
				Check:  testAccCheckHarborRegistryExists("harbor_registry.registry"),
			},
			{
				Config: testHarborRegistryBasic(registryName, "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHarborRegistryExists("harbor_registry.registry"),
					resource.TestCheckResourceAttr("harbor_registry.registry", "verify_remote_cert", "false"),
				),
			},
		},
	})
}

func testHarborRegistryBasic(registryName string, verify_cert string) string {
	return fmt.Sprintf(`
resource "harbor_registry" "registry" {
	name = "%s"
	type = "docker-hub"
	endpoint_url = "http://hub.docker.com"
	verify_remote_cert = %s
  }`, registryName, verify_cert)
}

func testAccCheckHarborRegistryExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getRegistryFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func getRegistryFromState(s *terraform.State, resourceName string) (*harbor.Registry, error) {
	client := testAccProvider.Meta().(*harbor.Client)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID

	registry, err := client.GetRegistry(id)
	if err != nil {
		return nil, fmt.Errorf("error getting group with id %s: %s", id, err)
	}

	return registry, nil
}

func testAccCheckHarborRegistryDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "harbor_registry" {
				continue
			}

			id := rs.Primary.ID

			client := testAccProvider.Meta().(*harbor.Client)

			registry, _ := client.GetRegistry(id)
			if registry != nil {
				return fmt.Errorf("registry with id %s still exists", id)
			}
		}

		return nil
	}
}
