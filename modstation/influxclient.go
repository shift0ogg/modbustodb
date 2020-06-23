package modstation

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

//InfluxManager ---
type InfluxManager struct {
	pChannelMgr *ChannelManager
	influx      client.Client
}

//NewInfluxManager ---
func NewInfluxManager(pChannel *ChannelManager) (*InfluxManager, error) {
	pInfluxManager := &InfluxManager{pChannelMgr: pChannel}
	err := pInfluxManager.init()
	return pInfluxManager, err
}

func (influxMgr *InfluxManager) ping() error {
	_, _, err := influxMgr.influx.Ping(5 * time.Second)
	return err
}

//init ok
func (influxMgr *InfluxManager) init() error {
	var err error
	influxMgr.influx, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     influxMgr.pChannelMgr.Influxdb,
		Username: influxMgr.pChannelMgr.Username,
		Password: influxMgr.pChannelMgr.Password,
		Timeout:  time.Second * 5,
	})
	if err != nil {
		return err
	}
	err = influxMgr.ping()
	return err

}

//ToDb --
func (influxMgr *InfluxManager) ToDb() error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Precision: "ms",
		Database:  influxMgr.pChannelMgr.MyDB,
	})
	if err != nil {
		return err
	}

	for _, ch := range influxMgr.pChannelMgr.Channels {
		for _, dev := range ch.Device {
			for _, dp := range dev.Datapoint {
				point, _ := client.NewPoint(
					dp.Dpid,
					map[string]string{
						"addr": strconv.Itoa(int(dev.Addr)), "alarm": strconv.Itoa(int(dp.Alarm)), "connected": strconv.FormatBool(dp.Connected), "valid": strconv.FormatBool(dp.Valid)},
					map[string]interface{}{"t": dp.Val}, dp.updatetime)
				bp.AddPoint(point)
			}

			err = influxMgr.influx.Write(bp)
			if err != nil {
				fmt.Println("write error", err)
				log.Println("write error", err)
			}
		}
	}
	return nil

}
