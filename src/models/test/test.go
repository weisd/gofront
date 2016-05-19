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

	m := &TestGorm{}
	m.Name = name
	m.Passwd = pass

	err := models.Model().Create(m).Error
	if err != nil {
		return nil, models.Error(err)
	}

	return m, nil
}

func InfoByName(name string) (*TestGorm, error) {

	m := &TestGorm{}

	err := models.Model().Where("name = ?", name).First(m).Error
	if err != nil {
		return nil, models.Error(err)
	}

	return m, nil
}

func List() ([]TestGorm, error) {
	var list []TestGorm
	err := models.Model().Find(&list).Error
	if err != nil {
		return list, models.Error(err)
	}
	return list, nil
}
