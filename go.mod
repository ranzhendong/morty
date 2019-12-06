module morty

go 1.13

replace k8sapi v1.11.3 => ./src/public/k8sapi

replace dpimageupdate v1.11.3 => ./src/public/dpimageupdate

replace user v1.11.3 => ./src/public/user

replace datastructure v1.11.3 => ./src/public/datastructure

replace alert v1.11.3 => ./src/public/alert

replace configanalysis v1.11.3 => ./src/public/configanalysis

require (
	configanalysis v1.11.3
	datastructure v1.11.3
	dpimageupdate v1.11.3
)
