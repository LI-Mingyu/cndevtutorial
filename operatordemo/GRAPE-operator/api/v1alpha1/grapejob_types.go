/*
Copyright 2021.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GrapeJobSpec defines the desired state of GrapeJob
type GrapeJobSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// alpha1: Parallelism、 AppExec、 AppArgs
	// alpha2: ttl
	// alpha3: Volumes、 VolumeMounts
	// alpha4: Resources requests and limits
	// ...: FailedJobsHistoryLimit, SuccessfulJobsHistoryLimit ...

	// +optional
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:default:=1
	Parallelism *int32 `json:"parallelism,omitempty"`

	AppExec string   `json:"appExec"`
	AppArgs []string `json:"appArgs"`
}

// GrapeJobStatus defines the observed state of GrapeJob
type GrapeJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	JobStatus JobStatus `json:"jobStatus"`
}
type JobStatus string

const (
	Pending   JobStatus = "Pending"
	Running   JobStatus = "Running"
	Failed    JobStatus = "Failed"
	Completed JobStatus = "Completed"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.jobStatus`
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// GrapeJob is the Schema for the grapejobs API
type GrapeJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GrapeJobSpec   `json:"spec,omitempty"`
	Status GrapeJobStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GrapeJobList contains a list of GrapeJob
type GrapeJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GrapeJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GrapeJob{}, &GrapeJobList{})
}
