package restore

import (
	"atlas-cli-plugin/internal/utils"
	"bytes"
	"os"
	"os/exec"

	"log"

	"github.com/mongodb/mongo-tools/common/signals"
	"github.com/mongodb/mongo-tools/mongorestore"
	"github.com/spf13/cobra"
)

func restoreArchive(connectionString string, archive string) {

	// I am not proud of the next few lines of code, but this seems to be the only way to
	// instantiate the required mongorestore options properly
	args := []string{}
	args = append(args, "--uri="+connectionString)
	args = append(args, "--archive="+archive)
	opts, _ := mongorestore.ParseOptions(args, "", "")
	opts.NormalizeOptionsAndURI()

	restore, err := mongorestore.New(opts)

	if err != nil {
		log.Fatalf("Error creating mongorestore: %v", err)
	}

	defer restore.Close()

	finishedChan := signals.HandleWithInterrupt(restore.HandleInterrupt)
	defer close(finishedChan)

	result := restore.Restore()
	if result.Err != nil {
		log.Fatalf("Failed: %v", result.Err)
	}
}

func Builder() *cobra.Command {
	restoreCmd := &cobra.Command{
		Use:   "restore",
		Short: "Restores archived dump into a running MongoDB cluster",
		RunE: func(cmd *cobra.Command, _ []string) error {
			// deploymentName is the first argument
			deploymentName := cmd.Flags().Arg(0)
			archive, _ := cmd.Flags().GetString("archive")

			atlasCliExe := utils.AtlasCliExe()

			log.Printf("Using Atlas CLI executable: %s\n", atlasCliExe)
			log.Printf("Restoring archive %s to deployment: %s\n", archive, deploymentName)

			atlasCmd := exec.Command(atlasCliExe, "deployments", "connect", deploymentName, "--connectWith", "connectionString")
			atlasCmd.Env = os.Environ()
			var stdout bytes.Buffer
			atlasCmd.Stdout = &stdout

			if err := atlasCmd.Run(); err != nil {
				log.Fatalf("Error running command: %v", err)
			}

			connectionString := stdout.String()
			connectionString = connectionString[:len(connectionString)-1]

			log.Printf("Connection String: %s\n", connectionString)

			restoreArchive(connectionString, archive)

			return nil
		},
	}

	restoreCmd.Flags().StringP("archive", "a", "", "Path to the archive to restore")

	return restoreCmd
}
