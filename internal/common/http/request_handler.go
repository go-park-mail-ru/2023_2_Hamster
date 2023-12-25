package http

import (
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetIDFromRequest(id string, r *http.Request) (uuid.UUID, error) {
	uuidString := mux.Vars(r)[id]
	parsedUUID, err := uuid.Parse(uuidString)
	if err != nil {
		return parsedUUID, errors.New("invalid uuid parameter")
	}
	return parsedUUID, nil
}

func GetloginFromRequest(login string, r *http.Request) string {
	userLogin := mux.Vars(r)[login]

	return userLogin
}

var ErrUnauthorized = &models.UnathorizedError{}

func GetUserFromRequest(r *http.Request) (*models.User, error) {
	user, ok := r.Context().Value(models.ContextKeyUserType{}).(*models.User)
	if !ok {
		return nil, ErrUnauthorized
	}
	if user == nil {
		return nil, ErrUnauthorized
	}

	return user, nil
}

func GetQueryParam(r *http.Request) (*models.QueryListOptions, error) {
	values := r.URL.Query()
	params := &models.QueryListOptions{}
	var err error
	categoryStr := values.Get("category")
	if categoryStr != "" {
		params.Category, err = uuid.Parse(categoryStr)
		if err != nil {
			return nil, errors.New("invalid category UUID")
		}
	}

	accountStr := values.Get("account")
	if accountStr != "" {
		params.Account, err = uuid.Parse(accountStr)
		if err != nil {
			return nil, errors.New("invalid account UUID")
		}
	}

	incomeStr := values.Get("income")
	if incomeStr != "" {
		params.Income, err = strconv.ParseBool(incomeStr)
		if err != nil {
			return nil, errors.New("invalid value for income")
		}
	}

	outcomeStr := values.Get("outcome")
	if outcomeStr != "" {
		params.Outcome, err = strconv.ParseBool(outcomeStr)
		if err != nil {
			return nil, errors.New("invalid value for outcome")
		}
	}

	startDateStr := values.Get("start_date")
	endDateStr := values.Get("end_date")

	if startDateStr != "" {
		params.StartDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return nil, errors.New("invalid start date format")
		}
	}

	if endDateStr != "" {
		params.EndDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			return nil, errors.New("invalid end date format")
		}
	}

	return params, nil
}

func GetIpFromRequest(r *http.Request) (string, error) {
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		// get last IP in list since ELB prepends other user defined IPs, meaning the last one is the actual client IP.
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String(), nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	netIP := net.ParseIP(ip)
	if netIP != nil {
		ip := netIP.String()
		if ip == "::1" {
			return "127.0.0.1", nil
		}
		return ip, nil
	}

	return "", errors.New("IP not found")
}

// func GetGoalIdFromRequest(r *http.Request) {
// 	return
// }
