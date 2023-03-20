package domain

import (
	"errors"
	"net/url"
)

// SurvivalTime
type SurvivalTime interface {
	SurvivalTime() int64
}

func NewSurvivalTime(v int64) (SurvivalTime, error) {
	return survivalTime(v), nil
}

type survivalTime int64

func (r survivalTime) SurvivalTime() int64 {
	return int64(r)
}

// AccessURL
type AccessURL interface {
	AccessURL() string
}

func NewAccessURL(v string) (AccessURL, error) {
	if _, err := url.Parse(v); err != nil {
		return nil, errors.New("invalid url")
	}

	return accessURL(v), nil
}

type accessURL string

func (r accessURL) AccessURL() string {
	return string(r)
}

// ErrorMsg
type ErrorMsg interface {
	ErrorMsg() string
	IsGood() bool
}

func NewErrorMsg(v string) (ErrorMsg, error) {
	return errorMsg(v), nil
}

type errorMsg string

func (r errorMsg) ErrorMsg() string {
	return string(r)
}

func (p errorMsg) IsGood() bool {
	return p.ErrorMsg() == ""
}
