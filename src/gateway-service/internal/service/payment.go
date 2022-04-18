package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"rsoi2/src/gateway-service/internal/models"
)

func GetPayment(paymentServiceAddress string, paymentUID string) (*models.Payment, error) {
	requestURL := fmt.Sprintf(paymentServiceAddress+"/api/v1/payment/%s", paymentUID)

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Println("failed to create an http request")
		return nil, err
	}

	client := &http.Client{
		Timeout: 10 * time.Minute,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed request to payment service: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body")
		}
	}(res.Body)

	payment := &models.Payment{}
	err = json.NewDecoder(res.Body).Decode(payment)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return payment, nil
}

func CreatePayment(paymentServiceAddress string, price int) (string, error) {
	requestURL := fmt.Sprintf(paymentServiceAddress + "/api/v1/payment")

	uid := uuid.New().String()

	payment := &models.Payment{
		PaymentUID: uid,
		Status:     "PAID",
		Price:      price,
	}

	data, err := json.Marshal(payment)
	if err != nil {
		return "", fmt.Errorf("encoding error: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to create an http request: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Minute,
	}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed request to payment service: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body")
		}
	}(res.Body)

	return uid, nil
}

func CancelPayment(paymentServiceAddress, paymentUid string) error {
	requestURL := fmt.Sprintf(paymentServiceAddress+"/api/v1/payment/%s", paymentUid)

	req, err := http.NewRequest(http.MethodPatch, requestURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create an http request: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Minute,
	}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed request to rental service: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body")
		}
	}(res.Body)

	return nil
}
