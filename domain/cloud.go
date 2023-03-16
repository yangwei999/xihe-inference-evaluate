package domain

type CloudPod struct {
	PodId string
	SurvivalTime SurvivalTime
}

type CloudContainerDetail struct {
	AccessUrl AccessURL
	ErrorMsg  ErrorMsg
}
