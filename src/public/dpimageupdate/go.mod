module dpimageupdate

require (
	alert v1.11.3
	datastructure v1.11.3
	k8sapi v1.11.3
	user v1.11.3
)

replace k8sapi v1.11.3 => ../k8sapi

replace datastructure v1.11.3 => ../datastructure

replace user v1.11.3 => ../user

replace alert v1.11.3 => ../alert

go 1.13
