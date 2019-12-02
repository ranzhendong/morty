module morty

go 1.13

require (
	dpimageupdate v1.11.3
	k8sapi v1.11.3
)

replace k8sapi v1.11.3 => ./src/public/k8sapi

replace dpimageupdate v1.11.3 => ./src/public/dpimageupdate
