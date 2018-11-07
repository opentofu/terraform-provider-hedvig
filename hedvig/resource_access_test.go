package hedvig

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestAccHedvigAccess(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccHedvigAccessConfig,
				Check:  resource.ComposeTestCheckFunc(testAccCheckHedvigAccessExists("hedvig_access.test-access1"), testAccCheckHedvigAccessExists("hedvig_access.test-access2")),
			},
		},
	})
}

var testAccHedvigAccessConfig = fmt.Sprintf(`
provider "hedvig" {
  node = "%s"
  username = "%s"
  password = "%s"
}

resource "hedvig_vdisk" "test-access-vdisk1" {
  cluster = "%s"
  name = "%s"
  size = 9
  type = "BLOCK"
}

resource "hedvig_vdisk" "test-access-vdisk2" {
  cluster = "%s"
  name = "%s"
  size = 14
  type = "NFS"
}

resource "hedvig_lun" "test-access-lun" {
  cluster = "%s"
  vdisk = "${hedvig_vdisk.test-access-vdisk1.name}"
  controller = "%s"
}

resource "hedvig_mount" "test-access-mount" {
  cluster = "%s"
  vdisk = "${hedvig_vdisk.test-access-vdisk2.name}"
  controller = "%s"
}

resource "hedvig_access" "test-access1" {
  cluster = "%s"
  vdisk = "${hedvig_vdisk.test-access-vdisk1.name}"
  host = "${hedvig_lun.test-access-lun.controller}"
  address = "%s"
  type = "host"
}

resource "hedvig_access" "test-access2" {
  cluster = "%s"
  vdisk = "${hedvig_vdisk.test-access-vdisk2.name}"
  host = "${hedvig_mount.test-access-mount.controller}"
  address = "%s"
  type = "host"
}
`, os.Getenv("HV_TESTNODE"), os.Getenv("HV_TESTUSER"), os.Getenv("HV_TESTPASS"),
	os.Getenv("HV_TESTCLUST"), genRandomVdiskName(),
	os.Getenv("HV_TESTCLUST"), genRandomVdiskName(),
	os.Getenv("HV_TESTCLUST"), os.Getenv("HV_TESTCONT"),
	os.Getenv("HV_TESTCLUST"), os.Getenv("HV_TESTCONT"),
	os.Getenv("HV_TESTCLUST"), os.Getenv("HV_TESTADDR"),
	os.Getenv("HV_TESTCLUST"), os.Getenv("HV_TESTADDR2"))

func genRandomVdiskName() string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, 10)
	for i := 0; i < 10; i++ {
		bytes[i] = byte(65 + rand.Intn(25))
	}
	return fmt.Sprintf("HV-Test-%s", string(bytes))
}

func testAccCheckHedvigAccessExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No lun ID is set")
		}

		return nil
	}
}

func testAccCheckHedvigAccessCheckDestroyed(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return nil
		}
		return errors.New("Access Resource not Destroyed")
	}
}
