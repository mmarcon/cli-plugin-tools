package restore

import (
	"atlas-cli-plugin/internal/spinner"
	"atlas-cli-plugin/internal/utils"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"log"

	"github.com/mongodb/mongo-tools/common/signals"
	"github.com/mongodb/mongo-tools/mongorestore"
	"github.com/spf13/cobra"
)

func restoreArchive(connectionString string, archive string, debug bool) {

	// I am not proud of the next few lines of code, but this seems to be the only way to
	// instantiate the required mongorestore options properly
	args := []string{}
	args = append(args, "--uri="+connectionString)
	args = append(args, "--archive="+archive)
	if !debug {
		args = append(args, "--quiet")
	}
	opts, _ := mongorestore.ParseOptions(args, "", "")
	opts.NormalizeOptionsAndURI()

	spin := spinner.New("Restoring archive")
	defer spin.Stop()

	restore, err := mongorestore.New(opts)

	if err != nil {
		log.Fatalf("Error creating mongorestore: %v", err)
	}

	defer restore.Close()

	finishedChan := signals.HandleWithInterrupt(restore.HandleInterrupt)
	defer close(finishedChan)

	result := restore.Restore()
	spin.Stop()
	if result.Err != nil {
		log.Fatalf("Failed: %v", result.Err)
	}
}

func isUri(uri string) bool {
	parsedUrl, err := url.ParseRequestURI(uri)
	if err != nil {
		return false
	}

	// Check if the scheme is http or https
	return strings.HasPrefix(parsedUrl.Scheme, "http")
}

func downloadArchive(uri string) (string, error) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "*.tmp.archive")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	spin := spinner.New("Downloading archive")

	// Perform the HTTP GET request
	resp, err := http.Get(uri)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file: %s", resp.Status)
	}

	// Write the response body to the temporary file
	_, err = io.Copy(tempFile, resp.Body)
	spin.Stop()
	if err != nil {
		return "", err
	}
	// Return the path to the temporary file
	return tempFile.Name(), nil
}

func deleteFile(filePath string) error {
	return os.Remove(filePath)
}

func Builder() *cobra.Command {
	restoreCmd := &cobra.Command{
		Use:   "restore",
		Short: "Restores archived dump into a running MongoDB cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			// deploymentName is the first argument
			deploymentName := cmd.Flags().Arg(0)
			archive, _ := cmd.Flags().GetString("archive")
			debug, _ := cmd.Flags().GetBool("debug")

			if !debug {
				log.SetOutput(io.Discard)
			}

			atlasCliExe := utils.AtlasCliExe()

			log.Printf("Using Atlas CLI executable: %s\n", atlasCliExe)
			log.Printf("Restoring archive %s to deployment: %s\n", archive, deploymentName)

			if isUri(archive) {
				log.Printf("Downloading archive from %s\n", archive)
				downloadedArchive, err := downloadArchive(archive)
				if err != nil {
					log.Fatalf("Error downloading archive: %v", err)
				}
				defer deleteFile(downloadedArchive)
				archive = downloadedArchive
			}

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

			restoreArchive(connectionString, archive, debug)

			return nil
		},
	}

	restoreCmd.Flags().StringP("archive", "a", "", "Path to the archive to restore")
	restoreCmd.Flags().Bool("debug", false, "Enable debug mode")

	return restoreCmd
}
