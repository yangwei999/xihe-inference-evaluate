package domain

type CloudPod struct {
	PodId        string
	User         string
	SurvivalTime SurvivalTime
}

type CloudContainerDetail struct {
	AccessUrl AccessURL
	ErrorMsg  ErrorMsg
}
