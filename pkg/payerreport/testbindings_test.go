package payerreport

import "github.com/xmtp/xmtpd/pkg/merkle"

func GenerateMerkleTreeTestBinding(payerMap PayerMap) (*merkle.MerkleTree, error) {
	return generateMerkleTree(payerMap)
}

func ValidateReportTransitionTestBinding(prevReport *PayerReport, newReport *PayerReport) error {
	return validateReportTransition(prevReport, newReport)
}

func ValidateReportStructureTestBinding(report *PayerReport) error {
	return validateReportStructure(report)
}
