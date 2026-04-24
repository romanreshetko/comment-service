package serviceIntegrations

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func UpdateUserPoints(userID int64, points int64) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered: ", r)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	baseURL := "http://auth-service:8080/user/points/update"
	params := url.Values{}
	params.Add("user_id", strconv.FormatInt(userID, 10))
	params.Add("points", strconv.FormatInt(points, 10))
	fullURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "PATCH", fullURL, nil)
	if err != nil {
		log.Printf("error creating request: %v", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SERVICE_TOKEN"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error executing request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("auth-service returned error: %d", resp.StatusCode)
		return
	}

	log.Println("request to auth-service succeeded")
}
