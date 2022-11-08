package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var psql *gorm.DB

func initPSql(host, port, user, password, db string) error {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, password, db, port)
	var err error
	if psql, err = gorm.Open(postgres.Open(dsn)); err != nil {
		return err
	}
	return nil
}

type Drive struct {
	ID              int64 `gorm:"column:id"`
	StartPositionID int64 `gorm:"column:start_position_id"`
	EndPositionID   int64 `gorm:"column:end_position_id"`
	StartAddressID  int64 `gorm:"column:start_address_id"`
	EndAddressID    int64 `gorm:"column:end_address_id"`
}

type ChargingProcess struct {
	ID         int64 `gorm:"column:id"`
	PositionID int64 `gorm:"column:position_id"`
	AddressID  int64 `gorm:"column:address_id"`
}

type Position struct {
	ID        int64   `gorm:"column:id"`
	Latitude  float64 `gorm:"column:latitude"`
	Longitude float64 `gorm:"column:longitude"`
}

type Address struct {
	ID            int64          `gorm:"column:id"`
	DisplayName   string         `gorm:"column:display_name"`
	Latitude      float64        `gorm:"column:latitude"`
	Longitude     float64        `gorm:"column:longitude"`
	Name          string         `gorm:"column:name"`
	HouseNumber   sql.NullString `gorm:"column:house_number"`
	Road          sql.NullString `gorm:"column:road"`
	Neighbourhood sql.NullString `gorm:"column:neighbourhood"`
	City          sql.NullString `gorm:"column:city"`
	County        sql.NullString `gorm:"column:county"`
	Postcode      sql.NullString `gorm:"column:postcode"`
	State         sql.NullString `gorm:"column:state"`
	StateDistrict sql.NullString `gorm:"column:state_district"`
	Country       sql.NullString `gorm:"column:country"`
	Raw           []byte         `gorm:"column:raw"`
	InsertedAt    time.Time      `gorm:"column:inserted_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
	OsmID         int64          `gorm:"column:osm_id"`
	OsmType       string         `gorm:"column:osm_type"`
}

func saveBrokenAddr() error {
	return psql.Transaction(func(tx *gorm.DB) error {

		// find drive graph broken addresses
		var drives []*Drive
		err := tx.Table("drives").Where("start_address_id IS NULL").Or("end_address_id IS NULL").Find(&drives).Error
		if err != nil {
			return err
		}

		positions := []*Position{}
		for _, d := range drives {
			startPos, endPos := &Position{}, &Position{}

			if err := tx.Table("positions").Where("id = ?", d.StartPositionID).First(startPos).Error; err != nil {
				log.Printf("start position not found, id=%v, drive=%+v", d.StartPositionID, d)
			} else {
				positions = append(positions, startPos)
			}

			if err := tx.Table("positions").Where("id = ?", d.EndPositionID).First(endPos).Error; err != nil {
				log.Printf("end position not found, id=%v, drive=%+v", d.EndPositionID, d)
			} else {
				positions = append(positions, endPos)
			}
		}

		// find charge graph broken addresses
		var charges []*ChargingProcess
		err = tx.Table("charging_processes").Where("address_id IS NULL").Find(&charges).Error
		if err != nil {
			return err
		}

		for _, c := range charges {
			pos := &Position{}
			if err := tx.Table("positions").Where("id = ?", c.PositionID).First(pos).Error; err != nil {
				log.Printf("charge position not found, id=%v, charge=%+v", c.PositionID, c)
			} else {
				positions = append(positions, pos)
			}
		}

		// fix addresses by positions
		for _, p := range positions {
			osmAddr, err := getAddressByProxy(p.Latitude, p.Longitude)
			if err != nil {
				log.Printf("get address from osm failed, lat=%v, lon=%v, err=%#v", p.Latitude, p.Longitude, err)
				continue
			}

			var exist int64
			tx.Table("addresses").Where("osm_id = ?", osmAddr.OsmID).Where("osm_type = ?", osmAddr.OsmType).Count(&exist)
			if exist > 0 {
				continue
			}

			dnames := strings.Split(osmAddr.DisplayName, ",")
			var name string
			if len(dnames) > 0 {
				name = strings.TrimSpace(dnames[0])
			}
			raw, _ := json.Marshal(osmAddr.Address)
			newAddr := &Address{
				DisplayName:   osmAddr.DisplayName,
				Latitude:      p.Latitude,
				Longitude:     p.Longitude,
				Name:          name,
				HouseNumber:   getOrNull(osmAddr.Address, "housenumber"),
				Road:          getOrNull(osmAddr.Address, "road"),
				Neighbourhood: getOrNull(osmAddr.Address, "neighbourhood"),
				City:          getOrNull(osmAddr.Address, "city"),
				County:        getOrNull(osmAddr.Address, "county"),
				Postcode:      getOrNull(osmAddr.Address, "postcode"),
				State:         getOrNull(osmAddr.Address, "state"),
				StateDistrict: getOrNull(osmAddr.Address, "state_district"),
				Country:       getOrNull(osmAddr.Address, "country"),
				Raw:           raw,
				InsertedAt:    time.Now(),
				UpdatedAt:     time.Now(),
				OsmID:         int64(osmAddr.OsmID),
				OsmType:       osmAddr.OsmType,
			}

			err = tx.Table("addresses").Create(newAddr).Error
			if err == nil {
				log.Printf("save address success, addr=%+v", newAddr)
			}
		}
		return nil
	})
}

func getOrNull(m map[string]interface{}, key string) sql.NullString {
	v, ok := m[key]
	if !ok {
		return sql.NullString{}
	}
	cv, ok := v.(string)
	if !ok {
		return sql.NullString{}
	}
	return sql.NullString{
		String: cv,
		Valid:  true,
	}
}

func fixAddrBroken() error {
	return psql.Transaction(func(tx *gorm.DB) error {

		// fix drives
		var drives []*Drive
		err := tx.Table("drives").Where("start_address_id IS NULL").Or("end_address_id IS NULL").Find(&drives).Error
		if err != nil {
			return err
		}
		for _, d := range drives {
			startPos, endPos := &Position{}, &Position{}
			tx.Table("positions").Where("id = ?", d.StartPositionID).First(startPos)
			tx.Table("positions").Where("id = ?", d.EndPositionID).First(endPos)

			osmStartAddr, err := getAddressByProxy(startPos.Latitude, startPos.Longitude)
			if err != nil {
				log.Printf("get address from osm failed, lat=%v, lon=%v, err=%#v", startPos.Latitude, startPos.Longitude, err)
				continue
			}
			osmEndAddr, err := getAddressByProxy(endPos.Latitude, endPos.Longitude)
			if err != nil {
				log.Printf("get address from osm failed, lat=%v, lon=%v, err=%#v", endPos.Latitude, endPos.Longitude, err)
				continue
			}

			startAddr, endAddr := &Address{}, &Address{}
			tx.Table("addresses").Where("osm_id = ?", osmStartAddr.OsmID).Where("osm_type = ?", osmStartAddr.OsmType).First(startAddr)
			if startAddr.ID > 0 {
				err := tx.Table("drives").Where("id = ?", d.ID).Update("start_address_id", startAddr.ID).Error
				if err == nil {
					log.Printf("fix address success, drives id=%v, fix start addr=%v", d.ID, startAddr.DisplayName)
				}
			}

			tx.Table("addresses").Where("osm_id = ?", osmEndAddr.OsmID).Where("osm_type = ?", osmEndAddr.OsmType).First(endAddr)
			if endAddr.ID > 0 {
				err := tx.Table("drives").Where("id = ?", d.ID).Update("end_address_id", endAddr.ID).Error
				if err == nil {
					log.Printf("fix address success, drives id=%v, fix end addr=%v", d.ID, endAddr.DisplayName)
				}
			}
		}

		// fix charges
		var charges []*ChargingProcess
		err = tx.Table("charging_processes").Where("address_id IS NULL").Find(&charges).Error
		if err != nil {
			return err
		}

		for _, c := range charges {
			pos := &Position{}
			tx.Table("positions").Where("id = ?", c.PositionID).First(pos)

			osmAddr, err := getAddressByProxy(pos.Latitude, pos.Longitude)
			if err != nil {
				log.Printf("get address from osm failed, lat=%v, lon=%v, err=%#v", pos.Latitude, pos.Longitude, err)
				continue
			}

			addr := &Address{}
			tx.Table("addresses").Where("osm_id = ?", osmAddr.OsmID).Where("osm_type = ?", osmAddr.OsmType).First(addr)
			if addr.ID > 0 {
				err := tx.Table("charging_processes").Where("id = ?", c.ID).Update("address_id", addr.ID).Error
				if err == nil {
					log.Printf("fix address success, charge id=%v, fix addr=%v", c.ID, addr.DisplayName)
				}
			}
		}
		return nil
	})
}
