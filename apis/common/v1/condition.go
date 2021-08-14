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

	"github.com/netw-device-driver/ndd-runtime/pkg/yang/leafref"
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
	ConditionKindSynced ConditionKind = "Synced"
	// handled per target per resource
	ConditionKindReady ConditionKind = "Ready"
)

// A ConditionReason represents the reason a resource is in a condition.
type ConditionReason string

// Reasons a resource validation is or is not ok
// applicable to leafref validation and target validation
const (
	ConditionReasonSuccess ConditionReason = "Success"
	ConditionReasonFailed  ConditionReason = "Failed"
)

// Reasons a resource is or is not ready
const (
	ConditionReasonUnknown     ConditionReason = "Unknown"
	ConditionReasonCreating    ConditionReason = "Creating"
	ConditionReasonDeleting    ConditionReason = "Deleting"
	ConditionReasonUpdating    ConditionReason = "Updating"
	ConditionReasonUnavailable ConditionReason = "UnAvailable"
	ConditionReasonAvailable   ConditionReason = "Available"
)

// Reasons a resource is or is not synced.
const (
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

	// A Resolved leafref
	ResolvedLeafRefs []*leafref.ResolvedLeafRef `json:"resolvedLeafRefs,omitempty"`

	// ExternalResourceNames for external leafrefs
	ExternalResourceNames []string `json:"externalResourceNames,omitempty"`
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

// WithResolvedLeafRefs returns a condition by adding the provided resolvedLeafRefs 
// to an existing condition.
func (c Condition) WithResolvedLeafRefs(r []*leafref.ResolvedLeafRef) Condition {
	c.ResolvedLeafRefs = r
	return c
}

// WithResolvedLeafRefs returns a condition by adding the provided resolvedLeafRefs 
// to an existing condition.
func (c Condition) WithExternalResourceNames(r []string) Condition {
	c.ExternalResourceNames = r
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

// Unknown returns a condition that indicates the resource is in an
// unknown status.
func Unknown() Condition {
	return Condition{
		Kind:               ConditionKindReady,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonUnknown,
	}
}

// Creating returns a condition that indicates the resource is currently
// being created.
func Creating() Condition {
	return Condition{
		Kind:               ConditionKindReady,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonCreating,
	}
}

// Updating returns a condition that indicates the resource is currently
// being updated.
func Updating() Condition {
	return Condition{
		Kind:               ConditionKindReady,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonUpdating,
	}
}

// Deleting returns a condition that indicates the resource is currently
// being deleted.
func Deleting() Condition {
	return Condition{
		Kind:               ConditionKindReady,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonDeleting,
	}
}

// Available returns a condition that indicates the resource is
// currently observed to be available for use.
func Available() Condition {
	return Condition{
		Kind:               ConditionKindReady,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonAvailable,
	}
}

// Unavailable returns a condition that indicates the resource is not
// currently available for use. Unavailable should be set only when ndd
// expects the resource to be available but knows it is not, for example
// because its API reports it is unhealthy.
func Unavailable() Condition {
	return Condition{
		Kind:               ConditionKindReady,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonUnavailable,
	}
}

// ReconcileSuccess returns a condition indicating that ndd successfully
// completed the most recent reconciliation of the resource.
func ReconcileSuccess() Condition {
	return Condition{
		Kind:               ConditionKindSynced,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonReconcileSuccess,
	}
}

// ReconcileError returns a condition indicating that ndd encountered an
// error while reconciling the resource. This could mean ndd was
// unable to update the resource to reflect its desired state, or that
// ndd was unable to determine the current actual state of the resource.
func ReconcileError(err error) Condition {
	return Condition{
		Kind:               ConditionKindSynced,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonReconcileFailure,
		Message:            err.Error(),
	}
}

// TargetFound returns a condition that indicates the resource has
// target(s) available for use.
func TargetFound() Condition {
	return Condition{
		Kind:               ConditionKindTarget,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonSuccess,
	}
}

// TargetNotFound returns a condition that indicates the resource has no
// target(s) available for use.
func TargetNotFound() Condition {
	return Condition{
		Kind:               ConditionKindTarget,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonFailed,
	}
}

// LocalLeafRefValidationSuccess returns a condition that indicates
// the resource leafreference(s) are found or no leafrefs exist
func LocalLeafRefValidationSuccess() Condition {
	return Condition{
		Kind:               ConditionKindInternalLeafRef,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonSuccess,
	}
}

// LocalLeafRefValidationFailure returns a condition that indicates
// the resource leafreference(s) are missing
func LocalLeafRefValidationFailure() Condition {
	return Condition{
		Kind:               ConditionKindInternalLeafRef,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonFailed,
	}
}

// LocalLeafRefValidationUnknown returns a condition that indicates
// the internal leafref validation is unknown
func LocalLeafRefValidationUnknown() Condition {
	return Condition{
		Kind:               ConditionKindInternalLeafRef,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonUnknown,
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

// ExternalLeafRefValidationUnknown returns a condition that indicates
// the external leafref validation is unknown
func ExternalLeafRefValidationUnknown() Condition {
	return Condition{
		Kind:               ConditionKindExternalLeafRef,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonUnknown,
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

// ParentValidationUnknown returns a condition that indicates
// the parent validation is unknown
func ParentValidationUnknown() Condition {
	return Condition{
		Kind:               ConditionKindParent,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonUnknown,
	}
}
