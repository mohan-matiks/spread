package cmd

import (
	"os"

	"github.com/SwishHQ/spread/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "spread",
	Short: "OTA update tool for React Native apps",
	Long: `
	                                            888 
                                                888 
                                                888 
.d8888b  88888b.  888d888 .d88b.   8888b.   .d88888 
88K      888 "88b 888P"  d8P  Y8b     "88b d88" 888 
"Y8888b. 888  888 888    88888888 .d888888 888  888 
     X88 888 d88P 888    Y8b.     888  888 Y88b 888 
 88888P' 88888P"  888     "Y8888  "Y888888  "Y88888 
         888                                        
         888                                        
         888                                        
	`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.L.Error("Failed to execute root command", zap.Error(err))
		os.Exit(1)
	}
}
