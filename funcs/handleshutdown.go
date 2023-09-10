package funcs

import (
	"log"
	"os"
	"os/signal"
)

func HandleShutdown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	// Block until a signal is received
	sig := <-sigCh
	log.Printf("Received signal %s. Shutting down...\n", sig)
	close(shutdownCh) // Signal other goroutines to stop
	// Perform any cleanup here.

	os.Exit(0)
}
