// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

import (
	configv1 "github.com/openshift/api/config/v1"
)

// NodeSpecApplyConfiguration represents a declarative configuration of the NodeSpec type for use
// with apply.
type NodeSpecApplyConfiguration struct {
	CgroupMode            *configv1.CgroupMode               `json:"cgroupMode,omitempty"`
	WorkerLatencyProfile  *configv1.WorkerLatencyProfileType `json:"workerLatencyProfile,omitempty"`
	MinimumKubeletVersion *string                            `json:"minimumKubeletVersion,omitempty"`
}

// NodeSpecApplyConfiguration constructs a declarative configuration of the NodeSpec type for use with
// apply.
func NodeSpec() *NodeSpecApplyConfiguration {
	return &NodeSpecApplyConfiguration{}
}

// WithCgroupMode sets the CgroupMode field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CgroupMode field is set to the value of the last call.
func (b *NodeSpecApplyConfiguration) WithCgroupMode(value configv1.CgroupMode) *NodeSpecApplyConfiguration {
	b.CgroupMode = &value
	return b
}

// WithWorkerLatencyProfile sets the WorkerLatencyProfile field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the WorkerLatencyProfile field is set to the value of the last call.
func (b *NodeSpecApplyConfiguration) WithWorkerLatencyProfile(value configv1.WorkerLatencyProfileType) *NodeSpecApplyConfiguration {
	b.WorkerLatencyProfile = &value
	return b
}

// WithMinimumKubeletVersion sets the MinimumKubeletVersion field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MinimumKubeletVersion field is set to the value of the last call.
func (b *NodeSpecApplyConfiguration) WithMinimumKubeletVersion(value string) *NodeSpecApplyConfiguration {
	b.MinimumKubeletVersion = &value
	return b
}
