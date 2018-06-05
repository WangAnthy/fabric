/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"github.com/hyperledger/fabric/core/common/ccprovider"
	"github.com/hyperledger/fabric/core/ledger"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
)

// InstantiatedChaincodeStore returns information on chaincodes which are instantiated
type InstantiatedChaincodeStore interface {
	ChaincodeDeploymentSpec(channelID, chaincodeName string) (*pb.ChaincodeDeploymentSpec, error)
	ChaincodeDefinition(chaincodeName string, txSim ledger.QueryExecutor) (ccprovider.ChaincodeDefinition, error)
}

// Lifecycle provides methods to invoke the lifecycle system chaincode.
type Lifecycle struct {
	InstantiatedChaincodeStore InstantiatedChaincodeStore
}

// GetChaincodeDeploymentSpec retrieves a chaincode deployment spec for the specified chaincode.
func (l *Lifecycle) ChaincodeContainerInfo(channelID, chaincodeName string) (*ccprovider.ChaincodeContainerInfo, error) {
	cds, err := l.InstantiatedChaincodeStore.ChaincodeDeploymentSpec(channelID, chaincodeName)
	if err != nil {
		return nil, errors.Wrapf(err, "could not retrieve deployment spec for %s/%s", channelID, chaincodeName)
	}

	return DeploymentSpecToChaincodeContainerInfo(cds), nil
}

// GetChaincodeDefinition returns a ccprovider.ChaincodeDefinition for the chaincode
// associated with the provided txsim and name.
func (l *Lifecycle) GetChaincodeDefinition(chaincodeName string, txSim ledger.QueryExecutor) (ccprovider.ChaincodeDefinition, error) {
	return l.InstantiatedChaincodeStore.ChaincodeDefinition(chaincodeName, txSim)
}

func DeploymentSpecToChaincodeContainerInfo(cds *pb.ChaincodeDeploymentSpec) *ccprovider.ChaincodeContainerInfo {
	return &ccprovider.ChaincodeContainerInfo{
		Name:          cds.Name(),
		Version:       cds.Version(),
		Path:          cds.Path(),
		Type:          cds.CCType(),
		ContainerType: cds.ExecEnv.String(),
	}
}
