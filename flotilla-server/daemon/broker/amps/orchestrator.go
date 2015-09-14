package amps

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

const (
	ampsDocker   = "amps"
	internalPort = "9007"
)

// Broker implements the broker interface for AMPS.
type Broker struct {
	containerID  string
	DockerExtras string
}

// Start will start the message broker and prepare it for testing.
func (a *Broker) Start(host, port string) (interface{}, error) {
	containerID, err := exec.Command("/bin/sh", "-c",
		fmt.Sprintf("docker run %s -d -p %s:%s %s /root/AMPS-develop-Release-Linux/bin/ampServer /root/config.xml", a.DockerExtras, port, internalPort, ampsDocker)).Output()
	if err != nil {
		log.Printf("Failed to start container %s: %s", ampsDocker, err.Error())
		return "", err
	}
	time.Sleep(2 * time.Second)

	log.Printf("Started container %s: %s", ampsDocker, containerID)
	a.containerID = string(containerID)
	return string(containerID), nil
}

// Stop will stop the message broker.
func (a *Broker) Stop() (interface{}, error) {
	containerID, err := exec.Command("/bin/sh", "-c",
		fmt.Sprintf("docker kill %s", a.containerID)).Output()
	if err != nil {
		log.Printf("Failed to stop container %s: %s", ampsDocker, err.Error())
		return "", err
	}

	log.Printf("Stopped container %s: %s", ampsDocker, a.containerID)
	return string(containerID), nil
}
