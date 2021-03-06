/*
Copyright 2020.

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

package managementseldonio

import (
	"context"

	"github.com/go-logr/logr"
    "k8s.io/apimachinery/pkg/types"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	managementseldoniov1alpha1 "github.com/SeldonIO/mlops-operator/apis/management.seldon.io/v1alpha1"
)

// SeldonClusterReconciler reconciles a SeldonCluster object
type SeldonClusterReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=management.seldon.io,resources=seldonclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=management.seldon.io,resources=seldonclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=management.seldon.io,resources=seldonclusters/finalizers,verbs=update
// +kubebuilder:rbac:groups=machinelearning.seldon.io,resources=seldondeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=machinelearning.seldon.io,resources=seldondeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=machinelearning.seldon.io,resources=seldondeployments/finalizers,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=v1,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=v1,resources=services/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=networking.istio.io,resources=virtualservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.istio.io,resources=virtualservices/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=networking.istio.io,resources=destinationrules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.istio.io,resources=destinationrules/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=v1,resources=namespaces,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *SeldonClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("seldoncluster", req.NamespacedName)

	instance := &managementseldoniov1alpha1.SeldonCluster{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("SeldonCluster resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get SeldonCluster")
		return ctrl.Result{}, err
	}

    // Seldon Deploy cluster section
    sdName := r.getSeldonDeployName(instance)
    found := &appsv1.Deployment{}
    err = r.Get(ctx, types.NamespacedName{Name: sdName, Namespace: instance.Namespace}, found)
    if err != nil && errors.IsNotFound(err) {
        // Define a new deployment
        dep := r.getDeploymentForSeldonDeploy(instance)
        log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "DEployment.Name", dep.Name)
        err = r.Create(ctx, dep)
        if err != nil {
            log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
            return ctrl.Result{}, err
        }
        // Deployment created successfully - return and requeue
        return ctrl.Result{Requeue: true}, nil
    } else if err != nil {
        log.Error(err, "Failed to get Deployment")
        return ctrl.Result{}, err
    }

    // Seldon Deploy cluster section

	return ctrl.Result{}, nil
}

func (r *SeldonClusterReconciler) getDeploymentForSeldonDeploy(instance *managementseldoniov1alpha1.SeldonCluster) *appsv1.Deployment {
    return &appsv1.Deployment{}
}

func (r *SeldonClusterReconciler) getSeldonDeployName(instance *managementseldoniov1alpha1.SeldonCluster) string {
    return instance.Name + "-deploy"
}

// SetupWithManager sets up the controller with the Manager.
func (r *SeldonClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&managementseldoniov1alpha1.SeldonCluster{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
