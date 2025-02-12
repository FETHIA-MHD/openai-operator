package controllers

import (
	"context"
	oaiv1 "github.com/FETHIA-MHD/openai-operator/api/v1"
	"k8s.io/api/apps/v1"                   // Pour le package apps/v1 (Déploiement)
	"k8s.io/api/core/v1"                   // Pour le package core/v1 (Conteneurs et PodSpec)
	"k8s.io/apimachinery/pkg/apis/meta/v1" // Pour le package meta/v1 (Métadonnées)
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// createUPFDeployment crée un déploiement pour UPF
func (r *OaiCoreReconciler) createUPFDeployment(ctx context.Context, oaiCore *oaiv1.OaiCore) error {
	log := log.FromContext(ctx)

	labels := map[string]string{
		"app":        "oaicore",
		"component":  "upf",
		"controller": oaiCore.Name,
	}

	// Configuration UPF
	upfConfig := map[string]string{
		"host":           "oai-upf",
		"sbi_port":       "80",
		"api_version":    "v1",
		"interface_name": "eth0",
		"n3_port":        "2152",
		"n4_port":        "8805",
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      oaiCore.Name + "-upf",
			Namespace: oaiCore.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &oaiCore.Spec.UPF.Replicas,
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
							Name:  "upf",
							Image: oaiCore.Spec.UPF.Image,
							Env: []corev1.EnvVar{
								{
									Name:  "UPF_HOST",
									Value: upfConfig["host"],
								},
								{
									Name:  "UPF_SBI_PORT",
									Value: upfConfig["sbi_port"],
								},
								{
									Name:  "UPF_API_VERSION",
									Value: upfConfig["api_version"],
								},
								{
									Name:  "UPF_INTERFACE_NAME",
									Value: upfConfig["interface_name"],
								},
								{
									Name:  "UPF_N3_PORT",
									Value: upfConfig["n3_port"],
								},
								{
									Name:  "UPF_N4_PORT",
									Value: upfConfig["n4_port"],
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
			log.Info("UPF Deployment already exists.", "Deployment", deployment.Name)
			return nil
		}
		log.Error(err, "Failed to create UPF Deployment.", "Deployment", deployment.Name)
		return err
	}

	log.Info("Created UPF Deployment.", "Deployment", deployment.Name)
	return nil
}
