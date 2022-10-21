package cmd

import (
	"fmt"
	civoHandler "github.com/kubesimplify/Kubesimpctl/src/api/civo"
	"github.com/kubesimplify/Kubesimpctl/src/api/payload"
	"github.com/spf13/cobra"
)

var createClusterCivo = &cobra.Command{
	Use:   "civo",
	Short: "Use to create a CIVO k3s cluster",
	Long: `It is used to create cluster with the given name from user. For example:

kubesimpctl create-cluster civo <arguments to civo cloud provider>
`,
	Run: func(cmd *cobra.Command, args []string) {
		clusterConfig := civoHandler.ClusterInfoInjecter(
			cclusterName,
			cregion,
			cspec.Disk,
			cspec.Nodes,
			apps,
			cni)
		err := civoHandler.CreateCluster(clusterConfig)
		if err != nil {
			fmt.Printf("\033[31;40m%v\033[0m\n", err)
			return
		}
		fmt.Printf("\033[32;40mCREATED!\033[0m\n")
	},
}

var (
	cclusterName string
	cspec        payload.Machine
	cregion      string
	apps         string
	cni          string
)

func init() {
	createClusterCmd.AddCommand(createClusterCivo)
	createClusterCivo.Flags().StringVarP(&cclusterName, "name", "C", "demo", "Cluster name")
	createClusterCivo.Flags().StringVarP(&cspec.Disk, "nodeSize", "s", "500M", "Node size")
	createClusterCivo.Flags().StringVarP(&cregion, "region", "r", "", "Region")
	createClusterCivo.Flags().StringVarP(&apps, "apps", "a", "", "PreInstalled Apps with comma seperated string")
	createClusterCivo.Flags().StringVarP(&cni, "cni", "c", "", "CNI Plugin to be installed")
	createClusterCivo.Flags().IntVarP(&cspec.Nodes, "nodes", "n", 1, "Number of Nodes")
	createClusterCivo.MarkFlagRequired("name")
	createClusterCivo.MarkFlagRequired("nodes")
	createClusterCivo.MarkFlagRequired("region")
}
