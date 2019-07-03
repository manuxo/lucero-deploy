package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/otiai10/copy"
)

type CopyOption int8
type EnvOption int8

const (
	COPY_WEBSITE  CopyOption = 1
	COPY_SERVICES CopyOption = 2
	COPY_BOTH     CopyOption = 3
	DEV_ENV       EnvOption  = 1
	TEST_ENV      EnvOption  = 2
	PROD_ENV      EnvOption  = 3
	ALL_ENV       EnvOption  = 4
)

type DeployConfig struct {
	WebSiteSourcePath  string `json:"WEBSITE_SOURCE_PATH"`
	WebSiteDestPath    string `json:"WEBSITE_DEST_PATH"`
	WebSiteBackupPath  string `json:"WEBSITE_BACKUP_PATH"`
	ServicesSourcePath string `json:"SERVICES_SOURCE_PATH"`
	ServicesDestPath   string `json:"SERVICES_DEST_PATH"`
	ServicesBackupPath string `json:"SERVICES_BACKUP_PATH"`
}

type Config struct {
	Dev     DeployConfig `json:"dev"`
	Testing DeployConfig `json:"test"`
	Prod    DeployConfig `json:"prod"`
}

var config Config

func ReadConfig() {
	jsonFile, err := os.Open("deploy.config.json")
	if err != nil {
		log.Panic("Error at reading file.")
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)
}

func Deploy(deployConfig DeployConfig, copyOption CopyOption,chRead chan<- bool) {
	currentDate := time.Now()
	minute := currentDate.Minute()
	hour := currentDate.Hour()
	strDate := currentDate.Format("02-01-2006")
	webBackupPath := fmt.Sprintf("%s/%s-%dhrs%dmin", deployConfig.WebSiteBackupPath, strDate, hour, minute)
	servicesBackupPath := fmt.Sprintf("%s/%s-%dhrs%dmin", deployConfig.ServicesBackupPath, strDate, hour, minute)
	if copyOption == COPY_WEBSITE {
		_ = copy.Copy(deployConfig.WebSiteDestPath, webBackupPath)
		_ = copy.Copy(deployConfig.WebSiteSourcePath, deployConfig.WebSiteDestPath)
	} else if copyOption == COPY_SERVICES {
		_ = copy.Copy(deployConfig.ServicesDestPath, servicesBackupPath)
		_ = copy.Copy(deployConfig.ServicesSourcePath, deployConfig.ServicesDestPath)
	} else if copyOption == COPY_BOTH {
		_ = copy.Copy(deployConfig.WebSiteDestPath, webBackupPath)
		_ = copy.Copy(deployConfig.ServicesDestPath, servicesBackupPath)
		_ = copy.Copy(deployConfig.WebSiteSourcePath, deployConfig.WebSiteDestPath)
		_ = copy.Copy(deployConfig.ServicesSourcePath, deployConfig.ServicesDestPath)
	}
	chRead<-true
}

func main() {
	chRead := make(chan bool)
	for {
		ReadConfig()
		var envOption EnvOption
		var copyOption CopyOption
		fmt.Println("--- Lucero Deploy ---")
		fmt.Print("Ingrese opción: [Desarrollo: 1, Testing: 2, Produccion: 3, Todos: 4, Salir: 0]: ")
		fmt.Scanf("%d\n", &envOption)
		if envOption == 0 {
			fmt.Println("--- --- --- ---")
			break
		}
		fmt.Print("Ingrese opción de copiado: [Web: 1, Servicios: 2, Ambos: 3]: ")
		fmt.Scanf("%d\n", &copyOption)
		fmt.Println("Copiando archivos...")
		if envOption == DEV_ENV {
			go Deploy(config.Dev, copyOption,chRead)
		} else if envOption == TEST_ENV {
			go Deploy(config.Testing, copyOption,chRead)
		} else if envOption == PROD_ENV {
			go Deploy(config.Prod, copyOption,chRead)
		} else if envOption == ALL_ENV {
			go Deploy(config.Dev, copyOption,chRead)
			go Deploy(config.Testing, copyOption,chRead)
			go Deploy(config.Prod, copyOption,chRead)
		}
		<-chRead
		fmt.Println("Se han copiado los archivos satisfactoriamente!!!")
	}
}
