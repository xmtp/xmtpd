package payerreport

func ValidateReportTransitionTestBinding(prevReport *PayerReport, newReport *PayerReport) error {
	return validateReportTransition(prevReport, newReport)
}

func ValidateReportStructureTestBinding(report *PayerReport) error {
	return validateReportStructure(report)
}
