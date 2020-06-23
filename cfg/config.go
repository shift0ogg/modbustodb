package cfg

// import (
// 	"encoding/json"
// 	"io/ioutil"
// )

// type Datapoint struct {
// 	Alarmcondition string `json:"alarmcondition"`
// 	Bitoffset      int64  `json:"bitoffset"`
// 	Circle         int64  `json:"circle"`
// 	Datatype       int64  `json:"datatype"`
// 	Dpid           string `json:"dpid"`
// 	Regaddr        uint16 `json:"regaddr"`
// 	Reserved       string `json:"reserved"`
// }

// type Device struct {
// 	Addr      string      `json:"addr"`
// 	Datapoint []Datapoint `json:"datapoint"`
// }

// //Channel model
// type Channel struct {
// 	Device   []Device `json:"device"`
// 	ID       int64    `json:"id"`
// 	Portname string   `json:"portname"`
// }

// //ChannelSlice --
// type Channels struct {
// 	MyDB     string     `json:"mydb"`
// 	Username string     `json:"username"`
// 	Password string     `json:"password"`
// 	Influxdb string     `json:"influxdb"`
// 	Channels []*Channel `json:"channels"`
// }

// func Test(cfgfile string) {
// 	var s Channels
// 	data, err := ioutil.ReadFile(cfgfile)
// 	if err != nil {
// 		println(err.Error())
// 		return
// 	}

// 	//读取的数据为json格式，需要进行解码
// 	err = json.Unmarshal(data, &s)
// 	if err != nil {
// 		println(err.Error())
// 		return
// 	}
// 	println(s.Channels[0].Portname)
// }
