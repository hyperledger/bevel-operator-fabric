package tests

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func getDefaultResources() (corev1.ResourceRequirements, error) {
	requestCpu, err := resource.ParseQuantity("10m")
	if err != nil {
		return corev1.ResourceRequirements{}, err
	}
	requestMemory, err := resource.ParseQuantity("128Mi")
	if err != nil {
		return corev1.ResourceRequirements{}, err
	}
	limitsCpu, err := resource.ParseQuantity("300m")
	if err != nil {
		return corev1.ResourceRequirements{}, err
	}
	limitsMemory, err := resource.ParseQuantity("256Mi")
	if err != nil {
		return corev1.ResourceRequirements{}, err
	}
	return corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    requestCpu,
			corev1.ResourceMemory: requestMemory,
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    limitsCpu,
			corev1.ResourceMemory: limitsMemory,
		},
	}, nil
}
