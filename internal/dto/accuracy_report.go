package dto

type TestCaseResult struct {
	UserName  string `json:"user_name"`
	EventName string `json:"event_name"`

	DijkstraDistanceKM float64 `json:"-"`
	OSRMDistanceKM     float64 `json:"-"`
	PercentageErrorPct float64 `json:"-"`
	DijkstraExecTimeMs float64 `json:"-"`

	DijkstraDistance string `json:"dijkstra_distance"`  
	OSRMDistance     string `json:"osrm_distance"`      
	PercentageError  string `json:"percentage_error"`   
	DijkstraExecTime string `json:"dijkstra_exec_time"` 
}


type UserMAPEResult struct {
	UserName string           `json:"user_name"`
	Cases    []TestCaseResult `json:"cases"`

	MAPEPercentValue     float64 `json:"-"`
	TotalExecTimeMsValue float64 `json:"-"`

	MAPEPercent   string `json:"mape_percent"`    
	TotalExecTime string `json:"total_exec_time"` 
}

type AccuracyReport struct {
	GeneratedAt        string           `json:"generated_at"`
	TotalUsers         int              `json:"total_users"`
	TotalEvents        int              `json:"total_events"`
	Results            []UserMAPEResult `json:"results"`
	OverallMAPEPercent float64          `json:"overall_mape_percent"`
}

type MAPESummaryRequest struct {
	MAPEValues []string `json:"mape_values"`
}

type MAPESummaryResponse struct {
	TotalDataUji       int       `json:"total_data_uji"`
	MAPEValues         []float64 `json:"mape_values"`
	AverageMAPEPercent float64   `json:"average_mape_percent"`
}
