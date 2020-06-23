package modstation

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

//PostParams 设置报警条件
func (cm *ChannelManager) PostParams(c *gin.Context) {
	fmt.Println(">>>> 设置报警条件 <<<<")
	dpid := c.PostForm("dpid")
	alarmcon := c.PostForm("alarmcon")
	err := cm.setAlarmCondition(dpid, alarmcon)
	//cm.writeConfigfile(false)
	c.JSON(http.StatusOK, gin.H{
		"status": err == nil,
	})
}

//GetJSON 设置报警条件
func (cm *ChannelManager) GetJSON(c *gin.Context) {
	lockerData.RLock()
	defer lockerData.RUnlock()
	c.JSON(http.StatusOK, cm)
}

//GetData 设置报警条件
func (cm *ChannelManager) GetData(c *gin.Context) {
	c.Status(200)
	tmpl, err := template.New("message").Parse(tpldata)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}
	lockerData.RLock()
	defer lockerData.RUnlock()
	tmpl.Execute(c.Writer, cm)
}

//SetAlarmCondition 设置报警条件
func (cm *ChannelManager) setAlarmCondition(dpid string, alarmCondition string) error {
	lockerData.Lock()
	defer lockerData.Unlock()
	for _, ch := range cm.Channels {
		for _, dev := range ch.Device {
			for _, dp := range dev.Datapoint {
				if dp.Dpid == dpid {
					vstr := strings.Split(alarmCondition, ",")
					if len(vstr) != 2 {
						return errors.New("Format Error")
					}
					f2, err := strconv.ParseFloat(vstr[0], 32) //黄色
					if err != nil {
						return errors.New("Format Error")
					}
					f1, err := strconv.ParseFloat(vstr[1], 32) //红色
					if err != nil {
						return errors.New("Format Error")
					}
					dp.Alarmcondition = alarmCondition
					dp.alarmLevels = []float32{float32(f2), float32(f1)}
					return nil
				}
			}
		}
	}
	return errors.New("Dpid not found")
}

//ReloadModbus  ---
func (cm *ChannelManager) ReloadModbus(c *gin.Context) {
	fmt.Println(">>>> ReloadModbus <<<<")
	log.Println(">>>> ReloadModbus <<<<")
	err := cm.reInitFromDB()
	c.JSON(http.StatusOK, gin.H{
		"status": err == nil,
	})
}
