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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/davecgh/go-spew/spew"
	pequodv1alpha1 "github.com/soal-one/pequod/v2/api/v1alpha1"
)

// ClaimReconciler reconciles a Claim object
type ClaimReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=pequod.soal.one,resources=claims,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pequod.soal.one,resources=claims/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pequod.soal.one,resources=claims/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Claim object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *ClaimReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	_ = log.FromContext(ctx)

	var claim = pequodv1alpha1.Claim{}

	if err := r.Get(ctx, req.NamespacedName, &claim); err != nil {
		logger.Info(fmt.Sprintf("unable to fetch claim %s, assuming it has been deleted.", req.NamespacedName))
	}

	var claimList = pequodv1alpha1.ClaimList{}

	if err := r.List(ctx, &claimList, client.InNamespace(req.Namespace)); err != nil {
		logger.Error(err, "unable to list claims")
	}

	// Ensure concretion type exists
	var concretionList = pequodv1alpha1.ConcretionList{}

	if err := r.List(ctx, &concretionList, client.MatchingFields{".metadata.name": claim.Spec.Type}); err != nil {
		logger.Error(err, "Unable to list concretions")
	}
	spew.Dump(claim)
	spew.Dump(concretionList)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClaimReconciler) SetupWithManager(mgr ctrl.Manager) error {

	claimGVStr := pequodv1alpha1.GroupVersion.String()

	if err := mgr.GetFieldIndexer().IndexField(context.Background(),
		&pequodv1alpha1.Claim{},
		".metadata.name",
		func(rawObj client.Object) []string {
			claim := rawObj.(*pequodv1alpha1.Claim)
			name := claim.ObjectMeta.Name
			if name == "" {
				fmt.Println("REEE")
				return nil
			}

			// TODO Probably superfluous
			if claim.APIVersion != claimGVStr || claim.Kind != "Claim" {
				fmt.Println("REEEEEEEEEEEE")
				return nil
			}

			return []string{name}
		}); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&pequodv1alpha1.Claim{}).
		Complete(r)
}
