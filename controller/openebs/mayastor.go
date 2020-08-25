package openebs

import (
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"mayadata.io/openebs-upgrade/types"
	"mayadata.io/openebs-upgrade/unstruct"
	"strings"
)

const (
	MayastorVersion010EE     string = "0.1.0-ee"
	MayastorVersion020       string = "0.2.0"
	MayastorVersion020EE     string = "0.2.0-ee"
	MayastorVersion030       string = "v0.3.0"
	MayastorVersion030EE     string = "v0.3.0-ee"
	MayastorCSIVersion030    string = "v0.3.0"
	MayastorCSIVersion030EE  string = "v0.3.0-ee"
	NATSVersion21Alpine311   string = "2.1-alpine3.11"
	NATSVersion21Alpine311EE string = "2.1-alpine3.11-ee"
	DefaultMoacReplicaCount  int32  = 1
	DefaultNATSReplicaCount  int32  = 1
)

// supportedMayastorVersionForOpenEBSVersion stores the mapping for
// Mayastor to OpenEBS version i.e., a Mayastor version for each of the
// supported OpenEBS versions.
var supportedMayastorVersionForOpenEBSVersion = map[string]string{
	types.OpenEBSVersion1100EE: MayastorVersion010EE,
	types.OpenEBSVersion1110:   MayastorVersion020,
	types.OpenEBSVersion1110EE: MayastorVersion020EE,
	types.OpenEBSVersion1120:   MayastorVersion020,
	types.OpenEBSVersion1120EE: MayastorVersion020EE,
	types.OpenEBSVersion200:    MayastorVersion030,
	types.OpenEBSVersion200EE:  MayastorVersion030EE,
}

// supportedNATSVersionForOpenEBSVersion stores the mapping for
// NATS to OpenEBS version i.e., a NATS version for each of the
// supported OpenEBS versions.
var supportedNATSVersionForOpenEBSVersion = map[string]string{
	types.OpenEBSVersion200:   NATSVersion21Alpine311,
	types.OpenEBSVersion200EE: NATSVersion21Alpine311EE,
}

// supportedMayastorCSIVersionForOpenEBSVersion stores the mapping for
// mayastor-csi to OpenEBS version i.e., a mayastor-csi version for each of the
// supported OpenEBS versions.
var supportedMayastorCSIVersionForOpenEBSVersion = map[string]string{
	types.OpenEBSVersion200:   MayastorCSIVersion030,
	types.OpenEBSVersion200EE: MayastorCSIVersion030EE,
}

var (
	// List of images which are by default fetched from quay.io/k8scsi registry.
	CSIProvisionerForMOACImage                string
	CSIAttacherForMOACImage                   string
	CSINodeDriverRegistrarForMayastorImage    string
	CSINodeDriverRegistrarForMayastorCSIImage string
)

// SupportedCSIProvisionerVersionForMOACVersion stores the mapping for
// CSI provisioner to moac(Mayastor) version.
var SupportedCSIProvisionerVersionForMOACVersion = map[string]string{
	types.OpenEBSVersion1100:   types.CSIProvisionerVersion150,
	types.OpenEBSVersion1100EE: types.CSIProvisionerVersion111,
	types.OpenEBSVersion1110:   types.CSIProvisionerVersion160,
	types.OpenEBSVersion1110EE: types.CSIProvisionerVersion160,
	types.OpenEBSVersion1120:   types.CSIProvisionerVersion160,
	types.OpenEBSVersion1120EE: types.CSIProvisionerVersion160,
	types.OpenEBSVersion200:    types.CSIProvisionerVersion160,
	types.OpenEBSVersion200EE:  types.CSIProvisionerVersion160,
}

// SupportedCSIAttacherVersionForMOACVersion stores the mapping for
// CSI provisioner to MOAC(mayastor) version.
var SupportedCSIAttacherVersionForMOACVersion = map[string]string{
	types.OpenEBSVersion1100:   types.CSIAttacherVersion111,
	types.OpenEBSVersion1100EE: types.CSIAttacherVersion111,
	types.OpenEBSVersion1110:   types.CSIAttacherVersion220,
	types.OpenEBSVersion1110EE: types.CSIAttacherVersion220,
	types.OpenEBSVersion1120:   types.CSIAttacherVersion220,
	types.OpenEBSVersion1120EE: types.CSIAttacherVersion220,
	types.OpenEBSVersion200:    types.CSIAttacherVersion220,
	types.OpenEBSVersion200EE:  types.CSIAttacherVersion220,
}

// SupportedCSINodeDriverRegistrarVersionForMayastorVersion stores the mapping for
// CSINodeDriverRegistrar to mayastor version.
var SupportedCSINodeDriverRegistrarVersionForMayastorVersion = map[string]string{
	types.OpenEBSVersion1100EE: types.CSINodeDriverRegistrarVersion110,
	types.OpenEBSVersion1110:   types.CSINodeDriverRegistrarVersion130,
	types.OpenEBSVersion1110EE: types.CSINodeDriverRegistrarVersion130,
	types.OpenEBSVersion1120:   types.CSINodeDriverRegistrarVersion130,
	types.OpenEBSVersion1120EE: types.CSINodeDriverRegistrarVersion130,
	types.OpenEBSVersion200:    types.CSINodeDriverRegistrarVersion130,
	types.OpenEBSVersion200EE:  types.CSINodeDriverRegistrarVersion130,
}

// SupportedCSINodeDriverRegistrarVersionForMayastorCSIVersion stores the mapping for
// CSINodeDriverRegistrar to mayastor-csi version.
var SupportedCSINodeDriverRegistrarVersionForMayastorCSIVersion = map[string]string{
	types.OpenEBSVersion200:   types.CSINodeDriverRegistrarVersion130,
	types.OpenEBSVersion200EE: types.CSINodeDriverRegistrarVersion130,
}

// Set the default values for Mayastor if not already given.
func (p *Planner) setMayastorDefaultsIfNotSet() error {
	var (
		defaultImageRegistryForMayastor string
		// List of images which are by default fetched from quay.io/k8scsi registry.
		CSIProvisionerForMOACImageTag                string
		CSIAttacherForMOACImageTag                   string
		CSINodeDriverRegistrarForMayastorImageTag    string
		CSINodeDriverRegistrarForMayastorCSIImageTag string
	)
	isMayastorSupported, err := p.isMayastorSupported()
	// Do not return the error as not to block installing other components.
	if err != nil {
		isMayastorSupported = false
		glog.Errorf("Failed to set Mayastor defaults, error: %v", err)
	}
	// set the default image registry value.
	defaultImageRegistryForMayastor = p.ObservedOpenEBS.Spec.ImagePrefix
	// If OpenEBS version is lower than 2.0.0 then use the default imageRegistry if no custom registry
	// is provided otherwise use mayadata/ as default registry for Mayastor components from OpenEBS version
	// 2.0.0 onwards for community edition.
	comp, err := compareVersion(p.ObservedOpenEBS.Spec.Version, types.OpenEBSVersion200)
	if err != nil {
		glog.Errorf("Error setting default image registry for Mayastor based on OpenEBS version: %+v", err)
	}
	if comp >= 0 && defaultImageRegistryForMayastor == types.QUAYIOOPENEBSREGISTRY {
		// use mayadata/ only in case of community edition, for enterprise mayadataio/ will
		// be used by default.
		defaultImageRegistryForMayastor = "mayadata/"
	}
	if !isMayastorSupported {
		glog.V(5).Infof("Skipping Mayastor installation.")
	}

	if p.ObservedOpenEBS.Spec.MayastorConfig == nil {
		p.ObservedOpenEBS.Spec.MayastorConfig = &types.MayastorConfig{}
	}

	if p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Enabled == nil {
		p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Enabled = false
	}

	if isMayastorSupported && *p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Enabled == true {
		if len(p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Name) == 0 {
			p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Name = types.MoacDeploymentNameKey
		}
		// set the defaults for MOAC service
		if p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Service == nil {
			p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Service = &types.MOACService{}
		}
		if len(p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Service.Name) == 0 {
			p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Service.Name = types.MoacServiceNameKey
		}
		if p.ObservedOpenEBS.Spec.MayastorConfig.Moac.ImageTag == "" {
			if moacVersion, exist :=
				supportedMayastorVersionForOpenEBSVersion[p.ObservedOpenEBS.Spec.Version]; exist {
				p.ObservedOpenEBS.Spec.MayastorConfig.Moac.ImageTag = moacVersion +
					p.ObservedOpenEBS.Spec.ImageTagSuffix
			} else {
				return errors.Errorf("Failed to get moac version for the given OpenEBS version: %s",
					p.ObservedOpenEBS.Spec.Version)
			}
		}
		p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Image = defaultImageRegistryForMayastor +
			"moac:" + p.ObservedOpenEBS.Spec.MayastorConfig.Moac.ImageTag

		if p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Replicas == nil {
			p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Replicas = new(int32)
			*p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Replicas = DefaultMoacReplicaCount
		}

		// form the CSI provisioner image for MOAC(Mayastor)
		if csiProvisionerForMOAC, exist :=
			SupportedCSIProvisionerVersionForMOACVersion[p.ObservedOpenEBS.Spec.Version]; exist {
			CSIProvisionerForMOACImageTag = "csi-provisioner:" +
				csiProvisionerForMOAC
		} else {
			return errors.Errorf(
				"Failed to get csi-provisioner version for moac(mayastor) for the given OpenEBS version: %s",
				p.ObservedOpenEBS.Spec.Version)
		}

		// form the CSI attacher image for MOAC(mayastor)
		if csiAttacherForMOAC, exist :=
			SupportedCSIAttacherVersionForMOACVersion[p.ObservedOpenEBS.Spec.Version]; exist {
			CSIAttacherForMOACImageTag = "csi-attacher:" +
				csiAttacherForMOAC
		} else {
			return errors.Errorf(
				"Failed to get csi-attacher version for moac(mayastor) for the given OpenEBS version: %s",
				p.ObservedOpenEBS.Spec.Version)
		}
	}

	if p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Enabled == nil {
		p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Enabled = false
	}

	if isMayastorSupported && *p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Enabled == true {
		if len(p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Name) == 0 {
			p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Name = types.MayastorDaemonsetNameKey
		}
		if p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Mayastor.ImageTag == "" {
			if mayastorVersion, exist :=
				supportedMayastorVersionForOpenEBSVersion[p.ObservedOpenEBS.Spec.Version]; exist {
				p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Mayastor.ImageTag = mayastorVersion +
					p.ObservedOpenEBS.Spec.ImageTagSuffix
			} else {
				return errors.Errorf("Failed to get mayastor version for the given OpenEBS version: %s",
					p.ObservedOpenEBS.Spec.Version)
			}
		}
		if p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.MayastorGRPC.ImageTag == "" {
			if mayastorVersion, exist :=
				supportedMayastorVersionForOpenEBSVersion[p.ObservedOpenEBS.Spec.Version]; exist {
				p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.MayastorGRPC.ImageTag = mayastorVersion +
					p.ObservedOpenEBS.Spec.ImageTagSuffix
			} else {
				return errors.Errorf("Failed to get mayastor grpc version for the given OpenEBS version: %s",
					p.ObservedOpenEBS.Spec.Version)
			}
		}

		p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Mayastor.Image = defaultImageRegistryForMayastor +
			"mayastor:" + p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Mayastor.ImageTag
		p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.MayastorGRPC.Image = p.ObservedOpenEBS.Spec.ImagePrefix +
			"mayastor-grpc:" + p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.MayastorGRPC.ImageTag

		// form the csi-node-driver-registrar image for mayastor for the given OpenEBS version
		if csiNodeDriverRegistrar, exist :=
			SupportedCSINodeDriverRegistrarVersionForMayastorVersion[p.ObservedOpenEBS.Spec.Version]; exist {
			CSINodeDriverRegistrarForMayastorImageTag = "csi-node-driver-registrar:" +
				csiNodeDriverRegistrar
		} else {
			return errors.Errorf(
				"Failed to get csi-node-driver-registrar version for mayastor for the given OpenEBS version: %s",
				p.ObservedOpenEBS.Spec.Version)
		}
	}

	if p.ObservedOpenEBS.Spec.MayastorConfig.NATS.Enabled == nil {
		p.ObservedOpenEBS.Spec.MayastorConfig.NATS.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.MayastorConfig.NATS.Enabled = false
	}
	// set the defaults for NATS service
	if p.ObservedOpenEBS.Spec.MayastorConfig.NATS.Service == nil {
		p.ObservedOpenEBS.Spec.MayastorConfig.NATS.Service = &types.NATSService{}
	}
	// Update NATS default only if it supported
	isNATSSupported, err := p.isNATSSupported()
	// Do not return the error as not to block installing other components.
	if err != nil {
		isNATSSupported = false
		glog.Errorf("Failed to set NATS defaults, error: %v", err)
	}
	if isMayastorSupported && isNATSSupported &&
		*p.ObservedOpenEBS.Spec.MayastorConfig.NATS.Enabled == true {
		if len(p.ObservedOpenEBS.Spec.MayastorConfig.NATS.Name) == 0 {
			p.ObservedOpenEBS.Spec.MayastorConfig.NATS.Name = types.NATSDeploymentNameKey
		}
		if len(p.ObservedOpenEBS.Spec.MayastorConfig.NATS.Service.Name) == 0 {
			p.ObservedOpenEBS.Spec.MayastorConfig.NATS.Service.Name = types.NATSServiceNameKey
		}
		if p.ObservedOpenEBS.Spec.MayastorConfig.NATS.ImageTag == "" {
			if natsVersion, exist :=
				supportedNATSVersionForOpenEBSVersion[p.ObservedOpenEBS.Spec.Version]; exist {
				p.ObservedOpenEBS.Spec.MayastorConfig.NATS.ImageTag = natsVersion
			} else {
				return errors.Errorf("Failed to get nats version for the given OpenEBS version: %s",
					p.ObservedOpenEBS.Spec.Version)
			}
		}
		p.ObservedOpenEBS.Spec.MayastorConfig.NATS.Image =
			"nats:" + p.ObservedOpenEBS.Spec.MayastorConfig.NATS.ImageTag

		if p.ObservedOpenEBS.Spec.MayastorConfig.NATS.Replicas == nil {
			p.ObservedOpenEBS.Spec.MayastorConfig.NATS.Replicas = new(int32)
			*p.ObservedOpenEBS.Spec.MayastorConfig.NATS.Replicas = DefaultNATSReplicaCount
		}
	}

	if p.ObservedOpenEBS.Spec.MayastorConfig.MayastorCSI.Enabled == nil {
		p.ObservedOpenEBS.Spec.MayastorConfig.MayastorCSI.Enabled = new(bool)
		*p.ObservedOpenEBS.Spec.MayastorConfig.MayastorCSI.Enabled = false
	}
	// Update mayastor-csi default only if it supported
	isMayastorCSISupported, err := p.isMayastorCSISupported()
	// Do not return the error as not to block installing other components.
	if err != nil {
		isMayastorCSISupported = false
		glog.Errorf("Failed to set mayastor-csi defaults, error: %v", err)
	}
	if isMayastorSupported && isMayastorCSISupported &&
		*p.ObservedOpenEBS.Spec.MayastorConfig.MayastorCSI.Enabled == true {
		if len(p.ObservedOpenEBS.Spec.MayastorConfig.MayastorCSI.Name) == 0 {
			p.ObservedOpenEBS.Spec.MayastorConfig.MayastorCSI.Name = types.MayastorCSIDaemonsetNameKey
		}
		if p.ObservedOpenEBS.Spec.MayastorConfig.MayastorCSI.ImageTag == "" {
			if mayastorCSIVersion, exist :=
				supportedMayastorCSIVersionForOpenEBSVersion[p.ObservedOpenEBS.Spec.Version]; exist {
				p.ObservedOpenEBS.Spec.MayastorConfig.MayastorCSI.ImageTag = mayastorCSIVersion +
					p.ObservedOpenEBS.Spec.ImageTagSuffix
			} else {
				return errors.Errorf("Failed to get mayastor-csi version for the given OpenEBS version: %s",
					p.ObservedOpenEBS.Spec.Version)
			}
		}
		p.ObservedOpenEBS.Spec.MayastorConfig.MayastorCSI.Image = defaultImageRegistryForMayastor +
			"mayastor-csi:" + p.ObservedOpenEBS.Spec.MayastorConfig.Moac.ImageTag

		// form the csi-node-driver-registrar image for mayastor-csi for the given OpenEBS version
		if csiNodeDriverRegistrar, exist :=
			SupportedCSINodeDriverRegistrarVersionForMayastorCSIVersion[p.ObservedOpenEBS.Spec.Version]; exist {
			CSINodeDriverRegistrarForMayastorCSIImageTag = "csi-node-driver-registrar:" +
				csiNodeDriverRegistrar
		} else {
			return errors.Errorf(
				"Failed to get csi-node-driver-registrar version for mayastor-csi for the given OpenEBS version: %s",
				p.ObservedOpenEBS.Spec.Version)
		}
	}

	// check if the image registry is the default ones i.e., quay.io/openebs/, openebs/ or mayadataio/,
	// if not then form the k8s repositories related images also so that they can also be pulled from
	// the specified repository only.
	if !(p.ObservedOpenEBS.Spec.ImagePrefix == types.QUAYIOOPENEBSREGISTRY ||
		p.ObservedOpenEBS.Spec.ImagePrefix == types.MAYADATAIOREGISTRY ||
		p.ObservedOpenEBS.Spec.ImagePrefix == types.OPENEBSREGISTRY) {
		CSIProvisionerForMOACImage = p.ObservedOpenEBS.Spec.ImagePrefix + CSIProvisionerForMOACImageTag
		CSIAttacherForMOACImage = p.ObservedOpenEBS.Spec.ImagePrefix + CSIAttacherForMOACImageTag
		CSINodeDriverRegistrarForMayastorImage = p.ObservedOpenEBS.Spec.ImagePrefix + CSINodeDriverRegistrarForMayastorImageTag
		CSINodeDriverRegistrarForMayastorCSIImage = p.ObservedOpenEBS.Spec.ImagePrefix + CSINodeDriverRegistrarForMayastorCSIImageTag
	} else {
		CSIProvisionerForMOACImage = types.QUAYIOK8SCSI + CSIProvisionerForMOACImageTag
		CSIAttacherForMOACImage = types.QUAYIOK8SCSI + CSIAttacherForMOACImageTag
		CSINodeDriverRegistrarForMayastorImage = types.QUAYIOK8SCSI + CSINodeDriverRegistrarForMayastorImageTag
		CSINodeDriverRegistrarForMayastorCSIImage = types.QUAYIOK8SCSI + CSINodeDriverRegistrarForMayastorCSIImageTag
	}

	return nil
}

// isMayastorSupported checks if mayastor is supported or not in the current kubernetes cluster of openebs version,
// if not it will return false else true.
func (p *Planner) isMayastorSupported() (bool, error) {
	// compare the openebs version with the supported version of mayastor.
	comp, err := compareVersion(p.ObservedOpenEBS.Spec.Version, types.MayastorSupportedVersion)
	if err != nil {
		return false, errors.Errorf("Error comparing versions, error: %v", err)
	}

	// if the versions are equal, check if that contains the "ee" in the image. As mayastor installation
	// is supported from 1.10.0-ee.
	if comp < 0 || (comp == 0 && !strings.Contains(p.ObservedOpenEBS.Spec.Version, "ee")) {
		glog.Warningf("Mayastor is not supported in %s openebs version. "+
			"Mayastor is supported from %s openebs version.", p.ObservedOpenEBS.Spec.Version, types.MayastorSupportedVersion)
		return false, nil
	}

	return true, nil
}

// isNATSSupported checks if nats is supported or not for the given OpenEBS version.
func (p *Planner) isNATSSupported() (bool, error) {
	// compare the openebs version with the supported version of nats.
	comp, err := compareVersion(p.ObservedOpenEBS.Spec.Version, types.NATSSupportedVersion)
	if err != nil {
		return false, errors.Errorf("Error comparing versions, error: %v", err)
	}
	if comp < 0 {
		return false, nil
	}

	return true, nil
}

// isMayastorCSISupported checks if mayastor-csi is supported or not for the given OpenEBS version.
func (p *Planner) isMayastorCSISupported() (bool, error) {
	// compare the openebs version with the supported version of nats.
	comp, err := compareVersion(p.ObservedOpenEBS.Spec.Version, types.MayastorCSISupportedVersion)
	if err != nil {
		return false, errors.Errorf("Error comparing versions, error: %v", err)
	}
	if comp < 0 {
		return false, nil
	}

	return true, nil
}

// updateMoac updates the moac manifest as per the reconcile.ObservedOpenEBS values.
func (p *Planner) updateMoac(deploy *unstructured.Unstructured) error {
	deploy.SetName(p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Name)
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := deploy.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for moac deploy
	// 1. openebs-upgrade.dao.mayadata.io/component-group: mayastor
	// 2. openebs-upgrade.dao.mayadata.io/component-name: moac
	desiredLabels[types.OpenEBSComponentGroupLabelKey] = types.OpenEBSMayastorComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.MoacDeploymentNameKey
	// set the desired labels
	deploy.SetLabels(desiredLabels)

	// Overwrite the namespace to mayastor for mayastor based components.
	// Note: mayastor based components will be installed only in mayastor namespace only.
	deploy.SetNamespace(types.MayastorNamespaceNameKey)

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
		if containerName == types.MoacContainerKey {
			// Set the image of the container.
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Image,
				"spec", "image")
		} else if containerName == ContainerCSIProvisionerName {
			// Set the image of the container.
			err = unstructured.SetNestedField(obj.Object, CSIProvisionerForMOACImage,
				"spec", "image")
		} else if containerName == ContainerCSIAttacherName {
			// Set the image of the container.
			err = unstructured.SetNestedField(obj.Object, CSIAttacherForMOACImage,
				"spec", "image")
		}
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

// updateMoacService updates the moac service manifest as per the
// reconcile.ObservedOpenEBS values.
func (p *Planner) updateMoacService(svc *unstructured.Unstructured) error {
	svc.SetName(p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Service.Name)
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := svc.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for moac service
	// 1. openebs-upgrade.dao.mayadata.io/component-group: mayastor
	// 2. openebs-upgrade.dao.mayadata.io/component-name: moac
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSMayastorComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.MoacServiceNameKey
	// set the desired labels
	svc.SetLabels(desiredLabels)

	// Overwrite the namespace to mayastor for mayastor based components.
	// Note: mayastor based components will be installed only in mayastor namespace only.
	svc.SetNamespace(types.MayastorNamespaceNameKey)

	return nil
}

// updateMayastor updates the values of mayastor daemonset as per given configuration.
func (p *Planner) updateMayastor(daemonset *unstructured.Unstructured) error {
	daemonset.SetName(p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Name)
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := daemonset.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for mayastor daemonset:
	// 1. openebs-upgrade.dao.mayadata.io/component-group: mayastor
	// 2. openebs-upgrade.dao.mayadata.io/component-name: mayastor
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSMayastorComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.MayastorDaemonsetNameKey
	// set the desired labels
	daemonset.SetLabels(desiredLabels)

	// Overwrite the namespace to mayastor for mayastor based components.
	// Note: mayastor based components will be installed only in mayastor namespace only.
	daemonset.SetNamespace(types.MayastorNamespaceNameKey)

	containers, err := unstruct.GetNestedSliceOrError(daemonset, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}
	updateContainer := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}

		if containerName == types.MayastorContainerKey {
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Mayastor.Image,
				"spec", "image")
			if err != nil {
				return err
			}

			// Set the resource of the container.
			resources := map[string]interface{}{
				"limits": map[string]interface{}{
					"cpu":           "1",
					"memory":        "500Mi",
					"hugepages-2Mi": "1Gi",
				},
				"requests": map[string]interface{}{
					"cpu":           "1",
					"memory":        "500Mi",
					"hugepages-2Mi": "1Gi",
				},
			}
			err = unstructured.SetNestedField(obj.Object, resources, "spec", "resources")
		}
		if containerName == types.MayastorGRPCContainerKey {
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.MayastorGRPC.Image,
				"spec", "image")
		}
		if err != nil {
			return err
		}
		// set the image of csi-driver-registrar container
		if containerName == ContainerCSIDriverRegistrarName {
			// Set the image of the container.
			err = unstructured.SetNestedField(obj.Object, CSINodeDriverRegistrarForMayastorImage,
				"spec", "image")
			if err != nil {
				return err
			}
		}

		if p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Resources != nil &&
			containerName != types.MayastorContainerKey {
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Resources,
				"spec", "resources")
		} else if p.ObservedOpenEBS.Spec.Resources != nil &&
			containerName != types.MayastorContainerKey {
			err = unstructured.SetNestedField(obj.Object,
				p.ObservedOpenEBS.Spec.Resources, "spec", "resources")
		}
		if err != nil {
			return err
		}
		return nil
	}
	err = unstruct.SliceIterator(containers).ForEachUpdate(updateContainer)
	if err != nil {
		return err
	}
	err = unstructured.SetNestedSlice(daemonset.Object, containers, "spec",
		"template", "spec", "containers")
	if err != nil {
		return err
	}

	return nil
}

// updateMayastorNamespace updates the mayastor namespace
// structure as per the provided values otherwise default values.
func (p *Planner) updateMayastorNamespace(namespace *unstructured.Unstructured) error {
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := namespace.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Set some component specific labels in order to identify specific components.
	// These labels will be only set by openebs-upgrade and will help the end-users
	// identify a particular or a set of OpenEBS components.
	//
	// Component specific labels for mayastor daemonset:
	// 1. openebs-upgrade.dao.mayadata.io/component-group: mayastor
	// 2. openebs-upgrade.dao.mayadata.io/component-name: mayastor
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSMayastorComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.MayastorNamespaceNameKey
	// set the desired labels
	namespace.SetLabels(desiredLabels)

	namespace.SetName(types.MayastorNamespaceNameKey)

	return nil
}

// removeMayastorManifests removes the manifests of mayastor if disabled.
func (p *Planner) removeMayastorManifests() {
	if *p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Enabled == false &&
		*p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Enabled == false {
		delete(p.ComponentManifests, types.MayastorNamespaceManifestKey)
		delete(p.ComponentManifests, types.MayastorPoolsCRDManifestKey)
	}

	if *p.ObservedOpenEBS.Spec.MayastorConfig.Moac.Enabled == false {
		delete(p.ComponentManifests, types.MoacSAManifestKey)
		delete(p.ComponentManifests, types.MoacClusterRoleManifestKey)
		delete(p.ComponentManifests, types.MoacClusterRoleBindingManifestKey)
		delete(p.ComponentManifests, types.MoacDeploymentManifestKey)
		delete(p.ComponentManifests, types.MoacServiceManifestKey)
	}

	if *p.ObservedOpenEBS.Spec.MayastorConfig.Mayastor.Enabled == false {
		delete(p.ComponentManifests, types.MayastorDaemonsetManifestKey)
	}

	if *p.ObservedOpenEBS.Spec.MayastorConfig.MayastorCSI.Enabled == false {
		delete(p.ComponentManifests, types.MayastorCSIDaemonsetManifestKey)
	}

	if *p.ObservedOpenEBS.Spec.MayastorConfig.NATS.Enabled == false {
		delete(p.ComponentManifests, types.NATSDeploymentManifestKey)
		delete(p.ComponentManifests, types.NATSServiceManifestKey)
	}
}

// updateNATS updates the nats manifest as per the reconcile.ObservedOpenEBS values.
func (p *Planner) updateNATS(deploy *unstructured.Unstructured) error {
	deploy.SetName(p.ObservedOpenEBS.Spec.MayastorConfig.NATS.Name)
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := deploy.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for nats deploy
	// 1. openebs-upgrade.dao.mayadata.io/component-group: mayastor
	// 2. openebs-upgrade.dao.mayadata.io/component-name: nats
	desiredLabels[types.OpenEBSComponentGroupLabelKey] = types.OpenEBSMayastorComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.NATSDeploymentNameKey
	// set the desired labels
	deploy.SetLabels(desiredLabels)

	// Overwrite the namespace to mayastor for mayastor based components.
	// Note: mayastor based components will be installed only in mayastor namespace only.
	deploy.SetNamespace(types.MayastorNamespaceNameKey)

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
		if containerName == types.NATSContainerKey {
			// Set the image of the container.
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.MayastorConfig.NATS.Image,
				"spec", "image")
		}
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

// updateNATSService updates the NATS service manifest as per the
// reconcile.ObservedOpenEBS values.
func (p *Planner) updateNATSService(svc *unstructured.Unstructured) error {
	svc.SetName(p.ObservedOpenEBS.Spec.MayastorConfig.NATS.Service.Name)
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := svc.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for nats service
	// 1. openebs-upgrade.dao.mayadata.io/component-group: mayastor
	// 2. openebs-upgrade.dao.mayadata.io/component-name: nats
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSMayastorComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.NATSServiceNameKey
	// set the desired labels
	svc.SetLabels(desiredLabels)

	// Overwrite the namespace to mayastor for mayastor based components.
	// Note: mayastor based components will be installed only in mayastor namespace only.
	svc.SetNamespace(types.MayastorNamespaceNameKey)

	return nil
}

// updateMayastorCSI updates the values of mayastor-csi daemonset as per given configuration.
func (p *Planner) updateMayastorCSI(daemonset *unstructured.Unstructured) error {
	daemonset.SetName(p.ObservedOpenEBS.Spec.MayastorConfig.MayastorCSI.Name)
	// desiredLabels is used to form the desired labels of a particular OpenEBS component.
	desiredLabels := daemonset.GetLabels()
	if desiredLabels == nil {
		desiredLabels = make(map[string]string, 0)
	}
	// Component specific labels for mayastor-csi daemonset:
	// 1. openebs-upgrade.dao.mayadata.io/component-group: mayastor
	// 2. openebs-upgrade.dao.mayadata.io/component-name: mayastor-csi
	desiredLabels[types.OpenEBSComponentGroupLabelKey] =
		types.OpenEBSMayastorComponentGroupLabelValue
	desiredLabels[types.OpenEBSComponentNameLabelKey] = types.MayastorCSIDaemonsetNameKey
	// set the desired labels
	daemonset.SetLabels(desiredLabels)

	// Overwrite the namespace to mayastor for mayastor based components.
	// Note: mayastor based components will be installed only in mayastor namespace only.
	daemonset.SetNamespace(types.MayastorNamespaceNameKey)

	containers, err := unstruct.GetNestedSliceOrError(daemonset, "spec", "template", "spec", "containers")
	if err != nil {
		return err
	}
	updateContainer := func(obj *unstructured.Unstructured) error {
		containerName, _, err := unstructured.NestedString(obj.Object, "spec", "name")
		if err != nil {
			return err
		}

		if containerName == types.MayastorCSIContainerKey {
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.MayastorConfig.MayastorCSI.Image,
				"spec", "image")
			if err != nil {
				return err
			}
		}
		// set the image of csi-driver-registrar container
		if containerName == ContainerCSIDriverRegistrarName {
			// Set the image of the container.
			err = unstructured.SetNestedField(obj.Object, CSINodeDriverRegistrarForMayastorCSIImage,
				"spec", "image")
			if err != nil {
				return err
			}
		}

		if p.ObservedOpenEBS.Spec.MayastorConfig.MayastorCSI.Resources != nil {
			err = unstructured.SetNestedField(obj.Object, p.ObservedOpenEBS.Spec.MayastorConfig.MayastorCSI.Resources,
				"spec", "resources")
		} else if p.ObservedOpenEBS.Spec.Resources != nil {
			err = unstructured.SetNestedField(obj.Object,
				p.ObservedOpenEBS.Spec.Resources, "spec", "resources")
		}
		if err != nil {
			return err
		}
		return nil
	}
	err = unstruct.SliceIterator(containers).ForEachUpdate(updateContainer)
	if err != nil {
		return err
	}
	err = unstructured.SetNestedSlice(daemonset.Object, containers, "spec",
		"template", "spec", "containers")
	if err != nil {
		return err
	}

	return nil
}
