package service

import (
	"regexp"
	"testing"

	"github.com/my/app/model"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUserRepository(t *testing.T) {

	Convey("User Repository", t, func() {
		s := SetupService()

		Convey(".FindOne should fail", func() {
			name := "test"
			sql := `SELECT * FROM "users"`

			s.mock.ExpectQuery(regexp.QuoteMeta(sql)).WithArgs(name)

			res, err := s.service.User.FindOne(&UserFilter{Name: name})
			So(err, ShouldNotBeNil)
			So(res, ShouldResemble, &model.User{})
		})

		Convey(".FindOne should success", func() {
			var (
				id   = 1
				name = "test"
			)
			sql := `SELECT * FROM "users"`
			rows := sqlmock.NewRows([]string{"id", "name"}).
				AddRow(id, name)

			s.mock.ExpectQuery(regexp.QuoteMeta(sql)).
				WithArgs(name).
				WillReturnRows(rows)

			res, err := s.service.User.FindOne(&UserFilter{Name: name})
			So(err, ShouldBeNil)
			So(res, ShouldResemble, &model.User{ID: id, Name: name})
		})

		Convey(".FindAll should success", func() {
			var (
				id     = 1
				name   = "test"
				limit  = 1
				offset = 1
				users  = []*model.User{}
			)
			sql := `SELECT * FROM "users"`
			rows := sqlmock.NewRows([]string{"id", "name"}).
				AddRow(id, name)

			s.mock.ExpectQuery(regexp.QuoteMeta(sql)).
				WithArgs(name).
				WillReturnRows(rows)

			res, err := s.service.User.FindAll(&UserFilter{Name: name, Limit: &limit, Offset: &offset})
			So(err, ShouldBeNil)
			So(res, ShouldResemble, append(users, &model.User{ID: id, Name: name}))
		})

		// This reset is run after each `Convey` at the same scope.
		Reset(func() {
			err := s.mock.ExpectationsWereMet() // make sure all expectations were met
			So(err, ShouldBeNil)
		})
	})
}
