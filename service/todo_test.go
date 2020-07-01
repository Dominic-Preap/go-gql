package service

import (
	"regexp"
	"testing"

	"github.com/my/app/model"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestTodoRepository(t *testing.T) {

	Convey("Todo Repository", t, func() {
		s := SetupService()

		Convey(".FindOne should fail", func() {
			text := "test"
			sql := `SELECT * FROM "todos"`

			s.mock.ExpectQuery(regexp.QuoteMeta(sql)).WithArgs(text)

			res, err := s.service.Todo.FindOne(&TodoFilter{Text: &text})
			So(err, ShouldNotBeNil)
			So(res, ShouldResemble, &model.Todo{})
		})

		Convey(".FindOne should success", func() {
			var (
				id   = 1
				text = "test"
			)
			sql := `SELECT * FROM "todos"`
			rows := sqlmock.NewRows([]string{"id", "text"}).
				AddRow(id, text)

			s.mock.ExpectQuery(regexp.QuoteMeta(sql)).
				WithArgs(text).
				WillReturnRows(rows)

			res, err := s.service.Todo.FindOne(&TodoFilter{Text: &text})
			So(err, ShouldBeNil)
			So(res, ShouldResemble, &model.Todo{ID: id, Text: text})
		})

		Convey(".FindAll should success", func() {
			var (
				id     = 1
				text   = "test"
				limit  = 1
				offset = 1
				todos  = []*model.Todo{}
			)
			sql := `SELECT * FROM "todos"`
			rows := sqlmock.NewRows([]string{"id", "text"}).
				AddRow(id, text)

			s.mock.ExpectQuery(regexp.QuoteMeta(sql)).
				WithArgs(text).
				WillReturnRows(rows)

			res, err := s.service.Todo.FindAll(&TodoFilter{Text: &text, Limit: &limit, Offset: &offset})
			So(err, ShouldBeNil)
			So(res, ShouldResemble, append(todos, &model.Todo{ID: id, Text: text}))
		})

		Convey(".Create should success", func() {
			var (
				id   = 1
				text = "test"
			)
			s.mock.ExpectBegin()
			s.mock.ExpectQuery(regexp.QuoteMeta(
				`INSERT INTO "todos"`)).
				WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg()).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
			s.mock.ExpectCommit()

			res, err := s.service.Todo.Create(&model.Todo{ID: id, Text: text})

			So(err, ShouldBeNil)
			So(res.ID, ShouldEqual, id)
		})

		// This reset is run after each `Convey` at the same scope.
		Reset(func() {
			err := s.mock.ExpectationsWereMet() // make sure all expectations were met
			So(err, ShouldBeNil)
		})
	})
}
