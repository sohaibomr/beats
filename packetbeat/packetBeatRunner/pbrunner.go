package main

import "github.com/gin-gonic/gin"
import "nvisible/packetBeatRunner/handlers"

func main() {
	router := gin.Default()
	pb := &handlers.PbDuration{StartTime: 0, StopTime: 0, Running: false}

	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "World"})
	})

	router.GET("/start-packetbeat", pb.PacketbeatStart)
	// swagger:operation GET /start-packetbeat Packetbeat pb
	// ---
	// summary: Starts the packetbeat process by making new process session
	// description: Starts the packetbeat.
	// responses:
	//   "200":
	//    description: Count of field in path as response
	//    schema:
	//     "$ref": "#/definitions/resPBStatus"

	router.GET("/stop-packetbeat", pb.PacketbeatStop)
	// swagger:operation GET /stop-packetbeat Packetbeat pb
	// ---
	// summary: Stops the packetbeat process by making new process session
	// description: Stops the packetbeat.
	// responses:
	//   "200":
	//    description: Count of field in path as response
	//    schema:
	//     "$ref": "#/definitions/resPBStatus"

	router.Run()
}
