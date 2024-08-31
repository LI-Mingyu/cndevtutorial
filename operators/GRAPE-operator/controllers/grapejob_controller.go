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

package controllers

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kbatchv1 "k8s.io/api/batch/v1"
	kcorev1 "k8s.io/api/core/v1"
	kmetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	grapeJobv1alpha2 "limingyu.io/GRAPE-operator/api/v1alpha2"
)

// GrapeJobReconciler reconciles a GrapeJob object
type GrapeJobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=graph.limingyu.io,resources=grapejobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=graph.limingyu.io,resources=grapejobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=graph.limingyu.io,resources=grapejobs/finalizers,verbs=update

//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GrapeJob object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *GrapeJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// TODO(user): your logic here
	//var logger = log.Log.WithValues("GrapeJob", req.NamespacedName)
	logger.Info("In the reconcile loop.")

	var grapejob grapeJobv1alpha2.GrapeJob
	if err := r.Get(ctx, req.NamespacedName, &grapejob); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.Info("Resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get GrapeJob")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var foundCtlJob kbatchv1.Job
	var foundCtlSvc kcorev1.Service
	var foundWorkersJob kbatchv1.Job

	//获取该 GrapeJob 对应的主控 Job，若不存在，则创建它
	if err := r.Get(ctx, types.NamespacedName{Namespace: grapejob.Namespace, Name: nameForCtl(grapejob.Name)}, &foundCtlJob); err != nil {
		if errors.IsNotFound(err) {
			// 根据获取的 GrapeJob，构造主控 Job
			ctlJob := r.constructCtlJob(&grapejob)
			// 创建主控 Job
			if err := r.Create(ctx, ctlJob); err != nil && !errors.IsAlreadyExists(err) {
				logger.Error(err, "Unable to create controller job for GrapeJob", "Job", ctlJob)
				return ctrl.Result{}, err //若主控 Job 创建失败，而且失败原因也不是“资源已存在”，报错并返回错误
			}
			logger.Info("Controller job created", "job", ctlJob)
			return ctrl.Result{Requeue: true}, nil //主控 Job 创建成功，重新排进队列以便后续操作
		}
		logger.Error(err, "Failed to get controller job")
		return ctrl.Result{}, err //主控 Job 获取失败，且原因也不是“该资源不存在”，报错并返回错误
	}

	//获取该 GrapeJob 对应的主控 Service，若不存在，则创建它
	if err := r.Get(ctx, types.NamespacedName{Namespace: grapejob.Namespace, Name: nameForCtl(grapejob.Name)}, &foundCtlSvc); err != nil {
		if errors.IsNotFound(err) {
			ctlSvc := r.constructCtlSvc(&grapejob)
			if err := r.Create(ctx, ctlSvc); err != nil && !errors.IsAlreadyExists(err) {
				logger.Error(err, "Unable to create controller service for GrapeJob", "Service", ctlSvc)
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}
		logger.Error(err, "Failed to get controller service")
		return ctrl.Result{}, err
	}

	//获取该 GrapeJob 对应的 WorkersJob，若不存在，则创建它
	if err := r.Get(ctx, types.NamespacedName{Namespace: grapejob.Namespace, Name: nameForWorkersJob(grapejob.Name)}, &foundWorkersJob); err != nil {
		if errors.IsNotFound(err) {
			workersJob := r.constructWorkersJob(&grapejob)
			if err := r.Create(ctx, workersJob); err != nil && !errors.IsAlreadyExists(err) {
				logger.Error(err, "Unable to create workers job for GrapeJob", "Job", workersJob)
				return ctrl.Result{}, err
			}
			logger.Info("Workers job created", "job", workersJob)
			return ctrl.Result{Requeue: true}, nil
		}
		logger.Error(err, "Failed to get workers job")
		return ctrl.Result{}, err
	}

	//更新 status
	if foundWorkersJob.Status.Failed+foundCtlJob.Status.Failed > 0 {
		grapejob.Status.JobStatus = grapeJobv1alpha2.Failed
	} else if foundWorkersJob.Status.Succeeded == *grapejob.Spec.Parallelism && foundCtlJob.Status.Succeeded == 1 {
		grapejob.Status.JobStatus = grapeJobv1alpha2.Completed
	} else {
		grapejob.Status.JobStatus = grapeJobv1alpha2.Running
	}
	if err := r.Status().Update(ctx, &grapejob); err != nil {
		logger.Error(err, "Failed to update GrapeJob status")
		return ctrl.Result{}, err
	}

	if *grapejob.Spec.TTL > 0 {
		ttl_left := grapejob.CreationTimestamp.Add(time.Second * time.Duration(*grapejob.Spec.TTL)).Sub(time.Now())
		if ttl_left > 0 {
			return ctrl.Result{RequeueAfter: ttl_left}, nil
		}
		if ttl_left < 0 {
			r.Delete(ctx, &grapejob)
		}
	}

	return ctrl.Result{}, nil
}

// 主控 Job 和 Service 的名称
func nameForCtl(grapeJobName string) string {
	return grapeJobName + "-controller"
}

//主控 Pod 的 Labels，也是 Service 的 Selector
func labelsForCtlPod(grapeJobName string) map[string]string {
	return map[string]string{"app": grapeJobName, "type": "controller"}
}

func nameForWorkersJob(grapeJobName string) string {
	return grapeJobName + "-workers"
}

// 构造主控Job
func (r *GrapeJobReconciler) constructCtlJob(gj *grapeJobv1alpha2.GrapeJob) *kbatchv1.Job {
	name := nameForCtl(gj.Name)
	controllerPodLabels := labelsForCtlPod(gj.Name)

	// 构造主控进程的命令行
	hostsIDs := "0"
	for i := 1; i < int(*gj.Spec.Parallelism); i++ {
		hostsIDs = fmt.Sprintf("%s, %d", hostsIDs, i)
	}
	cmdArgs := append([]string{"-launcher", "manual", "-verbose", "-disable-hostname-propagation", "-hosts", hostsIDs,
		gj.Spec.AppExec}, gj.Spec.AppArgs...)

	// 构造主控 Job
	job := &kbatchv1.Job{
		ObjectMeta: kmetav1.ObjectMeta{
			Name:      name,
			Namespace: gj.Namespace,
		},
		Spec: kbatchv1.JobSpec{
			Template: kcorev1.PodTemplateSpec{
				ObjectMeta: kmetav1.ObjectMeta{
					Labels: controllerPodLabels,
				},
				Spec: kcorev1.PodSpec{
					Containers: []kcorev1.Container{{
						// 目前这里的镜像名是写死的，这个后面要改一下
						Image:   "limingyu007/run_app:v1.1",
						Name:    "grape-mpi-controller",
						Command: []string{"mpirun"},
						Args:    cmdArgs,
						Env: []kcorev1.EnvVar{{
							Name:  "MPICH_PORT_RANGE",
							Value: "20000:20100",
						}},
						Ports: []kcorev1.ContainerPort{{
							ContainerPort: 20000,
						}},
					}},
					RestartPolicy: "Never",
				},
			},
		},
	}
	ctrl.SetControllerReference(gj, job, r.Scheme)

	return job
}

// 构造主控Job对应的Service（worker要根据这个service name寻址主控）
func (r *GrapeJobReconciler) constructCtlSvc(gj *grapeJobv1alpha2.GrapeJob) *kcorev1.Service {
	name := nameForCtl(gj.Name)
	controllerPodLabels := labelsForCtlPod(gj.Name)

	// 构造主控 Pod 对应的 Service
	svc := &kcorev1.Service{
		ObjectMeta: kmetav1.ObjectMeta{
			Name:      name,
			Namespace: gj.Namespace,
		},
		Spec: kcorev1.ServiceSpec{
			Selector: controllerPodLabels,
			Ports: []kcorev1.ServicePort{{
				Port: 20000,
			}},
		},
	}
	ctrl.SetControllerReference(gj, svc, r.Scheme)

	return svc
}

// 构造 Workers Job
func (r *GrapeJobReconciler) constructWorkersJob(gj *grapeJobv1alpha2.GrapeJob) *kbatchv1.Job {
	name := nameForWorkersJob(gj.Name)
	completionMode := "Indexed"
	cmdArgs := []string{"-c", "hydra_pmi_proxy --control-port " + gj.Name +
		"-controller:20000 --debug --demux poll --pgid 0 --retries 10 --usize -2 --proxy-id $JOB_COMPLETION_INDEX"}

	// 构造 Workers Job
	job := &kbatchv1.Job{
		ObjectMeta: kmetav1.ObjectMeta{
			Name:      name,
			Namespace: gj.Namespace,
		},
		Spec: kbatchv1.JobSpec{
			CompletionMode: (*kbatchv1.CompletionMode)(&completionMode),
			Completions:    gj.Spec.Parallelism,
			Parallelism:    gj.Spec.Parallelism,
			Template: kcorev1.PodTemplateSpec{
				Spec: kcorev1.PodSpec{
					Containers: []kcorev1.Container{{
						// 目前这里的镜像名是写死的，这个后面要改一下
						Image:   "limingyu007/run_app:v1.1",
						Name:    "grape-mpi-worker",
						Command: []string{"bash"},
						Args:    cmdArgs,
					}},
					RestartPolicy: "Never",
				},
			},
		},
	}
	ctrl.SetControllerReference(gj, job, r.Scheme)

	return job
}

// SetupWithManager sets up the controller with the Manager.
func (r *GrapeJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&grapeJobv1alpha2.GrapeJob{}).
		Owns(&kbatchv1.Job{}).
		Owns(&kcorev1.Service{}).
		Complete(r)
}
