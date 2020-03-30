package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	rt "github.com/qiangxue/fasthttp-routing"
)

type DeviceConfig struct {
	Resolution       int
	WhiteBalanceMode int
	NormalDelay      int
	RetakeDelay      int
}

func (s *DeviceConfig) String() string {
	return fmt.Sprintf("resolution=%d\nwb_mode=%d\nnormal_delay=%d\nretake_delay=%d", s.Resolution, s.WhiteBalanceMode, s.NormalDelay, s.RetakeDelay)
}

func (s *DeviceConfig) ScanRow(rows *sql.Rows) error {
	return rows.Scan(&s.Resolution, &s.WhiteBalanceMode, &s.NormalDelay, &s.RetakeDelay)
}

type DeviceConfigCenter struct {
	dbPath string
	db     *sql.DB
}

func (d *DeviceConfigCenter) connect(path string) (reterr error) {
	d.dbPath = path
	connectStr := fmt.Sprintf("file:%s?cache=shared&mode=ro", path)
	d.db, reterr = sql.Open("sqlite3", connectStr)
	return
}

func (d *DeviceConfigCenter) fetchConfigueByID(id string) (ret DeviceConfig, reterr error) {
	cmd := fmt.Sprintf("select Resolution, WBMode, NormalDelay, RetakeDelay from devices where ID = \"%s\"", id)
	res, err := d.db.Query(cmd)
	if err != nil {
		reterr = err
		return
	}
	defer res.Close()

	if !res.Next() {
		reterr = fmt.Errorf("No device named %s", id)
		return
	}

	reterr = ret.ScanRow(res)
	return
}

var GDevCenter DeviceConfigCenter

func handleConfigQuery(ctx *rt.Context) (err error) {
	log.Printf("Getting config...")

	deviceName := string(ctx.QueryArgs().Peek("device_id"))
	log.Printf("deviceName = %s", deviceName)

	if deviceName == "" {
		log.Printf("empty device name")
		ctx.SetStatusCode(500)
	} else {
		if cfg, err := GDevCenter.fetchConfigueByID(deviceName); err != nil {
			log.Printf("Error in fetching configuration, %v", err)
			ctx.SetStatusCode(500)
		} else {
			log.Printf("Configuration dispatched!")
			ctx.SetStatusCode(200)
			ctx.SetBody([]byte(cfg.String()))
		}
	}

	return
}
