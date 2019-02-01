package main

// Tests individual Libvirt actuator actions

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/praveenkumar/cluster-api-provider-hyperkit/cmd/hyperkit-actuator/utils"

	flag "github.com/spf13/pflag"

	goflag "flag"

	"github.com/spf13/cobra"
)

const (
	pollInterval        = 5 * time.Second
	timeoutPoolInterval = 20 * time.Minute
)

var rootCmd = &cobra.Command{
	Use:   "libvirt-actuator-test",
	Short: "Test for Cluster API Libvirt actuator",
}

func createCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create machine instance for specified cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkFlags(cmd); err != nil {
				return err
			}
			cluster, machine, userData, err := utils.ReadClusterResources(
				cmd.Flag("cluster").Value.String(),
				cmd.Flag("machine").Value.String(),
				cmd.Flag("userdata").Value.String(),
			)
			if err != nil {
				return err
			}

			actuator := utils.CreateActuator(machine, userData)
			err = actuator.Create(context.TODO(), cluster, machine)
			if err != nil {
				return err
			}
			fmt.Printf("Machine creation was successful!\n")
			return nil
		},
	}
}

func deleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete INSTANCE-ID",
		Short: "Delete machine instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkFlags(cmd); err != nil {
				return err
			}

			cluster, machine, userData, err := utils.ReadClusterResources(
				cmd.Flag("cluster").Value.String(),
				cmd.Flag("machine").Value.String(),
				cmd.Flag("userdata").Value.String(),
			)
			if err != nil {
				return err
			}

			actuator := utils.CreateActuator(machine, userData)
			err = actuator.Delete(context.TODO(), cluster, machine)
			if err != nil {
				return err
			}
			fmt.Printf("Machine delete operation was successful.\n")
			return nil
		},
	}
}

func existsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "exists",
		Short: "Determine if underlying machine instance exists",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkFlags(cmd); err != nil {
				return err
			}

			cluster, machine, userData, err := utils.ReadClusterResources(
				cmd.Flag("cluster").Value.String(),
				cmd.Flag("machine").Value.String(),
				cmd.Flag("userdata").Value.String(),
			)
			if err != nil {
				return err
			}

			actuator := utils.CreateActuator(machine, userData)
			exists, err := actuator.Exists(context.TODO(), cluster, machine)
			if err != nil {
				return err
			}
			if exists {
				fmt.Printf("Underlying machine's instance exists.\n")
			} else {
				fmt.Printf("Underlying machine's instance not found.\n")
			}
			return nil
		},
	}
}

func bootstrapCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstrap kubernetes cluster with kubeadm",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	// To bootstrap the master guest that will run the cluster api stack
	cmd.PersistentFlags().StringP("libvirt-uri", "", "", "Libvirt URI. E.g. qemu//system")
	// libvirt URI for actuator running inside a container (as part of the cluster API stack)
	cmd.PersistentFlags().StringP("in-cluster-libvirt-uri", "", "", "Libvirt URI for docker container. E.g. qemu+ssh://root@IP/system")
	// ssh private key so the actuator running inside the container can talk to libvirt running in libvirt instance (in case qemu+ssh is used)
	cmd.PersistentFlags().StringP("libvirt-private-key", "", "", "Private key file for libvirt qemu+ssh URI")
	// ssh private key to pull kubeconfig from master guest
	cmd.PersistentFlags().StringP("master-guest-private-key", "", "", "Private key file of the master guest to pull kubeconfig")

	return cmd
}

func init() {
	rootCmd.PersistentFlags().StringP("machine", "m", "", "Machine manifest")
	rootCmd.PersistentFlags().StringP("cluster", "c", "", "Cluster manifest")
	rootCmd.PersistentFlags().StringP("userdata", "u", "", "User data manifest")

	cUser, err := user.Current()
	if err != nil {
		rootCmd.PersistentFlags().StringP("environment-id", "p", "", "Directory with bootstrapping manifests")
	} else {
		rootCmd.PersistentFlags().StringP("environment-id", "p", cUser.Username, "Machine prefix, by default set to the current user")
	}

	rootCmd.AddCommand(createCommand())
	rootCmd.AddCommand(deleteCommand())
	rootCmd.AddCommand(existsCommand())
	rootCmd.AddCommand(bootstrapCommand())

	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	// the following line exists to make glog happy, for more information, see: https://github.com/kubernetes/kubernetes/issues/17162
	flag.CommandLine.Parse([]string{})
}

func checkFlags(cmd *cobra.Command) error {
	if cmd.Flag("cluster").Value.String() == "" {
		return fmt.Errorf("--%v/-%v flag is required", cmd.Flag("cluster").Name, cmd.Flag("cluster").Shorthand)
	}
	if cmd.Flag("machine").Value.String() == "" {
		return fmt.Errorf("--%v/-%v flag is required", cmd.Flag("machine").Name, cmd.Flag("machine").Shorthand)
	}
	return nil
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred: %v\n", err)
		os.Exit(1)
	}
}
