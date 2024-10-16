package cmd

import (
	"github.com/spf13/cobra"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"go.uber.org/zap"
)

var Run func() error

var rootCmd = &cobra.Command{
	Use:   "uttarawave-backend",
	Short: "Uttarawave Backend Application",
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the Uttarawave backend server",
	Long:  `This command starts the Uttarawave backend server and all its components.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := Run(); err != nil {
			logger.Fatal("Failed to run the application", zap.Error(err))
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
