package ksyun

const (
	// statusPending is the general status when remote resource is not completed
	statusPending = "pending"

/*
// statusInitialized is the general status when remote resource is completed
statusInitialized = "initialized"

// statusRunning is the general status when remote resource is running
statusRunning = "running"

// statusStopped is the general status when remote resource is stopped
statusStopped = "stopped"
*/
)

// trove front
const (
	tActiveStatus    = "ACTIVE"
	tDeletedStatus   = "DELETED"
	tError           = "ERROR"
	tFailedStatus    = "FAILED"
	tStopedStatus    = "STOPPED"
	tDeleting        = "DELETING"
	tCreatingStatus  = "CREATING"
	tModifyingSpec   = "MODIFYING_SPEC"
	tRuningTask      = "RUNNING_TASK"
	tBackingUp       = "BACKING_UP"
	tModifyType      = "MODIFYING_TYPE"
	tRebooting       = "REBOOTING"
	tRestartRequired = "RESTART_REQUIRED"
	tInvalid         = "INVALID"
	tRestoring       = "RESTORING"
	tBuildingRr      = "BUILDING_RR"
	tUpgrading       = "UPGRADING"
	tExpiringSoon    = "EXPIRING_SOON"
	tLocked          = "LOCKED" //欠费状态，实例正常执行，但是不允许操作
)

// 不需要等待的状态
var finalStatus = []string{
	tActiveStatus,
	tDeletedStatus,
	tFailedStatus,
	tStopedStatus,
	tDeleting,
	tError,
	tLocked,
}
var waitStatus = []string{
	tCreatingStatus,
	tModifyingSpec,
	tRuningTask,
	tBackingUp,
	tModifyType,
	tRebooting,
	tRestartRequired,
	tInvalid,
	tRestoring,
	tBuildingRr,
	tUpgrading,
	tExpiringSoon,
}
