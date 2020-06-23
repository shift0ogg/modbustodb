package modstation

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/goburrow/modbus"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/sasha-s/go-deadlock"
)

//Debug 同时记录到日志和Console
func Debug(a ...interface{}) {
	fmt.Println(a...)
	log.Println(a...)
}

var (
	chPoint chan *Datapoint = make(chan *Datapoint, 1000)

	lockerData    deadlock.RWMutex = deadlock.RWMutex{} //lock for channels data
	lockerQuit    deadlock.RWMutex = deadlock.RWMutex{} //lock Quit
	lockerRequest deadlock.RWMutex = deadlock.RWMutex{} //lock Request
	lockerWrite   deadlock.RWMutex = deadlock.RWMutex{} //lock Write
	requesting    bool             = false              //是否还在循环485请求中
	writing       bool             = false              //是否写influxdb中
	quitRun       bool             = false
)

func isQuit() bool {
	lockerQuit.RLock()
	defer lockerQuit.RUnlock()
	return quitRun
}

func setQuit(b bool) {
	lockerQuit.Lock()
	defer lockerQuit.Unlock()
	quitRun = b
}

func isRequesting() bool {
	lockerRequest.RLock()
	defer lockerRequest.RUnlock()
	return requesting
}

func setRequest(b bool) {
	lockerRequest.Lock()
	defer lockerRequest.Unlock()
	requesting = b
}

func isWriting() bool {
	lockerWrite.RLock()
	defer lockerWrite.RUnlock()
	return writing
}
func setWrite(b bool) {
	lockerWrite.Lock()
	defer lockerWrite.Unlock()
	writing = b
}

//BytesToInt16 转换
func BytesToInt16(buf []byte) int16 {
	return int16(binary.BigEndian.Uint16(buf))
}

type onConn func(dp1 *Datapoint)

func (ch *Channel) getTemp(dev *Device, data *Datapoint, fun1 onConn) error {

	d, e := ch.readTemperature(dev.Addr, data.Regaddr)
	if e != nil {
		fmt.Printf("addr:[%s][%d][%d],[%v]\n", ch.Portname, dev.Addr, data.Regaddr, e)
		if data.Valid == true {
			data.updatetime = time.Now()
			data.UpdateTime = data.updatetime.Format("2006-01-02 15:04:05")
			data.Valid = false
		}
		data.validNum = 0
		data.inValidNum++
		//超过5分钟秒未更新，认为断线了
		if (time.Now().Sub(data.updatetime) > 5*time.Minute || data.inValidNum >= 2) && data.Connected {
			data.updatetime = time.Now()
			data.UpdateTime = data.updatetime.Format("2006-01-02 15:04:05")
			data.Connected = false
			/*记录系统日志*/
			if fun1 != nil {
				fmt.Println(data.Dpid, " lost")
				fun1(data)
			}
		}
	} else {
		data.Val = d
		data.Valid = true

		data.inValidNum = 0 //请求失败计数清零
		data.validNum++
		if (data.validNum >= 2) && !data.Connected {
			data.Connected = true
			/*记录系统日志*/
			if fun1 != nil {
				fmt.Println(data.Dpid, " connected")
				fun1(data)
			}
		}
		data.updatetime = time.Now()
		data.UpdateTime = data.updatetime.Format("2006-01-02 15:04:05")
		/*历史温度存储到切片，并做突变消除判断*/
		if len(data.valList) > 2 {
			data.valList[0] = data.valList[1]
			data.valList[1] = data.valList[2]
			data.valList[2] = data.Val
		} else { //
			data.valList = append(data.valList, data.Val, data.Val, data.Val)
		}
		// //变化大于10摄氏度,滤除畸变
		if math.Abs(float64(data.valList[2]-data.valList[0])) > 10 && math.Abs(float64(data.valList[2]-data.valList[1])) > 10 {
			data.Val = data.valList[1]
		}
		if math.Abs(float64(data.valList[2]-data.valList[0])) > 10 && math.Abs(float64(data.valList[1]-data.valList[0])) > 10 {
			data.Val = data.valList[2]
		}
		//判断报警
		data.Alarm = 0 //初始化无报警
		if len(data.alarmLevels) == 2 {
			f1 := data.alarmLevels[0]
			f2 := data.alarmLevels[1]
			if data.Val > f1 && data.valLast > f1 {
				data.Alarm = 1 //红色
			}
			if data.Val > f2 && data.valLast > f2 {
				data.Alarm = 2 //黄色
			}
		}
		data.valLast = data.Val
	}

	return e

}

//initParameters ---
func (dev *Device) initParameters() {
	for _, vdata := range dev.Datapoint {
		vstr := strings.Split(vdata.Alarmcondition, ",")
		if len(vstr) == 2 {
			f2, _ := strconv.ParseFloat(vstr[0], 32) //黄色
			f1, _ := strconv.ParseFloat(vstr[1], 32) //红色
			vdata.alarmLevels = []float32{float32(f2), float32(f1)}
		} else {
			vdata.alarmLevels = []float32{40, 80}
		}

		vdata.updatetime = time.Now()
		vdata.UpdateTime = vdata.updatetime.Format("2006-01-02 15:04:05")
	}
}

//更新配置文件
func (cm *ChannelManager) writeConfigfile(bBak bool) {
	//备份原有配置文件
	if bBak {
		bytes, err := ioutil.ReadFile(cm.configFile)
		if err != nil {
			return
		}
		dir, fileWithExt := path.Split(cm.configFile)
		strs := strings.Split(fileWithExt, ".")
		bakFile := path.Join(dir, strs[0]+"_"+strconv.FormatInt(time.Now().Unix(), 10)+"."+strs[1])
		err = ioutil.WriteFile(bakFile, bytes, 0666)
		if err != nil {
			return
		}
	}

	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	bytes, err := json.Marshal(cm)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(cm.configFile, bytes, 0666)
	if err != nil {
		return
	}
}

//RequestData --- 一次数据请求
func (cm *ChannelManager) RequestData() {
	if isQuit() {
		return
	}
	fmt.Println("query start...")

	setRequest(true)
	lockerData.RLock()
	var group sync.WaitGroup
	group.Add(len(cm.Channels))
	for _, v := range cm.Channels {
		go func(vp *Channel) {
			fmt.Println(vp.Portname + " query start!")
			err := vp.requestData(func(dp1 *Datapoint) {
				cm.toSysLog(dp1)
			})
			if err != nil {
				fmt.Println(vp.Portname + " quit...")
			}
			fmt.Println(vp.Portname + " query finished!")
			group.Done()
		}(v)
	}
	group.Wait()
	lockerData.RUnlock()
	fmt.Println("query end...")
	setRequest(false)
}

//Close ---
func (cm *ChannelManager) Close() {
	fmt.Println("Close...")

	setQuit(true)
	for isRequesting() {
		time.Sleep(time.Millisecond * 10) //等待ms
	}
	fmt.Println("quit requesting")
	cm.closeChannels()
	for isWriting() {
		time.Sleep(time.Millisecond * 10) //等待ms
	}
	fmt.Println("quit writing")

	cm.closeChannels()
	cm.influx.Close()
	cm.influx = nil
}

func (cm *ChannelManager) closeChannels() {
	for _, v := range cm.Channels {
		v.closeCom()
	}
}

//Init 通道管理初始化
func (cm *ChannelManager) Init(jsonFile string) error {

	//读取配置文件
	cm.configFile = jsonFile
	data, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	/**
	读取配置文件
	*/
	err = json.Unmarshal(data, cm)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	//初始化influx连接参数
	err = cm.initInflux()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	//初始化写库通道
	chPoint = make(chan *Datapoint, 100)
	go cm.writePointThread() //启动写库通道

	chs, err := cm.getDpsFromDB()
	if err != nil {
		log.Println(err.Error())
	} else {
		cm.Channels = chs
	}

	//初始化485读取
	for _, v1 := range cm.Channels {
		v1.initParameters()
		v1.initCom()
	}

	return nil
}

//ReInitDataPointFromDB 从数据库初始化数据点
func (cm *ChannelManager) reInitFromDB() error {
	setQuit(true)
	for isRequesting() {
		time.Sleep(time.Millisecond * 10) //等待ms
	}
	fmt.Println("quit requesting")
	cm.closeChannels()
	for isWriting() {
		time.Sleep(time.Millisecond * 10) //等待ms
	}
	fmt.Println("quit writing")

	setQuit(false)

	lockerData.Lock()
	if chPoint != nil {
		close(chPoint)
	}
	//初始化写库通道
	chPoint = make(chan *Datapoint, 100)
	go cm.writePointThread() //启动写库通道

	//初始化influx连接参数
	err := cm.initInflux()
	if err != nil {
		Debug("initInflux", err.Error())
		lockerData.Unlock()
		return err
	}

	chs, err := cm.getDpsFromDB()
	if err != nil {
		Debug("读取测点库失败", err.Error())
		lockerData.Unlock()
		return err
	}
	cm.Channels = chs
	for _, v1 := range cm.Channels {
		v1.initParameters()
		v1.initCom()
	}
	lockerData.Unlock()
	//485读取

	return err
}

func (cm *ChannelManager) writePointThread() error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Precision: "ms",
		Database:  cm.MyDB,
	})
	if err != nil {
		Debug("NewBatchPoints", err.Error())
		return err
	}
	setWrite(true)
	for !isQuit() {
		select {
		case dp := <-chPoint: //写数据
			if dp == nil {
				break
			}
			point, _ := client.NewPoint(
				dp.Dpid,
				map[string]string{
					"addr":      strconv.Itoa(1),
					"alarm":     strconv.Itoa(int(dp.Alarm)),
					"connected": strconv.FormatBool(dp.Connected),
					"valid":     strconv.FormatBool(dp.Valid),
					"alarmcon":  dp.Alarmcondition,
				},
				map[string]interface{}{"t": dp.Val}, dp.updatetime)
			bp.AddPoint(point)
			if dp.Alarm != 0 { //有报警
				cm.toAlarmLog(dp)
			}

		case <-time.After(time.Second * 1):
			len1 := len(bp.Points())
			if len1 != 0 {
				fmt.Println("LEN:", len1)
				err = cm.influx.Write(bp)
				if err != nil {
					Debug("influx.Write error", err)
				}

				bp, err = client.NewBatchPoints(client.BatchPointsConfig{
					Precision: "ms",
					Database:  cm.MyDB,
				})
				if err != nil {
					Debug("NewBatchPoints", err.Error())
					break
				}
			}
		}
	}
	setWrite(false)
	return nil
}

func (cm *ChannelManager) ping() error {
	_, _, err := cm.influx.Ping(5 * time.Second)
	return err
}

//初始化Inlufx库连接
func (cm *ChannelManager) initInflux() error {
	if cm.influx != nil {
		cm.influx.Close()
	}
	var err error
	cm.influx, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     cm.Influxdb,
		Username: cm.Username,
		Password: cm.Password,
		Timeout:  time.Second * 5,
	})
	if err != nil {
		return err
	}
	err = cm.ping()
	return err

}

//NewChannelManager ---
func NewChannelManager(jsonFile string) (*ChannelManager, error) {
	pChannelManager := &ChannelManager{influx: nil}
	err := pChannelManager.Init(jsonFile)
	return pChannelManager, err
}

//初始化 串口连接
func (ch *Channel) initCom() {
	if ch.client != nil {
		ch.handler.Close()
		ch.client = nil
	}

	//
	ch.handler = modbus.NewRTUClientHandler(ch.Portname)
	ch.handler.BaudRate = 9600
	ch.handler.DataBits = 8
	ch.handler.Parity = "N"
	ch.handler.StopBits = 1
	ch.handler.Timeout = 500 * time.Millisecond
	ch.client = modbus.NewClient(ch.handler)
}

//initParameters 初始化报警参数
func (ch *Channel) initParameters() {
	for _, v := range ch.Device {
		v.initParameters()
	}
}

//Close 关闭串口
func (ch *Channel) closeCom() {
	ch.handler.Close()
	ch.client = nil
}

//ReadTemperature 获取温度
func (ch *Channel) readTemperature(slaveID byte, regaddr uint16) (float32, error) {
	if ch.client == nil {
		return 0.0, errors.New("client error")
	}
	ch.handler.SlaveId = slaveID
	adu, err := ch.client.ReadHoldingRegisters(regaddr, 1)
	var vv float32 = 0.0
	if err != nil {
		str1 := err.Error()
		if strings.Contains(str1, "timeout") {
			Debug("reconnect:", ch.Portname)
			ch.initCom()
		}
		return 0.0, err
	}
	v16 := BytesToInt16(adu[0:2])
	if v16 == 32767 { //没有接设备
		v16 = 0
		vv = float32(v16) / 10.0
		return vv, errors.New("Not Valid")
	}
	vv = float32(v16) / 10.0
	return vv, err
}

//requestData ---
func (ch *Channel) requestData(fun1 onConn) error {
	for _, dev := range ch.Device {
		for _, vdata := range dev.Datapoint {
			if isQuit() {
				fmt.Println("requestData - quit")
				return errors.New("quit")
			}
			ch.getTemp(dev, vdata, fun1)
			chPoint <- vdata
		}
	}
	return nil
}
