package handlers

import (
	"fmt"
	servicePB "nvisible/packetBeatRunner/services"
	"nvisible/packetBeatRunner/utils"

	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

//Packet beat starting, stopping

//PbDuration stores start and endtime of packetbeat runtime(I do not belong here, find me a good place)
type PbDuration struct {
	StartTime int64 `json:"start_time"`
	StopTime  int64 `json:"end_time"`
	Running   bool  `json:"status"`
}

// PacketbeatStart starts the packcetbeat
func (pb *PbDuration) PacketbeatStart(c *gin.Context) {

	if pb.Running == false {
		res := servicePB.StartPB()
		pb.Running = true
		pb.StartTime = utils.MakeTimestamp()

		c.JSON(utils.StatusCodeOK, res)
		return
	}

	if pb.Running == true {
		res := servicePB.PBRunning()
		c.JSON(utils.StatusCodeOK, res)
		return
	}

}

// PacketbeatStop stops the packetbeat and sends the duration into ES
func (pb *PbDuration) PacketbeatStop(c *gin.Context) {

	if pb.Running == true {
		res := servicePB.KillPB()
		pb.Running = false
		pb.StopTime = utils.MakeTimestamp()
		//send to elastic search here
		// pb.SendToEs()
		c.JSON(utils.StatusCodeOK, res)
		return
	}
	if pb.Running == false {
		res := servicePB.PBStopped()
		c.JSON(utils.StatusCodeOK, res)
		return
	}

}

func (pb *PbDuration) toESObj() string {
	str := `{"start_time":%s, "stop_time":%s}`
	return fmt.Sprintf(str, pb.StartTime, pb.StopTime)

}

//SendToEs sends packetbeat runtime to ES
func (pb *PbDuration) SendToEs() {
	insertString := fmt.Sprintf("http://%s:%s/pb_alive_duration/_doc/", "localhost", "9200")
	request := gorequest.New()
	request.AppendHeader("Content-Type", "application/json")

	status, _, _ := request.Post(insertString).Send(pb.toESObj()).End()

	fmt.Println(status.StatusCode)

}
