// Package crutils ...
package crutils

import (
	"github.com/IBM/ibm-object-csi-driver-operator/pkg/config"
	"github.com/IBM/ibm-object-csi-driver-operator/pkg/util/boolptr"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateCSIDriver ...
func (c *IBMObjectCSI) GenerateCSIDriver() *storagev1.CSIDriver {
	defaultFSGroupPolicy := storagev1.FileFSGroupPolicy
	return &storagev1.CSIDriver{
		ObjectMeta: metav1.ObjectMeta{
			Name: config.DriverName,
			Labels: map[string]string{
				"app.kubernetes.io/name":       "ibm-object-csi",
				"app.kubernetes.io/part-of":    config.CSIDriverName,
				"app.kubernetes.io/managed-by": config.CSIOperatorName,
			},
		},
		Spec: storagev1.CSIDriverSpec{
			AttachRequired: boolptr.False(),
			PodInfoOnMount: boolptr.True(),
			FSGroupPolicy:  &defaultFSGroupPolicy,
		},
	}
}

// GenerateControllerServiceAccount ...
func (c *IBMObjectCSI) GenerateControllerServiceAccount() *corev1.ServiceAccount {
	return getServiceAccount(c, config.CSIControllerServiceAccount)
}

// GenerateNodeServiceAccount ...
func (c *IBMObjectCSI) GenerateNodeServiceAccount() *corev1.ServiceAccount {
	return getServiceAccount(c, config.CSINodeServiceAccount)
}

func getServiceAccount(c *IBMObjectCSI, serviceAccountResourceName config.ResourceName) *corev1.ServiceAccount {
	secrets := GetImagePullSecrets(c.Spec.ImagePullSecrets)
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.GetNameForResource(serviceAccountResourceName, c.Name),
			Namespace: c.Namespace,
			Labels:    c.GetLabels(),
		},
		ImagePullSecrets: secrets,
	}
}

// GenerateExternalProvisionerClusterRole ...
func (c *IBMObjectCSI) GenerateExternalProvisionerClusterRole() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:   config.GetNameForResource(config.ExternalProvisionerClusterRole, c.Name),
			Labels: config.CommonCSIResourceLabels,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{config.SecretsResource},
				Verbs:     []string{config.VerbGet, config.VerbList},
			},
			{
				APIGroups: []string{""},
				Resources: []string{config.PersistentVolumesResource},
				Verbs:     []string{config.VerbGet, config.VerbList, config.VerbWatch, config.VerbCreate, config.VerbDelete},
			},
			{
				APIGroups: []string{""},
				Resources: []string{config.PersistentVolumeClaimsResource},
				Verbs:     []string{config.VerbGet, config.VerbList, config.VerbWatch, config.VerbUpdate},
			},
			{
				APIGroups: []string{config.StorageAPIGroup},
				Resources: []string{config.StorageClassesResource},
				Verbs:     []string{config.VerbGet, config.VerbList, config.VerbWatch},
			},
			{
				APIGroups: []string{""},
				Resources: []string{config.EventsResource},
				Verbs:     []string{config.VerbList, config.VerbWatch, config.VerbCreate, config.VerbUpdate, config.VerbPatch},
			},
			{
				APIGroups: []string{config.StorageAPIGroup},
				Resources: []string{config.CsiNodesResource},
				Verbs:     []string{config.VerbGet, config.VerbList, config.VerbWatch},
			},
			{
				APIGroups: []string{""},
				Resources: []string{config.NodesResource},
				Verbs:     []string{config.VerbGet, config.VerbList, config.VerbWatch},
			},
		},
	}
}

// GenerateExternalProvisionerClusterRoleBinding ...
func (c *IBMObjectCSI) GenerateExternalProvisionerClusterRoleBinding() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:   config.GetNameForResource(config.ExternalProvisionerClusterRoleBinding, c.Name),
			Labels: config.CommonCSIResourceLabels,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      config.GetNameForResource(config.CSIControllerServiceAccount, c.Name),
				Namespace: c.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     config.GetNameForResource(config.ExternalProvisionerClusterRole, c.Name),
			APIGroup: config.RbacAuthorizationAPIGroup,
		},
	}
}

// GenerateSCCForControllerClusterRole ...
func (c *IBMObjectCSI) GenerateSCCForControllerClusterRole() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:   config.GetNameForResource(config.CSIControllerSCCClusterRole, c.Name),
			Labels: config.CommonCSIResourceLabels,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups:     []string{config.SecurityOpenshiftAPIGroup},
				Resources:     []string{config.SecurityContextConstraintsResource},
				ResourceNames: []string{"anyuid"},
				Verbs:         []string{"use"},
			},
		},
	}
}

// GenerateSCCForControllerClusterRoleBinding ...
func (c *IBMObjectCSI) GenerateSCCForControllerClusterRoleBinding() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:   config.GetNameForResource(config.CSIControllerSCCClusterRoleBinding, c.Name),
			Labels: config.CommonCSIResourceLabels,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      config.GetNameForResource(config.CSIControllerServiceAccount, c.Name),
				Namespace: c.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     config.GetNameForResource(config.CSIControllerSCCClusterRole, c.Name),
			APIGroup: config.RbacAuthorizationAPIGroup,
		},
	}
}

// GenerateSCCForNodeClusterRole ...
func (c *IBMObjectCSI) GenerateSCCForNodeClusterRole() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:   config.GetNameForResource(config.CSINodeSCCClusterRole, c.Name),
			Labels: config.CommonCSIResourceLabels,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups:     []string{config.SecurityOpenshiftAPIGroup},
				Resources:     []string{config.SecurityContextConstraintsResource},
				ResourceNames: []string{"privileged"},
				Verbs:         []string{"use"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{config.NodesResource},
				Verbs:     []string{config.VerbGet},
			},
			{
				APIGroups: []string{""},
				Resources: []string{
					config.PersistentVolumesResource,
					config.SecretsResource},
				Verbs: []string{config.VerbGet},
			},
		},
	}
}

// GenerateSCCForNodeClusterRoleBinding ...
func (c *IBMObjectCSI) GenerateSCCForNodeClusterRoleBinding() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:   config.GetNameForResource(config.CSINodeSCCClusterRoleBinding, c.Name),
			Labels: config.CommonCSIResourceLabels,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      config.GetNameForResource(config.CSINodeServiceAccount, c.Name),
				Namespace: c.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     config.GetNameForResource(config.CSINodeSCCClusterRole, c.Name),
			APIGroup: config.RbacAuthorizationAPIGroup,
		},
	}
}

// Generates3fsSC ...
func (c *IBMObjectCSI) Generates3fsSC() *storagev1.StorageClass {
	reclaimPolicy := corev1.PersistentVolumeReclaimRetain
	return &storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:   config.GetNameForResource(config.S3fsStorageClass, c.Name),
			Labels: config.CommonCSIResourceLabels,
		},
		Provisioner:   config.DriverName,
		ReclaimPolicy: &reclaimPolicy,
		MountOptions: []string{
			"multipart_size=62",
			"max_dirty_data=51200",
			"parallel_count=8",
			"max_stat_cache_size=100000",
			"retries=5",
			"kernel_cache",
		},
		Parameters: map[string]string{
			"mounter": "s3fs",
			"client":  "awss3",
			"csi.storage.k8s.io/provisioner-secret-name":       "${pvc.name}",
			"csi.storage.k8s.io/provisioner-secret-namespace":  "${pvc.namespace}",
			"csi.storage.k8s.io/node-publish-secret-name":      "${pvc.name}",
			"csi.storage.k8s.io/node-publish-secret-namespace": "${pvc.namespace}",
		},
	}
}

// GenerateRcloneSC ...
func (c *IBMObjectCSI) GenerateRcloneSC() *storagev1.StorageClass {
	reclaimPolicy := corev1.PersistentVolumeReclaimRetain
	return &storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:   config.GetNameForResource(config.RcloneStorageClass, c.Name),
			Labels: config.CommonCSIResourceLabels,
		},
		Provisioner:   config.DriverName,
		ReclaimPolicy: &reclaimPolicy,
		MountOptions: []string{
			"acl=private",
			"bucket_acl=private",
			"upload_cutoff=256Mi",
			"chunk_size=64Mi",
			"max_upload_parts=64",
			"upload_concurrency=20",
			"copy_cutoff=1Gi",
			"memory_pool_flush_time=30s",
			"disable_checksum=true",
		},
		Parameters: map[string]string{
			"mounter": "rclone",
			"client":  "awss3",
			"csi.storage.k8s.io/provisioner-secret-name":       "${pvc.name}",
			"csi.storage.k8s.io/provisioner-secret-namespace":  "${pvc.namespace}",
			"csi.storage.k8s.io/node-publish-secret-name":      "${pvc.name}",
			"csi.storage.k8s.io/node-publish-secret-namespace": "${pvc.namespace}",
		},
	}
}
