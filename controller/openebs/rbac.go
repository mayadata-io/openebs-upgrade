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
	"mayadata.io/openebs-upgrade/unstruct"
)

// getDesiredNamespace updates the namespace manifest as per the given configuration
// in OpenEBS CR.
func (p *Planner) getDesiredNamespace(namespace *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	namespace.SetNamespace(p.ObservedOpenEBS.Namespace)
	return namespace, nil
}

// getDesiredServiceAccount updates the service account manifest as per the
// given configuration in OpenEBS CR.
func (p *Planner) getDesiredServiceAccount(sa *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	sa.SetNamespace(p.ObservedOpenEBS.Namespace)
	return sa, nil
}

// getDesiredClusterRoleBinding updates the clusterRoleBinding manifest as per the
// given configuration in OpenEBS CR.
func (p *Planner) getDesiredClusterRoleBinding(crb *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	setNamespaceOfEachSubject := func(obj *unstructured.Unstructured) error {
		unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Namespace, "spec", "namespace")
		return nil
	}
	crbSubjects, _, err := unstruct.GetSlice(crb, "subjects")
	if err != nil {
		return crb, err
	}
	err = unstruct.SliceIterator(crbSubjects).ForEachUpdate(setNamespaceOfEachSubject)
	if err != nil {
		return crb, err
	}
	err = unstructured.SetNestedSlice(crb.Object, crbSubjects, "subjects")
	if err != nil {
		return crb, err
	}
	return crb, nil
}
