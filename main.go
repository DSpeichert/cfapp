package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/urfave/cli.v1"
)

var (
	db         *gorm.DB
	notify_url string
)

// Generate documentation:
// - go get -u github.com/go-swagger/go-swagger/cmd/swagger
// - swagger generate spec -o ./swagger.json
func main() {
	app := cli.NewApp()
	app.Name = "cfapp"
	app.Usage = "A naive, insecure HTTP REST API server"
	app.Version = "0.0.1-dev"
	app.Author = "Daniel Speichert"
	app.Email = "daniel@speichert.pl"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "mysql, m", Value: "user:pass@/dbname?charset=utf8&parseTime=True&loc=Local",
			Usage: "The MySQL connection URI, see https://github.com/go-sql-driver/mysql#dsn-data-source-name", EnvVar: "MYSQL_URI"},
		cli.StringFlag{Name: "notify, n", Value: "",
			Usage: "The URL to notify on certificate updates", EnvVar: "NOTIFY_URL"},
		cli.StringFlag{Name: "listen, l", Value: ":80",
			Usage: "IP:PORT to listen on for HTTP interface", EnvVar: "HTTP_LISTEN"},
		cli.StringFlag{Name: "salt, s", Value: "defaultsalt",
			Usage: "Salt used for password", EnvVar: "HTTP_LISTEN"},
	}

	app.Action = func(c *cli.Context) error {
		if c.String("mysql") == "" {
			cli.ShowSubcommandHelp(c)
			os.Exit(1)
		}
		notify_url = c.String("notify")

		// open MySQL connection
		var err error
		db, err = gorm.Open("mysql", c.String("mysql"))
		if err != nil {
			return cli.NewExitError("can't open database "+c.String("mysql")+": "+err.Error(), 2)
		}
		defer db.Close()
		db.LogMode(true)

		// sync database schema
		db.AutoMigrate(&Customer{}, &Certificate{})
		db.Model(&Certificate{}).AddForeignKey("customer_id", "customers(id)", "CASCADE", "CASCADE")

		// for password salt generation
		rand.Seed(time.Now().UTC().UnixNano())

		// build HTTP server
		e := echo.New()

		// Middleware
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
		e.Use(middleware.CORS()) // allows all - for Swagger UI

		// Routes
		e.GET("/certificates", CertificateIndex)
		e.POST("/certificates", CertificateCreate)
		e.GET("/certificates/:id", CertificateShow)
		e.PUT("/certificates/:id", CertificateUpdate)
		e.GET("/customers", CustomerIndex)
		e.POST("/customers", CustomerCreate)
		e.GET("/customers/:id", CustomerShow)
		e.DELETE("/customers/:id", CustomerDelete)
		e.File("/swagger.json", "swagger.json")

		// Start server
		println("Test with http://petstore.swagger.io/?url=http://localhost/swagger.json")
		e.Logger.Fatal(e.Start(c.String("listen")))
		return nil
	}

	app.Run(os.Args)
}
