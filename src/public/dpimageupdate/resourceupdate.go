package dpimageupdate

import (
	"datastructure"
	"encoding/json"
	"fmt"
	"log"
)

func eliminateStatus(a datastructure.Request, bodyContentByte []byte) (err error, newDeploymentByte []byte) {
	var (
		deploymentMap map[string]interface{}
	)
	//Unmarshal the body
	if err = json.Unmarshal(bodyContentByte, &deploymentMap); err != nil {
		log.Printf("[EliminateStatus] Json TO DeploymentMap Json Change ERR: %v", err)
		return
	}

	//delete status from deployment
	delete(deploymentMap, "status")

	//Marshal the new body
	if newDeploymentByte, err = json.Marshal(deploymentMap); err != nil {
		log.Println("[EliminateStatus] DeploymentByte TO Json Change ERR: ", err)
		return
	}
	return
}

func replaceResource(a datastructure.Request, bodyContentByte []byte) (err error, newDeploymentByte []byte) {
	var (
		deploymentMap map[string]interface{}
		strategy      = make(map[string]interface{})
		rollingUpdate = make(map[string]interface{})
	)
	//Unmarshal the body
	if err = json.Unmarshal(bodyContentByte, &deploymentMap); err != nil {
		log.Printf("[ReplaceImage] Json TO DeploymentMap Json Change ERR: %v", err)
		return
	}

	//exchange the image from body
	deploymentMap["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})[0].(map[string]interface{})["image"] = a.Image

	//replace replicas from deployment
	if a.Replicas != 0 {
		deploymentMap["spec"].(map[string]interface{})["replicas"] = a.Replicas
	}

	//replace minReadySeconds from deployment
	if a.MinReadySeconds != 0 {
		deploymentMap["spec"].(map[string]interface{})["minReadySeconds"] = a.MinReadySeconds
	}

	//replace rollingUpdate from deployment
	if a.RollingUpdate.MaxSurge != "" {
		//must be equal
		if a.RollingUpdate.MaxSurge != a.RollingUpdate.MaxUnavailable {
			log.Printf("[ReplaceImage] MaxSurge:%v ≠ MaxUnavailable:%v \n"+
				"MaxSurge and MaxUnavailable Must Be Equal ",
				a.RollingUpdate.MaxSurge, a.RollingUpdate.MaxUnavailable)
			err = fmt.Errorf("[ReplaceImage] MaxSurge:%v ≠ MaxUnavailable:%v \n"+
				"MaxSurge and MaxUnavailable Must Be Equal ",
				a.RollingUpdate.MaxSurge, a.RollingUpdate.MaxUnavailable)
			return
		}
		//init the strategy
		rollingUpdate["maxUnavailable"] = a.RollingUpdate.MaxUnavailable
		rollingUpdate["maxSurge"] = a.RollingUpdate.MaxSurge
		strategy["rollingUpdate"] = rollingUpdate
		strategy["type"] = "RollingUpdate"
		deploymentMap["spec"].(map[string]interface{})["strategy"] = strategy
	}

	//Marshal the new body
	if newDeploymentByte, err = json.Marshal(deploymentMap); err != nil {
		log.Println("[ReplaceImage] DeploymentByte TO Json Change ERR: ", err)
		return
	}
	return
}
