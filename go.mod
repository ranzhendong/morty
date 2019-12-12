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
	dpimageupdate v1.11.3
	github.com/spf13/viper v1.5.0
	golang.org/x/oauth2 v0.0.0-20191202225959-858c2ad4c8b6 // indirect
)
