package database

import (
	"time"
)

type Clone struct {
	Day        time.Time `gorm:"primaryKey"`
	Count      int
	Uniques    int
	Repository string `gorm:"primaryKey"`
}

type View struct {
	Day        time.Time `gorm:"primaryKey"`
	Count      int
	Uniques    int
	Repository string `gorm:"primaryKey"`
}

type Path struct {
	Path       string `gorm:"primaryKey"`
	Title      string
	Count      int
	Uniques    int
	Day        time.Time `gorm:"primaryKey"`
	Repository string    `gorm:"primaryKey"`
}

type Referral struct {
	Referrer   string `gorm:"primaryKey"`
	Count      int
	Uniques    int
	Day        time.Time `gorm:"primaryKey"`
	Repository string    `gorm:"primaryKey"`
}
