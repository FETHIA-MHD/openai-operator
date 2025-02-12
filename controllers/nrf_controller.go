package controllers

import (
	"context"
	oaiv1 "github.com/FETHIA-MHD/openai-operator/api/v1"
	"k8s.io/api/apps/v1"                   // Pour le package apps/v1 (Déploiement)
	"k8s.io/api/core/v1"                   // Pour le package core/v1 (Conteneurs et PodSpec)
	"k8s.io/apimachinery/pkg/apis/meta/v1" // Pour le package meta/v1 (Métadonnées)
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// createNRFDeployment crée un déploiement pour NRF
func (r *OaiCoreReconciler) createNRFDeployment(ctx context.Context, oaiCore *oaiv1.OaiCore) error {
	log := log.FromContext(ctx)

	labels := map[string]string{
		"app":        "oaicore",
		"component":  "nrf",
		"controller": oaiCore.Name,
	}

	// Configuration NRF
	nrfConfig := map[string]string{
		"host":           "oai-nrf",
		"sbi_port":       "80",
		"api_version":    "v1",
		"interface_name": "eth0",
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      oaiCore.Name + "-nrf",
			Namespace: oaiCore.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &oaiCore.Spec.NRF.Replicas,
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
							Name:  "nrf",
							Image: oaiCore.Spec.NRF.Image,
							Env: []corev1.EnvVar{
								{
									Name:  "NRF_HOST",
									Value: nrfConfig["host"],
								},
								{
									Name:  "NRF_SBI_PORT",
									Value: nrfConfig["sbi_port"],
								},
								{
									Name:  "NRF_API_VERSION",
									Value: nrfConfig["api_version"],
								},
								{
									Name:  "NRF_INTERFACE_NAME",
									Value: nrfConfig["interface_name"],
								},
							},
						},
					},
				},
			},
		},
	}

	if err := r.Create(ctx, deployment); err != nil {
		if errors.IsAlreadyExists(err) {
			log.Info("NRF Deployment already exists.", "Deployment", deployment.Name)
			return nil
		}
		log.Error(err, "Failed to create NRF Deployment.", "Deployment", deployment.Name)
		return err
	}

	log.Info("Created NRF Deployment.", "Deployment", deployment.Name)
	return nil
}
