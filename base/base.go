package base

type BaseResp struct {
    ErrorCode string `json:"errorCode,omitempty"`
    ErrorDesc string `json:"errorDesc,omitempty"`
}
