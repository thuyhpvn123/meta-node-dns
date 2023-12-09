package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	"github.com/meta-node-blockchain/meta-node/pkg/storage"
	"github.com/spf13/cobra"
)

var (
	dnsStoragePath    string
	address           string
	connectionAddress string
	verbose           bool
)

var rootCmd = &cobra.Command{
	Use:   "dns-cli",
	Short: "Metanode DNS CLI application",
	Long:  `This is a Metanode DNS CLI application. It has two subcommands: create and list.`,
}

func init() {
	// Create the create subcommand
	createCmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create dns record",
		Long:    `This subcommand creates dns record.`,
		Run: func(cmd *cobra.Command, args []string) {
			dnsStorage, err := storage.NewLevelDB(dnsStoragePath)
			if err != nil {
				panic(err)
			}
			dnsStorage.Put(common.FromHex(address), []byte(connectionAddress))
			logger.Info("Success create dns record for " + address + ":" + connectionAddress)
		},
	}

	// Add the --name flag to the create subcommand
	createCmd.PersistentFlags().StringVarP(&dnsStoragePath, "storage-path", "s", "", "The path of dns storage")
	createCmd.PersistentFlags().StringVarP(&address, "address", "a", "", "The address of record")
	createCmd.PersistentFlags().StringVarP(&connectionAddress, "connection-address", "c", "", "The connection address of record")

	// Add the create subcommand to the root command
	rootCmd.AddCommand(createCmd)

	// Create the list subcommand
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all record",
		Long:  `This subcommand lists all record.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	listCmd.PersistentFlags().StringVarP(&dnsStoragePath, "storage-path", "s", "", "The path of dns storage")

	// Add the list subcommand to the root command
	rootCmd.AddCommand(listCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
