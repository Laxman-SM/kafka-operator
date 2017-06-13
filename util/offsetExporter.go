package util

import (
	"github.com/krallistic/kafka-operator/spec"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1Beta1 "k8s.io/client-go/pkg/apis/apps/v1beta1"

	"k8s.io/client-go/pkg/api/v1"
)

const (
	deplyomentPrefix = "kafka-offset-checker"
)

func (c *ClientUtil) getOffsetMonitorName(cluster spec.KafkaCluster) string {
	return deplyomentPrefix + "-" + cluster.Metadata.Name
}

// Deploys the OffsetMonitor as an extra Pod inside the Cluster
func (c *ClientUtil) DeployOffsetMonitor(cluster spec.KafkaCluster) error {
	deployment, err := c.KubernetesClient.AppsV1beta1().Deployments(cluster.Metadata.Namespace).Get(c.getOffsetMonitorName(cluster), c.DefaultOption)


	if err != nil {
		fmt.Println("error while talking to k8s api: ", err)
		//TODO better error handling, global retry module?
		return err
	}
	if len(deployment.Name) == 0 {
		//Service dosnt exist, creating new.
		fmt.Println("Deployment dosnt exist, creating new")
		replicas := int32(1)

		objectMeta := metav1.ObjectMeta{
			Name: c.getOffsetMonitorName(cluster),
			Annotations: map[string]string{
				"component": "kafka",
				"name":      cluster.Metadata.Name,
				"role":      "data",
				"type":      "service",
			},
		}

		podObjectMeta := metav1.ObjectMeta{
			Name: c.getOffsetMonitorName(cluster),
			Annotations: map[string]string{
				"component": "kafka",
				"name":      cluster.Metadata.Name,
				"role":      "data",
				"type":      "service",
				//TODO Prometheus Annotations
			},
		}
		deploy := appsv1Beta1.Deployment{
			ObjectMeta: objectMeta,
			Spec: appsv1Beta1.DeploymentSpec{
				Replicas: &replicas,
				Template: v1.PodTemplateSpec{
					ObjectMeta: podObjectMeta,
					Spec: v1.PodSpec{

					},
				},

			},
		}



		_, err := c.KubernetesClient.AppsV1beta1().Deployments(cluster.Metadata.Namespace).Create(&deploy)
		if err != nil {
			fmt.Println("Error while creating Deployment: ", err)
		}
	} else {
		//Service exist
		fmt.Println("Deployment already exist: ", deployment)
		//TODO maybe check for correct service?
	}

	return nil

	return nil
}


//Deletes the offset checker for the given kafka cluster.
// Return error if any problems occurs. (Except if monitor dosnt exist)
//
func (c *ClientUtil) DeleteOffsetMonitor(cluster spec.KafkaCluster) error {


	return nil
}

