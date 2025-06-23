package workers_test

import (
	"crypto/ecdsa"
	"database/sql"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/payerreport/workers"
	protoEnvelopes "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/server"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
	apiTestUtils "github.com/xmtp/xmtpd/pkg/testutils/api"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	networkTestUtils "github.com/xmtp/xmtpd/pkg/testutils/network"
	registryTestUtils "github.com/xmtp/xmtpd/pkg/testutils/registry"
	serverTestUtils "github.com/xmtp/xmtpd/pkg/testutils/server"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

var (
	server1NodeID = uint32(100)
	server2NodeID = uint32(200)
)

type multiNodeTestScaffold struct {
	servers            []*server.ReplicationServer
	nodeIDs            []uint32
	nodePrivateKeys    []*ecdsa.PrivateKey
	payerPrivateKeys   []*ecdsa.PrivateKey
	payerAddresses     []string
	clients            []message_api.ReplicationApiClient
	dbs                []*sql.DB
	reportGenerators   []*workers.GeneratorWorker
	attestationWorkers []*workers.AttestationWorker
	payerReportStores  []payerreport.IPayerReportStore
	log                *zap.Logger
	registrants        []*registrant.Registrant
	registry           registry.NodeRegistry
	reportsManager     blockchain.PayerReportsManager
}

func setupBlockchain(
	t *testing.T,
	nodes []registry.Node,
) (registry.NodeRegistry, blockchain.PayerReportsManager, config.ContractsOptions) {
	log := testutils.NewLog(t)
	rpcUrl := anvil.StartAnvil(t, false)
	contractsOptions := testutils.NewContractsOptions(t, rpcUrl)

	signer, err := blockchain.NewPrivateKeySigner(
		testutils.GetPayerOptions(t).PrivateKey,
		contractsOptions.SettlementChain.ChainID,
	)
	require.NoError(t, err)

	client, err := blockchain.NewClient(t.Context(), contractsOptions.SettlementChain.WssURL)
	require.NoError(t, err)

	registryAdmin, err := blockchain.NewNodeRegistryAdmin(
		log,
		client,
		signer,
		contractsOptions,
	)
	require.NoError(t, err)

	for _, node := range nodes {
		_, err := registryAdmin.AddNode(
			t.Context(),
			ethcrypto.PubkeyToAddress(*node.SigningKey).String(),
			node.SigningKey,
			node.HttpAddress,
		)
		require.NoError(t, err)

		err = registryAdmin.AddToNetwork(t.Context(), node.NodeID)
		require.NoError(t, err)
	}

	registry, err := registry.NewSmartContractRegistry(
		t.Context(),
		client,
		testutils.NewLog(t),
		contractsOptions,
	)
	require.NoError(t, err)

	err = registry.Start()
	require.NoError(t, err)
	t.Cleanup(registry.Stop)

	fromRegistry, err := registry.GetNodes()
	require.NoError(t, err)
	require.Equal(t, len(fromRegistry), len(nodes))

	reportsManager, err := blockchain.NewReportsManager(
		log, client, signer, contractsOptions.SettlementChain,
	)
	require.NoError(t, err)

	return registry, reportsManager, contractsOptions
}

func setupMultiNodeTest(t *testing.T) multiNodeTestScaffold {
	ctx := t.Context()
	dbs := testutils.NewDBs(t, ctx, 2)
	log := testutils.NewLog(t)
	privateKey1 := testutils.RandomPrivateKey(t)
	privateKey2 := testutils.RandomPrivateKey(t)

	payerPrivateKey1 := testutils.RandomPrivateKey(t)
	payerPrivateKey2 := testutils.RandomPrivateKey(t)

	server1Port := networkTestUtils.FindFreePort(t)
	server2Port := networkTestUtils.FindFreePort(t)

	httpServer1Port := networkTestUtils.FindFreePort(t)
	httpServer2Port := networkTestUtils.FindFreePort(t)

	nodes := []registry.Node{
		registryTestUtils.CreateNode(server1NodeID, server1Port, privateKey1),
		registryTestUtils.CreateNode(server2NodeID, server2Port, privateKey2),
	}

	registry, reportsManager, contractsOptions := setupBlockchain(t, nodes)

	domainSeparator, err := reportsManager.GetDomainSeparator(t.Context())
	require.NoError(t, err)

	server1 := serverTestUtils.NewTestServer(
		t,
		serverTestUtils.TestServerCfg{
			Port:             server1Port,
			HttpPort:         httpServer1Port,
			Db:               dbs[0],
			Registry:         registry,
			PrivateKey:       privateKey1,
			ContractsOptions: contractsOptions,
			Services: serverTestUtils.EnabledServices{
				Replication: true,
				Sync:        true,
			},
		},
	)
	server2 := serverTestUtils.NewTestServer(
		t,
		serverTestUtils.TestServerCfg{
			Port:             server2Port,
			HttpPort:         httpServer2Port,
			Db:               dbs[1],
			Registry:         registry,
			PrivateKey:       privateKey2,
			ContractsOptions: contractsOptions,
			Services: serverTestUtils.EnabledServices{
				Replication: true,
				Sync:        true,
			},
		},
	)

	require.NotEqual(t, server1.Addr(), server2.Addr())

	client1 := apiTestUtils.NewReplicationAPIClient(t, server1.Addr().String())
	client2 := apiTestUtils.NewReplicationAPIClient(t, server2.Addr().String())

	registrant1, err := registrant.NewRegistrant(
		t.Context(),
		log,
		queries.New(dbs[0]),
		registry,
		utils.EcdsaPrivateKeyToString(privateKey1),
		nil,
	)
	require.NoError(t, err)
	registrant2, err := registrant.NewRegistrant(
		t.Context(),
		log,
		queries.New(dbs[1]),
		registry,
		utils.EcdsaPrivateKeyToString(privateKey2),
		nil,
	)
	require.NoError(t, err)

	payerReportStore1 := payerreport.NewStore(dbs[0], log)
	payerReportStore2 := payerreport.NewStore(dbs[1], log)
	reportGenerator1 := workers.NewGeneratorWorker(
		t.Context(),
		log.With(zap.Uint32("node_id", server1NodeID)),
		payerReportStore1,
		registry,
		registrant1,
		1*time.Hour,
		domainSeparator,
	)
	reportGenerator2 := workers.NewGeneratorWorker(
		t.Context(),
		log.With(zap.Uint32("node_id", server2NodeID)),
		payerReportStore2,
		registry,
		registrant2,
		1*time.Hour,
		domainSeparator,
	)

	attestationWorker1 := workers.NewAttestationWorker(
		t.Context(),
		log,
		registrant1,
		payerReportStore1,
		1*time.Hour,
		domainSeparator,
	)
	attestationWorker2 := workers.NewAttestationWorker(
		t.Context(),
		log,
		registrant2,
		payerReportStore2,
		1*time.Hour,
		domainSeparator,
	)

	t.Cleanup(func() {
		log.Info("Shutting down servers")
		server1.Shutdown(0)
		server2.Shutdown(0)
	})

	return multiNodeTestScaffold{
		servers:          []*server.ReplicationServer{server1, server2},
		nodeIDs:          []uint32{server1NodeID, server2NodeID},
		nodePrivateKeys:  []*ecdsa.PrivateKey{privateKey1, privateKey2},
		payerPrivateKeys: []*ecdsa.PrivateKey{payerPrivateKey1, payerPrivateKey2},
		payerAddresses: []string{
			utils.EcdsaPublicKeyToAddress(&payerPrivateKey1.PublicKey),
			utils.EcdsaPublicKeyToAddress(&payerPrivateKey2.PublicKey),
		},
		clients:          []message_api.ReplicationApiClient{client1, client2},
		dbs:              dbs,
		reportGenerators: []*workers.GeneratorWorker{reportGenerator1, reportGenerator2},
		registrants:      []*registrant.Registrant{registrant1, registrant2},
		attestationWorkers: []*workers.AttestationWorker{
			attestationWorker1,
			attestationWorker2,
		},
		payerReportStores: []payerreport.IPayerReportStore{payerReportStore1, payerReportStore2},
		log:               log,
		registry:          registry,
		reportsManager:    reportsManager,
	}
}

func (s *multiNodeTestScaffold) publishRandomMessage(
	t *testing.T,
	topic []byte,
	nodeIndex int,
	payerIndex int,
) {
	payerEnv := envelopeTestUtils.CreatePayerEnvelopeWithSigner(
		t,
		s.nodeIDs[nodeIndex],
		s.payerPrivateKeys[payerIndex],
		10,
		envelopeTestUtils.CreateClientEnvelope(&protoEnvelopes.AuthenticatedData{
			TargetTopic: topic,
		}),
	)

	_, err := s.clients[nodeIndex].PublishPayerEnvelopes(
		t.Context(),
		&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*protoEnvelopes.PayerEnvelope{payerEnv},
		},
	)
	require.NoError(t, err)
}

func (s *multiNodeTestScaffold) getMessagesFromTopic(
	t *testing.T,
	nodeIndex int,
	topic []byte,
) []*protoEnvelopes.OriginatorEnvelope {
	client := s.clients[nodeIndex]

	response, err := client.QueryEnvelopes(t.Context(), &message_api.QueryEnvelopesRequest{
		Query: &message_api.EnvelopesQuery{
			Topics: [][]byte{topic},
		},
	})
	require.NoError(t, err)

	return response.Envelopes
}

func TestValidSignature(t *testing.T) {
	scaffold := setupMultiNodeTest(t)
	reportsManager := scaffold.reportsManager

	domainSeparator, err := reportsManager.GetDomainSeparator(t.Context())
	require.NoError(t, err)

	payerReport, err := payerreport.BuildPayerReport(payerreport.BuildPayerReportParams{
		OriginatorNodeID:    100,
		StartSequenceID:     0,
		EndSequenceID:       1,
		EndMinuteSinceEpoch: 10,
		Payers:              map[common.Address]currency.PicoDollar{},
		NodeIDs:             []uint32{100, 200},
		DomainSeparator:     domainSeparator,
	})

	require.NoError(t, err)

	signatures := make([]payerreport.NodeSignature, len(scaffold.registrants))

	for idx, registrant := range scaffold.registrants {
		signature, err := registrant.SignPayerReportAttestation(payerReport.ID)
		require.NoError(t, err)
		signatures[idx] = *signature
	}

	reportWithStatus := payerreport.PayerReportWithStatus{
		PayerReport:           payerReport.PayerReport,
		AttestationSignatures: signatures,
	}

	// Ensure the report ID matches the one we built
	reportID, err := reportsManager.GetReportID(t.Context(), &reportWithStatus)
	require.NoError(t, err)
	require.Equal(t, reportID, payerReport.ID)

	// Submit the report to the blockchain
	err = reportsManager.SubmitPayerReport(t.Context(), &reportWithStatus)
	require.NoError(t, err)
}

func TestCanGenerateReport(t *testing.T) {
	scaffold := setupMultiNodeTest(t)
	groupID := testutils.RandomGroupID()
	messageTopic := topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, groupID[:]).Bytes()

	scaffold.publishRandomMessage(t, messageTopic, 0, 0)
	scaffold.publishRandomMessage(t, messageTopic, 1, 1)

	// Confirm that both nodes have received the messages
	require.Eventually(t, func() bool {
		messagesOnNode1 := scaffold.getMessagesFromTopic(t, 0, messageTopic)
		messagesOnNode2 := scaffold.getMessagesFromTopic(t, 1, messageTopic)
		return len(messagesOnNode1) == 2 && len(messagesOnNode2) == 2
	}, 2*time.Second, 50*time.Millisecond)

	err := scaffold.reportGenerators[0].GenerateReports()
	require.NoError(t, err)

	node1ReportTopic := topic.NewTopic(topic.TOPIC_KIND_PAYER_REPORTS_V1, utils.Uint32ToBytes(scaffold.nodeIDs[0])).
		Bytes()

	require.Eventually(t, func() bool {
		messagesOnNode1 := scaffold.getMessagesFromTopic(t, 0, node1ReportTopic)
		messagesOnNode2 := scaffold.getMessagesFromTopic(t, 1, node1ReportTopic)
		return len(messagesOnNode1) == 1 && len(messagesOnNode2) == 1
	}, 2*time.Second, 50*time.Millisecond)

	// Try and generate a report again. This should be a no-op.
	err = scaffold.reportGenerators[0].GenerateReports()
	require.NoError(t, err)

	// Make sure there is still only one report after generating again
	time.Sleep(100 * time.Millisecond)
	messagesOnNode1 := scaffold.getMessagesFromTopic(t, 0, node1ReportTopic)
	require.Len(t, messagesOnNode1, 1)
}

func TestCanGenerateAndAttestReport(t *testing.T) {
	scaffold := setupMultiNodeTest(t)
	groupID := testutils.RandomGroupID()
	messageTopic := topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, groupID[:]).Bytes()

	scaffold.publishRandomMessage(t, messageTopic, 0, 0)
	scaffold.publishRandomMessage(t, messageTopic, 1, 1)

	// Confirm that both nodes have received the messages
	require.Eventually(t, func() bool {
		messagesOnNode1 := scaffold.getMessagesFromTopic(t, 0, messageTopic)
		messagesOnNode2 := scaffold.getMessagesFromTopic(t, 1, messageTopic)
		return len(messagesOnNode1) == 2 && len(messagesOnNode2) == 2
	}, 2*time.Second, 50*time.Millisecond)

	err := scaffold.reportGenerators[0].GenerateReports()
	require.NoError(t, err)

	node1ReportTopic := topic.NewTopic(topic.TOPIC_KIND_PAYER_REPORTS_V1, utils.Uint32ToBytes(scaffold.nodeIDs[0])).
		Bytes()

	require.Eventually(t, func() bool {
		messagesOnNode1 := scaffold.getMessagesFromTopic(t, 0, node1ReportTopic)
		messagesOnNode2 := scaffold.getMessagesFromTopic(t, 1, node1ReportTopic)
		return len(messagesOnNode1) == 1 && len(messagesOnNode2) == 1
	}, 2*time.Second, 50*time.Millisecond)

	// Make both node's attestation workers try and attest reports. Do this multiple times to ensure no dupes
	for range 5 {
		err = scaffold.attestationWorkers[0].AttestReports()
		require.NoError(t, err)
		err = scaffold.attestationWorkers[1].AttestReports()
		require.NoError(t, err)
	}

	attestationTopic := topic.NewTopic(topic.TOPIC_KIND_PAYER_REPORT_ATTESTATIONS_V1, utils.Uint32ToBytes(scaffold.nodeIDs[0])).
		Bytes()

	require.Eventually(t, func() bool {
		messagesOnNode1 := scaffold.getMessagesFromTopic(t, 0, attestationTopic)
		messagesOnNode2 := scaffold.getMessagesFromTopic(t, 1, attestationTopic)
		// We are expecting 2 attestations total. One from each node. Each node's attestation should have synced from the other node
		return len(messagesOnNode1) == 2 && len(messagesOnNode2) == 2
	}, 2*time.Second, 50*time.Millisecond)

	// Get the attestations of the two reports from both nodes
	for nodeIndex := range 2 {
		// See all the reports from the perspective of node1
		node1Reports, err := scaffold.payerReportStores[0].FetchReports(
			t.Context(),
			payerreport.NewFetchReportsQuery().WithOriginatorNodeID(scaffold.nodeIDs[nodeIndex]),
		)
		require.NoError(t, err)
		require.Len(t, node1Reports, 1)
		for _, report := range node1Reports {
			require.Len(t, report.AttestationSignatures, 2)
		}

		// See all the reports from the perspective of node1
		node2Reports, err := scaffold.payerReportStores[1].FetchReports(
			t.Context(),
			payerreport.NewFetchReportsQuery().WithOriginatorNodeID(scaffold.nodeIDs[nodeIndex]),
		)
		require.NoError(t, err)
		require.Len(t, node2Reports, 1)
		for _, report := range node2Reports {
			require.Len(t, report.AttestationSignatures, 2)
		}
	}

	submitterWorker := workers.NewSubmitterWorker(
		t.Context(),
		scaffold.log,
		scaffold.payerReportStores[0],
		scaffold.registry,
		scaffold.reportsManager,
		scaffold.nodeIDs[0],
	)
	err = submitterWorker.SubmitReports(t.Context())
	require.NoError(t, err)

	fetchedReports, err := scaffold.payerReportStores[0].FetchReports(
		t.Context(),
		payerreport.NewFetchReportsQuery().WithOriginatorNodeID(scaffold.nodeIDs[0]),
	)
	require.NoError(t, err)
	require.Len(t, fetchedReports, 1)
	fetchedReportID := fetchedReports[0].ID

	report, err := scaffold.reportsManager.GetReport(t.Context(), scaffold.nodeIDs[0], 0)
	require.NoError(t, err)
	require.Equal(
		t,
		report.ID,
		fetchedReportID,
	)
}
