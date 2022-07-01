package main

import (
	"fmt"
	"github.com/NextronSystems/ransomware-simulator/lib/encrypt"
	"github.com/NextronSystems/ransomware-simulator/lib/note"
	"github.com/NextronSystems/ransomware-simulator/lib/persistance"
	"github.com/NextronSystems/ransomware-simulator/lib/shadowcopy"
	"github.com/NextronSystems/ransomware-simulator/lib/simulatemacro"
	"github.com/secDre4mer/go-parseflags"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

func init() {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run Ransomware Simulator",
		Run:   run,
	}
	runCmd.Flags().AddFlagSet(parseflags.CreateFlagset(&runOptions))
	rootCmd.AddCommand(runCmd)
}

var runOptions = struct {
	DisableMacroSimulation bool `flag:"disable-macro-simulation" description:"Don't simulate start from a macro by building the following process chain: winword.exe -> cmd.exe -> ransomware-simulator.exe"`

	DisableShadowCopyDeletion bool `flag:"disable-shadow-copy-deletion" description:"Don't simulate volume shadow copy deletion"`

	EnableAddRegistryAutoRun bool `flag:"enable-registry-autorun" description:"Add a value to the registry startup and after reboot execute ransomware functional"`

	DisableFileEncryption bool   `flag:"disable-file-encryption" description:"Don't simulate document encryption"`
	EncryptionDirectory   string `flag:"dir" description:"Directory where files that will be encrypted should be staged"`

	DisableNoteDrop bool   `flag:"disable-note-drop" description:"Don't drop pseudo ransomware note"`
	NoteLocation    string `flag:"note-location" description:"Ransomware note location"`
}{
	EncryptionDirectory: `./encrypted-files`,
	NoteLocation:        filepath.Join(homeDir, "Desktop", "ransomware-simulator-note.txt"),
}

var homeDir, _ = os.UserHomeDir()

func run(cmd *cobra.Command, args []string) {
	if runOptions.EnableAddRegistryAutoRun {
		if err := persistance.RegistryAutoRun(); err != nil {
			log.Fatal(err)
		}
	}
	// Simulate Macro execution of this executable with current parameters, including --disable-macro-simulation
	if !runOptions.DisableMacroSimulation && !runOptions.EnableAddRegistryAutoRun {
		if err := simulatemacro.Run(append(os.Args, "--disable-macro-simulation")); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(asciiArt)
	if !runOptions.DisableShadowCopyDeletion && !runOptions.EnableAddRegistryAutoRun {
		if err := shadowcopy.Delete(); err != nil {
			log.Fatal(err)
		}
	}
	if !runOptions.DisableFileEncryption && !runOptions.EnableAddRegistryAutoRun {
		if err := encrypt.StageFiles(runOptions.EncryptionDirectory); err != nil {
			log.Fatal(err)
		}
		if err := encrypt.EncryptFiles(runOptions.EncryptionDirectory); err != nil {
			log.Fatal(err)
		}
	}
	if !runOptions.DisableNoteDrop && !runOptions.EnableAddRegistryAutoRun {
		if err := note.Write(runOptions.NoteLocation); err != nil {
			log.Fatal(err)
		}
	}
}
