// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package subscriptions

// [START pubsub_subscriber_concurrency_control]
import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
)

func pullMsgsConcurrenyControl(w io.Writer, projectID, subName string, numGoroutines int) ([]string, error) {
	// projectID := "my-project-id"
	// subName := projectID + "-example-sub"
	// numGoroutines := 4
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(subName)
	// Must set ReceiveSettings.Synchronous to false to enable concurrency settings.
	// Otherwise, NumGoroutines will be set to 1.
	sub.ReceiveSettings.Synchronous = false
	// NumGoroutines is the number of goroutines sub.Receive will spawn to pull messages concurrently.
	sub.ReceiveSettings.NumGoroutines = numGoroutines

	// Receive messages for 10 seconds.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var msgs []string
	var lock sync.Mutex
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		lock.Lock()
		defer lock.Unlock()
		msgs = append(msgs, string(msg.Data))
		fmt.Fprintf(w, "Got message: %s\n", string(msg.Data))
		msg.Ack()
	})
	if err != nil {
		return nil, fmt.Errorf("Receive: %v", err)
	}
	return msgs, nil
}

// [END pubsub_subscriber_concurrency_control]