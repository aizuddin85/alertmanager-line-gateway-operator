package linemessengergateway

import (
	"context"

	redhatcopv1alpha1 "github.com/aizuddin85/alertmanager-line-gateway-operator/pkg/apis/redhatcop/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	//"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_linemessengergateway")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new LineMessengerGateway Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileLineMessengerGateway{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("linemessengergateway-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource LineMessengerGateway
	err = c.Watch(&source.Kind{Type: &redhatcopv1alpha1.LineMessengerGateway{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner LineMessengerGateway
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &redhatcopv1alpha1.LineMessengerGateway{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileLineMessengerGateway implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileLineMessengerGateway{}

// ReconcileLineMessengerGateway reconciles a LineMessengerGateway object
type ReconcileLineMessengerGateway struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a LineMessengerGateway object and makes changes based on the state read
// and what is in the LineMessengerGateway.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileLineMessengerGateway) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling LineMessengerGateway")

	// Fetch the LineMessengerGateway instance
	instance := &redhatcopv1alpha1.LineMessengerGateway{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new Deployment object
	deploymentconf := newDeploymentForCR(instance)

	// Set LineMessengerGateway instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, deploymentconf, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Deployment already exists
	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: deploymentconf.Name, Namespace: deploymentconf.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new deployment object", "Deployname namespace", deploymentconf.Namespace, "Deployment name", deploymentconf.Name)
		err = r.client.Create(context.TODO(), deploymentconf)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Deployment created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Deployment already exists - don't requeue
	reqLogger.Info("Skip reconcile: Deployment already exists", "Deployment namespace", found.Namespace, "Deployment name", found.Name)
	return reconcile.Result{}, nil
}

// newDeploymentForCR returns a deployment config
func newDeploymentForCR(cr *redhatcopv1alpha1.LineMessengerGateway) *appsv1.Deployment {
    labels := map[string]string{
		"apps": cr.Name,
	}

	DeploymentPodSpec := corev1.PodSpec {
		Containers: []corev1.Container{{
			Image: cr.Spec.Image,
			Name: cr.Name,
			Ports: []corev1.ContainerPort{{
				ContainerPort:	8080,
				Name:			"http",
			},{
				ContainerPort:	8443,
				Name:			"https",
			}},
		}},
	}
	var RepSize int32 = cr.Spec.Size

	DeploymentConf := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: appsv1.SchemeGroupVersion.String(),
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: cr.Name,
			Namespace: cr.Namespace,
			Labels: labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &RepSize,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: cr.Name,
					Namespace: cr.Namespace,
					Labels: labels,
				},
				Spec: DeploymentPodSpec,
			},
		},
	}
	return DeploymentConf
}

	