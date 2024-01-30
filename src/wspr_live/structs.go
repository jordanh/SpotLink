package wspr_live

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
