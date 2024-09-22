package dump

import (
	"atlas-cli-plugin/internal/spinner"
	"atlas-cli-plugin/internal/utils"
	"bytes"
	"io"
	"net/url"
	"os"
	"os/exec"

	"log"

	"github.com/mongodb/mongo-tools/common/signals"
	"github.com/mongodb/mongo-tools/mongodump"
	"github.com/spf13/cobra"
)

func dumpToArchive(connectionString string, archive string, db string, debug bool) {

	// I am not proud of the next few lines of code, but this seems to be the only way to
	// instantiate the required mongoredump options properly
	args := []string{}
	args = append(args, "--uri="+connectionString)
	args = append(args, "--archive="+archive)
	if db != "" {
		args = append(args, "--db="+db)
	}
	if !debug {
		args = append(args, "--quiet")
	}
	opts, _ := mongodump.ParseOptions(args, "", "")
	opts.NormalizeOptionsAndURI()

	if !debug {
		spin := spinner.New("Dumping data to archive")
		defer spin.Stop()
	}

	dump := mongodump.MongoDump{
		ToolOptions:   opts.ToolOptions,
		OutputOptions: opts.OutputOptions,
		InputOptions:  opts.InputOptions,
	}

	finishedChan := signals.HandleWithInterrupt(dump.HandleInterrupt)
	defer close(finishedChan)

	if err := dump.Init(); err != nil {
		log.Fatalf("Failed: %v", err)
	}

	if err := dump.Dump(); err != nil {
		log.Fatalf("Failed: %v", err)
	}
}

func Builder() *cobra.Command {
	dumpCmd := &cobra.Command{
		Use:   "dump",
		Short: "Dumps a MongoDB deployment into an archive",
		RunE: func(cmd *cobra.Command, args []string) error {
			// deploymentName is the first argument
			deploymentName := cmd.Flags().Arg(0)
			archive, _ := cmd.Flags().GetString("archive")
			debug, _ := cmd.Flags().GetBool("debug")
			dbuser, _ := cmd.Flags().GetString("dbuser")
			dbpass, _ := cmd.Flags().GetString("dbpass")
			db, _ := cmd.Flags().GetString("db")

			if !debug {
				log.SetOutput(io.Discard)
			}

			atlasCliExe := utils.AtlasCliExe()

			log.Printf("Using Atlas CLI executable: %s\n", atlasCliExe)
			log.Printf("Dumping to archive %s deployment: %s with database %s \n", archive, deploymentName, db)

			atlasCmd := exec.Command(atlasCliExe, "deployments", "connect", deploymentName, "--connectWith", "connectionString")
			atlasCmd.Env = os.Environ()
			var stdout bytes.Buffer
			atlasCmd.Stdout = &stdout

			if err := atlasCmd.Run(); err != nil {
				log.Fatalf("Error running command: %v", err)
			}

			connectionString := stdout.String()
			connectionString = connectionString[:len(connectionString)-1]

			// convert connection string to a URI object so we can edit its different parts
			if dbuser != "" || dbpass != "" {
				uri, err := url.Parse(connectionString)
				if err != nil {
					log.Fatalf("Error parsing connection string: %v", err)
				}

				uri.User = url.UserPassword(dbuser, dbpass)
				connectionString = uri.String()
			}

			log.Printf("Connection String: %s\n", connectionString)

			dumpToArchive(connectionString, archive, db, debug)

			return nil
		},
	}

	dumpCmd.Flags().String("db", "", "Database to dump")
	dumpCmd.Flags().StringP("archive", "a", "", "Path to the archive where the dump will be stored")
	dumpCmd.MarkFlagRequired("archive")
	dumpCmd.Flags().String("dbuser", "", "Database user")
	dumpCmd.Flags().String("dbpass", "", "Database password")
	dumpCmd.Flags().Bool("debug", false, "Enable debug mode")

	return dumpCmd
}
