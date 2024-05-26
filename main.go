package main

import (
	"ZED-Magdy/Delivery-go/Handlers"
	"ZED-Magdy/Delivery-go/Middlewares"
	"ZED-Magdy/Delivery-go/Models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	db := initDb()
	trans, validate := initValidator()

	r := mux.NewRouter()
	r.Use(Middlewares.CORSMethodMiddleware)
	// * Routes
	registerRoutes(r, db, validate, trans)

	srv := &http.Server{
		Addr:         "127.0.0.1:8000",
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func registerRoutes(r *mux.Router, db *gorm.DB, validate *validator.Validate, trans ut.Translator) {
	/**
	* Handlers initialization
	 */
	initInfoHandler := Handlers.NewInitInfoHandler(*db)
	userHandler := Handlers.NewUserHandler(*db, validate, trans)

	/**
	 * Routes registration
	 */

	authRoutes := mux.NewRouter().PathPrefix("/api").Subrouter()
	authRoutes.Use(Middlewares.AuthMiddleware)

	authRoutes.HandleFunc("/init-info", initInfoHandler.GetInitialInformation).Methods("GET")
	authRoutes.HandleFunc("/user", userHandler.GetCurrentUser).Methods("GET")

	r.HandleFunc("/api/register", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/api/login", userHandler.Login).Methods("POST")

	r.PathPrefix("/api").Handler(authRoutes)
	http.Handle("/", r)
}

func initValidator() (ut.Translator, *validator.Validate) {
	en := en.New()
	uni := ut.New(en, en)

	trans, _ := uni.GetTranslator("en")

	validate := validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)
	return trans, validate
}

func initDb() *gorm.DB {
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Models.User{}, &Models.PriceList{}, &Models.Region{})
	seed(db)
	return db
}
func seed(db *gorm.DB) {
	user := &Models.User{}
	db.First(user)
	if user.ID == 0 {
		db.Create(&Models.PriceList{Name: "Standard", KmCost: 0.5, CancellationCost: 0.5})
		db.Exec("INSERT INTO regions (name, geofence, price_list_id) VALUES (?, ST_GeomFromText(?), ?)", "Cairo", "POLYGON((30.0444 31.2357, 30.0444 30.0444, 31.2357 30.0444, 31.2357 31.2357, 30.0444 31.2357))", 1)
	}

}
