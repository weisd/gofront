package test

import (
	"models"
	"time"
)

// AUTO_INCREMENT

type TestGorm struct {
	Id        int64  `gorm:"primary_key"`
	Name      string `gorm:"index"`
	Passwd    string
	Status    int `gorm:"default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// 自动创建表
func AutoMigrate() {
	models.Model().AutoMigrate(&TestGorm{})
}

func Add(name, pass string) (*TestGorm, error) {
	db := models.Model()

	m := &TestGorm{}
	m.Name = name
	m.Passwd = pass

	err := db.Create(m).Error
	if err != nil {
		return nil, models.Error(err)
	}

	return m, nil
}

func InfoByName(name string) (*TestGorm, error) {
	db := models.Model()

	m := &TestGorm{}

	err := db.Where("name = ?", name).First(m).Error
	if err != nil {
		return nil, models.Error(err)
	}

	return m, nil
}
