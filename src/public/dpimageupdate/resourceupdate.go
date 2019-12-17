package dpimageupdate

import (
	"datastructure"
	"encoding/json"
	"fmt"
	"log"
)

func eliminateStatus(bodyContentByte []byte) (err error, newDeploymentByte []byte) {
	var (
		deploymentMap, metadataMap map[string]interface{}
	)
	//Unmarshal the body
	if err = json.Unmarshal(bodyContentByte, &deploymentMap); err != nil {
		log.Printf("[EliminateStatus] Json TO DeploymentMap Json Change ERR: %v", err)
		err = fmt.Errorf("[EliminateStatus] Json TO DeploymentMap Json Change ERR: %v", err)
		return
	}

	//delete status and metadata from deployment
	//second apiserver put will filed if not delete key resourceVersion of metadataMap
	delete(deploymentMap, "status")
	metadataMap = deploymentMap["metadata"].(map[string]interface{})
	delete(metadataMap, "resourceVersion")
	delete(metadataMap, "annotations")
	delete(metadataMap, "creationTimestamp")
	delete(metadataMap, "generation")
	delete(metadataMap, "uid")
	delete(metadataMap, "annotations")
	delete(metadataMap, "selfLink")

	//Marshal the new body
	if newDeploymentByte, err = json.Marshal(deploymentMap); err != nil {
		log.Printf("[EliminateStatus] DeploymentByte TO Json Change ERR: %v", err)
		err = fmt.Errorf("[EliminateStatus] DeploymentByte TO Json Change ERR: %v", err)
		return
	}
	return
}

func replaceResource(a datastructure.Request, bodyContentByte []byte) (err error, newDeploymentByte []byte, newA datastructure.Request) {
	var (
		deploymentMap map[string]interface{}
		strategy      = make(map[string]interface{})
		rollingUpdate = make(map[string]interface{})
	)
	//Unmarshal the body
	if err = json.Unmarshal(bodyContentByte, &deploymentMap); err != nil {
		log.Printf("[ReplaceResource] Json TO DeploymentMap Json Change ERR: %v", err)
		err = fmt.Errorf("[ReplaceResource] Json TO DeploymentMap Json Change ERR: %v", err)
		return
	}

	//must be equal
	if a.RollingUpdate.MaxSurge != a.RollingUpdate.MaxUnavailable {
		log.Printf("[ReplaceResource] MaxSurge:%v ≠ MaxUnavailable:%v \n"+
			"MaxSurge and MaxUnavailable Must Be Equal ",
			a.RollingUpdate.MaxSurge, a.RollingUpdate.MaxUnavailable)
		err = fmt.Errorf("[ReplaceResource] MaxSurge:%v ≠ MaxUnavailable:%v \n"+
			"MaxSurge and MaxUnavailable Must Be Equal ",
			a.RollingUpdate.MaxSurge, a.RollingUpdate.MaxUnavailable)
		return
	}

	//exchange the image from body
	deployImage := deploymentMap["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})[0].(map[string]interface{})["image"]
	if a.Image != deployImage.(string) {
		deploymentMap["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})[0].(map[string]interface{})["image"] = a.Image
	} else {
		log.Printf("Deployment Image %v == %v, ImageTag Can Not Be Same! ", deployImage.(string), a.Image)
		err = fmt.Errorf("Deployment Image %v == %v, ImageTag Can Not Be Same! ", deployImage.(string), a.Image)
		return
	}

	//replace replicas from deployment
	if t := a.Replicas.String(); t == "" || a.Name == "GrayDeployment" {
		//a.Replicas = json.Number(deploymentMap["spec"].(map[string]interface{})["replicas"].(string))
		deploymentMap["spec"].(map[string]interface{})["replicas"] = 1
		newA = a
	} else if a.Name == "InstantDeployment" {
		deploymentMap["spec"].(map[string]interface{})["replicas"] = a.Replicas
	}

	//replace minReadySeconds from deployment
	if a.MinReadySeconds.String() != string(0) {
		deploymentMap["spec"].(map[string]interface{})["minReadySeconds"] = a.MinReadySeconds
	}

	if a.Name == "InstantDeployment" {
		deploymentMap["metadata"].(map[string]interface{})["name"] = a.Deployment
	} else if a.Name == "GrayDeployment" {
		deploymentMap["metadata"].(map[string]interface{})["name"] = "temp-" + a.Deployment
	} else {
		log.Printf("[ReplaceResource] Name %v Is Not Right", a.Name)
		err = fmt.Errorf("[ReplaceResource] Name %v Is Not Right", a.Name)
		return
	}

	////replace rollingUpdate from deployment
	rollingUpdate["maxUnavailable"] = a.RollingUpdate.MaxUnavailable
	rollingUpdate["maxSurge"] = a.RollingUpdate.MaxSurge
	strategy["rollingUpdate"] = rollingUpdate
	strategy["type"] = "RollingUpdate"
	deploymentMap["spec"].(map[string]interface{})["strategy"] = strategy

	//Marshal the new body
	if newDeploymentByte, err = json.Marshal(deploymentMap); err != nil {
		log.Printf("[ReplaceResource] DeploymentByte TO Json Change ERR: %v", err)
		err = fmt.Errorf("[ReplaceResource] DeploymentByte TO Json Change ERR: %v", err)
		return
	}
	return
}

func ReplaceResourceName(a datastructure.Request, bodyContentByte []byte) (err error, newDeploymentByte []byte) {
	var (
		deploymentMap map[string]interface{}
	)
	//Unmarshal the body
	if err = json.Unmarshal(bodyContentByte, &deploymentMap); err != nil {
		log.Printf("[ReplaceResourceName] Json TO DeploymentMap Json Change ERR: %v", err)
		err = fmt.Errorf("[ReplaceResourceName] Json TO DeploymentMap Json Change ERR: %v", err)
		return
	}

	if a.Name == "InstantDeployment" {
		deploymentMap["metadata"].(map[string]interface{})["name"] = a.Deployment
	} else if a.Name == "GrayDeployment" {
		deploymentMap["metadata"].(map[string]interface{})["name"] = "temp-" + a.Deployment
	} else {
		log.Printf("[ReplaceResourceName] Name %v Is Not Right", a.Name)
		err = fmt.Errorf("[ReplaceResourceName] Name %v Is Not Right", a.Name)
		return
	}

	//Marshal the new body
	if newDeploymentByte, err = json.Marshal(deploymentMap); err != nil {
		log.Printf("[ReplaceResourceName] DeploymentByte TO Json Change ERR: %v", err)
		err = fmt.Errorf("[ReplaceResourceName] DeploymentByte TO Json Change ERR: %v", err)
		return
	}
	return
}

func ReplaceResourceReplicas(spec datastructure.MySpec) (err error, newDeploymentByte []byte) {
	//Marshal the new body
	if newDeploymentByte, err = json.Marshal(spec); err != nil {
		log.Printf("[ReplaceResourceReplicas] DeploymentByte TO Json Change ERR: %v", err)
		err = fmt.Errorf("[ReplaceResourceReplicas] DeploymentByte TO Json Change ERR: %v", err)
		return
	}
	return
}
