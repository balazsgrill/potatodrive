package main

import (
	"log"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

func addWritePermissionToCurrentUser(keyPath string) error {
	// Open the registry key with WRITE_DAC access to modify its ACL
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer key.Close()

	// Get the current user's SID
	token := windows.Token(0) // Use the current process token
	userSID, err := token.GetTokenUser()
	if err != nil {
		return err
	}

	// Create an explicit access entry for the current user
	explicitAccess := windows.EXPLICIT_ACCESS{
		AccessPermissions: windows.KEY_WRITE,
		AccessMode:        windows.GRANT_ACCESS,
		Inheritance:       windows.NO_INHERITANCE,
		Trustee: windows.TRUSTEE{
			TrusteeForm:  windows.TRUSTEE_IS_SID,
			TrusteeType:  windows.TRUSTEE_IS_USER,
			TrusteeValue: windows.TrusteeValue(unsafe.Pointer(userSID.User.Sid)),
		},
	}

	// Get the current ACL of the registry key
	var oldACL *windows.ACL
	sd, err := windows.GetNamedSecurityInfo(
		keyPath,
		windows.SE_REGISTRY_KEY,
		windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		return err
	}
	oldACL, _, err = sd.DACL()
	if err != nil {
		return err
	}

	// Create a new ACL with the added permission
	newACL, err := windows.ACLFromEntries([]windows.EXPLICIT_ACCESS{explicitAccess}, oldACL)
	if err != nil {
		return err
	}

	// Set the new ACL to the registry key
	err = windows.SetNamedSecurityInfo(
		keyPath,
		windows.SE_REGISTRY_KEY,
		windows.DACL_SECURITY_INFORMATION,
		nil,
		nil,
		newACL,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	keyPath := `SOFTWARE\MyApp`

	err := addWritePermissionToCurrentUser(keyPath)
	if err != nil {
		log.Fatalf("Failed to add write permission: %v", err)
	}

	log.Println("Write permission granted to the current user.")
}
