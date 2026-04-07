package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const version = "0.1.0-dev"

func main() {
	addr := flag.String("addr", "127.0.0.1:8087", "listen address")
	domain := flag.String("trust-domain", "domain-a.test", "server trust domain (issued principals must belong to it)")
	store := flag.String("store", "", "path to a durable state file (default: in-memory only)")
	keyPath := flag.String("key", "", "path to the authority key (PEM); created if absent. Needed for records to survive restarts")
	apiKey := flag.String("api-key", envOr("ATLAS_API_KEY", ""), "if set, /issue and /revoke require Authorization: Bearer <key>")
	grant := flag.String("grant", "", "comma-separated scopes a principal may delegate (default: a built-in demo set)")
	flag.Parse()

	var grantSet []string
	if *grant != "" {
		for _, s := range strings.Split(*grant, ",") {
			if s = strings.TrimSpace(s); s != "" {
				grantSet = append(grantSet, s)
			}
		}
	}
	app, err := NewApp(Config{Domain: *domain, StorePath: *store, KeyPath: *keyPath, APIKey: *apiKey, Grant: grantSet}, systemClock{})
	if err != nil {
		log.Fatalf("atlas-server: %v", err)
	}
	if *store != "" {
		log.Printf("atlas-server: durable store at %s", *store)
	}
	if *keyPath != "" {
		log.Printf("atlas-server: authority key at %s", *keyPath)
	}
	if *apiKey != "" {
		log.Printf("atlas-server: mutating endpoints require a bearer token")
	}

	// Freshness refresher: re-sign the current revoked set every second so the
	// held snapshot never ages past the freshness bound R (= 2s) even when no
	// revocations occur.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if err := app.Refresh(); err != nil {
					log.Printf("atlas-server: snapshot refresh failed: %v", err)
				}
				if err := app.Flush(); err != nil {
					log.Printf("atlas-server: store flush failed: %v", err)
				}
			}
		}
	}()

	srv := &http.Server{
		Addr:              *addr,
		Handler:           app.Router(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("atlas-server %s listening on http://%s  (trust domain %s, R=%s)", version, *addr, *domain, app.revWindow)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("atlas-server: %v", err)
		}
	}()

	// Graceful shutdown.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("atlas-server: shutting down…")
	shutCtx, shutCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutCancel()
	if err := srv.Shutdown(shutCtx); err != nil {
		log.Printf("atlas-server: shutdown: %v", err)
	}
	if err := app.Flush(); err != nil {
		log.Printf("atlas-server: final flush: %v", err)
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
