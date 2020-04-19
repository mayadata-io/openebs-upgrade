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

	// CSINodeInfoCRDNameKey is the name of the CSINodeInfo CRD.
	CSINodeInfoCRDNameKey string = "csinodeinfos.csi.storage.k8s.io"
	// CSIVolumeCRDNameKey is the name of the CSIVolume CRD.
	CSIVolumeCRDNameKey string = "csivolumes.openebs.io"
	// VolumeSnapshotClassCRDNameKey is the name of the VolumeSnapshotClass CRD.
	VolumeSnapshotClassCRDNameKey string = "volumesnapshotclasses.snapshot.storage.k8s.io"
	// VolumeSnapshotContentCRDNameKey is the name of the VolumeSnapshotContent CRD.
	VolumeSnapshotContentCRDNameKey string = "volumesnapshotcontents.snapshot.storage.k8s.io"
	// VolumeSnapshotCRDNameKey is the name of the VolumeSnapshot CRD.
	VolumeSnapshotCRDNameKey string = "volumesnapshots.snapshot.storage.k8s.io"
	// CStorCSISnapshottterBindingNameKey is the name of the cstor csi snapshotter cluster role binding.
	CStorCSISnapshottterBindingNameKey string = "openebs-cstor-csi-snapshotter-binding"
	// CStorCSISnapshottterRoleNameKey is the name of the cstor csi snapshotter cluster role.
	CStorCSISnapshottterRoleNameKey string = "openebs-cstor-csi-snapshotter-role"
	// CStorCSIControllerSANameKey is the name of the cstor csi controller service account.
	CStorCSIControllerSANameKey string = "openebs-cstor-csi-controller-sa"
	// CStorCSIProvisionerRoleNameKey is the name of the cstor csi provisioner cluster role.
	CStorCSIProvisionerRoleNameKey string = "openebs-cstor-csi-provisioner-role"
	// CStorCSIProvisionerBindingNameKey is the name of the cstor csi provisioner cluster role binding.
	CStorCSIProvisionerBindingNameKey string = "openebs-cstor-csi-provisioner-binding"
	// CStorCSIControllerNameKey is the name of the cstor csi controller statefulset.
	CStorCSIControllerNameKey string = "openebs-cstor-csi-controller"
	// CStorCSIAttacherRoleNameKey is the name of the cstor csi attacher cluster role.
	CStorCSIAttacherRoleNameKey string = "openebs-cstor-csi-attacher-role"
	// CStorCSIAttacherBindingNameKey is the name of the cstor csi attacher cluster role binding.
	CStorCSIAttacherBindingNameKey string = "openebs-cstor-csi-attacher-binding"
	// CStorCSIClusterRegistrarRoleNameKey is the name of the cstor csi cluster registrar cluster role.
	CStorCSIClusterRegistrarRoleNameKey string = "openebs-cstor-csi-cluster-registrar-role"
	// CStorCSIClusterRegistrarBindingNameKey is the name of the cstor csi cluster registrar cluster role binding.
	CStorCSIClusterRegistrarBindingNameKey string = "openebs-cstor-csi-cluster-registrar-binding"
	// CStorCSINodeSANameKey is the name of the cstor csi node service account.
	CStorCSINodeSANameKey string = "openebs-cstor-csi-node-sa"
	// CStorCSIRegistrarRoleNameKey is the name of the cstor csi registrar cluster role.
	CStorCSIRegistrarRoleNameKey string = "openebs-cstor-csi-registrar-role"
	// CStorCSIRegistrarBindingNameKey is the name of the cstor csi registrar cluster role binding.
	CStorCSIRegistrarBindingNameKey string = "openebs-cstor-csi-registrar-binding"
	// CStorCSINodeNameKey is the name of the cstor csi node daemonset.
	CStorCSINodeNameKey string = "openebs-cstor-csi-node"
	// CStorCSIDriverNameKey is the name of the cstor csi csidriver.
	CStorCSIDriverNameKey string = "cstor.csi.openebs.io"

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
	// KindCustomResourceDefinition is the k8s kind of CustomResourceDefinition.
	KindCustomResourceDefinition string = "CustomResourceDefinition"
	// KindStatefulset is the k8s kind of statefulset.
	KindStatefulset string = "StatefulSet"
	// KindCSIDriver is the k8s kind of csidriver.
	KindCSIDriver string = "CSIDriver"

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

	// Below constants are used to get the manifests of cstor csi operator and driver.
	CSINodeInfoCRDManifestKey                  string = CSINodeInfoCRDNameKey + "_" + KindCustomResourceDefinition
	CSIVolumeCRDManifestKey                    string = CSIVolumeCRDNameKey + "_" + KindCustomResourceDefinition
	VolumeSnapshotClassCRDManifestKey          string = VolumeSnapshotClassCRDNameKey + "_" + KindCustomResourceDefinition
	VolumeSnapshotContentCRDManifestKey        string = VolumeSnapshotContentCRDNameKey + "_" + KindCustomResourceDefinition
	VolumeSnapshotCRDManifestKey               string = VolumeSnapshotCRDNameKey + "_" + KindCustomResourceDefinition
	CStorCSISnapshottterBindingManifestKey     string = CStorCSISnapshottterBindingNameKey + "_" + KindClusterRoleBinding
	CStorCSISnapshottterRoleManifestKey        string = CStorCSISnapshottterRoleNameKey + "_" + KindClusterRole
	CStorCSIControllerSAManifestKey            string = CStorCSIControllerSANameKey + "_" + KindServiceAccount
	CStorCSIProvisionerRoleManifestKey         string = CStorCSIProvisionerRoleNameKey + "_" + KindClusterRole
	CStorCSIProvisionerBindingManifestKey      string = CStorCSIProvisionerBindingNameKey + "_" + KindClusterRoleBinding
	CStorCSIControllerManifestKey              string = CStorCSIControllerNameKey + "_" + KindStatefulset
	CStorCSIAttacherRoleManifestKey            string = CStorCSIAttacherRoleNameKey + "_" + KindClusterRole
	CStorCSIAttacherBindingManifestKey         string = CStorCSIAttacherBindingNameKey + "_" + KindClusterRoleBinding
	CStorCSIClusterRegistrarRoleManifestKey    string = CStorCSIClusterRegistrarRoleNameKey + "_" + KindClusterRole
	CStorCSIClusterRegistrarBindingManifestKey string = CStorCSIClusterRegistrarBindingNameKey + "_" + KindClusterRoleBinding
	CStorCSINodeSAManifestKey                  string = CStorCSINodeSANameKey + "_" + KindServiceAccount
	CStorCSIRegistrarRoleManifestKey           string = CStorCSIRegistrarRoleNameKey + "_" + KindClusterRole
	CStorCSIRegistrarBindingManifestKey        string = CStorCSIRegistrarBindingNameKey + "_" + KindClusterRoleBinding
	CStorCSINodeManifestKey                    string = CStorCSINodeNameKey + "_" + KindDaemonSet
	CStorCSIDriverManifestKey                  string = CStorCSIDriverNameKey + "_" + KindCSIDriver

	// CVCOperatorNameKey is the name of cvc-operator deployment.
	CVCOperatorNameKey = "cvc-operator"
	// CSPCOperatorNameKey is the name of cspc-operator deployment.
	CSPCOperatorNameKey = "cspc-operator"
	// CSPCCRDNameKey is the name of CSPC CRD
	CSPCCRDNameKey = "cstorpoolclusters.cstor.openebs.io"

	// CVCOperatorManifestKey is used to get the manifest of CVC operator
	CVCOperatorManifestKey string = CVCOperatorNameKey + "_" + KindDeployment
	// CSPCOperatorManifestKey is used to get the manifest of CSPC operator
	CSPCOperatorManifestKey string = CSPCOperatorNameKey + "_" + KindDeployment
	// CSPCCRDManifestKey is used to get the manifest of CSPC CRD
	CSPCCRDManifestKey string = CSPCCRDNameKey + "_" + KindCustomResourceDefinition

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

	// OSImageUbuntu1804 is the OS Image value of a Node.
	OSImageUbuntu1804 string = "Ubuntu 18.04"
	// OSImageSLES12 is the OS Image value of a Node.
	OSImageSLES12 string = "SUSE Linux Enterprise Server 12"
	// OSImageSLES15 is the OS Image value of a Node.
	OSImageSLES15 string = "SUSE Linux Enterprise Server 15"

	// NamespaceKubeSystem is the value of kube-system namespace
	NamespaceKubeSystem string = "kube-system"

	// CSISupportedVersion is the k8s version from where csi is supported.
	CSISupportedVersion string = "v1.14.0"
)
