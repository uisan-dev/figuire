package amiami

import (
	"encoding/json"
	"errors"
	"figurebot/figuire/merch"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	APIBase   = "https://api.amiami.com/api/v1.0"
	ImageBase = "https://img.amiami.com"
	UserKey   = "amiami_dev"
	UserAgent = "Mozilla/5.0 (compatible; figuire/0.1)"
)

type Client struct {
	http    *http.Client
	limiter *time.Ticker
}

func NewClient() *Client {
	return &Client{
		http:    &http.Client{Timeout: 10 * time.Second},
		limiter: time.NewTicker(2 * time.Second),
	}
}

func (c *Client) Close() {
	c.limiter.Stop()
}

func (c *Client) fetch(gCode string) (merch.Item, error) {
	<-c.limiter.C

	params := url.Values{"gcode": {gCode}, "lang": {"eng"}}
	req, err := http.NewRequest("GET", APIBase+"/item?"+params.Encode(), nil)
	if err != nil {
		return merch.Item{}, err
	}

	req.Header.Set("X-User-Key", UserKey)
	req.Header.Set("User-Agent", UserAgent)

	resp, err := c.http.Do(req)
	if err != nil {
		return merch.Item{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		se := &StatusError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			GCode:      gCode,
			Body:       strings.TrimSpace(string(body)),
		}
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			if secs, err := strconv.Atoi(ra); err == nil {
				se.RetryAfter = time.Duration(secs) * time.Second
			}
		}
		return merch.Item{}, se
	}

	if ct := resp.Header.Get("Content-Type"); !strings.Contains(ct, "json") {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return merch.Item{}, fmt.Errorf("Expected JSON for fetch of %s, recieved %s: %s", gCode, ct, body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return merch.Item{}, err
	}

	var ar apiResponse
	if err := json.Unmarshal(body, &ar); err != nil {
		return merch.Item{}, err
	}

	if !ar.RSuccess || ar.Item == nil {
		return merch.Item{}, fmt.Errorf("%s: %w", gCode, ErrItemNotFound)
	}

	return ar.Item.toMerchItem(), nil
}

func (c *Client) FetchItem(gCode string) (merch.Item, error) {
	var lastErr error

	for attempt := range 3 {
		if attempt > 0 {
			wait := time.Duration(1<<attempt) * time.Second

			var se *StatusError
			if errors.As(lastErr, &se) && se.RetryAfter > 0 {
				wait = se.RetryAfter + time.Second
			}
			log.Printf("Retry %d/3 for %s in %s", attempt, gCode, wait)
			time.Sleep(wait)
		}

		item, err := c.fetch(gCode)
		if err == nil {
			return item, nil
		}
		lastErr = err

		if errors.Is(err, ErrItemNotFound) {
			break
		}

		var se *StatusError
		if errors.As(err, &se) && !se.Retryable() {
			break
		}
	}

	return merch.Item{}, lastErr
}
