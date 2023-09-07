package ksyun

import "time"

type TransformType int

const (
	TransformDefault TransformType = iota
	TransformWithN
	TransformWithFilter
	TransformListFilter
	TransformListUnique
	TransformListN
	TransformSingleN
)

const (
	ResourceKrdsParameterGroup = "krds_parameters_group"
)

const (
	RetryTimeoutMinute = 10 * time.Minute
)

const (
	ApiCallBeforeProcess = iota
	ApiCallExecuteProcess
	ApiCallAfterProcess
)

type ParamsList []interface{}
