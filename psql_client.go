package main

import (
	"fmt"

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
	EndPositionID   int64 `gorm:"column:end_position_id "`
	StartAddressID  int64 `gorm:"column:start_address_id"`
	EndAddressID    int64 `gorm:"column:end_address_id"`
}

type Position struct {
	ID        int64  `gorm:"column:id"`
	Latitude  string `gorm:"column:latitude"`
	Longitude string `gorm:"column:longitude"`
}

func fixAddrBroken() error {
	return psql.Transaction(func(tx *gorm.DB) error {
		var drives []*Drive
		err := tx.Table("drives").Where("start_address_id IS NULL").Or("end_address_id IS NULL").Find(&drives).Error
		if err != nil {
			return err
		}
		for _, d := range drives {
			startPos, endPos := &Position{}, &Position{}
			tx.Table("positions").Where("id = ?", d.StartPositionID).First(startPos)
			tx.Table("positions").Where("id = ?", d.EndPositionID).First(endPos)

			tx.Table("addresses").Where("latitude = ?", startPos.Latitude).Where("longitude = ?", startPos.Longitude)
			tx.Table("addresses").Where("latitude = ?", endPos.Latitude).Where("longitude = ?", endPos.Longitude)
		}
		return nil
	})
}
