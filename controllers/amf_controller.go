package controllers

import (
	"context"
	oaiv1 "github.com/FETHIA-MHD/openai-operator/api/v1"
	"k8s.io/api/apps/v1"                   // Pour le package apps/v1 (Déploiement)
	"k8s.io/api/core/v1"                   // Pour le package core/v1 (Conteneurs et PodSpec)
	"k8s.io/apimachinery/pkg/apis/meta/v1" // Pour le package meta/v1 (Métadonnées)
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// createAMFDeployment crée un déploiement pour AMF
func (r *OaiCoreReconciler) createAMFDeployment(ctx context.Context, oaiCore *oaiv1.OaiCore) error {
	log := log.FromContext(ctx)

	labels := map[string]string{
		"app":        "oaicore",
		"component":  "amf",
		"controller": oaiCore.Name,
	}

	// Configuration AMF
	amfConfig := map[string]string{
		"host":           "oai-amf",
		"sbi_port":       "80",
		"api_version":    "v1",
		"interface_name": "eth0",
		"n2_port":        "38412",
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      oaiCore.Name + "-amf",
			Namespace: oaiCore.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &oaiCore.Spec.AMF.Replicas,
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
							Name:  "amf",
							Image: oaiCore.Spec.AMF.Image,
							Env: []corev1.EnvVar{
								{
									Name:  "AMF_HOST",
									Value: amfConfig["host"],
								},
								{
									Name:  "AMF_SBI_PORT",
									Value: amfConfig["sbi_port"],
								},
								{
									Name:  "AMF_API_VERSION",
									Value: amfConfig["api_version"],
								},
								{
									Name:  "AMF_INTERFACE_NAME",
									Value: amfConfig["interface_name"],
								},
								{
									Name:  "AMF_N2_PORT",
									Value: amfConfig["n2_port"],
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
			log.Info("AMF Deployment already exists.", "Deployment", deployment.Name)
			return nil
		}
		log.Error(err, "Failed to create AMF Deployment.", "Deployment", deployment.Name)
		return err
	}

	log.Info("Created AMF Deployment.", "Deployment", deployment.Name)
	return nil
}
