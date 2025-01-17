package cmd

import (
	"log"

	"github.com/Appelfeldt/osfmonitor/internal/osfm"

	"github.com/spf13/cobra"
)

var BuildVersion string

var rootCmd = &cobra.Command{
	Use:     "osfmonitor",
	Version: BuildVersion,
	Short:   "osfmonitor - View incoming OpenSeeFace tracking data",
	Long:    "osfmonitor is a tool for viewing OpenSeeFace tracking data received over the network",
	Args:    cobra.MaximumNArgs(0),
	Run:     command,
}

func Execute() {
	cobra.MousetrapHelpText = "" //Disables cobra functionality that blocks the program from being run through the GUI (explorer.exe) on Microsoft Windows.
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.PersistentFlags().Uint16P("port", "p", 11573, "Listening port")
}

func command(cmd *cobra.Command, args []string) {
	port, err := cmd.Flags().GetUint16("port")
	if err != nil {
		log.Fatalf("invalid port value\n%v", err)
	}

	settings := osfm.Settings{
		Port: port,
	}

	osfm.Run(settings)
}
