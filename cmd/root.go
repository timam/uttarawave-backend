package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"go.uber.org/zap"
)

var Run func() error
var Migrate func() error

var rootCmd = &cobra.Command{
	Use:   "uttarawave-backend",
	Short: "Uttarawave Backend Application",
	Long:  `Uttarawave Backend Application is a server application with various functionalities.`,
	Run: func(cmd *cobra.Command, args []string) {
		displayBanner()
		cmd.Usage()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(migrateCmd)
}

func Execute() error {
	return rootCmd.Execute()
}

func displayBanner() {
	banner := `
 _   _ _   _                                        
| | | | | | |                                       
| | | | |_| |_ __ _ _ __ __ ___      ____ ___   ___ 
| | | | __| __/ _` + "`" + ` | '__/ _` + "`" + ` \ \ /\ / / _` + "`" + ` \ \ / / _ \
| |_| | |_| || (_| | | | (_| |\ V  V / (_| |\ V /  __/
 \___/ \__|\__\__,_|_|  \__,_| \_/\_/ \__,_| \_/ \___|
                                                    
Uttarawave Backend Application
======================================================
Available commands:
  run     : Run the backend server
  migrate : Run database migrations 
======================================================
`
	fmt.Println(banner)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the Uttarawave backend server",
	Long:  `This command starts the Uttarawave backend server.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := Run(); err != nil {
			logger.Fatal("Failed to run the application", zap.Error(err))
		}
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migration for Uttarawave backend server",
	Long:  `This command migrates database of Uttarawave backend server.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := Migrate(); err != nil {
			logger.Fatal("Failed to run the application", zap.Error(err))
		}
	},
}
