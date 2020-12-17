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
	// CStorCSIISCSIADMConfigmapNameKey is the name of openebs-cstor-csi-iscsiadm configmap
	CStorCSIISCSIADMConfigmapNameKey string = "openebs-cstor-csi-iscsiadm"
	// SnapshotOperatorNameKey is the name of snapshot operator deployment
	SnapshotOperatorNameKey string = "openebs-snapshot-operator"
	// AdmissionServerNameKey is the name of admission server deployment
	AdmissionServerNameKey string = "openebs-admission-server"
	// AdmissionServerSVCNameKey is the name of admission server service
	AdmissionServerSVCNameKey string = "admission-server-svc"

	// APIServerContainerKey is the name of the container of maya-apiserver.
	APIServerContainerKey string = "maya-apiserver"
	// NodeDiskOperatorContainerKey is the name of the container of ndm-operator.
	NodeDiskOperatorContainerKey string = "node-disk-operator"
	// NodeDiskManagerContainerKey is the name of the container of ndm daemon.
	NodeDiskManagerContainerKey string = "node-disk-manager"
	// SnapshotControllerContainerKey is one of the container of snapshot operator
	SnapshotControllerContainerKey string = "snapshot-controller"
	// SnapshotProvisionerContainerKey is also one of the container of snapshot operator
	SnapshotProvisionerContainerKey string = "snapshot-provisioner"
	// OpenEBSProvisionerContainerKey is the name of the container of openebs-provisioner.
	OpenEBSProvisionerContainerKey string = "openebs-provisioner"
	// LocalPVProvisionerContainerKey is the container name of openebs-provisioner-hostpath
	LocalPVProvisionerContainerKey string = "openebs-provisioner-hostpath"
	// AdmissionServerContainerKey is the container name of admission server
	AdmissionServerContainerKey string = "admission-webhook"
	// CSPCOperatorContainerKey is the container name of cspc-operator
	CSPCOperatorContainerKey string = "cspc-operator"
	// CVCOperatorContainerKey is the container name of cvc-operator
	CVCOperatorContainerKey string = "cvc-operator"
	// CStorCSIControllerCSIPluginContainerKey is the openebs-csi-plugin container name.
	CStorCSIControllerCSIPluginContainerKey string = "openebs-csi-plugin"
	// CStorCSINodeCSIPluginContainerKey is the openebs-csi-plugin container name.
	CStorCSINodeCSIPluginContainerKey string = "openebs-csi-plugin"
	// NDMDaemonContainerKey is the node-disk-manager container.
	NDMDaemonContainerKey string = "node-disk-manager"

	// MoacContainerKey is one of the container of moac deployment.
	MoacContainerKey string = "moac"
	// MayastorContainerKey is one of the container of mayastor daemonset.
	MayastorContainerKey string = "mayastor"
	// MayastorGRPCContainerKey is one of the container of mayastor daemonset.
	MayastorGRPCContainerKey string = "mayastor-grpc"
	// NATSContainerKey is one of the container of nats deployment.
	NATSContainerKey string = "nats"
	// MayastorCSIContainerKey is one of the container of mayastor-csi daemonset.
	MayastorCSIContainerKey string = "mayastor-csi"

	// CSINodeInfoCRDNameKey is the name of the CSINodeInfo CRD.
	CSINodeInfoCRDNameKey string = "csinodeinfos.csi.storage.k8s.io"
	// CSIVolumeCRDNameKey is the name of the CSIVolume CRD.
	CSIVolumeCRDNameKey string = "csivolumes.openebs.io"
	// VolumeSnapshotClassCRDNameKey is the name of the VolumeSnapshotClass CRD.
	VolumeSnapshotClassCRDNameKey string = "volumesnapshotclasses.snapshot.storage.k8s.io"
	// VolumeSnapshotContentCRDNameKey is the name of the VolumeSnapshotContent CRD.
	VolumeSnapshotContentCRDNameKey string = "volumesnapshotcontents.snapshot.storage.k8s.io"
	// CStorVolumeAttachmentsCRDNameKey is the name of the CStorVolumeAttachments CRD.
	CStorVolumeAttachmentsCRDNameKey string = "cstorvolumeattachments.cstor.openebs.io"
	// VolumeSnapshotCRDNameKey is the name of the VolumeSnapshot CRD.
	VolumeSnapshotCRDNameKey string = "volumesnapshots.snapshot.storage.k8s.io"
	// CStorVolumesCRDV1NameKey is the name of CStorVolumes V1 CRD
	CStorVolumesCRDV1NameKey string = "cstorvolumes.cstor.openebs.io"
	// CStorVolumeConfigsCRDV1NameKey is the name of CStorVolumeConfigs V1 CRD
	CStorVolumeConfigsCRDV1NameKey string = "cstorvolumeconfigs.cstor.openebs.io"
	// CStorVolumePoliciesCRDV1NameKey is the name of CStorVolumePolicies V1 CRD
	CStorVolumePoliciesCRDV1NameKey string = "cstorvolumepolicies.cstor.openebs.io"
	// CStorVolumeReplicasCRDV1NameKey is the name of CStorVolumeReplicas V1 CRD
	CStorVolumeReplicasCRDV1NameKey string = "cstorvolumereplicas.cstor.openebs.io"
	// CStorRestoresCRDV1alpha1NameKey is the name of CStorRestores v1alpha1 CRD
	CStorRestoresCRDV1alpha1NameKey string = "cstorrestores.openebs.io"
	// MayastorPoolsCRDV1alpha1NameKey is the name of Mayastor pools v1alpha1 CRD
	MayastorPoolsCRDV1alpha1NameKey string = "mayastorpools.openebs.io"
	// CStorCompletedBackupsCRDV1alpha1NameKey is the name of cstorcompletedbackups v1alpha1 CRD
	CStorCompletedBackupsCRDV1alpha1NameKey string = "cstorcompletedbackups.openebs.io"
	// CStorBackupsCRDV1alpha1NameKey is the name of CStorBackups v1alpha1 CRD
	CStorBackupsCRDV1alpha1NameKey string = "cstorbackups.openebs.io"
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

	// MoacSANameKey is the name of the moac service account.
	MoacSANameKey string = "moac"
	// MoacClusterRoleNameKey is the name of the moac cluster role.
	MoacClusterRoleNameKey string = "moac"
	// MoacClusterRoleBindingNameKey is the name of the moac cluster role binding.
	MoacClusterRoleBindingNameKey string = "moac"
	// MoacDeploymentNameKey is the name of the moac deployment.
	MoacDeploymentNameKey string = "moac"
	// MoacServiceNameKey is the name of the moac service.
	MoacServiceNameKey string = "moac"
	// NATSDeploymentNameKey is the name of the nats deployment.
	NATSDeploymentNameKey string = "nats"
	// NATSServiceNameKey is the name of the nats service.
	NATSServiceNameKey string = "nats"
	// MayastorCSIDaemonsetNameKey is the name of the mayastor-csi daemonset
	MayastorCSIDaemonsetNameKey string = "mayastor-csi"
	// MayastorDaemonsetNameKey is the name of the mayastor daemonset
	MayastorDaemonsetNameKey string = "mayastor"
	// MayastorNamespaceNameKey is the name of the mayastor namespace
	MayastorNamespaceNameKey string = "mayastor"
	// MayastorMOACSVCNameKey is the key to identify moac service
	MayastorMOACSVCNameKey string = "moac-svc"

	// OpenEBSNodeSetupDaemonsetNameKey is the name of daemonset which is launched to
	// setup ISCSI client on nodes prior to OpenEBS installation.
	OpenEBSNodeSetupDaemonsetNameKey string = "openebs-node-setup"
	// OpenEBSNodeSetupConfigmapNameKey is the name of configmap which contains the
	// configuration that is run to install ISCSI client on the nodes.
	OpenEBSNodeSetupConfigmapNameKey string = "node-setup"

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
	// KindPriorityClass is the k8s kind of PriorityClass.
	KindPriorityClass string = "PriorityClass"
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
	CStorVolumeAttachmentCRDManifestKey        string = CStorVolumeAttachmentsCRDNameKey + "_" + KindCustomResourceDefinition
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
	CStorCSIISCSIADMManifestKey                string = CStorCSIISCSIADMConfigmapNameKey + "_" + KindConfigMap
	CStorCSIDriverManifestKey                  string = CStorCSIDriverNameKey + "_" + KindCSIDriver

	// CVCOperatorNameKey is the name of cvc-operator deployment.
	CVCOperatorNameKey = "cvc-operator"
	// CVCOperatorServiceNameKey is the name of cvc-operator service.
	CVCOperatorServiceNameKey = "cvc-operator-service"
	// CSPCOperatorNameKey is the name of cspc-operator deployment.
	CSPCOperatorNameKey = "cspc-operator"
	// CSPCCRDV1NameKey is the name of CSPC V1 CRD
	CSPCCRDV1NameKey = "cstorpoolclusters.cstor.openebs.io"
	// CSPCCRDV1alpha1NameKey is the name of v1alpha1 CSPC CRD
	CSPCCRDV1alpha1NameKey = "cstorpoolclusters.openebs.io"
	// CSPICRDV1NameKey is the name of v1 CSPI CRD
	CSPICRDV1NameKey = "cstorpoolinstances.cstor.openebs.io"
	// CStorAdmissionServerNameKey is the name of CStor admission server
	CStorAdmissionServerNameKey = "openebs-cstor-admission-server"
	// OpenEBSCSINodePriorityClassNameKey is the name of OpenEBS CSI node priority class.
	OpenEBSCSINodePriorityClassNameKey = "openebs-csi-node-critical"
	// OpenEBSCSIControllerPriorityClassNameKey is the name of OpenEBS CSI controller priority class.
	OpenEBSCSIControllerPriorityClassNameKey = "openebs-csi-controller-critical"

	// CStorAdmissionServerManifestKey is used to get the manifest of CStor admission server
	CStorAdmissionServerManifestKey string = CStorAdmissionServerNameKey + "_" + KindDeployment
	// CVCOperatorManifestKey is used to get the manifest of CVC operator
	CVCOperatorManifestKey string = CVCOperatorNameKey + "_" + KindDeployment
	// CSPCOperatorManifestKey is used to get the manifest of CSPC operator
	CSPCOperatorManifestKey string = CSPCOperatorNameKey + "_" + KindDeployment
	// CSPCCRDV1ManifestKey is used to get the manifest of CSPC CRD
	CSPCCRDV1ManifestKey string = CSPCCRDV1NameKey + "_" + KindCustomResourceDefinition
	// CSPCCRDV1alpha1ManifestKey is used to get the manifest of CSPC CRD v1alpha1
	CSPCCRDV1alpha1ManifestKey string = CSPCCRDV1alpha1NameKey + "_" + KindCustomResourceDefinition
	// CSPICRDV1ManifestKey is used to get the manifest of CSPI CRD v1
	CSPICRDV1ManifestKey string = CSPICRDV1NameKey + "_" + KindCustomResourceDefinition

	// MoacSAManifestKey is used to get the manifest of moac service account.
	MoacSAManifestKey string = MoacSANameKey + "_" + KindServiceAccount
	// MoacClusterRoleManifestKey is used to get the manifest of moac cluster role.
	MoacClusterRoleManifestKey string = MoacClusterRoleBindingNameKey + "_" + KindClusterRole
	// MoacClusterRoleBindingManifestKey is used to get the manifest of moac cluster role binding.
	MoacClusterRoleBindingManifestKey string = MoacClusterRoleBindingNameKey + "_" + KindClusterRoleBinding
	// MoacDeploymentManifestKey is used to get the manifest of moac deployment.
	MoacDeploymentManifestKey string = MoacDeploymentNameKey + "_" + KindDeployment
	// MoacServiceManifestKey is used to get the manifest of moac service.
	MoacServiceManifestKey string = MoacServiceNameKey + "_" + KindService
	// MayastorDaemonsetManifestKey is used to get the manifest of mayastor daemonset.
	MayastorDaemonsetManifestKey string = MayastorDaemonsetNameKey + "_" + KindDaemonSet
	// MayastorNamespaceManifestKey is used to get the manifest of mayastor namespace.
	MayastorNamespaceManifestKey string = MayastorNamespaceNameKey + "_" + KindNamespace
	// NATSDeploymentManifestKey is used to get the manifest of nats deployment.
	NATSDeploymentManifestKey string = NATSDeploymentNameKey + "_" + KindDeployment
	// NATSServiceManifestKey is used to get the manifest of nats service.
	NATSServiceManifestKey string = NATSServiceNameKey + "_" + KindService
	// MayastorCSIDaemonsetManifestKey is used to get the manifest of mayastor-csi daemonset.
	MayastorCSIDaemonsetManifestKey string = MayastorCSIDaemonsetNameKey + "_" + KindDaemonSet
	// MayastorPoolsCRDManifestKey is used to get the manifest of mayastorpools CRD.
	MayastorPoolsCRDManifestKey string = MayastorPoolsCRDV1alpha1NameKey + "_" + KindCustomResourceDefinition

	// MayastorSupportedVersion is the openebs version from where mayastor is supported.
	MayastorSupportedVersion string = "1.10.0-ee" // MayastorSupportedVersion is the openebs version from where mayastor is supported.
	// NATSSupportedVersion is the openebs version from where NATS is supported.
	NATSSupportedVersion string = "2.0.0"
	// MayastorCSISupportedVersion is the openebs version from where mayastor-csi is supported.
	MayastorCSISupportedVersion string = "2.0.0"

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
	// OpenEBSVersion190EE is the OpenEBS version 1.9.0-ee
	OpenEBSVersion190EE string = "1.9.0-ee"
	// OpenEBSVersion1100 is the OpenEBS version 1.10.0
	OpenEBSVersion1100 string = "1.10.0"
	// OpenEBSVersion1100EE is the OpenEBS version 1.10.0-ee
	OpenEBSVersion1100EE string = "1.10.0-ee"
	// OpenEBSVersion1110 is the OpenEBS version 1.11.0
	OpenEBSVersion1110 string = "1.11.0"
	// OpenEBSVersion1110EE is the OpenEBS version 1.11.0-ee
	OpenEBSVersion1110EE string = "1.11.0-ee"
	// OpenEBSVersion1120 is the OpenEBS version 1.12.0
	OpenEBSVersion1120 string = "1.12.0"
	// OpenEBSVersion1120EE is the OpenEBS version 1.12.0-ee
	OpenEBSVersion1120EE string = "1.12.0-ee"
	// OpenEBSVersion200 is the OpenEBS version 2.0.0
	OpenEBSVersion200 string = "2.0.0"
	// OpenEBSVersion200EE is the OpenEBS version 2.0.0-ee
	OpenEBSVersion200EE string = "2.0.0-ee"
	// OpenEBSVersion210 is the OpenEBS version 2.1.0
	OpenEBSVersion210 string = "2.1.0"
	// OpenEBSVersion210EE is the OpenEBS version 2.1.0-ee
	OpenEBSVersion210EE string = "2.1.0-ee"
	// OpenEBSVersion220 is the OpenEBS version 2.2.0
	OpenEBSVersion220 string = "2.2.0"
	// OpenEBSVersion220EE is the OpenEBS version 2.2.0-ee
	OpenEBSVersion220EE string = "2.2.0-ee"
	// OpenEBSVersion240 is the OpenEBS version 2.4.0
	OpenEBSVersion240 string = "2.4.0"

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
	// CSISupportedVersionFromOpenEBS200 is the k8s version from where csi is supported for
	// OpenEBS version 2.0.0 or greater.
	CSISupportedVersionFromOpenEBS200 string = "v1.17.0"
	// K8sVersion1170 is the const for kubernetes version v1.17.0
	K8sVersion1170 string = "v1.17.0"

	// OpenEBSMayaOperatorSANameKey is the name of OpenEBS service account.
	OpenEBSMayaOperatorSANameKey string = "openebs-maya-operator"
	// OpenEBSMayaOperatorRoleNameKey is the name of OpenEBS cluster role.
	OpenEBSMayaOperatorRoleNameKey string = "openebs-maya-operator"
	// OpenEBSMayaOperatorBindingNameKey is the name of OpenEBS cluster role binding.
	OpenEBSMayaOperatorBindingNameKey string = "openebs-maya-operator"

	// openebs-upgrade dao specific label, this label will be present across all the
	// OpenEBS components whether created or adopted by openebs-upgrade.

	// OpenEBSUpgradeDAOManagedLabelKey is the key for openebs-upgrade dao managed label
	OpenEBSUpgradeDAOManagedLabelKey string = "openebs-upgrade.dao.mayadata.io/managed"
	// OpenEBSUpgradeDAOManagedLabelValue is the value for openebs-upgrade dao managed label
	OpenEBSUpgradeDAOManagedLabelValue string = "true"

	// OpenEBSUpgradeDAOAdoptLabelKey is the key for openebs-upgrade dao adopt label
	OpenEBSUpgradeDAOAdoptLabelKey string = "openebs-upgrade.dao.mayadata.io/adopt"
	// OpenEBSUpgradeDAOAdoptLabelValue is the value for openebs-upgrade dao adopt label
	OpenEBSUpgradeDAOAdoptLabelValue string = "true"

	// OpenEBSComponentGroupLabelKey is the label key which helps in identifying
	// a particular group of OpenEBS components i.e., the component-type value for
	// CSI related cluster roles can still be "cluster-role" but the value for
	// component-group can be "csi" which will help in identifying all the
	// OpenEBS related cluster roles as well as cluster roles specific to CSI also.
	OpenEBSComponentGroupLabelKey string = "openebs-upgrade.dao.mayadata.io/component-group"
	// OpenEBSComponentSubGroupLabelKey is the label key which helps in identifying
	// a particular subgroup of OpenEBS components i.e., the value for
	// component-group for NDM components can be "ndm" which will help in identifying all the
	// OpenEBS NDM related components, but it will have a component-subgroup too like
	// operator or daemon which will tell that this particular NDM components is an operator or
	// NDM daemon or something else.
	OpenEBSComponentSubGroupLabelKey string = "openebs-upgrade.dao.mayadata.io/component-subgroup"
	// OpenEBSComponentNameLabelKey is the label key which helps in
	// identifying a particular OpenEBS component i.e., openebs-ndm will be the label value
	// for ndm daemonset while openebs-ndm-operator will be the value for NDM operator.
	OpenEBSComponentNameLabelKey string = "openebs-upgrade.dao.mayadata.io/component-name"

	// OpenEBSSAComponentNameLabelValue is the value of the component-name label
	// of OpenEBS service account.
	OpenEBSSAComponentNameLabelValue string = "openebs-maya-operator"
	// CStorCSICtrlSAComponentNameLabelValue is the value of the component-name label
	// of OpenEBS CStor CSI controller service account.
	CStorCSICtrlSAComponentNameLabelValue string = "cstor-csi-controller"
	// CStorCSINodeSAComponentNameLabelValue is the value of the component-name label
	// of OpenEBS CStor CSI node service account.
	CStorCSINodeSAComponentNameLabelValue string = "cstor-csi-node"

	// CSIComponentGroupLabelValue is the value of the component-group label
	// of CSI components.
	CSIComponentGroupLabelValue string = "csi"
	// OpenEBSCStorCSIComponentGroupLabelValue is the value of the component-group label
	// of CSI components.
	OpenEBSCStorCSIComponentGroupLabelValue string = "cstor-csi"
	// OpenEBSRoleComponentNameLabelValue is the value of the component-name label
	// of OpenEBS cluster role.
	OpenEBSRoleComponentNameLabelValue string = "openebs-maya-operator"

	// OpenEBSRoleBindingComponentNameLabelValue is the value of the component-name label
	// of OpenEBS cluster role.
	OpenEBSRoleBindingComponentNameLabelValue string = "openebs-maya-operator"

	// OpenEBSMayaAPIServerComponentGroupLabelValue is the value of the component-group label
	// of OpenEBS apiservers.
	OpenEBSMayaAPIServerComponentGroupLabelValue string = "maya-apiserver"
	// OpenEBSProvisionerComponentGroupLabelValue is the value of the component-group label
	// of OpenEBS provisioner.
	OpenEBSProvisionerComponentGroupLabelValue string = "openebs-provisioner"
	// OpenEBSOperatorComponentSubGroupLabelValue is the value of the component-subgroup label
	// of OpenEBS operators.
	OpenEBSOperatorComponentSubGroupLabelValue string = "operator"
	// OpenEBSDaemonComponentSubGroupLabelValue is the value of the component-subgroup label
	// of OpenEBS daemons.
	OpenEBSDaemonComponentSubGroupLabelValue string = "daemon"
	// OpenEBSAdmissionServerComponentGroupLabelValue is the value of the component-group label
	// of OpenEBS admission server.
	OpenEBSAdmissionServerComponentGroupLabelValue string = "admission-server"
	// OpenEBSNDMComponentGroupLabelValue is the value of the component-group label
	// of NDM components.
	OpenEBSNDMComponentGroupLabelValue string = "ndm"
	// CSPCComponentGroupLabelValue is the value of the component-group label
	// of CSPC components.
	CSPCComponentGroupLabelValue string = "cspc"
	// CSPIComponentGroupLabelValue is the value of the component-group label
	// of CSPI components.
	CSPIComponentGroupLabelValue string = "cspi"
	// CVCComponentGroupLabelValue is the value of the component-group label
	// of CVC components.
	CVCComponentGroupLabelValue string = "cvc"
	// CStorAdmissionServerComponentNameLabelValue is the value of the component-name label
	// of CStor admission server component.
	CStorAdmissionServerComponentNameLabelValue string = "cstor-admission-webhook"
	// OpenEBSVersionLabelKey is the label that can be used to get the OpenEBS version.
	OpenEBSVersionLabelKey string = "openebs.io/version"
	// OpenEBSMayastorComponentGroupLabelValue is the value of the component-group label
	// of mayastor components.
	OpenEBSMayastorComponentGroupLabelValue string = "mayastor"

	// ComponentNameLabelKey is the label key which is found in OpenEBS components.
	// These labels and their values already exists in the OpenEBS components even
	// if these are not installed via openebs-upgrade.
	//
	// NOTE: These keys and values can be used to identify a particular OpenEBS component.
	ComponentNameLabelKey                            string = "openebs.io/component-name"
	NDMComponentNameLabelValue                       string = "ndm"
	NDMOperatorComponentNameLabelValue               string = "ndm-operator"
	NDMConfigComponentNameLabelValue                 string = "ndm-config"
	CSPCOperatorComponentNameLabelValue              string = "cspc-operator"
	CVCOperatorComponentNameLabelValue               string = "cvc-operator"
	CVCOperatorServiceComponentNameLabelValue        string = "cvc-operator-svc"
	MayaAPIServerComponentNameLabelValue             string = "maya-apiserver"
	MayaAPIServerSVCComponentNameLabelValue          string = "maya-apiserver-svc"
	AdmissionServerComponentNameLabelValue           string = "admission-webhook"
	AdmissionServerSVCComponentNameLabelValue        string = "admission-webhook-svc"
	LocalPVProvisionerComponentNameLabelValue        string = "openebs-localpv-provisioner"
	OpenEBSProvisionerComponentNameLabelValue        string = "openebs-provisioner"
	SnapshotOperatorComponentNameLabelValue          string = "openebs-snapshot-operator"
	CStorCSINodeComponentNameLabelValue              string = "openebs-cstor-csi-node"
	CStorCSIControllerComponentNameLabelValue        string = "openebs-cstor-csi-controller"
	CStorCSIISCSIADMConfigmapComponentNameLabelValue string = "openebs-cstor-csi-iscsiadm"
	MayastorMOACComponentNameLabelValue              string = "moac"
	MayastorMOACServiceComponentNameLabelValue       string = "moac-svc"
	MayastorMayastorComponentNameLabelValue          string = "mayastor"

	KeyName              string = "name"
	KeyEnabled           string = "enabled"
	KeyReplicas          string = "replicas"
	KeyNodeSelector      string = "nodeSelector"
	KeyAffinity          string = "affinity"
	KeyTolerations       string = "tolerations"
	KeyMatchLabels       string = "matchLabels"
	KeyPodTemplateLabels string = "podTemplateLabels"
	KeyResources         string = "resources"
	KeyImage             string = "image"
	KeyImageTag          string = "imageTag"
	KeyContainerName     string = "containerName"
	KeyISCSIPath         string = "iscsiPath"
	KeyAdoptionJobID     string = "mayadata.io/openebsAdoptionJobId"
	KeyMicroK8s          string = "microk8s"

	QUAYIOOPENEBSREGISTRY string = "quay.io/openebs/"
	MAYADATAIOREGISTRY    string = "mayadataio/"
	OPENEBSREGISTRY       string = "openebs/"
	QUAYIOK8SCSI          string = "quay.io/k8scsi/"

	CSIResizerVersion010                string = "v0.1.0"
	CSIResizerVersion040                string = "v0.4.0"
	CSISnapshotterVersion201            string = "v2.0.1"
	CSISnapshotControllerVersion201     string = "v2.0.1"
	CSIProvisionerVersion111            string = "v1.1.1"
	CSIProvisionerVersion150            string = "v1.5.0"
	CSIProvisionerVersion160            string = "v1.6.0"
	CSIAttacherVersion200               string = "v2.0.0"
	CSIAttacherVersion220               string = "v2.2.0"
	CSIAttacherVersion111               string = "v1.1.1"
	CSIClusterDriverRegistrarVersion101 string = "v1.0.1"
	CSINodeDriverRegistrarVersion101    string = "v1.0.1"
	CSINodeDriverRegistrarVersion110    string = "v1.1.0"
	CSINodeDriverRegistrarVersion130    string = "v1.3.0"
)
