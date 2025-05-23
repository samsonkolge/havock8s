//go:build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Havock8sExperiment) DeepCopyInto(out *Havock8sExperiment) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Havock8sExperiment.
func (in *Havock8sExperiment) DeepCopy() *Havock8sExperiment {
	if in == nil {
		return nil
	}
	out := new(Havock8sExperiment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Havock8sExperiment) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Havock8sExperimentList) DeepCopyInto(out *Havock8sExperimentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Havock8sExperiment, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Havock8sExperimentList.
func (in *Havock8sExperimentList) DeepCopy() *Havock8sExperimentList {
	if in == nil {
		return nil
	}
	out := new(Havock8sExperimentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Havock8sExperimentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Havock8sExperimentSpec) DeepCopyInto(out *Havock8sExperimentSpec) {
	*out = *in
	in.Target.DeepCopyInto(&out.Target)
	if in.Parameters != nil {
		in, out := &in.Parameters, &out.Parameters
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Schedule != nil {
		in, out := &in.Schedule, &out.Schedule
		*out = new(ScheduleSpec)
		**out = **in
	}
	if in.Safety != nil {
		in, out := &in.Safety, &out.Safety
		*out = new(SafetySpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Havock8sExperimentSpec.
func (in *Havock8sExperimentSpec) DeepCopy() *Havock8sExperimentSpec {
	if in == nil {
		return nil
	}
	out := new(Havock8sExperimentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Havock8sExperimentStatus) DeepCopyInto(out *Havock8sExperimentStatus) {
	*out = *in
	if in.StartTime != nil {
		in, out := &in.StartTime, &out.StartTime
		*out = (*in).DeepCopy()
	}
	if in.EndTime != nil {
		in, out := &in.EndTime, &out.EndTime
		*out = (*in).DeepCopy()
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.TargetResources != nil {
		in, out := &in.TargetResources, &out.TargetResources
		*out = make([]TargetResourceStatus, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Havock8sExperimentStatus.
func (in *Havock8sExperimentStatus) DeepCopy() *Havock8sExperimentStatus {
	if in == nil {
		return nil
	}
	out := new(Havock8sExperimentStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HealthCheckSpec) DeepCopyInto(out *HealthCheckSpec) {
	*out = *in
	if in.Command != nil {
		in, out := &in.Command, &out.Command
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HealthCheckSpec.
func (in *HealthCheckSpec) DeepCopy() *HealthCheckSpec {
	if in == nil {
		return nil
	}
	out := new(HealthCheckSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PauseConditionSpec) DeepCopyInto(out *PauseConditionSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PauseConditionSpec.
func (in *PauseConditionSpec) DeepCopy() *PauseConditionSpec {
	if in == nil {
		return nil
	}
	out := new(PauseConditionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProtectionSpec) DeepCopyInto(out *ProtectionSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProtectionSpec.
func (in *ProtectionSpec) DeepCopy() *ProtectionSpec {
	if in == nil {
		return nil
	}
	out := new(ProtectionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SafetySpec) DeepCopyInto(out *SafetySpec) {
	*out = *in
	if in.HealthChecks != nil {
		in, out := &in.HealthChecks, &out.HealthChecks
		*out = make([]HealthCheckSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.PauseConditions != nil {
		in, out := &in.PauseConditions, &out.PauseConditions
		*out = make([]PauseConditionSpec, len(*in))
		copy(*out, *in)
	}
	if in.ResourceProtections != nil {
		in, out := &in.ResourceProtections, &out.ResourceProtections
		*out = make([]ProtectionSpec, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SafetySpec.
func (in *SafetySpec) DeepCopy() *SafetySpec {
	if in == nil {
		return nil
	}
	out := new(SafetySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScheduleSpec) DeepCopyInto(out *ScheduleSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScheduleSpec.
func (in *ScheduleSpec) DeepCopy() *ScheduleSpec {
	if in == nil {
		return nil
	}
	out := new(ScheduleSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TargetResourceStatus) DeepCopyInto(out *TargetResourceStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TargetResourceStatus.
func (in *TargetResourceStatus) DeepCopy() *TargetResourceStatus {
	if in == nil {
		return nil
	}
	out := new(TargetResourceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TargetSpec) DeepCopyInto(out *TargetSpec) {
	*out = *in
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TargetSpec.
func (in *TargetSpec) DeepCopy() *TargetSpec {
	if in == nil {
		return nil
	}
	out := new(TargetSpec)
	in.DeepCopyInto(out)
	return out
}
