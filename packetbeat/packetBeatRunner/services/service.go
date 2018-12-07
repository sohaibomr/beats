package services

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/elastic/beats/packetbeat/packetBeatRunner/models"
)

const (
	pbStartMsg = "PacketBeat started"
	pbStopMsg  = "PacketBeat stoped"
)

// StartPB starts the packetbeat
func StartPB() models.PacketBeatStatus {
	var res models.PacketBeatStatus
	cmd := exec.Command("./packetbeat", "-e", "-c", "configs/packetbeat.yml")
	cmd.Dir = "/packetbeat"

	go func() {
		fmt.Println("Starting PB new process in goroutine")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("StartPB", err)
			log.Fatal(err)
		}
	}()

	res.Msg = pbStartMsg
	res.PID = getPBID()

	return res
}

// PBRunning constructs if PB is already running
func PBRunning() models.PacketBeatStatus {
	res := models.PacketBeatStatus{Msg: "PackcetBeats is running already, stop it first.", PID: getPBID()}
	return res
}

// PBStopped constructs if PB is already stopped
func PBStopped() models.PacketBeatStatus {
	res := models.PacketBeatStatus{Msg: "PackcetBeats is stopped already.", PID: "0"}
	return res
}

// GetPBID gets the process ID of PB
func getPBID() string {
	cmd := exec.Command("/bin/sh", "-c",
		"pgrep packetbeat")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("getPBID", err)
		log.Fatal(err)
	}

	return fmt.Sprintf("%s", out)
}

// KillPB kills the the packetbeat by finding it's id
func KillPB() models.PacketBeatStatus {
	var res models.PacketBeatStatus
	res.PID = getPBID()
	killCmd := fmt.Sprintf("kill -15 %s", res.PID)
	fmt.Print(killCmd)

	err := exec.Command("/bin/sh", "-c",
		killCmd).Start()
	if err != nil {
		fmt.Println("killpid", err)
		log.Fatal(err)
	}
	res.Msg = pbStopMsg
	return res

}
