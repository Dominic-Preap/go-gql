package service

import (
	"database/sql"

	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type Suite struct {
	DB      *gorm.DB
	mock    sqlmock.Sqlmock
	service *Service
}

func SetupService() Suite {
	var (
		db  *sql.DB
		err error
	)
	s := Suite{}

	db, s.mock, err = sqlmock.New()
	So(err, ShouldBeNil)

	s.DB, err = gorm.Open("postgres", db)
	So(err, ShouldBeNil)

	s.DB.LogMode(true)
	s.service = Init(s.DB)

	return s
}

// func TestServiceFunction(t *testing.T) {
// 	Convey("Service Package", t, func() {
// 		Convey("", func() {
// 		})
// 	})
// }
