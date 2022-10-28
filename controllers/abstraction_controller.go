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

	"cuelang.org/go/cue/cuecontext"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	pequodv1alpha1 "github.com/soal-one/pequod/v2/api/v1alpha1"
)

// AbstractionReconciler reconciles a Abstraction object
type AbstractionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=pequod.soal.one,resources=abstractions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pequod.soal.one,resources=abstractions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pequod.soal.one,resources=abstractions/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Abstraction object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *AbstractionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	logger := log.FromContext(ctx)
	abstraction := pequodv1alpha1.Abstraction{}

	if err := r.Get(ctx, req.NamespacedName, &abstraction); err != nil {
		logger.Info(fmt.Sprintf("unable to fetch abstraction %s, assuming it has been deleted.", req.NamespacedName))
	}

	var absList = pequodv1alpha1.AbstractionList{}

	if err := r.List(ctx, &absList, client.InNamespace(req.Namespace)); err != nil {
		logger.Error(err, "unable to list abstractions")
	}

	qctx := cuecontext.New()
	val := qctx.CompileString(abstraction.Spec.Template)
	err := val.Validate()

	if err != nil {
		logger.Error(err, "Error with Abstraction Template")
	}

	// rt.Breakpoint()

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AbstractionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pequodv1alpha1.Abstraction{}).
		Complete(r)
}
