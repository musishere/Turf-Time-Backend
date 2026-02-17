package types

// UpdateTurfRequest contains optional fields for updating a turf
type UpdateTurfRequest struct {
    Name       *string  `json:"name" form:"name"`
    StartTime  *int     `json:"startTime" form:"startTime"`
    EndTime    *int     `json:"endTime" form:"endTime"`
    Status     *string  `json:"status" form:"status"`
    NoOfFields *int     `json:"noOfFields" form:"noOfFields"`
    Address    *string  `json:"address" form:"address"`
}
