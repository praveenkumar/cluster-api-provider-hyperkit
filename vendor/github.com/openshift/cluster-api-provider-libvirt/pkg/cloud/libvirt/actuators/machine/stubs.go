package machine

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	providerconfigv1 "github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1alpha1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

const (
	defaultNamespace   = "default"
	userDataSecretName = "libvirt-actuator-user-data-secret"

	clusterID = "libvirt-actuator-cluster"
)

func stubProviderConfig() *providerconfigv1.LibvirtMachineProviderConfig {
	return &providerconfigv1.LibvirtMachineProviderConfig{
		DomainMemory: 2048,
		DomainVcpu:   1,
		CloudInit: &providerconfigv1.CloudInit{
			SSHAccess: true,
		},
		Volume: &providerconfigv1.Volume{
			PoolName:     "default",
			BaseVolumeID: "/var/lib/libvirt/images/fedora_base",
		},
		NetworkInterfaceName:    "default",
		NetworkInterfaceAddress: "192.168.124.12/24",
		Autostart:               false,
		URI:                     "http://localhost",
	}
}

func stubMachine() (*clusterv1.Machine, error) {
	machinePc := stubProviderConfig()

	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		return nil, fmt.Errorf("failed creating codec: %v", err)
	}
	config, err := codec.EncodeToProviderSpec(machinePc)
	if err != nil {
		return nil, fmt.Errorf("encodeToProviderConfig failed: %v", err)
	}

	machine := &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "libvirt-actuator-testing-machine",
			Namespace: defaultNamespace,
			Labels: map[string]string{
				providerconfigv1.ClusterIDLabel:   clusterID,
				providerconfigv1.MachineRoleLabel: "infra",
				providerconfigv1.MachineTypeLabel: "master",
			},
		},

		Spec: clusterv1.MachineSpec{
			ProviderSpec: *config,
		},
	}

	return machine, nil
}

func stubCluster() *clusterv1.Cluster {
	return &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID,
			Namespace: defaultNamespace,
		},
	}
}
