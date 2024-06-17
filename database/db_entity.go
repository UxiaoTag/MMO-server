package database

import "gorm.io/gorm"

type DbPlayer struct {
	gorm.Model
	Username string
	Password string
	Coin     int
}

type DbCharacter struct {
	gorm.Model
	JobId    int
	Name     string
	Hp       int
	Mp       int
	Level    int
	Exp      int
	SpaceId  int
	X        int
	Y        int
	Z        int
	Gold     int64
	PlayerId int
}

func NewDbCharacter() *DbCharacter {
	return &DbCharacter{
		Hp:    100,
		Mp:    100,
		Level: 1,
	}
}
