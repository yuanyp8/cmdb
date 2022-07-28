package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yuanyp8/cmdb/version"
)

var flagVersion bool

var RootCmd = &cobra.Command{
	Use:   "cmdb-api",
	Short: "cmdb 项目后端API",
	Long:  "cmdb 项目后端API，基于Restful API实现资源的管理",
	RunE: func(cmd *cobra.Command, args []string) error {
		if flagVersion {
			fmt.Println(version.FullVersion())
			return nil
		}
		return nil
	},
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&flagVersion, "version", "v", false, "print cmdb-api version")
}
