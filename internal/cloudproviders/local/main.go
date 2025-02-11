package local

import (
	"encoding/json"
	"fmt"

	"github.com/kubesimplify/ksctl/pkg/utils"

	"github.com/kubesimplify/ksctl/pkg/resources"
	"github.com/kubesimplify/ksctl/pkg/resources/controllers/cloud"
	cloudControlRes "github.com/kubesimplify/ksctl/pkg/resources/controllers/cloud"
	. "github.com/kubesimplify/ksctl/pkg/utils/consts"
)

type StateConfiguration struct {
	ClusterName string `json:"cluster_name"`
	Distro      string `json:"distro"`
	Version     string `json:"version"`
	Nodes       int    `json:"nodes"`
}

type Metadata struct {
	ResName string
	Version string

	Cni string
}

type LocalProvider struct {
	ClusterName string `json:"cluster_name"`
	NoNodes     int    `json:"no_nodes"`
	Metadata
}

// GetSecretTokens implements resources.CloudFactory.
func (*LocalProvider) GetSecretTokens(resources.StorageFactory) (map[string][]byte, error) {
	return nil, nil
}

// GetStateFile implements resources.CloudFactory.
func (*LocalProvider) GetStateFile(resources.StorageFactory) (string, error) {
	cloudstate, err := json.Marshal(localState)
	if err != nil {
		return "", err
	}
	return string(cloudstate), nil
}

var (
	localState *StateConfiguration
)

const (
	STATE_FILE = "kind-state.json"
	KUBECONFIG = "kubeconfig"
)

func ReturnLocalStruct(metadata resources.Metadata) (*LocalProvider, error) {
	return &LocalProvider{
		ClusterName: metadata.ClusterName,
	}, nil
}

// InitState implements resources.CloudFactory.
func (cloud *LocalProvider) InitState(storage resources.StorageFactory, operation KsctlOperation) error {
	switch operation {
	case OperationStateCreate:
		if isPresent(storage, cloud.ClusterName) {
			return fmt.Errorf("[local] already present")
		}
		localState = &StateConfiguration{
			ClusterName: cloud.ClusterName,
			Distro:      "kind",
		}
		var err error
		err = storage.Path(utils.GetPath(UtilClusterPath, CloudLocal, ClusterTypeMang, cloud.ClusterName)).
			Permission(0750).CreateDir()
		if err != nil {
			return err
		}

		err = saveStateHelper(storage, utils.GetPath(UtilClusterPath, CloudLocal, ClusterTypeMang, cloud.ClusterName, STATE_FILE))
		if err != nil {
			return err
		}
	case OperationStateDelete, OperationStateGet:
		err := loadStateHelper(storage, utils.GetPath(UtilClusterPath, CloudLocal, ClusterTypeMang, cloud.ClusterName, STATE_FILE))
		if err != nil {
			return err
		}
	}
	storage.Logger().Success("[local] initialized the state")
	return nil
}

// it will contain the name of the resource to be created
func (cloud *LocalProvider) Name(resName string) resources.CloudFactory {
	cloud.Metadata.ResName = resName
	return cloud
}

func (cloud *LocalProvider) Application(s string) (externalApps bool) {
	return true
}

func (client *LocalProvider) CNI(s string) (externalCNI bool) {

	switch KsctlValidCNIPlugin(s) {
	case CNIKind, "":
		client.Metadata.Cni = string(CNIKind)
	default:
		client.Metadata.Cni = string(CNINone)
		return true
	}

	return false
}

// Version implements resources.CloudFactory.
func (cloud *LocalProvider) Version(ver string) resources.CloudFactory {
	// TODO: validation of version
	cloud.Metadata.Version = ver
	return cloud
}

func GetRAWClusterInfos(storage resources.StorageFactory) ([]cloudControlRes.AllClusterData, error) {
	var data []cloudControlRes.AllClusterData

	managedFolders, err := storage.Path(utils.GetPath(UtilClusterPath, CloudLocal, ClusterTypeMang)).GetFolders()
	if err != nil {
		return nil, err
	}

	for _, folder := range managedFolders {

		path := utils.GetPath(UtilClusterPath, CloudLocal, ClusterTypeMang, folder[0], STATE_FILE)
		raw, err := storage.Path(path).Load()
		if err != nil {
			return nil, err
		}
		var clusterState *StateConfiguration
		if err := json.Unmarshal(raw, &clusterState); err != nil {
			return nil, err
		}

		data = append(data,
			cloudControlRes.AllClusterData{
				Provider:   CloudLocal,
				Name:       folder[0],
				Region:     "N/A",
				Type:       ClusterTypeMang,
				K8sDistro:  KsctlKubernetes(clusterState.Distro),
				K8sVersion: clusterState.Version,
				NoMgt:      clusterState.Nodes,
			})
	}
	return data, nil
}

// //// NOT IMPLEMENTED //////

// it will contain whether the resource to be created belongs for controlplane component or loadbalancer...
func (cloud *LocalProvider) Role(KsctlRole) resources.CloudFactory {
	return nil
}

// it will contain which vmType to create
func (cloud *LocalProvider) VMType(string) resources.CloudFactory {
	return nil
}

// whether to have the resource as public or private (i.e. VMs)
func (cloud *LocalProvider) Visibility(bool) resources.CloudFactory {
	return nil
}

func (*LocalProvider) GetHostNameAllWorkerNode() []string {
	return nil
}

// CreateUploadSSHKeyPair implements resources.CloudFactory.
func (*LocalProvider) CreateUploadSSHKeyPair(state resources.StorageFactory) error {
	return nil

}

// DelFirewall implements resources.CloudFactory.
func (*LocalProvider) DelFirewall(state resources.StorageFactory) error {
	return nil
}

// DelNetwork implements resources.CloudFactory.
func (*LocalProvider) DelNetwork(state resources.StorageFactory) error {
	return nil
}

// DelSSHKeyPair implements resources.CloudFactory.
func (*LocalProvider) DelSSHKeyPair(state resources.StorageFactory) error {
	return nil
}

// DelVM implements resources.CloudFactory.
func (*LocalProvider) DelVM(resources.StorageFactory, int) error {
	return nil
}

// GetStateForHACluster implements resources.CloudFactory.
func (*LocalProvider) GetStateForHACluster(state resources.StorageFactory) (cloud.CloudResourceState, error) {
	return cloud.CloudResourceState{}, fmt.Errorf("[local] should not be implemented")
}

// NewFirewall implements resources.CloudFactory.
func (*LocalProvider) NewFirewall(state resources.StorageFactory) error {
	return nil
}

// NewNetwork implements resources.CloudFactory.
func (*LocalProvider) NewNetwork(state resources.StorageFactory) error {
	return nil
}

// NewVM implements resources.CloudFactory.
func (*LocalProvider) NewVM(resources.StorageFactory, int) error {
	return nil
}

// NoOfControlPlane implements resources.CloudFactory.
func (cloud *LocalProvider) NoOfControlPlane(int, bool) (int, error) {
	return -1, fmt.Errorf("[local] unsupported operation")
}

// NoOfDataStore implements resources.CloudFactory.
func (cloud *LocalProvider) NoOfDataStore(int, bool) (int, error) {
	return -1, fmt.Errorf("[local] unsupported operation")
}

// NoOfWorkerPlane implements resources.CloudFactory.
func (cloud *LocalProvider) NoOfWorkerPlane(resources.StorageFactory, int, bool) (int, error) {
	return -1, fmt.Errorf("[local] unsupported operation")
}

func (obj *LocalProvider) SwitchCluster(storage resources.StorageFactory) error {

	if isPresent(storage, obj.ClusterName) {

		printKubeconfig(storage, OperationStateCreate, obj.ClusterName)
		return nil
	}
	return fmt.Errorf("[local] Cluster not found")
}
