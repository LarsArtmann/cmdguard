// Integration test for di example
package main

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func TestDIExample_CreateCLI(t *testing.T) {
	t.Parallel()

	cli, err := v2.NewCLI[Config]("di-app", "DI Example App", Config{})
	if err != nil {
		t.Fatalf("Failed to create CLI: %v", err)
	}

	if cli == nil {
		t.Fatal("CLI is nil")
	}

	cfg := cli.Config()
	if cfg == nil {
		t.Fatal("Config is nil")
	}
}

func TestDIExample_ServiceRegistration(t *testing.T) {
	t.Parallel()

	cli, err := v2.NewCLI[Config]("di-app", "DI Example App", Config{})
	if err != nil {
		t.Fatalf("Failed to create CLI: %v", err)
	}

	// Register database service
	scope := cli.Scope()

	err = v2.Provide(scope, NewDatabaseService)
	if err != nil {
		t.Fatalf("Failed to register database service: %v", err)
	}

	// Register API service
	err = v2.Provide(scope, NewAPIService)
	if err != nil {
		t.Fatalf("Failed to register API service: %v", err)
	}

	// Verify database service can be invoked
	database, err := v2.Invoke[*DatabaseService](scope)
	if err != nil {
		t.Fatalf("Failed to invoke database service: %v", err)
	}database	if db == nil {
		t.Fatal("Database service is nildatabase
	}

	if !db.IsConnected() {
		t.Error("Database should be connected")
	}

	// Verify API service can be invoked
	api, err := v2.Invoke[*APIService](scope)
	if err != nil {
		t.Fatalf("Failed to invoke API service: %v", err)
	}

	if api == nil {
		t.Fatal("API service is nil")
	}
}

func TestDIExample_MustInvoke(t *testing.T) {
	t.Parallel()

	cli, err := v2.NewCLI[Config]("di-app", "DI Example App", Config{})
	if err != nil {
		t.Fatalf("Failed to create CLI: %v", err)
	}

	// Register services
	scope := cli.Scope()

	err = v2.Provide(scope, NewDatabaseService)
	if err != nil {
		t.Fatalf("Failed to register database service: %v", err)
	}

	// Test Invoke (Musdatabasenvoke removed)
	db, err := v2.Invoke[*DatabaseService](scope)
	if err != nil {
		t.Fatalf("Failed to invoke database servdatabasee: %v", err)
	}

	if !db.IsConnected() {
		t.Error("Database should be connected")
	}
}

func TestDIExample_HealthCheck(t *testing.T) {
	t.Parallel()

	cli, err := v2.NewCLI[Config]("di-app", "DI Example App", Config{})
	if err != nil {
		t.Fatalf("Failed to create CLI: %v", err)
	}

	// Register database service
	scope := cli.Scope()

	err = v2.Provide(scope, NewDatabaseService)
	if err != nil {
		t.Fatalf("Failed to register database service: %v", err)
	}

	// Health check should pass when service is registered and connected
	ctx := context.Background()

	err = scope.HealthCheckWithContext(ctx)
	if err != nil {
		t.Errorf("Health check failed: %v", err)
	}
}

func TestDIExample_ConfigAccess(t *testing.T) {
	t.Parallel()

	cfg := Config{
		LogLevel:  v2.LogLevel{},
		ServerURL: "http://localhost:8080",
	}

	cli, err := v2.NewCLI[Config]("di-app", "DI Example App", cfg)
	if err != nil {
		t.Fatalf("Failed to create CLI: %v", err)
	}

	retrievedCfg := cli.Config()
	if retrievedCfg == nil {
		t.Fatal("Config is nil")
	}

	if retrievedCfg.ServerURL != cfg.ServerURL {
		t.Errorf("Expected ServerURL %q, got %q", cfg.ServerURL, retrievedCfg.ServerURL)
	}
}
