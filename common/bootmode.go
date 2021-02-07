package common

import (
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
)

// BootMode is the boot mode of the system
// +kubebuilder:validation:Enum=UEFI;legacy
type BootMode string

// Allowed boot mode
const (
	UEFI            BootMode = "UEFI"
	SecureUEFI      BootMode = "UEFI secure boot"
	Legacy          BootMode = "legacy"
	DefaultBootMode BootMode = UEFI
)

func (bm BootMode) IronicBootMode() string {
	switch bm {
	case Legacy:
		return "bios"
	default:
		return "uefi"
	}
}

func (bm BootMode) MakeCapabilities() string {
	switch bm {
	case SecureUEFI:
		return "boot_mode:uefi,secure_boot:true"
	default:
		return fmt.Sprintf("boot_mode:%s", bm.IronicBootMode())
	}
}

func (bm BootMode) MakeInstanceUpdate(node *nodes.Node) (updates nodes.UpdateOpts) {
	value := map[string]string{
		"boot_mode": bm.IronicBootMode(),
	}
	if bm == SecureUEFI {
		value["secure_boot"] = "true"
	}

	updates = append(updates, nodes.UpdateOperation{
		Op:    nodes.AddOp,
		Path:  "/instance_info/capabilities",
		Value: value,
	})
	return
}

func (bm BootMode) ReplaceInCapabilities(capabilities string) string {
	if capabilities == "" {
		// The existing value is empty so we can replace the whole
		// thing.
		return bm.MakeCapabilities()
	}

	var filteredCapabilities []string
	for _, item := range strings.Split(capabilities, ",") {
		if !strings.HasPrefix(item, "boot_mode:") && !strings.HasPrefix(item, "secure_boot:") {
			filteredCapabilities = append(filteredCapabilities, item)
		}
	}
	filteredCapabilities = append(filteredCapabilities, bm.MakeCapabilities())

	return strings.Join(filteredCapabilities, ",")

}
