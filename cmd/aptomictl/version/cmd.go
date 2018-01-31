package version

import (
	"fmt"
	"github.com/Aptomi/aptomi/cmd/common"
	"github.com/Aptomi/aptomi/pkg/client/rest"
	"github.com/Aptomi/aptomi/pkg/client/rest/http"
	"github.com/Aptomi/aptomi/pkg/config"
	"github.com/Aptomi/aptomi/pkg/version"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewCommand returns instance of cobra command that shows version from version package (injected at build tome)
func NewCommand(cfg *config.Client) *cobra.Command {
	var client, server bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the Aptomi Client (and Server) versions",
		Run: func(cmd *cobra.Command, args []string) {
			if !client && !server {
				client = true
				server = true
			}

			if client {
				info := version.GetBuildInfo()

				data, err := common.Format(cfg.Output, false, info)
				if err != nil {
					panic(fmt.Sprintf("Error while formating policy: %s", err))
				}

				log.Infof("Client: ")
				fmt.Println(string(data))
			}
			if server {
				info, err := rest.New(cfg, http.NewClient(cfg)).Version().Show()
				if err != nil {
					panic(fmt.Sprintf("Error while getting server version: %s", err))
				}

				data, err := common.Format(cfg.Output, false, info)
				if err != nil {
					panic(fmt.Sprintf("Error while formating server version: %s", err))
				}

				log.Infof("Server: ")
				fmt.Println(string(data))
			}
		},
	}

	cmd.Flags().BoolVarP(&client, "client", "c", false, "Show client version")
	cmd.Flags().BoolVarP(&server, "server", "s", false, "Show server version")

	return cmd
}
