package ksyun

import "time"

type TransformType int

const (
	TransformDefault TransformType = iota

	// TransformWithN output: Attribute.N
	TransformWithN

	// TransformWithFilter output: Filter.N.Name Filter.N.Value.1
	TransformWithFilter

	// TransformListFilter output:
	TransformListFilter

	// TransformListUnique output: Object.Attribute
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
