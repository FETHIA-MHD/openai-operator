package controllers

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	oaiv1 "github.com/FETHIA-MHD/openai-operator/api/v1"
)

// OaiCoreReconciler réconcilie un objet OaiCore
type OaiCoreReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=oai.openai.com,resources=oaicores,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=oai.openai.com,resources=oaicores/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=oai.openai.com,resources=oaicores/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile est la fonction principale qui gère la logique de votre opérateur
func (r *OaiCoreReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Récupérer l'objet OaiCore
	oaiCore := &oaiv1.OaiCore{}
	if err := r.Get(ctx, req.NamespacedName, oaiCore); err != nil {
		if errors.IsNotFound(err) {
			log.Info("OaiCore resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get OaiCore.")
		return ctrl.Result{}, err
	}

	// Créer les Deployments pour AMF, SMF, UPF, et NRF
	if err := r.createDeployment(ctx, oaiCore, "amf", oaiCore.Spec.AMF.Replicas, oaiCore.Spec.AMF.Image); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.createDeployment(ctx, oaiCore, "smf", oaiCore.Spec.SMF.Replicas, oaiCore.Spec.SMF.Image); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.createDeployment(ctx, oaiCore, "upf", oaiCore.Spec.UPF.Replicas, oaiCore.Spec.UPF.Image); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.createDeployment(ctx, oaiCore, "nrf", oaiCore.Spec.NRF.Replicas, oaiCore.Spec.NRF.Image); err != nil {
		return ctrl.Result{}, err
	}

	// Mettre à jour le statut de l'OaiCore
	oaiCore.Status.AMFReady = true
	oaiCore.Status.SMFReady = true
	oaiCore.Status.UPFReady = true
	oaiCore.Status.NRFReady = true
	if err := r.Status().Update(ctx, oaiCore); err != nil {
		log.Error(err, "Failed to update OaiCore status.")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// createDeployment crée un Deployment pour un composant donné
func (r *OaiCoreReconciler) createDeployment(ctx context.Context, oaiCore *oaiv1.OaiCore, component string, replicas int32, image string) error {
	log := log.FromContext(ctx)

	labels := map[string]string{
		"app":        "oaicore",
		"component":  component,
		"controller": oaiCore.Name,
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      oaiCore.Name + "-" + component,
			Namespace: oaiCore.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  component,
							Image: image,
						},
					},
				},
			},
		},
	}

	if err := r.Create(ctx, deployment); err != nil {
		if errors.IsAlreadyExists(err) {
			log.Info("Deployment already exists.", "Deployment", deployment.Name)
			return nil
		}
		log.Error(err, "Failed to create Deployment.", "Deployment", deployment.Name)
		return err
	}

	log.Info("Created Deployment.", "Deployment", deployment.Name)
	return nil
}

// SetupWithManager configure le contrôleur avec le Manager
func (r *OaiCoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&oaiv1.OaiCore{}).
		Complete(r)
}
