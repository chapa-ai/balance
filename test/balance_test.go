package test

import (
	"balance/pkg/models"
	"fmt"
	"github.com/go-resty/resty/v2"
	"testing"
)

var (
	url = "http://localhost:9999"
)

func TestSendMoney(t *testing.T) {
	t.Parallel()

	data1 := &models.Balance{
		UserId: "7",
		Sum:    50,
	}
	data2 := &models.Balance{
		UserId: "7",
		Sum:    -50,
	}

	balance, err := SelectBalance(data1)
	if err != nil {
		t.Fatal(err)
	}

	postiveTotal, err := SendMoney(data1)
	if err != nil {
		t.Fatal(err)
	}

	if balance.Balance == postiveTotal.Balance {
		t.Fatalf("balances must be different")
	}

	/// now checking for negative sum
	bal, err := SelectBalance(data2)
	if err != nil {
		t.Fatal(err)
	}

	negativeTotal, err := SendMoney(data2)
	if err != nil {
		t.Fatal(err)
	}

	if bal.Balance == negativeTotal.Balance {
		t.Fatalf("balances must be different")
	}

}

func TestPayBetweenUsers(t *testing.T) {
	t.Parallel()

	data := &models.Balance{
		UserId:  "7",
		UserId2: "9",
		Sum:     5,
	}

	usersBalance, err := PayBetweenUsers(data)
	if err != nil {
		t.Fatal(err)
	}

	bal := FindBalanceInArray(usersBalance, data.UserId2)
	if bal == nil {
		fmt.Print("empty balance")
	}
}

func TestSelectBalance(t *testing.T) {
	t.Parallel()

	data := &models.Balance{
		UserId: "7",
		Sum:    150,
	}

	_, err := SelectBalance(data)
	if err != nil {
		t.Fatal(err)
	}
}

func SelectBalance(data *models.Balance) (*models.Balance, error) {
	output := &models.Balance{}

	url := fmt.Sprintf("%v/balance/get", url)

	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		SetResult(output).
		Post(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("status code wrong. status: %v. body: %v", resp.StatusCode(), resp.String())
	}
	return output, nil

}

func PayBetweenUsers(data *models.Balance) ([]*models.Balance, error) {
	output := []*models.Balance{}

	url := fmt.Sprintf("%v/balance/users", url)

	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		SetResult(&output).
		Post(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("status code wrong. status: %v. body: %v", resp.StatusCode(), resp.String())
	}

	return output, nil
}

func SendMoney(data *models.Balance) (*models.Balance, error) {
	output := &models.Balance{}

	url := fmt.Sprintf("%v/balance/pay", url)

	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		SetResult(output).
		Post(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("status code wrong. status: %v. body: %v", resp.StatusCode(), resp.String())
	}
	return output, nil
}

func FindBalanceInArray(list []*models.Balance, userId string) *models.Balance {
	for _, item := range list {
		if item.UserId == userId {
			return item
		}
	}
	return nil
}
