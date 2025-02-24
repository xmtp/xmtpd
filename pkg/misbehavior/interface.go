package misbehavior

type MisbehaviorService interface {
	SafetyFailure(report *SafetyFailureReport) error
	// TODO:nm add liveness failures
}
