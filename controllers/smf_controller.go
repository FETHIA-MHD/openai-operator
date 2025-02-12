package controllers

import (
	"context"
	oaiv1 "github.com/FETHIA-MHD/openai-operator/api/v1"
	"k8s.io/api/apps/v1"                   // Pour le package apps/v1 (Déploiement)
	"k8s.io/api/core/v1"                   // Pour le package core/v1 (Conteneurs et PodSpec)
	"k8s.io/apimachinery/pkg/apis/meta/v1" // Pour le package meta/v1 (Métadonnées)
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// createSMFDeployment crée un déploiement pour SMF
func (r *OaiCoreReconciler) createSMFDeployment(ctx context.Context, oaiCore *oaiv1.OaiCore) error {
	log := log.FromContext(ctx)

	labels := map[string]string{
		"app":        "oaicore",
		"component":  "smf",
		"controller": oaiCore.Name,
	}

	// Configuration SMF
	smfConfig := map[string]string{
		"host":           "oai-smf",
		"sbi_port":       "80",
		"api_version":    "v1",
		"interface_name": "eth0",
		"n4_port":        "8805",
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      oaiCore.Name + "-smf",
			Namespace: oaiCore.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &oaiCore.Spec.SMF.Replicas,
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
							Name:  "smf",
							Image: oaiCore.Spec.SMF.Image,
							Env: []corev1.EnvVar{
								{
									Name:  "SMF_HOST",
									Value: smfConfig["host"],
								},
								{
									Name:  "SMF_SBI_PORT",
									Value: smfConfig["sbi_port"],
								},
								{
									Name:  "SMF_API_VERSION",
									Value: smfConfig["api_version"],
								},
								{
									Name:  "SMF_INTERFACE_NAME",
									Value: smfConfig["interface_name"],
								},
								{
									Name:  "SMF_N4_PORT",
									Value: smfConfig["n4_port"],
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
			log.Info("SMF Deployment already exists.", "Deployment", deployment.Name)
			return nil
		}
		log.Error(err, "Failed to create SMF Deployment.", "Deployment", deployment.Name)
		return err
	}

	log.Info("Created SMF Deployment.", "Deployment", deployment.Name)
	return nil
}
