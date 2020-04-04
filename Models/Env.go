package Models

import (
	"fmt"
	"os"
)

var envStruct *Env

type Env struct {
	AwsStorageUrl string
	OriginAllowed string
	Env string
	Port string
}

func InitEnv()  {
	envStruct = &Env{
		AwsStorageUrl:   			os.Getenv("AWS_STORAGE_URL"),
		OriginAllowed:  			os.Getenv("ORIGIN_ALLOWED"),
		Env: 			  			os.Getenv("ENV"),
		Port: 						os.Getenv("PORT"),
	}
	fmt.Println(envStruct)
}

func GetEnvStruct() *Env  {
	return  envStruct
}