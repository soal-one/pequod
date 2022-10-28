/*
Copyright 2022 soal.one.

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

	"cuelang.org/go/cue/cuecontext"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	pequodv1alpha1 "github.com/soal-one/pequod/v2/api/v1alpha1"
)

// ConcretionReconciler reconciles a Concretion object
type ConcretionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=pequod.soal.one,resources=concretions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pequod.soal.one,resources=concretions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pequod.soal.one,resources=concretions/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Concretion object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *ConcretionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	concretion := pequodv1alpha1.Concretion{}

	if err := r.Get(ctx, req.NamespacedName, &concretion); err != nil {
		logger.Error(err, "Concretion was likely deleted...")
	}

	conList := pequodv1alpha1.ConcretionList{}

	if err := r.List(ctx, &conList, client.InNamespace(req.Namespace)); err != nil {
		logger.Error(err, "unable to list concretions")
	}
	qctx := cuecontext.New()
	val := qctx.CompileString(concretion.Spec.Template)
	err := val.Validate()

	if err != nil {
		logger.Error(err, "Error with Concretion Template")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConcretionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	claimGVStr := pequodv1alpha1.GroupVersion.String()

	if err := mgr.GetFieldIndexer().IndexField(context.Background(),
		&pequodv1alpha1.Concretion{},
		".metadata.name",
		func(rawObj client.Object) []string {
			concretion := rawObj.(*pequodv1alpha1.Concretion)
			name := concretion.ObjectMeta.Name
			if name == "" {
				return nil
			}

			// TODO Probably superfluous
			if concretion.APIVersion != claimGVStr || concretion.Kind != "Concretion" {
				return nil
			}

			return []string{name}
		}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&pequodv1alpha1.Concretion{}).
		Complete(r)
}
