package client

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	hyperkit "github.com/moby/hyperkit/go"
	providerconfigv1 "github.com/praveenkumar/cluster-api-provider-hyperkit/pkg/apis/hyperkitproviderconfig/v1alpha1"
	"k8s.io/client-go/kubernetes"
	"github.com/pborman/uuid"
)

//go:generate mockgen -source=./client.go -destination=./mock/client_generated.go -package=mock

// CreateVMonHyperkit specifies input parameters for CreateVM operation
type CreateVMonHyperkit struct {
	QCow2DiskImage string
	CPU            int
	Memory         int
	Cmdline        string
	UUID           string
	BootKernel     string
	BootInitrd     string
	// Ignition configuration to be injected during bootstrapping
	Ignition *providerconfigv1.Ignition
	// KubeClient as kubernetes client
	KubeClient kubernetes.Interface
}

// HyperkitClientBuilderFuncType is function type for building aws client
type HyperkitClientBuilderFuncType func(binaryLocation string, statedir string) (Client, error)

// Client is a wrapper object for actual hyperkit library to allow for easier testing.
type Client interface {
	// CreateDomain creates VM based on CreateVMonHyperkit
	CreateVM(CreateVMonHyperkit) error

	// DeleteVM deletes a VM
	DeleteVM(name string) error

	// VMExists checks if VM exists
	VMExists(name string) bool
}

type hyperkitClient struct {
	connection *hyperkit.HyperKit
}

var _ Client = &hyperkitClient{}

// Client libvirt, generate libvirt client given URI
func NewClient(binaryLocation string, statedir string) (Client, error) {
	connection, err := hyperkit.New(binaryLocation, "", statedir)
	if err != nil {
		return nil, err
	}
	glog.Infof("Created hyperkit connection: %p", connection)
	return &hyperkitClient{
		connection: connection,
	}, nil
}


// CreateDomain creates domain based on CreateDomainInput
func (client *hyperkitClient) CreateVM(input CreateVMonHyperkit) error {
	client.connection.Kernel = input.BootKernel
	client.connection.Initrd = input.BootInitrd
	client.connection.VMNet  = true
	client.connection.Console = hyperkit.ConsoleFile
	client.connection.CPUs = input.CPU
	client.connection.Memory = input.Memory


	uid := uuid.NewUUID().String()
	glog.Infof("Using UUID %s", uid)
	mac, err := GetMACAddressFromUUID(uid)
	if err != nil {
		return err
	}

	// Need to strip 0's
	mac = trimMacAddress(mac)
	glog.Infof("Generated MAC %s", mac)
	glog.Infof("Starting with cmdline: %s", input.Cmdline)
	if _, err := client.connection.Start(input.Cmdline); err != nil {
		return err
	}

	getIP := func() error {
		_, err := GetIPAddressByMACAddress(mac)
		if err != nil {
			return &RetriableError{Err: err}
		}
		return nil
	}

	if err := RetryAfter(30, getIP, 2*time.Second); err != nil {
		return fmt.Errorf("IP address never found in dhcp leases file %v", err)
	}

	return nil
}

// DomainExists checks if domain exists
func (client *hyperkitClient) VMExists(name string) bool {
	return client.connection.IsRunning()
}

// DeleteDomain deletes a domain
func (client *hyperkitClient) DeleteVM(name string) error {
	return client.connection.Remove(true)
}
