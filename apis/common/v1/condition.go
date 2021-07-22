/*
Copyright 2020 Wim Henderickx.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"sort"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// A ConditionKind represents a condition kind for a resource
type ConditionKind string

// Condition Kinds.
const (
	// handled per resource
	ConditionKindInternalLeafRef ConditionKind = "InternalLeafrefValidationSuccess"
	// handled per target per resource
	ConditionKindExternalLeafRef ConditionKind = "ExternalLeafrefValidationSuccess"
	// handled per resource
	ConditionKindParent ConditionKind = "ParentValidationSuccess"
	// handled per resource
	ConditionKindTarget ConditionKind = "TargetFound"
	// handled per target per resource
	ConditionKindConfiguration ConditionKind = "ConfigurationSuccess"
)

// A ConditionReason represents the reason a resource is in a condition.
type ConditionReason string

// Reasons a resource validation is or is not ok
const (
	ConditionReasonSuccess ConditionReason = "Success"
	ConditionReasonFailed  ConditionReason = "Failed"
)

// Reasons a resource target is or is not ok
const (
	ConditionReasonFound    ConditionReason = "Target Found"
	ConditionReasonNotFound ConditionReason = "No valid target Found"
)

// Reasons a resource is or is not ready wrt configuration
const (
	ConditionReasonNone             ConditionReason = "None"
	ConditionReasonCreating         ConditionReason = "Creating"
	ConditionReasonDeleting         ConditionReason = "Deleting"
	ConditionReasonReconcileSuccess ConditionReason = "ReconcileSuccess"
	ConditionReasonReconcileFailure ConditionReason = "ReconcileFailure"
)

// A Condition that may apply to a resource
type Condition struct {
	// Type of this condition. At most one of each condition type may apply to
	// a resource at any point in time.
	Kind ConditionKind `json:"kind"`

	// Status of this condition; is it currently True, False, or Unknown?
	Status corev1.ConditionStatus `json:"status"`

	// LastTransitionTime is the last time this condition transitioned from one
	// status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`

	// A Reason for this condition's last transition from one status to another.
	Reason ConditionReason `json:"reason"`

	// A Message containing details about this condition's last transition from
	// one status to another, if any.
	// +optional
	Message string `json:"message,omitempty"`
}

// Equal returns true if the condition is identical to the supplied condition,
// ignoring the LastTransitionTime.
func (c Condition) Equal(other Condition) bool {
	return c.Kind == other.Kind &&
		c.Status == other.Status &&
		c.Reason == other.Reason &&
		c.Message == other.Message
}

// WithMessage returns a condition by adding the provided message to existing
// condition.
func (c Condition) WithMessage(msg string) Condition {
	c.Message = msg
	return c
}

// A ConditionedStatus reflects the observed status of a resource. Only
// one condition of each kind may exist.
type ConditionedStatus struct {
	// Conditions of the resource.
	// +optional
	Conditions []Condition `json:"conditions,omitempty"`
}

// NewConditionedStatus returns a stat with the supplied conditions set.
func NewConditionedStatus(c ...Condition) *ConditionedStatus {
	s := &ConditionedStatus{}
	s.SetConditions(c...)
	return s
}

// GetCondition returns the condition for the given ConditionKind if exists,
// otherwise returns nil
func (s *ConditionedStatus) GetCondition(ck ConditionKind) Condition {
	for _, c := range s.Conditions {
		if c.Kind == ck {
			return c
		}
	}
	return Condition{Kind: ck, Status: corev1.ConditionUnknown}
}

// SetConditions sets the supplied conditions, replacing any existing conditions
// of the same kind. This is a no-op if all supplied conditions are identical,
// ignoring the last transition time, to those already set.
func (s *ConditionedStatus) SetConditions(c ...Condition) {
	for _, new := range c {
		exists := false
		for i, existing := range s.Conditions {
			if existing.Kind != new.Kind {
				continue
			}

			if existing.Equal(new) {
				exists = true
				continue
			}

			s.Conditions[i] = new
			exists = true
		}
		if !exists {
			s.Conditions = append(s.Conditions, new)
		}
	}
}

// Equal returns true if the status is identical to the supplied status,
// ignoring the LastTransitionTimes and order of statuses.
func (s *ConditionedStatus) Equal(other *ConditionedStatus) bool {
	if s == nil || other == nil {
		return s == nil && other == nil
	}

	if len(other.Conditions) != len(s.Conditions) {
		return false
	}

	sc := make([]Condition, len(s.Conditions))
	copy(sc, s.Conditions)

	oc := make([]Condition, len(other.Conditions))
	copy(oc, other.Conditions)

	// We should not have more than one condition of each kind.
	sort.Slice(sc, func(i, j int) bool { return sc[i].Kind < sc[j].Kind })
	sort.Slice(oc, func(i, j int) bool { return oc[i].Kind < oc[j].Kind })

	for i := range sc {
		if !sc[i].Equal(oc[i]) {
			return false
		}
	}

	return true
}

// Creating returns a condition that indicates the resource is currently
// being created.
func Creating() Condition {
	return Condition{
		Kind:               ConditionKindConfiguration,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonCreating,
	}
}

// Deleting returns a condition that indicates the resource is currently
// being deleted.
func Deleting() Condition {
	return Condition{
		Kind:               ConditionKindConfiguration,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonDeleting,
	}
}

// ReconcileSuccess returns a condition that indicates the resource is
// currently successfully reconciled
func ReconcileSuccess() Condition {
	return Condition{
		Kind:               ConditionKindConfiguration,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonReconcileSuccess,
	}
}

// Unavailable returns a condition that indicates the resource is not
// currently available for use.
func ReconcileFailure() Condition {
	return Condition{
		Kind:               ConditionKindConfiguration,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonReconcileFailure,
	}
}

// TargetFound returns a condition that indicates the resource has
// target(s) available for use.
func TargetFound() Condition {
	return Condition{
		Kind:               ConditionKindTarget,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonFound,
	}
}

// TargetNotFound returns a condition that indicates the resource has no
// target(s) available for use.
func TargetNotFound() Condition {
	return Condition{
		Kind:               ConditionKindTarget,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonNotFound,
	}
}

// InternalLeafRefValidationSuccess returns a condition that indicates
// the resource leafreference(s) are found or no leafrefs exist
func InternalLeafRefValidationSuccess() Condition {
	return Condition{
		Kind:               ConditionKindInternalLeafRef,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonSuccess,
	}
}

// InternalLeafRefValidationFailure returns a condition that indicates
// the resource leafreference(s) are missing
func InternalLeafRefValidationFailure() Condition {
	return Condition{
		Kind:               ConditionKindInternalLeafRef,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonFailed,
	}
}

// ExternalLeafRefValidationSuccess returns a condition that indicates
// the resource leafreference(s) are found or no leafrefs exist
func ExternalLeafRefValidationSuccess() Condition {
	return Condition{
		Kind:               ConditionKindExternalLeafRef,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonSuccess,
	}
}

// ExternalLeafRefValidationFailure returns a condition that indicates
// the resource leafreference(s) are missing
func ExternalLeafRefValidationFailure() Condition {
	return Condition{
		Kind:               ConditionKindExternalLeafRef,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonFailed,
	}
}

// ParentValidationFailure returns a condition that indicates
// the resource parent is found
func ParentValidationSuccess() Condition {
	return Condition{
		Kind:               ConditionKindParent,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonSuccess,
	}
}

// ExternalLeafRefValidationFailure returns a condition that indicates
// the resource leafreference(s) are missing
func ParentValidationFailure() Condition {
	return Condition{
		Kind:               ConditionKindParent,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonFailed,
	}
}
