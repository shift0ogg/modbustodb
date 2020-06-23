package modstation

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

//refreshDpsFromDB 从数据刷新数据采集点信息
func (cm *ChannelManager) getDpsFromDB() ([]*Channel, error) {

	var (
		did int64
		pid int64
	)

	db, err := sql.Open("mysql", cm.MySQL)
	defer db.Close()
	if err != nil {
		Debug("sql.Open mysql error:", err)
		return nil, err
	}
	rows, err := db.Query("select id,connstr from iot_channel")
	defer rows.Close()
	if err != nil {
		Debug("sql.db.Query mysql error:", err)
		return nil, err
	}

	chs := []*Channel{}
	for rows.Next() {
		ch1 := Channel{client: nil}
		rows.Scan(&(ch1.ID), &(ch1.Portname))
		if err != nil {
			Debug("sql.rows.Scan mysql error:", err)
			return nil, err
		}
		rows1, err := db.Query("select id,addr from iot_device where channelid = ?", ch1.ID)
		defer rows1.Close()
		if err != nil {
			Debug("sql.db.Query mysql error:", err)
			return nil, err
		}
		devs := []*Device{}
		for rows1.Next() {
			dev1 := &Device{}
			err = rows1.Scan(&did, &(dev1.Addr))
			if err != nil {
				Debug("sql.rows2.Scan mysql error:", err)
				return nil, err
			}
			rows2, err := db.Query("select id,dataregaddr , dpid , alarmcondition as alc from iot_datapoint where deviceid = ?", did)
			defer rows2.Close()
			if err != nil {
				Debug("sql.db.Query mysql error:", err)
				return nil, err
			}
			dps := []*Datapoint{}
			for rows2.Next() {
				dp1 := Datapoint{}
				err = rows2.Scan(&pid, &dp1.Regaddr, &(dp1.Dpid), &(dp1.Alarmcondition))
				if err != nil {
					Debug("sql.rows.Scan mysql error:", err)
					return nil, err
				}
				dps = append(dps, &dp1)
			}
			dev1.Datapoint = dps
			devs = append(devs, dev1)
		}
		ch1.Device = devs
		chs = append(chs, &ch1)
	}
	return chs, nil

}

//toAlarmLog 对报警的判断是否已经处理，未处理则报警入库
func (cm *ChannelManager) toAlarmLog(dp1 *Datapoint) error {
	if dp1.Alarm == 0 {
		return nil
	}

	db, err := sql.Open("mysql", cm.MySQL)
	defer db.Close()
	if err != nil {
		Debug("sql.Open mysql error:", err)
		return err
	}
	/**判断是否已经有报警*/
	rows, err := db.Query("select * from iot_alarmlog where dpid = ? and state = 0 and  ctime between date_sub(now(),interval 1 day) and now() order by ctime desc", dp1.Dpid)
	defer rows.Close()
	if err != nil {
		Debug("sql.db.Query mysql error:", err)
	}
	if !rows.Next() { //有
		floatStr := fmt.Sprintf("%.1f", dp1.Val)
		_, err := db.Exec("insert into iot_alarmlog(dpid , alarmtype ,alarmval, alarmdesc) values(?,?,?,?)", dp1.Dpid, dp1.Alarm, floatStr, "")
		if err != nil {
			Debug("sql.db.Exec mysql error:", err)
		}
	}
	return nil
}

func (cm *ChannelManager) toSysLog(dp1 *Datapoint) error {
	db, err := sql.Open("mysql", cm.MySQL)
	defer db.Close()
	if err != nil {
		Debug("sql.Open mysql error:", err)
		return err
	}
	opdesc := "连接成功"
	if !dp1.Connected {
		opdesc = "连接失败"
	}
	_, err = db.Exec("insert into iot_syslog(dpid , opdesc ,memo) values(?,?,?)", dp1.Dpid, opdesc, "")
	if err != nil {
		Debug("sql.db.Exec mysql error:", err)
		return err
	}
	return nil
}
