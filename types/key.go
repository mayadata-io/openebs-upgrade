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

package types

const (
	// MayaAPIServerNameKey is the name of maya-apiserver deployment
	MayaAPIServerNameKey string = "maya-apiserver"
	// MayaAPIServerServiceNameKey is the name of maya-apiserver service
	MayaAPIServerServiceNameKey string = "maya-apiserver-service"
	// ProvisionerNameKey is the name of openebs provisioner deployment
	ProvisionerNameKey string = "openebs-provisioner"
	// LocalProvisionerNameKey is the name of openebs local pv provisioner deployment
	LocalProvisionerNameKey string = "openebs-localpv-provisioner"
	// NDMNameKey is the name of NDM daemonset
	NDMNameKey string = "openebs-ndm"
	// NDMOperatorNameKey is the name of NDM operator deployment
	NDMOperatorNameKey string = "openebs-ndm-operator"
	// NDMConfigNameKey is the name of NDM configmap
	NDMConfigNameKey string = "openebs-ndm-config"
	// SnapshotOperatorNameKey is the name of snapshot operator deployment
	SnapshotOperatorNameKey string = "openebs-snapshot-operator"
	// AdmissionServerNameKey is the name of admission server deployment
	AdmissionServerNameKey string = "openebs-admission-server"

	// SnapshotControllerContainerKey is one of the container of snapshot operator
	SnapshotControllerContainerKey string = "snapshot-controller"
	// SnapshotProvisionerContainerKey is also one of the container of snapshot operator
	SnapshotProvisionerContainerKey string = "snapshot-provisioner"

	// KindClusterRole is the k8s kind of cluster role
	KindClusterRole string = "ClusterRole"
	// KindClusterRoleBinding is the k8s kind of cluster role binding
	KindClusterRoleBinding string = "ClusterRoleBinding"
	// KindConfigMap is the k8s kind of configmap
	KindConfigMap string = "ConfigMap"
	// KindDaemonSet is the k8s kind of daemonset
	KindDaemonSet string = "DaemonSet"
	// KindDeployment is the k8s kind of  deployment
	KindDeployment string = "Deployment"
	// KindNamespace is the k8s kind of namespace
	KindNamespace string = "Namespace"
	// KindService is the k8s kind of service
	KindService string = "Service"
	// KindServiceAccount is the k8s kind of service account
	KindServiceAccount string = "ServiceAccount"
	// KindCustomResourceDefinition is the k8s kind of CustomResourceDefinition
	KindCustomResourceDefinition string = "CustomResourceDefinition"

	// MayaAPIServerManifestKey is used to get the manifest of maya-apiserver
	MayaAPIServerManifestKey string = MayaAPIServerNameKey + "_" + KindDeployment
	// MayaAPIServerServiceManifestKey is used to get the manifest of maya-apiserver-service
	MayaAPIServerServiceManifestKey string = MayaAPIServerServiceNameKey + "_" + KindService
	// ProvisionerManifestKey is used to get the manifest of openebs provisioner
	ProvisionerManifestKey string = ProvisionerNameKey + "_" + KindDeployment
	// SnapshotOperatorManifestKey is used to get the manifest of snapshot operator
	SnapshotOperatorManifestKey string = SnapshotOperatorNameKey + "_" + KindDeployment
	// NDMConfigManifestKey is used to get the manifest of NDM configmap
	NDMConfigManifestKey string = NDMConfigNameKey + "_" + KindConfigMap
	// NDMManifestKey is used to get the manifest of NDM
	NDMManifestKey string = NDMNameKey + "_" + KindDaemonSet
	// NDMOperatorManifestKey is used to get the manifest of NDM operator
	NDMOperatorManifestKey string = NDMOperatorNameKey + "_" + KindDeployment
	// LocalProvisionerManifestKey is used to get the manifest of local pv provisioner
	LocalProvisionerManifestKey string = LocalProvisionerNameKey + "_" + KindDeployment
	// AdmissionServerManifestKey is used to get the manifest of admission server
	AdmissionServerManifestKey string = AdmissionServerNameKey + "_" + KindDeployment

	// CVCOperatorNameKey is the name of cvc-operator deployment.
	CVCOperatorNameKey = "cvc-operator"
	// CSPCOperatorNameKey is the name of cspc-operator deployment.
	CSPCOperatorNameKey = "cspc-operator"
	// CstorOperatorNameKey is the name of cstor-operator service account,
	// cluster role, cluster role binding.
	CstorOperatorNameKey = "openebs-cstor-operator"
	// CSPCCRDNameKey is the name of CSPC CRD
	CSPCCRDNameKey = "cstorpoolclusters.cstor.openebs.io"
	// CSPICRDNameKey is the name of the CSPI CRD
	CSPICRDNameKey = "cstorpoolinstances.cstor.openebs.io"
	// CstorVolumesCRDNameKey is the name of the Cstor Volume CRD
	CstorVolumesCRDNameKey = "cstorvolumes.cstor.openebs.io"
	// CstorVolumeConfigsCRDNameKey is the name of Cstor Volume Configs CRD
	CstorVolumeConfigsCRDNameKey = "cstorvolumeconfigs.cstor.openebs.io"
	// CstorVolumeReplicasNameKey is the name of Cstor Volume Replicas CRD
	CstorVolumeReplicasNameKey = "cstorvolumereplicas.cstor.openebs.io"
	// CstorVolumePoliciesNameKey is the name of Cstor Volume Policies CRD
	CstorVolumePoliciesNameKey = "cstorvolumepolicies.cstor.openebs.io"

	// CstorOperatorServiceAccountManifestKey is used to get the manifest of cstor-operator
	// service account.
	CstorOperatorServiceAccountManifestKey string = CstorOperatorNameKey + "_" +
		KindServiceAccount
	// CstorOperatorClusterRoleManifestKey  is used to get the manifest of cstor-operator
	// cluster role.
	CstorOperatorClusterRoleManifestKey string = CstorOperatorNameKey + "_" + KindClusterRole
	// CstorOperatorClusterRoleBindingManifestKey is used to get the manifest of cstor-operator
	// cluster role binding.
	CstorOperatorClusterRoleBindingManifestKey string = CstorOperatorNameKey + "_" +
		KindClusterRoleBinding
	// CVCOperatorManifestKey is used to get the manifest of CVC operator
	CVCOperatorManifestKey string = CVCOperatorNameKey + "_" + KindDeployment
	// CSPCOperatorManifestKey is used to get the manifest of CSPC operator
	CSPCOperatorManifestKey string = CSPCOperatorNameKey + "_" + KindDeployment
	// CSPCCRDManifestKey is used to get the manifest of CSPC CRD
	CSPCCRDManifestKey string = CSPCCRDNameKey + "_" + KindCustomResourceDefinition
	// CSPICRDManifestKey is used to get the manifest of CSPI CRD
	CSPICRDManifestKey string = CSPICRDNameKey + "_" + KindCustomResourceDefinition
	// CstorVolumesCRDManifestKey is used to get the manifest of Cstor Volumes CRD
	CstorVolumesCRDManifestKey string = CstorVolumesCRDNameKey + "_" +
		KindCustomResourceDefinition
	// CstorVolumesConfigsCRDManifestKey is used to get the manifest of Cstor Volume Configs CRD
	CstorVolumesConfigsCRDManifestKey string = CstorVolumeConfigsCRDNameKey + "_" +
		KindCustomResourceDefinition
	// CstorVolumesPoliciesCRDManifestKey is used to get the manifest of Cstor Volume Policies CRD
	CstorVolumesPoliciesCRDManifestKey string = CstorVolumePoliciesNameKey + "_" +
		KindCustomResourceDefinition
	// CstorVolumesReplicasCRDManifestKey is used to get the manifest of Cstor Volume Replicas CRD
	CstorVolumesReplicasCRDManifestKey string = CstorVolumeReplicasNameKey + "_" +
		KindCustomResourceDefinition

	// OpenEBSVersion150 is the OpenEBS version 1.5.0
	OpenEBSVersion150 string = "1.5.0"
	// OpenEBSVersion160 is the OpenEBS version 1.6.0
	OpenEBSVersion160 string = "1.6.0"
	// OpenEBSVersion170 is the OpenEBS version 1.7.0
	OpenEBSVersion170 string = "1.7.0"
	// OpenEBSVersion180 is the OpenEBS version 1.8.0
	OpenEBSVersion180 string = "1.8.0"
	// OpenEBSVersion190 is the OpenEBS version 1.9.0
	OpenEBSVersion190 string = "1.9.0"
)
