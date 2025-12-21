// Package secretsmanager integration tests using LocalStack.
//go:build integration

package secretsmanager_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alessiosavi/GoGPUtils/aws/internal/testutil"
	"github.com/alessiosavi/GoGPUtils/aws/secretsmanager"
)

// testContext returns a context with timeout for tests.
func testContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)

	return ctx
}

// testConfig holds test configuration.
type testConfig struct {
	cfg    *testutil.TestConfig
	client *secretsmanager.Client
	prefix string
}

// setupTest creates a test environment with SecretsManager client.
func setupTest(t *testing.T) *testConfig {
	t.Helper()
	testutil.SkipIfNoLocalStack(t)

	cfg := testutil.MustLoadConfig(t)

	client, err := secretsmanager.NewClient(cfg.AWSConfig)
	if err != nil {
		t.Fatalf("failed to create SecretsManager client: %v", err)
	}

	// Use unique prefix for this test run
	prefix := testutil.UniqueSecretName()

	return &testConfig{
		cfg:    cfg,
		client: client,
		prefix: prefix,
	}
}

// cleanupSecret deletes a secret with force delete.
func cleanupSecret(t *testing.T, client *secretsmanager.Client, name string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client.DeleteSecret(ctx, name, secretsmanager.WithForceDelete())
}

func TestCreateSecretString(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	tests := []struct {
		name       string
		secretName string
		value      string
	}{
		{
			name:       "simple secret",
			secretName: "simple-secret",
			value:      "my-secret-value",
		},
		{
			name:       "json secret",
			secretName: "json-secret",
			value:      `{"username": "admin", "password": "secret123"}`,
		},
		{
			name:       "multiline secret",
			secretName: "multiline-secret",
			value:      "line1\nline2\nline3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secretName := tc.prefix + "-" + tt.secretName
			t.Cleanup(func() { cleanupSecret(t, tc.client, secretName) })

			err := tc.client.CreateSecretString(ctx, secretName, tt.value)
			if err != nil {
				t.Fatalf("CreateSecretString failed: %v", err)
			}

			// Verify by getting
			got, err := tc.client.GetSecretString(ctx, secretName)
			if err != nil {
				t.Fatalf("GetSecretString failed: %v", err)
			}

			if got != tt.value {
				t.Errorf("value mismatch: got %s, want %s", got, tt.value)
			}
		})
	}
}

func TestCreateSecretWithOptions(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	t.Run("with description", func(t *testing.T) {
		secretName := tc.prefix + "-with-desc"
		t.Cleanup(func() { cleanupSecret(t, tc.client, secretName) })

		err := tc.client.CreateSecretString(ctx, secretName, "value",
			secretsmanager.WithCreateDescription("Test secret description"),
		)
		if err != nil {
			t.Fatalf("CreateSecretString with description failed: %v", err)
		}

		info, err := tc.client.DescribeSecret(ctx, secretName)
		if err != nil {
			t.Fatalf("DescribeSecret failed: %v", err)
		}

		if info.Description != "Test secret description" {
			t.Errorf("description mismatch: got %s, want Test secret description", info.Description)
		}
	})

	t.Run("with tags", func(t *testing.T) {
		secretName := tc.prefix + "-with-tags"
		t.Cleanup(func() { cleanupSecret(t, tc.client, secretName) })

		tags := map[string]string{
			"env":     "test",
			"project": "integration",
		}

		err := tc.client.CreateSecretString(ctx, secretName, "value",
			secretsmanager.WithCreateTags(tags),
		)
		if err != nil {
			t.Fatalf("CreateSecretString with tags failed: %v", err)
		}

		info, err := tc.client.DescribeSecret(ctx, secretName)
		if err != nil {
			t.Fatalf("DescribeSecret failed: %v", err)
		}

		for k, v := range tags {
			if info.Tags[k] != v {
				t.Errorf("tag %s mismatch: got %s, want %s", k, info.Tags[k], v)
			}
		}
	})
}

func TestCreateSecretBinary(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	secretName := tc.prefix + "-binary"
	t.Cleanup(func() { cleanupSecret(t, tc.client, secretName) })

	binaryData := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}

	err := tc.client.CreateSecretBinary(ctx, secretName, binaryData)
	if err != nil {
		t.Fatalf("CreateSecretBinary failed: %v", err)
	}

	got, err := tc.client.GetSecretBinary(ctx, secretName)
	if err != nil {
		t.Fatalf("GetSecretBinary failed: %v", err)
	}

	if len(got) != len(binaryData) {
		t.Fatalf("binary data length mismatch: got %d, want %d", len(got), len(binaryData))
	}

	for i, b := range binaryData {
		if got[i] != b {
			t.Errorf("binary data mismatch at index %d: got %x, want %x", i, got[i], b)
		}
	}
}

func TestGetSecretString(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: create a secret
	secretName := tc.prefix + "-get-test"
	value := "test-value-for-get"
	t.Cleanup(func() { cleanupSecret(t, tc.client, secretName) })

	if err := tc.client.CreateSecretString(ctx, secretName, value); err != nil {
		t.Fatalf("setup CreateSecretString failed: %v", err)
	}

	t.Run("existing secret", func(t *testing.T) {
		got, err := tc.client.GetSecretString(ctx, secretName)
		if err != nil {
			t.Fatalf("GetSecretString failed: %v", err)
		}

		if got != value {
			t.Errorf("value mismatch: got %s, want %s", got, value)
		}
	})

	t.Run("non-existent secret", func(t *testing.T) {
		_, err := tc.client.GetSecretString(ctx, tc.prefix+"-non-existent")
		if !errors.Is(err, secretsmanager.ErrSecretNotFound) {
			t.Errorf("expected ErrSecretNotFound, got: %v", err)
		}
	})
}

func TestGetSecretJSON(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	type DBConfig struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	secretName := tc.prefix + "-json-secret"
	jsonValue := `{"host":"localhost","port":5432,"username":"admin","password":"secret"}`
	t.Cleanup(func() { cleanupSecret(t, tc.client, secretName) })

	if err := tc.client.CreateSecretString(ctx, secretName, jsonValue); err != nil {
		t.Fatalf("setup CreateSecretString failed: %v", err)
	}

	var config DBConfig
	err := tc.client.GetSecretJSON(ctx, secretName, &config)
	if err != nil {
		t.Fatalf("GetSecretJSON failed: %v", err)
	}

	if config.Host != "localhost" {
		t.Errorf("host mismatch: got %s, want localhost", config.Host)
	}

	if config.Port != 5432 {
		t.Errorf("port mismatch: got %d, want 5432", config.Port)
	}

	if config.Username != "admin" {
		t.Errorf("username mismatch: got %s, want admin", config.Username)
	}

	if config.Password != "secret" {
		t.Errorf("password mismatch: got %s, want secret", config.Password)
	}
}

func TestUpdateSecretString(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	secretName := tc.prefix + "-update-test"
	t.Cleanup(func() { cleanupSecret(t, tc.client, secretName) })

	// Create initial secret
	if err := tc.client.CreateSecretString(ctx, secretName, "original-value"); err != nil {
		t.Fatalf("CreateSecretString failed: %v", err)
	}

	// Verify original
	got, err := tc.client.GetSecretString(ctx, secretName)
	if err != nil {
		t.Fatalf("GetSecretString failed: %v", err)
	}

	if got != "original-value" {
		t.Errorf("original value mismatch: got %s, want original-value", got)
	}

	// Update
	if err := tc.client.UpdateSecretString(ctx, secretName, "updated-value"); err != nil {
		t.Fatalf("UpdateSecretString failed: %v", err)
	}

	// Verify update
	got, err = tc.client.GetSecretString(ctx, secretName)
	if err != nil {
		t.Fatalf("GetSecretString after update failed: %v", err)
	}

	if got != "updated-value" {
		t.Errorf("updated value mismatch: got %s, want updated-value", got)
	}
}

func TestUpdateSecretNotFound(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	err := tc.client.UpdateSecretString(ctx, tc.prefix+"-nonexistent", "value")
	if !errors.Is(err, secretsmanager.ErrSecretNotFound) {
		t.Errorf("expected ErrSecretNotFound, got: %v", err)
	}
}

func TestDeleteSecret(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	t.Run("force delete", func(t *testing.T) {
		secretName := tc.prefix + "-delete-force"

		if err := tc.client.CreateSecretString(ctx, secretName, "to-be-deleted"); err != nil {
			t.Fatalf("CreateSecretString failed: %v", err)
		}

		// Force delete
		if err := tc.client.DeleteSecret(ctx, secretName, secretsmanager.WithForceDelete()); err != nil {
			t.Fatalf("DeleteSecret failed: %v", err)
		}

		// Verify deleted
		_, err := tc.client.GetSecretString(ctx, secretName)
		if !errors.Is(err, secretsmanager.ErrSecretNotFound) {
			t.Errorf("expected ErrSecretNotFound after delete, got: %v", err)
		}
	})

	// t.Run("delete non-existent", func(t *testing.T) {
	// 	err := tc.client.DeleteSecret(ctx, tc.prefix+"-nonexistent2", secretsmanager.WithForceDelete())
	// 	if !errors.Is(err, secretsmanager.ErrSecretNotFound) {
	// 		t.Errorf("expected ErrSecretNotFound, got: %v", err)
	// 	}
	// })
}

func TestDescribeSecret(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	secretName := tc.prefix + "-describe-test"
	t.Cleanup(func() { cleanupSecret(t, tc.client, secretName) })

	err := tc.client.CreateSecretString(ctx, secretName, "value",
		secretsmanager.WithCreateDescription("Test description"),
	)
	if err != nil {
		t.Fatalf("CreateSecretString failed: %v", err)
	}

	info, err := tc.client.DescribeSecret(ctx, secretName)
	if err != nil {
		t.Fatalf("DescribeSecret failed: %v", err)
	}

	if info.Name != secretName {
		t.Errorf("name mismatch: got %s, want %s", info.Name, secretName)
	}

	if info.Description != "Test description" {
		t.Errorf("description mismatch: got %s, want Test description", info.Description)
	}

	if info.ARN == "" {
		t.Error("ARN should not be empty")
	}
}

func TestDescribeSecretNotFound(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	_, err := tc.client.DescribeSecret(ctx, tc.prefix+"-nonexistent3")
	if !errors.Is(err, secretsmanager.ErrSecretNotFound) {
		t.Errorf("expected ErrSecretNotFound, got: %v", err)
	}
}

func TestListSecrets(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: create multiple secrets
	secretNames := []string{
		tc.prefix + "-list-secret1",
		tc.prefix + "-list-secret2",
		tc.prefix + "-list-secret3",
	}

	for _, name := range secretNames {
		t.Cleanup(func() { cleanupSecret(t, tc.client, name) })
		if err := tc.client.CreateSecretString(ctx, name, "value"); err != nil {
			t.Fatalf("setup CreateSecretString failed: %v", err)
		}
	}

	secrets, err := tc.client.ListSecrets(ctx)
	if err != nil {
		t.Fatalf("ListSecrets failed: %v", err)
	}

	// Count our test secrets
	found := 0
	for _, s := range secrets {
		for _, name := range secretNames {
			if s.Name == name {
				found++
				break
			}
		}
	}

	if found != len(secretNames) {
		t.Errorf("expected to find %d secrets, found %d", len(secretNames), found)
	}
}

func TestListSecretsWithMaxResults(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: create some secrets
	for i := 0; i < 3; i++ {
		name := tc.prefix + "-maxresults-secret" + string(rune('0'+i))
		t.Cleanup(func() { cleanupSecret(t, tc.client, name) })
		if err := tc.client.CreateSecretString(ctx, name, "value"); err != nil {
			t.Fatalf("setup CreateSecretString failed: %v", err)
		}
	}

	secrets, err := tc.client.ListSecrets(ctx, secretsmanager.WithMaxResults(2))
	if err != nil {
		t.Fatalf("ListSecrets with max results failed: %v", err)
	}

	if len(secrets) > 2 {
		t.Errorf("expected at most 2 secrets, got %d", len(secrets))
	}
}

func TestListSecretNames(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: create some secrets
	secretNames := []string{
		tc.prefix + "-names-secret1",
		tc.prefix + "-names-secret2",
	}

	for _, name := range secretNames {
		t.Cleanup(func() { cleanupSecret(t, tc.client, name) })
		if err := tc.client.CreateSecretString(ctx, name, "value"); err != nil {
			t.Fatalf("setup CreateSecretString failed: %v", err)
		}
	}

	names, err := tc.client.ListSecretNames(ctx)
	if err != nil {
		t.Fatalf("ListSecretNames failed: %v", err)
	}

	// Verify our secrets are in the list
	found := 0
	for _, name := range names {
		for _, expected := range secretNames {
			if name == expected {
				found++
				break
			}
		}
	}

	if found != len(secretNames) {
		t.Errorf("expected to find %d secrets, found %d", len(secretNames), found)
	}
}

func TestValidationErrors(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	t.Run("empty name for create", func(t *testing.T) {
		err := tc.client.CreateSecretString(ctx, "", "value")
		if err == nil {
			t.Error("expected error for empty secret name")
		}
	})

	t.Run("empty name for get", func(t *testing.T) {
		_, err := tc.client.GetSecretString(ctx, "")
		if err == nil {
			t.Error("expected error for empty secret name")
		}
	})

	t.Run("empty name for update", func(t *testing.T) {
		err := tc.client.UpdateSecretString(ctx, "", "value")
		if err == nil {
			t.Error("expected error for empty secret name")
		}
	})

	t.Run("empty name for delete", func(t *testing.T) {
		err := tc.client.DeleteSecret(ctx, "")
		if err == nil {
			t.Error("expected error for empty secret name")
		}
	})

	t.Run("empty name for describe", func(t *testing.T) {
		_, err := tc.client.DescribeSecret(ctx, "")
		if err == nil {
			t.Error("expected error for empty secret name")
		}
	})
}

func TestSecretRoundTrip(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	secretName := tc.prefix + "-roundtrip"
	t.Cleanup(func() { cleanupSecret(t, tc.client, secretName) })

	// Create
	originalValue := "original-secret-value"
	err := tc.client.CreateSecretString(ctx, secretName, originalValue,
		secretsmanager.WithCreateDescription("Round trip test"),
	)
	if err != nil {
		t.Fatalf("CreateSecretString failed: %v", err)
	}

	// Read
	got, err := tc.client.GetSecretString(ctx, secretName)
	if err != nil {
		t.Fatalf("GetSecretString failed: %v", err)
	}
	if got != originalValue {
		t.Errorf("value mismatch: got %s, want %s", got, originalValue)
	}

	// Update
	updatedValue := "updated-secret-value"
	err = tc.client.UpdateSecretString(ctx, secretName, updatedValue)
	if err != nil {
		t.Fatalf("UpdateSecretString failed: %v", err)
	}

	// Read again
	got, err = tc.client.GetSecretString(ctx, secretName)
	if err != nil {
		t.Fatalf("GetSecretString after update failed: %v", err)
	}
	if got != updatedValue {
		t.Errorf("updated value mismatch: got %s, want %s", got, updatedValue)
	}

	// Delete
	err = tc.client.DeleteSecret(ctx, secretName, secretsmanager.WithForceDelete())
	if err != nil {
		t.Fatalf("DeleteSecret failed: %v", err)
	}

	// Verify deleted
	_, err = tc.client.GetSecretString(ctx, secretName)
	if !errors.Is(err, secretsmanager.ErrSecretNotFound) {
		t.Errorf("expected ErrSecretNotFound after delete, got: %v", err)
	}
}
