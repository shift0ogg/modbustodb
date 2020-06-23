package main

//Datapoint ---
// type Datapoint struct {
// 	Alarmcondition string `json:"alarmcondition"`
// 	Bitoffset      int64  `json:"bitoffset"`
// 	Circle         int64  `json:"circle"`
// 	Datatype       int64  `json:"datatype"`
// 	Dpid           string `json:"dpid"`
// 	Regaddr        uint16 `json:"regaddr"`
// 	Reserved       string `json:"reserved"`
// }

// //ChannelManager --
// type ChannelManager struct {
// 	MyDB     string     `json:"mydb"`
// 	Username string     `json:"username"`
// 	Password string     `json:"password"`
// 	Influxdb string     `json:"influxdb"`
// 	Channels []*Channel `json:"channels"`
// }

// //Channel 是对通道进行管理
// type Channel struct {
// 	ID       int64    `json:"id"`
// 	Portname string   `json:"portname"`
// 	Device   []Device `json:"device"`
// }

// //Device ---
// type Device struct {
// 	Addr      byte         `json:"addr"`
// 	Datapoint []*Datapoint `json:"datapoint"`
// }

// func test() {
// 	jsonParse := NewJsonStruct()
// 	v := ChannelManager{}
// 	//下面使用的是相对路径，config.json文件和main.go文件处于同一目录下
// 	jsonParse.Load("./config.json", &v)
// 	fmt.Println(v.MyDB)

// 	jsonParse.Write("./config1.json", v)

// }

// type JsonStruct struct {
// }

// func NewJsonStruct() *JsonStruct {
// 	return &JsonStruct{}
// }

// func (jst *JsonStruct) Load(filename string, v interface{}) {
// 	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
// 	data, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		return
// 	}

// 	//读取的数据为json格式，需要进行解码
// 	err = json.Unmarshal(data, v)
// 	if err != nil {
// 		return
// 	}
// }

// func (jst *JsonStruct) Write(filename string, dataobj interface{}) {
// 	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
// 	bytes, err := json.Marshal(dataobj)
// 	if err != nil {
// 		return
// 	}
// 	err = ioutil.WriteFile(filename, bytes, 0666)
// 	if err != nil {
// 		return
// 	}
// }

func test() {
	// var s1 = modstation.GetConfigFromDB()
	// fmt.Println(s1)
}
