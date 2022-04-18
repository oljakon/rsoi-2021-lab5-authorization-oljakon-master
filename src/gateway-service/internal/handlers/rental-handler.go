package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"rsoi2/src/gateway-service/internal/models"
	"rsoi2/src/gateway-service/internal/service"
)

var (
	keyJWT = "2IRMcTomN7eHrzqCkLa-zcK1FU5KR6nPCmm5W93anquAIxdL13XdNxIAq4NFQTJFztdmFNGwRvZEu9MS5i4JwUKV_leJTg0L5k3NAzLU_IZ039vzshd37Nl5dPIawhKd1b1vAiwQXjA_OtKDwI40MzCU0Q1bmrqMEt5sifMq6bUwIq6BH1k7EiRHuFsYuuIXhP-j3yx8L2Ba0LLSesP5Tyo3tFycZbw1IRKEVUuakXoK6AyEVj5kptvfZjIU9B0o7y3X6tjxOH15HjaAwLvP54DEHCD2PHmYZp2Py0y5K6CIRug44BiUU_G9XoEKxeEZoTy5JQ5DyB7GEAf2GYn4cQ"
)

type UserClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func (gs *GatewayService) GetUserRentals(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		log.Printf("authorization header is empty")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	s := strings.Split(authorizationHeader, " ")
	tokenString := s[1]

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(keyJWT), nil
	})

	username := claims["name"].(string)

	rentalServiceAddress := gs.Config.RentalServiceAddress
	paymentServiceAddress := gs.Config.PaymentServiceAddress
	carServiceAddress := gs.Config.CarServiceAddress
	rentalsInfo, err := service.UsersRentalWithPaymentController(rentalServiceAddress, paymentServiceAddress, carServiceAddress, username)
	if err != nil {
		log.Printf("failed to get response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(rentalsInfo)
	if err != nil {
		log.Printf("failed to encode response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (gs *GatewayService) GetRentalInfo(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		log.Printf("authorization header is empty")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	s := strings.Split(authorizationHeader, " ")
	tokenString := s[1]

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(keyJWT), nil
	})

	username := claims["name"].(string)

	params := mux.Vars(r)
	rentalUID := params["rentalUid"]

	rentalServiceAddress := gs.Config.RentalServiceAddress
	paymentServiceAddress := gs.Config.PaymentServiceAddress
	carServiceAddress := gs.Config.CarServiceAddress
	rentalsInfo, err := service.UsersRentalFullInfoController(rentalServiceAddress, paymentServiceAddress, carServiceAddress, username, rentalUID)
	if err != nil {
		log.Printf("failed to get response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(rentalsInfo)
	if err != nil {
		log.Printf("failed to encode response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (gs *GatewayService) RentCar(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		log.Printf("authorization header is empty")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	s := strings.Split(authorizationHeader, " ")
	tokenString := s[1]

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(keyJWT), nil
	})

	username := claims["name"].(string)

	var rentInfo models.RentCarRequest
	err = json.NewDecoder(r.Body).Decode(&rentInfo)
	if err != nil {
		fmt.Println("failed to decode post request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rentalServiceAddress := gs.Config.RentalServiceAddress
	paymentServiceAddress := gs.Config.PaymentServiceAddress
	carServiceAddress := gs.Config.CarServiceAddress
	rentedCar, err := service.RentCarController(rentalServiceAddress, paymentServiceAddress, carServiceAddress, username, &rentInfo)
	if err != nil {
		log.Printf("failed to get response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(rentedCar)
	if err != nil {
		log.Printf("failed to encode response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (gs *GatewayService) EndRental(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	rentalUID := params["rentalUid"]

	rentalServiceAddress := gs.Config.RentalServiceAddress
	carServiceAddress := gs.Config.CarServiceAddress
	err := service.EndRentalController(rentalServiceAddress, carServiceAddress, rentalUID)
	if err != nil {
		log.Printf("failed to get response: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (gs *GatewayService) CancelRental(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	rentalUID := params["rentalUid"]

	rentalServiceAddress := gs.Config.RentalServiceAddress
	carServiceAddress := gs.Config.CarServiceAddress
	paymentServiceAddress := gs.Config.PaymentServiceAddress
	err := service.CancelRentalController(rentalServiceAddress, carServiceAddress, paymentServiceAddress, rentalUID)
	if err != nil {
		log.Printf("failed to get response: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
