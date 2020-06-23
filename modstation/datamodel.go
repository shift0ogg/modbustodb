package modstation

import (
	"time"

	"github.com/goburrow/modbus"
	"github.com/influxdata/influxdb/client/v2"
)

//Datapoint ---
type Datapoint struct {
	Alarmcondition string `json:"alarmcondition"`
	Bitoffset      int64  `json:"bitoffset"`
	Circle         int64  `json:"circle"`
	Datatype       int64  `json:"datatype"`
	Dpid           string `json:"dpid"`
	Regaddr        uint16 `json:"regaddr"`
	Reserved       string `json:"reserved"`

	Connected  bool    `json:"connected"`  //是否连接正常
	Valid      bool    `json:"valid"`      //是否有效
	Val        float32 `json:"val"`        //温度值
	Alarm      uint16  `json:"alarm"`      //报警状态
	UpdateTime string  `json:"updatetime"` //最后数据时间【字符串】

	/*请求计数*/
	validNum   int //请求成功计数
	inValidNum int //请求失败计数
	/*温度缓存*/
	valLast     float32   //上次缓存的温度值
	valList     []float32 //历史缓存的温度值列表
	alarmLevels []float32 //报警条件
	updatetime  time.Time //最后数据时间
}

//ChannelManager --
type ChannelManager struct {
	MyDB       string        `json:"mydb"`
	Username   string        `json:"username"`
	Password   string        `json:"password"`
	Influxdb   string        `json:"influxdb"`
	Channels   []*Channel    `json:"channels"`
	MySQL      string        `json:"mysql"`      //mysql连接串
	RequestGap int64         `json:"requestgap"` //485请求时间间隔
	APIPort    uint16        `json:"apiport"`    //API端口
	influx     client.Client //influx连接
	configFile string        //配置文件名称

}

//Channel 是对通道进行管理
type Channel struct {
	Device   []*Device `json:"device"`
	ID       int64     `json:"id"`
	Portname string    `json:"portname"`
	handler  *modbus.RTUClientHandler
	client   modbus.Client
}

//Device ---
type Device struct {
	Addr      byte         `json:"addr"`
	Datapoint []*Datapoint `json:"datapoint"`
}
