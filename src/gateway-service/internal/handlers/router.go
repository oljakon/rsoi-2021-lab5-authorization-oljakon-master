package handlers

import (
	"fmt"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"rsoi2/src/gateway-service/pkg/authenticator"
	"rsoi2/src/gateway-service/pkg/jwt"
	"rsoi2/src/gateway-service/pkg/middleware"
)

type ServicesStruct struct {
	PaymentServiceAddress string
	RentalServiceAddress  string
	CarServiceAddress     string

	session *sessions.FilesystemStore
}

type GatewayService struct {
	Config ServicesStruct
}

func NewGatewayService(config *ServicesStruct) *GatewayService {
	return &GatewayService{Config: *config}
}

func Router(store *sessions.FilesystemStore) *mux.Router {
	carsHost := os.Getenv("CARS_SERVICE_SERVICE_HOST")
	carsPort := os.Getenv("CARS_SERVICE_SERVICE_PORT")

	paymentHost := os.Getenv("PAYMENT_SERVICE_SERVICE_HOST")
	paymentPort := os.Getenv("PAYMENT_SERVICE_SERVICE_PORT")

	rentalHost := os.Getenv("RENTAL_SERVICE_SERVICE_HOST")
	rentalPort := os.Getenv("RENTAL_SERVICE_SERVICE_PORT")

	servicesConfig := ServicesStruct{
		//PaymentServiceAddress: "http://payment-service:8082",
		//RentalServiceAddress:  "http://rental-service:8083",
		//CarServiceAddress:     "http://car-service:8081",
		//PaymentServiceAddress: "http://localhost:8082",
		//RentalServiceAddress:  "http://localhost:8083",
		//CarServiceAddress:     "http://localhost:8081",
		PaymentServiceAddress: fmt.Sprintf("http://%s:%s", paymentHost, paymentPort),
		RentalServiceAddress:  fmt.Sprintf("http://%s:%s", rentalHost, rentalPort),
		CarServiceAddress:     fmt.Sprintf("http://%s:%s", carsHost, carsPort),
	}

	httpClient := jwt.JWKHttpClient{}
	verifier := jwt.JwtTokenVerifier{
		JWKSUri:    "https://dev-jnrhuro0.us.auth0.com/.well-known/jwks.json",
		HTTPClient: &httpClient,
	}

	router := mux.NewRouter()

	gs := NewGatewayService(&servicesConfig)

	router.Handle("/api/v1/authorization", authenticator.LoginHandler(store))
	router.Handle("/api/v1/callback", authenticator.CallbackHandler(store))

	router.HandleFunc("/api/v1/cars", gs.GetAvailableCars).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/rental", gs.GetUserRentals).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/rental/{rentalUid}", gs.GetRentalInfo).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/rental", gs.RentCar).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/rental/{rentalUid}/finish", gs.EndRental).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/rental/{rentalUid}", gs.CancelRental).Methods("DELETE", "OPTIONS")

	amw := middleware.Auth{
		JwtTokenVerifier: verifier,
	}
	router.Use(amw.Middleware)

	return router
}
