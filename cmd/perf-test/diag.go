package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	messageApi "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/topic"
)

// runDiagnostic tests live vs catch-up delivery to isolate the bug.
// 1. Subscribe to a topic (live)
// 2. Publish messages to that topic
// 3. Open ANOTHER subscribe stream for catch-up verification
// 4. Report what each stream received
func runDiagnostic(cfg *config) error {
	conn, err := newGRPCConn(cfg)
	if err != nil {
		return fmt.Errorf("grpc connect: %w", err)
	}
	defer func() { _ = conn.Close() }()

	// Create a random topic
	topicID := randomBytes(16)
	tp := topic.NewTopic(topic.TopicKindGroupMessagesV1, topicID)
	topicBytes := tp.Bytes()
	fmt.Printf("Diagnostic: topic=%s (%d bytes: %x)\n", tp.String(), len(topicBytes), topicBytes)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Step 1: Open live subscribe stream
	queryClient := messageApi.NewQueryApiClient(conn)
	liveStream, err := queryClient.SubscribeTopics(ctx, &messageApi.SubscribeTopicsRequest{
		Filters: []*messageApi.SubscribeTopicsRequest_TopicFilter{
			{Topic: topicBytes},
		},
	})
	if err != nil {
		return fmt.Errorf("open live stream: %w", err)
	}
	fmt.Println("  Live stream opened")

	// Drain initial status messages
	for range 2 {
		resp, err := liveStream.Recv()
		if err != nil {
			return fmt.Errorf("recv initial: %w", err)
		}
		if su := resp.GetStatusUpdate(); su != nil {
			fmt.Printf("  Live recv: status=%s\n", su.GetStatus())
		}
	}

	time.Sleep(500 * time.Millisecond)

	// Step 2: Publish 5 messages to this topic
	key, err := ethcrypto.GenerateKey()
	if err != nil {
		return fmt.Errorf("gen key: %w", err)
	}

	publishClient := messageApi.NewPublishApiClient(conn)
	published := 0
	for range 5 {
		req, err := buildPublishRequestForTopic(cfg, key, tp, 256)
		if err != nil {
			fmt.Printf("  Publish build error: %v\n", err)
			continue
		}
		_, err = publishClient.PublishPayerEnvelopes(ctx, req)
		if err != nil {
			fmt.Printf("  Publish error: %v\n", err)
			continue
		}
		published++
	}
	fmt.Printf("  Published %d messages\n", published)

	// Step 3: Wait briefly, then check live stream (non-blocking with timeout)
	fmt.Println("  Waiting 3s for live delivery...")
	var liveReceived atomic.Int64
	liveCtx, liveCancel := context.WithTimeout(ctx, 3*time.Second)
	defer liveCancel()
	go func() {
		for {
			resp, err := liveStream.Recv()
			if err != nil {
				return
			}
			if su := resp.GetStatusUpdate(); su != nil {
				fmt.Printf("  Live recv: status=%s\n", su.GetStatus())
				continue
			}
			envs := resp.GetEnvelopes()
			if envs != nil && len(envs.GetEnvelopes()) > 0 {
				n := int64(len(envs.GetEnvelopes()))
				total := liveReceived.Add(n)
				fmt.Printf("  Live recv: %d envelopes (total %d)\n", n, total)
			}
		}
	}()
	<-liveCtx.Done()
	fmt.Printf("  Live stream received: %d messages\n", liveReceived.Load())

	// Step 4: Open catch-up stream on the same topic to verify messages are in DB
	fmt.Println("  Opening catch-up stream...")
	catchupCtx, catchupCancel := context.WithTimeout(ctx, 3*time.Second)
	defer catchupCancel()

	catchupStream, err := queryClient.SubscribeTopics(
		catchupCtx, &messageApi.SubscribeTopicsRequest{
			Filters: []*messageApi.SubscribeTopicsRequest_TopicFilter{
				{Topic: topicBytes},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("open catchup stream: %w", err)
	}

	catchupReceived := 0
	for {
		resp, err := catchupStream.Recv()
		if err != nil {
			if catchupCtx.Err() == nil {
				fmt.Printf("  Catch-up recv error: %v\n", err)
			}
			break
		}
		if su := resp.GetStatusUpdate(); su != nil {
			fmt.Printf("  Catch-up recv: status=%s\n", su.GetStatus())
			if su.GetStatus() == messageApi.SubscribeTopicsResponse_SUBSCRIPTION_STATUS_CATCHUP_COMPLETE {
				break
			}
			continue
		}
		envs := resp.GetEnvelopes()
		if len(envs.GetEnvelopes()) > 0 {
			catchupReceived += len(envs.GetEnvelopes())
			fmt.Printf("  Catch-up recv: %d envelopes (total %d)\n",
				len(envs.GetEnvelopes()), catchupReceived)
		}
	}
	fmt.Printf("  Catch-up stream received: %d messages\n", catchupReceived)

	fmt.Printf("\n=== DIAGNOSTIC RESULTS ===\n")
	live := liveReceived.Load()
	fmt.Printf("Published:      %d\n", published)
	fmt.Printf("Live received:  %d\n", live)
	fmt.Printf("Catchup received: %d\n", catchupReceived)

	if catchupReceived > 0 && live == 0 {
		fmt.Println(
			"\nDIAGNOSIS: Messages are in DB but subscribe worker is NOT delivering to live listeners.",
		)
		fmt.Println("This is a server-side bug in the subscribe worker dispatch path.")
	} else if catchupReceived == 0 && live == 0 {
		fmt.Println("\nDIAGNOSIS: Messages are NOT in DB at all. Publish path issue or wrong node.")
	} else if live > 0 {
		fmt.Println("\nDIAGNOSIS: Live delivery is working!")
	}

	return nil
}
