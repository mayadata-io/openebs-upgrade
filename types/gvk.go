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
	// GroupDAOMayaDataIO refers to the group for all
	// custom resources defined in this project
	GroupDAOMayaDataIO string = "dao.mayadata.io"

	// GroupOpenEBSIO refers to the group for all
	// custom resources defined in openebs
	GroupOpenEBSIO string = "openebs.io"

	// VersionV1Alpha1 refers to v1alpha1 version of the
	// custom resources used here
	VersionV1Alpha1 string = "v1alpha1"

	// APIVersionDAOMayaDataV1Alpha1 refers to v1alpha1 api
	// version of DAO based custom resources
	APIVersionDAOMayaDataV1Alpha1 string = GroupDAOMayaDataIO + "/" + VersionV1Alpha1

	// APIVersionOpenEBSV1Alpha1 refers to v1alpha1 api
	// version of openebs based custom resources
	APIVersionOpenEBSV1Alpha1 string = GroupOpenEBSIO + "/" + VersionV1Alpha1
)

// Kind is a custom datatype to refer to kubernetes native
// resource kind value
type Kind string

const (
	// KindOpenEBS refers to custom resource with
	// kind OpenEBS
	KindOpenEBS Kind = "OpenEBS"
)
