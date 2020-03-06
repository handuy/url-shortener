package main

import (
	"log"
	"os"

	"github.com/go-pg/pg/v9/orm"
	"github.com/go-pg/pg/v9"
	jsoniter "github.com/json-iterator/go"
	"github.com/kataras/iris/v12"
	"github.com/bradialabs/shortid"
)

type Config struct {
	Database struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
		Address  string `json:"address"`
	} `json:"database"`
}

type Url struct {
	Id         int
	OriginUrl  string
	ShortenUrl string
}

type ShortenReq struct {
	Path string
}

func DecodeDataFromJsonFile(f *os.File, data interface{}) error {
	jsonParser := jsoniter.NewDecoder(f)
	err := jsonParser.Decode(&data)
	if err != nil {
		return err
	}

	return nil
}

func SetupConfig() Config {
	var conf Config

	// Đọc file config.dev.json
	configFile, err := os.Open("config.default.json")
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	// Parse dữ liệu JSON và bind vào Controller
	err = DecodeDataFromJsonFile(configFile, &conf)
	if err != nil {
		log.Println("Không đọc được file config.")
		panic(err)
	}

	return conf
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*Url)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	app := iris.Default()
	tmpl := iris.HTML("./view", ".html")
	app.RegisterView(tmpl)

	// Kết nối CSDL
	config := SetupConfig()
	dbConfig := config.Database

	db := pg.Connect(&pg.Options{
		User:     dbConfig.User,
		Password: dbConfig.Password,
		Database: dbConfig.Database,
		Addr:     dbConfig.Address,
	})
	defer db.Close()

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	app.Get("/", func(ctx iris.Context) {
		ctx.View("index.html")
	})

	app.Post("/shorten", func(ctx iris.Context) {
		var req ShortenReq
		err := ctx.ReadJSON(&req)
		if err != nil {
			log.Println(err)
			return
		}

		var newUrl Url
		newUrl.OriginUrl = req.Path
		// Tạo unique string
		s := shortid.New()
		shortId := s.Generate()
		newUrl.ShortenUrl = "/" + shortId
		_, err = db.Model(&newUrl).Returning("*").Insert()
		if err != nil {
			log.Println(err)
			return
		}

		ctx.JSON(newUrl)
		ctx.StatusCode(iris.StatusOK)
	})

	app.Get("/{id}", func(ctx iris.Context) {
		shortId := "/" + ctx.Params().Get("id")

		var originUrl string
		_, err := db.Query(&originUrl, `SELECT origin_url FROM urls WHERE shorten_url = ?`, shortId)
		if err != nil {
			log.Println(err)
			return
		}

		ctx.Redirect(originUrl)
	})

	app.Run(iris.Addr(":8080"))
}
