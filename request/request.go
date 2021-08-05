package request

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jpillora/backoff"
	"github.com/juandiii/jetson-monitor/config"
	"github.com/juandiii/jetson-monitor/logging"
)

type Request struct {
	http http.Client
}

func RequestServer(c config.URL, log *logging.StandardLogger) (string, error) {
	b := &backoff.Backoff{
		Jitter: true,
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: false}
	if c.Timeout == 0 {
		c.Timeout = 15
	}

	client := &http.Client{
		Timeout: time.Duration(time.Duration(c.Timeout) * time.Second),
	}

	req, err := http.NewRequest("GET", c.URL, nil)

	if err != nil {
		log.Error(err)
		return "", errors.New("Received an invalid status code: 500 The service might be experiencing issues")
	}

	for tries := 0; tries < 3; tries++ {
		start := time.Now()
		resp, err := client.Do(req)
		t := time.Now()
		elapsed := t.Sub(start)

		if err != nil {
			d := b.Duration()
			time.Sleep(d)

			if tries == 2 {
				return "", fmt.Errorf("The server %s is down", c.URL)
			}

			continue

		}

		defer resp.Body.Close()

		if c.StatusCode != nil && *c.StatusCode != resp.StatusCode {
			log.Errorf("%s \n", c.URL)
			return "", errors.New("The server " + c.URL + " received an invalid status code: " + strconv.Itoa(resp.StatusCode) + " The service might be experiencing issues")
		}

		content, err := ioutil.ReadAll(resp.Body)
		if c.Match != "" && !strings.Contains(string(content), c.Match) {
			log.Errorf("%s \n", c.URL)
			return "", errors.New("The server " + c.URL + " received an invalid content.The response does not contain " + c.Match)
		}

		if c.ResponseTime != nil {
			responseTimeDuration := time.Duration(*c.ResponseTime) * time.Millisecond
			if responseTimeDuration-elapsed < 0 {
				responseTime := strconv.Itoa(*c.ResponseTime)
				log.Errorf("%s \n", c.URL)
				return "", errors.New("The server " + c.URL + ", Elapsed time: " + elapsed.String() + " instead of " + responseTime)
			}
		}

		log.Debugf("[OK] %s", c.URL)

		return fmt.Sprintf("[OK] %s", c.URL), nil
	}

	return "", errors.New("The request failed because it wasn't able to reach the service")

}
