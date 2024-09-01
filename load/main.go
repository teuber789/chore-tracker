package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

type result struct {
	// The number of calls made that succeeded
	numSuccesses int
	// The number of calls made that failed
	numFailures int
	err         error
	output      *string
}

func run(ctx context.Context, rawCmd string) (*result, error) {
	cmd := exec.CommandContext(ctx, "zsh", "-c", rawCmd)
	var output string
	combined, err := cmd.CombinedOutput() // Also executes the command
	if err != nil {
		// If we've recieved a SIGKILL around the time that the context is supposed to time out,
		// there isn't actually an error. This is just the program terminating when it's supposed to.
		now := time.Now()
		deadline, ok := ctx.Deadline()
		if ok && strings.HasSuffix(err.Error(), "signal: killed") && (now.Equal(deadline) || now.After(deadline)) {
			output = "Process completed normally"
			err = nil
		} else {
			output = fmt.Sprintf("Process terminated unexpectedly with error %s and has deadline %t", err, ok)
		}
	} else {
		output = "Process terminated unexpectedly without an error"
		err = errors.New("process terminated unexpectedly without an error")
	}
	output = fmt.Sprintf("%s\nCombined output from process:\n%s", output, string(combined))

	// TODO num successes and failures
	return &result{0, 0, err, &output}, err
}

func main() {
	// Read in args
	// serverType := flag.String("server", "", "Specifies whether to load test the GRPC or the HTTP server. Valid values are 'grpc' and 'http'.")
	users := flag.Int("users", 1, "The number of concurrent users to mimic. Must be an integer greater than 0.")
	seconds := flag.Int("seconds", 300, "The number of seconds the test should last. Must be an integer greater than 0.")
	flag.Parse()

	// Set up context with timeout
	duration := time.Duration(*seconds) * time.Second
	ctx, cancel := context.WithTimeout(context.TODO(), duration)
	defer cancel()

	// Launch a process per simulated user
	g, ctx := errgroup.WithContext(ctx)
	results := make([]*result, *users)
	rawCmd := "echo hello && sleep 10"

	for i := 0; i < *users; i++ {
		g.Go(func() error {
			result, err := run(ctx, rawCmd)
			if err != nil {
				return err
			}
			results[i] = result
			return nil
		})
	}

	// Wait for all to complete
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

	// All completed successfully
	log.Println("Load test completed successfully")
	log.Printf("\n\n")
	for i, result := range results {
		log.Printf("*******************************\n")
		log.Printf("********** RESULT %d ***********\n", i+1)
		log.Printf("*******************************\n")
		log.Printf(*result.output)
		log.Printf("\n\n\n")
	}
}
