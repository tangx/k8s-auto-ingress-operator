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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/sirupsen/logrus"
	myappv1 "github.com/tangx/k8s-auto-ingress-operator/api/v1"
	"github.com/tangx/k8s-auto-ingress-operator/controllers/helper"
	"github.com/tangx/k8s-auto-ingress-operator/controllers/util"
)

// AutoIngressReconciler reconciles a AutoIngress object
type AutoIngressReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=network.sodev.cc,resources=autoingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=network.sodev.cc,resources=autoingresses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=network.sodev.cc,resources=autoingresses/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources="services",verbs=get;list;watch
//+kubebuilder:rbac:groups="networking.k8s.io",resources="ingresses",verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AutoIngress object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *AutoIngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	logrus.Info("进入 Reconcile")
	defer logrus.Info("退出 Reconcile")

	// TODO(user): your logic here
	op := &myappv1.AutoIngress{}

	err := r.Client.Get(ctx, req.NamespacedName, op)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 删除
	if !op.DeletionTimestamp.IsZero() {
		autoIngressSet.Remove(*op)

		return ctrl.Result{}, nil
	}

	// 保存
	if len(op.Spec.ServicePrefixes) == 0 {
		op.Spec.ServicePrefixes = []string{"web-", "srv-"}
	}
	autoIngressSet.Add(*op)

	r.ReconcileServices(ctx, op)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AutoIngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&myappv1.AutoIngress{}).
		Watches(
			&source.Kind{
				Type: &corev1.Service{},
			},
			handler.Funcs{
				CreateFunc: r.onCreateService,
			},
		).
		Complete(r)
}

func (r *AutoIngressReconciler) onCreateService(e event.CreateEvent, q workqueue.RateLimitingInterface) {
	logrus.Info("新 Service 被创建")
	svc := r.getService(e.Object)
	if svc == nil {
		return
	}

	ctx := context.TODO()
	for _, op := range autoIngressSet.List() {
		r.HandleIngress(ctx, op, svc)
	}
}

// isExistInK8s 检查对象是否在 k8s 中存在
func (r *AutoIngressReconciler) isExistInK8s(obj client.Object) bool {

	key := r.objectKey(obj)
	err := r.Client.Get(context.TODO(), key, obj)
	if err != nil {
		return false
	}

	return true
}

func (r *AutoIngressReconciler) objectKey(e client.Object) types.NamespacedName {
	return types.NamespacedName{
		Namespace: e.GetNamespace(),
		Name:      e.GetName(),
	}
}

func (r *AutoIngressReconciler) getService(e client.Object) *corev1.Service {

	key := r.objectKey(e)
	svc := &corev1.Service{}

	err := r.Client.Get(context.TODO(), key, svc)
	if err != nil {
		return nil
	}

	return svc
}

func (r *AutoIngressReconciler) HandleIngress(ctx context.Context, op myappv1.AutoIngress, svc *corev1.Service) {

	if !util.IsValidServcieName(svc.Name, op.Spec.ServicePrefixes) {
		return
	}

	ing := helper.NewIngress(op, svc)
	_ing := helper.NewIngress(op, svc)
	action := "create"

	if r.isExistInK8s(_ing) {
		action = "update"
		ing.SetResourceVersion(_ing.ResourceVersion)
	}

	_ = controllerutil.SetOwnerReference(svc, ing, r.Scheme)
	_ = controllerutil.SetOwnerReference(&op, ing, r.Scheme)

	err := r.HandleObject(ctx, ing, action)
	if err != nil {
		logrus.Errorf("handle(%s) ingress (%s.%s) failed: %v", action, ing.Name, ing.Namespace, err)
		return
	}

	logrus.Infof("handle(%s) ingress (%s.%s) success", action, ing.Name, ing.Namespace)
}

func (r *AutoIngressReconciler) HandleObject(ctx context.Context, obj client.Object, action string) error {

	switch action {
	case "update":
		return r.Client.Update(ctx, obj)
	case "create":
		return r.Client.Create(ctx, obj)
	}

	return nil
}

func (r *AutoIngressReconciler) ReconcileServices(ctx context.Context, op *myappv1.AutoIngress) {
	svcs := &corev1.ServiceList{}
	err := r.Client.List(ctx, svcs)
	if err != nil {
		logrus.Errorf("list services failed: %v", err)
		return
	}

	for _, svc := range svcs.Items {
		r.HandleIngress(ctx, *op, &svc)
	}
}
