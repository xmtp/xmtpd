# Payer Reports

## Validating A Report

When a new report is picked up by the worker it is validated against the most recent report submitted to the smart contract for the originator.

### Validations

The following validations are performed. If any of these fail, the report can be deemed as invalid immediately.

1. The `originatorNodeID` in the new report and the previous report match
2. The `startSequenceID` of the new report is equal to the `endSequenceID` of the previous report
3. The `startSequenceID` of the new report is less than or equal to its `endSequenceID`
4. The `payersMerkleRoot` and `nodesHash` are exactly 32 bytes long
5. The `payersMerkleRoot` or `payersLeafCount` specified in the report are different than the values independently calculated by the validator based on the validator's own records of messages sent.
6. The `nodesHash` or `nodesCount` are different than the values independently calculated by the validator based on the validator's latest view of the `NodeRegistry` smart contract

### Errors validating a report

In some cases, the worker may encounter an error while validating a Payer Report. When an error occurs the report is not deemed either valid or invalid. Errors may be retried for up to 48 hours before the report is considered expired.

Some error conditions when evaluating a report:

1. The message referenced by the `startSequenceID` of the report is not found.
2. The message referenced by the `endSequenceID` of the report is not found.
3. Any other database error when validating the report

### Misbehaviour in reports

#### Withholding messages from the final minute

The system assumes that the messages used as the `startSequenceID` and `endSequenceID` of a report are both the last message processed by the originator node in the calendar minute (according to the originator's clock). Reports are always generated at least one full minute behind the current clock.

When calculating payer spend this assumption is used to query the database in the range `minuteFromStartSequenceID+1 -> minuteFromEndSequenceID`. A malicious node may choose to exploit this behaviour by selectively withholding messages from synchronization in the final minute in a report period. Once a report is confirmed, they can then sync those messages with the rest of the network and the messages will appear "between reports" and not be counted.

This behaviour can only be detected when the _next_ report is submitted for the originator. Each node keeps track of the last message for each minute. If the last message for the minute referenced in the `endSequenceID` of the previous report does not match what they have in their database they can submit a misbehaviour report including any messages with a greater sequence ID in the same minute. This proves the misbehaviour and is cause for a node to be removed from the network.
