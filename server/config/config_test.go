package config

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	file = ".env"

	MissingFile = "MISSING_FILE"
	MissingEnv  = "MISSING_ENV"
	WrongFormat = "WRONG_FORMAT"
)

func createEnvFile(s string) {
	os.Remove(file)

	if s == MissingFile {
		return
	}

	f, _ := os.Create(file)
	if s == WrongFormat {
		f.WriteString(`ENV+"xxxx"`)
		f.Close()
		return
	}

	if s == MissingEnv {
		f.WriteString("\nENVIRONMENT=development")
		f.WriteString("\nPORT=3000")
		f.Close()
		return
	}

	f.WriteString("\nENVIRONMENT=development")
	f.WriteString("\nPORT=3000")
	f.WriteString("\nSECRET_KEY=Key")

	f.WriteString("\nGORM_AUTOMIGRATE=false")
	f.WriteString("\nGORM_LOGMODE=false")
	f.WriteString("\nGORM_DIALECT=mysql")
	f.WriteString("\nGORM_CONNECTION_DSN=localhost")

	f.WriteString("\nREDIS_ADDRESS=localhost:123")
	f.WriteString("\nREDIS_PASSWORD=123")

	f.WriteString("\nMQTT_HOST=localhost")
	f.WriteString("\nMQTT_USER=user")
	f.WriteString("\nMQTT_PASS=pass")
	f.Close()
}

func TestConfig(t *testing.T) {

	Convey("Attempt to load environment variables", t, func() {

		Convey(`when ".env" is missing`, func() {
			createEnvFile(MissingFile)

			env, err := LoadEnv()
			So(env, ShouldBeNil)
			So(err, ShouldBeError)
		})

		Convey(`when ".env" is wrong format`, func() {
			createEnvFile(WrongFormat)

			env, err := LoadEnv()
			So(env, ShouldBeNil)
			So(err, ShouldBeError)
		})

		Convey("when some variables is missing", func() {
			createEnvFile(MissingEnv)

			env, err := LoadEnv()
			So(env, ShouldBeNil)
			So(err, ShouldBeError)
		})

		Convey("when all variables are provided", func() {
			createEnvFile("")

			env, _ := LoadEnv()
			So(env.Environment, ShouldEqual, DevelopmentEnv)
			So(env.Port, ShouldEqual, "3000")
			So(env.SecretKey, ShouldEqual, "Key")
		})

		Reset(func() {
			os.Remove(file)
		})
	})

}
