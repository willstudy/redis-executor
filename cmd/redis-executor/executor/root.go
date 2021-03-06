package executor

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const defaultDebugLevel uint32 = 4

var debugLevel uint32

var rootCmd = &cobra.Command{
	Use:           "redis-executor",
	Short:         "redis-executor provides a way to setup Redis Cluster.",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Please use -h to see usage")
	},
}

func init() {
	rootCmd.PersistentFlags().Uint32VarP(&debugLevel, "debuglevel", "l", defaultDebugLevel,
		"log debug level: 0[panic] 1[fatal] 2[error] 3[warn] 4[info] 5[debug]")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
