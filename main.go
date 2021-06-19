package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type conf struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Spec       spec     `yaml:"spec"`
	Metadata   metadata `yaml:"metadata"`
}

type metadata struct {
	Name string `yaml:"name"`
}

type spec struct {
	Template template `yaml:"template"`
}

type template struct {
	Spec templateSpec `yaml:"spec"`
}

type templateSpec struct {
	Volumes []volume `yaml:"volumes"`
}

type volume struct {
	Name      string    `yaml:"name"`
	Secret    secret    `yaml:"secret"`
	ConfigMap configMap `yaml:"configMap"`
}

type secret struct {
	SecretName string `yaml:"secretName"`
}

type configMap struct {
	Name string `yaml:"name"`
}

func (c *conf) getConf(fileName string) *conf {

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func findAndReplace(fileName string, find string, replace string) {
	input, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, find) {
			lines[i] = strings.Replace(line, find, replace, -1)
			log.Println(line)
		}
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(fileName, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	fileName := os.Args[1]
	operation := os.Args[2]

	var c conf
	c.getConf(fileName)

	if operation == "increment-secret-version" && c.Kind == "Secret" {
		secretName := c.Metadata.Name
		secretNameSplit := strings.Split(secretName, "-")
		currentSecretVersion := strings.Replace(secretNameSplit[len(secretNameSplit)-1], "v", "", -1)
		nextSecretVersionInt, _ := strconv.Atoi(currentSecretVersion)
		nextSecretName := strings.Replace(secretName, "v"+currentSecretVersion, "v"+strconv.Itoa(nextSecretVersionInt+1), -1)
		fmt.Print(nextSecretName)
		findAndReplace(fileName, secretName, nextSecretName)
	} else if operation == "increment-config-version" && c.Kind == "ConfigMap" {
		configName := c.Metadata.Name
		configNameSplit := strings.Split(configName, "-")
		currentConfigVersion := strings.Replace(configNameSplit[len(configNameSplit)-1], "v", "", -1)
		currentConfigVersionInt, _ := strconv.Atoi(currentConfigVersion)
		nextConfigName := strings.Replace(configName, "v"+currentConfigVersion, "v"+strconv.Itoa(currentConfigVersionInt+1), -1)
		fmt.Print(nextConfigName)
		findAndReplace(fileName, configName, nextConfigName)
	} else {
		// Defaults to Deployment Kind
		volumes := c.Spec.Template.Spec.Volumes
		for i := range volumes {
			if volumes[i].Secret.SecretName != "" {
				secretName := volumes[i].Secret.SecretName
				if operation == "get-secret-name" {
					fmt.Print(secretName)
				}
			} else if volumes[i].ConfigMap.Name != "" {
				configName := volumes[i].ConfigMap.Name
				if operation == "get-config-name" {
					fmt.Print(configName)
				}
			}
		}
	}
}
