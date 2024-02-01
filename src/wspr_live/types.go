package wspr_live

import (
	"time"

	"github.com/go-gota/gota/dataframe"
)

type QueryByCallsignOptions struct {
	FromTime time.Time
	ToTime   time.Time
	Limit    int
	MinSnr   int
}

type MetaData struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Statistics struct {
	Elapsed   float64 `json:"elapsed"`
	RowsRead  int     `json:"rows_read"`
	BytesRead int     `json:"bytes_read"`
}

type ApiResponse struct {
	Meta                   []MetaData      `json:"meta"`
	Data                   [][]interface{} `json:"data"`
	Rows                   int             `json:"rows"`
	RowsBeforeLimitAtLeast int             `json:"rows_before_limit_at_least"`
	Statistics             Statistics      `json:"statistics"`
}

func (apiResp ApiResponse) ToMaps() []map[string]interface{} {
	var result []map[string]interface{}

	for _, row := range apiResp.Data {
		rowMap := make(map[string]interface{})
		for i, value := range row {
			if i < len(apiResp.Meta) {
				rowMap[apiResp.Meta[i].Name] = value
			}
		}
		result = append(result, rowMap)
	}

	return result
}

func (apiResp ApiResponse) ToDataFrame() dataframe.DataFrame {
	df := dataframe.LoadMaps(apiResp.ToMaps())
	return df
}
