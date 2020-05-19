package types

const (
	// AnnotationPrefix is the prefix used across all
	// the annotations supported in this project
	AnnotationPrefix string = "openebs-upgrade.dao.mayadata.io"
	// AnnKeyOpenEBSUID is the annotation that refers to OpenEBS UID of openebs-upgrade
	AnnKeyOpenEBSUID string = AnnotationPrefix + "/openebs-uid"
)
