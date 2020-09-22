/*
Copyright 2020 The MayaData Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    https://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package openebs

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/types"
	"mayadata.io/openebs-upgrade/unstruct"
	"strconv"
	"strings"
)

const (
	// DefaultProvisionerReplicaCount is the default replica count for
	// openebs-k8s-provisioner.
	DefaultProvisionerReplicaCount int32 = 1
)

// setProvisionerDefaultsIfNotSet sets the default values for openebs-k8s-provisioner.
func (p *Planner) setProvisionerDefaultsIfNotSet() error {
	if p.ObservedOpenEBS.Spec.Provisioner == nil {
		p.ObservedOpenEBS.Spec.Provisioner = &types.Provisioner{}
	}
	if p.ObservedOpenEBS.Spec.Provisioner.Enabled == nil {
		p.ObservedOpenEBS.Spec.Provisioner.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.Provisioner.Enabled = true
	}
	// set the name with which openebs-provisioner will be deployed
	if len(p.ObservedOpenEBS.Spec.Provisioner.Name) == 0 {
		p.ObservedOpenEBS.Spec.Provisioner.Name = types.ProvisionerNameKey
	}
	if p.ObservedOpenEBS.Spec.Provisioner.ImageTag == "" {
		p.ObservedOpenEBS.Spec.Provisioner.ImageTag = p.ObservedOpenEBS.Spec.Version +
			p.ObservedOpenEBS.Spec.ImageTagSuffix
	}
	// form the image for openebs-k8s-provisioner.
	p.ObservedOpenEBS.Spec.Provisioner.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
		"openebs-k8s-provisioner:" + p.ObservedOpenEBS.Spec.Provisioner.ImageTag

	if p.ObservedOpenEBS.Spec.Provisioner.Replicas == nil {
		p.ObservedOpenEBS.Spec.Provisioner.Replicas = new(int32)
		*p.ObservedOpenEBS.Spec.Provisioner.Replicas = DefaultProvisionerReplicaCount
	}
	return nil
}

// updateOpenEBSProvisioner updates the openebs-provisioner manifest as per the
// reconcile.ObservedOpenEBS values.
func (p *Planner) updateOpenEBSProvisioner(deploy *unstructured.Unstructured) error {
	deploy.SetName(p.ObservedOpenEBS.Spec.Provisioner.Name)
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := deploy.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for openebs-provisioner deploy
	// 1. openebs-upgrade.dao.mayadata.io/component-name: openebs-provisioner
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.ProvisionerNameKey
	// set the desired labels
	deploy.SetLabels(desiredLabels)

	containers, err := unstruct.GetNestedSliceOrError(deploy, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}
	// update the containers
	updateContainer := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}
		envs, _, err := unstruct.GetSlice(obj, "spec", "env")
		if err != nil {
			return err
		}
		if containerName == types.OpenEBSProvisionerContainerKey {
			// Set the image of the container.
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.Provisioner.Image,
				"spec", "image")
			if err != nil {
				return err
			}
			// add ENVs to this container based on required conditions
			envs, err = p.addOpenEBSProvisionerEnvs(envs)
			if err != nil {
				return err
			}
		}
		err = unstructured.SetNestedSlice(obj.Object, envs, "spec", "env")
		if err != nil {
			return err
		}
		return nil
	}
	// Update the containers.
	err = unstruct.SliceIterator(containers).ForEachUpdate(updateContainer)
	if err != nil {
		return err
	}
	// Set back the value of the containers.
	err = unstructured.SetNestedSlice(deploy.Object,
		containers, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}

	return nil
}

// add env based on predefined conditions or values provided to the openebs-provisioner container.
func (p *Planner) addOpenEBSProvisionerEnvs(envs []interface{}) ([]interface{}, error) {
	// if leader election value is provided then insert the env for leader election
	if p.ObservedOpenEBS.Spec.Provisioner.EnableLeaderElection != nil {
		leaderElectionEnv := struct {
			name  string
			value string
		}{
			name:  "LEADER_ELECTION_ENABLED",
			value: strconv.FormatBool(*p.ObservedOpenEBS.Spec.Provisioner.EnableLeaderElection),
		}
		envs = append(envs, leaderElectionEnv)
	}
	return p.ignoreUpdatingImmutableEnvs(p.ObservedOpenEBS.Spec.Provisioner.ENV, envs)
}

// ignoreUpdatingImmutableEnvs returns the Envs without the ones which does not need update.
func (p *Planner) ignoreUpdatingImmutableEnvs(existingEnvs, envs []interface{}) ([]interface{}, error) {
	// check if we need not update some ENVs which are already present to avoid immutable
	// errors where update does not take place.
	var (
		elemIndexToReplace      string
		existingOpenEBSEnvValue map[string]interface{}
	)
	if !(existingEnvs == nil || len(existingEnvs) == 0) {
		env := func(obj *unstructured.Unstructured) error {
			envName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
			if err != nil {
				return err
			}
			if envName == "OPENEBS_NAMESPACE" {
				existingOpenEBSEnvValue, _, err = unstructured.NestedMap(obj.Object, "spec")
				if err != nil {
					return err
				}
			}
			return nil
		}
		err := unstruct.SliceIterator(existingEnvs).ForEach(env)
		if err != nil {
			return envs, err
		}
	}
	// get the index at which the env variable is present so that it can be deleted as
	// it is not required to be updated.
	if !(existingOpenEBSEnvValue == nil || len(existingOpenEBSEnvValue) == 0) {
		newENV := func(obj *unstructured.Unstructured) error {
			envName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
			if err != nil {
				return err
			}
			if envName == "OPENEBS_NAMESPACE" {
				elemIndexToReplace = strings.Split(obj.GetName(), "-")[1]
			}
			return nil
		}
		err := unstruct.SliceIterator(envs).ForEach(newENV)
		if err != nil {
			return envs, err
		}
		// replace the env variable if needs not be updated
		if len(elemIndexToReplace) > 0 {
			index, err := strconv.Atoi(elemIndexToReplace)
			if err != nil {
				return envs, err
			}
			envs[index] = existingOpenEBSEnvValue
		}
	}

	return envs, nil
}

func (p *Planner) fillProvisionerExistingValues(observedComponentDetails ObservedComponentDesiredDetails) error {
	var (
		containerName string
		err           error
	)
	p.ObservedOpenEBS.Spec.Provisioner.MatchLabels = observedComponentDetails.MatchLabels
	p.ObservedOpenEBS.Spec.Provisioner.PodTemplateLabels = observedComponentDetails.PodTemplateLabels
	if len(p.ObservedOpenEBS.Spec.Provisioner.ContainerName) > 0 {
		containerName = p.ObservedOpenEBS.Spec.Provisioner.ContainerName
	} else {
		containerName = types.OpenEBSProvisionerContainerKey
	}
	p.ObservedOpenEBS.Spec.Provisioner.ENV, err = fetchExistingContainerEnvs(
		observedComponentDetails.Containers, containerName)
	if err != nil {
		return err
	}

	return nil
}
