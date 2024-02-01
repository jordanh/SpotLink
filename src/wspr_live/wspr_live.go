package wspr_live

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"text/template"
)

func QueryByCallsign(callsign string, options QueryByCallsignOptions) (ApiResponse, error) {
	// Define your SQL query
	queryTmplStr := `SELECT *
FROM wspr.rx
WHERE tx_sign = '{{.Callsign}}'
AND time > '{{.FromTime}}' AND time <= '{{.ToTime}}'
AND snr > {{.MinSnr}}
ORDER BY time DESC, snr DESC
LIMIT {{.Limit}}
FORMAT JSONCompact
`

	queryTmpl, err := template.New("query").Parse(queryTmplStr)
	if err != nil {
		panic(err)
	}

	data := struct {
		Callsign string
		FromTime string
		ToTime   string
		Limit    int
		MinSnr   int
	}{
		Callsign: callsign,
		FromTime: options.FromTime.UTC().Format("2006-01-02 15:04:05"),
		ToTime:   options.ToTime.UTC().Format("2006-01-02 15:04:05"),
		Limit:    options.Limit,
		MinSnr:   options.MinSnr,
	}

	var queryBuf bytes.Buffer
	err = queryTmpl.Execute(&queryBuf, data)
	if err != nil {
		panic(err)
	}

	query := queryBuf.String()
	// fmt.Println("Query:\n", query)

	// Encode the query
	encodedQuery := url.QueryEscape(query)

	// Assuming you have a base URL for your request
	baseURL := "https://db1.wspr.live/"

	// Construct the full URL with the encoded query
	fullURL := fmt.Sprintf("%s?query=%s", baseURL, encodedQuery)

	// fmt.Println("Encoded URL:", fullURL)

	// Make the HTTP request
	resp, err := http.Get(fullURL)
	if err != nil {
		return ApiResponse{}, fmt.Errorf("error fetching URL: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ApiResponse{}, fmt.Errorf("error reading response body: %w", err)
	}

	// fmt.Println("Response from URL:")
	// fmt.Println(string(body))

	// Parse the JSON response into ApiResponse struct
	var response ApiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return ApiResponse{}, fmt.Errorf("error parsing JSON: %w", err)
	}

	// Return the ApiResponse struct
	return response, nil
}
