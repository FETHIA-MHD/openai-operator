package controllers

import (
	"context"
	oaiv1 "github.com/FETHIA-MHD/openai-operator/api/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// OaiCoreReconciler contient la logique principale pour le réconciliateur OaiCore
type OaiCoreReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile est la fonction principale qui gère la logique de réconciliation
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
	if err := r.createAMFDeployment(ctx, oaiCore); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.createSMFDeployment(ctx, oaiCore); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.createUPFDeployment(ctx, oaiCore); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.createNRFDeployment(ctx, oaiCore); err != nil {
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

// SetupWithManager enregistre le contrôleur avec le Manager
func (r *OaiCoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&oaiv1.OaiCore{}).
		Complete(r)
}
