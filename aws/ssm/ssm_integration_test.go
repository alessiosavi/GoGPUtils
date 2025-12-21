// Package ssm integration tests using LocalStack.
//go:build integration

package ssm_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alessiosavi/GoGPUtils/aws/internal/testutil"
	"github.com/alessiosavi/GoGPUtils/aws/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
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
	client *ssm.Client
	prefix string
}

// setupTest creates a test environment with SSM client.
func setupTest(t *testing.T) *testConfig {
	t.Helper()
	testutil.SkipIfNoLocalStack(t)

	cfg := testutil.MustLoadConfig(t)

	client, err := ssm.NewClient(cfg.AWSConfig)
	if err != nil {
		t.Fatalf("failed to create SSM client: %v", err)
	}

	// Use unique prefix for this test run
	prefix := testutil.UniqueParameterName()

	t.Cleanup(func() {
		// Clean up all parameters created during the test
		ctx := context.Background()
		params, _ := client.ListParametersByPath(ctx, prefix, ssm.WithRecursive(true))
		for _, p := range params {
			client.DeleteParameter(ctx, p.Name)
		}
	})

	return &testConfig{
		cfg:    cfg,
		client: client,
		prefix: prefix,
	}
}

func TestPutParameter(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	tests := []struct {
		name  string
		param string
		value string
	}{
		{
			name:  "simple parameter",
			param: tc.prefix + "/simple",
			value: "simple-value",
		},
		{
			name:  "nested parameter",
			param: tc.prefix + "/nested/deep/param",
			value: "nested-value",
		},
		{
			name:  "json value",
			param: tc.prefix + "/json",
			value: `{"key": "value", "number": 123}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tc.client.PutParameter(ctx, tt.param, tt.value)
			if err != nil {
				t.Fatalf("PutParameter failed: %v", err)
			}

			// Verify by getting
			got, err := tc.client.GetParameter(ctx, tt.param)
			if err != nil {
				t.Fatalf("GetParameter failed: %v", err)
			}

			if got != tt.value {
				t.Errorf("value mismatch: got %s, want %s", got, tt.value)
			}
		})
	}
}

func TestPutParameterWithOptions(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	t.Run("with description", func(t *testing.T) {
		param := tc.prefix + "/with-desc"
		err := tc.client.PutParameter(ctx, param, "value",
			ssm.WithDescription("Test parameter description"),
		)
		if err != nil {
			t.Fatalf("PutParameter with description failed: %v", err)
		}

		got, err := tc.client.GetParameter(ctx, param)
		if err != nil {
			t.Fatalf("GetParameter failed: %v", err)
		}

		if got != "value" {
			t.Errorf("value mismatch: got %s, want value", got)
		}
	})

	t.Run("with overwrite", func(t *testing.T) {
		param := tc.prefix + "/overwrite"

		// Create initial
		err := tc.client.PutParameter(ctx, param, "original")
		if err != nil {
			t.Fatalf("initial PutParameter failed: %v", err)
		}

		// Overwrite
		err = tc.client.PutParameter(ctx, param, "updated", ssm.WithOverwrite(true))
		if err != nil {
			t.Fatalf("PutParameter with overwrite failed: %v", err)
		}

		got, err := tc.client.GetParameter(ctx, param)
		if err != nil {
			t.Fatalf("GetParameter failed: %v", err)
		}

		if got != "updated" {
			t.Errorf("value mismatch: got %s, want updated", got)
		}
	})

	t.Run("secure string", func(t *testing.T) {
		param := tc.prefix + "/secure"
		err := tc.client.PutParameter(ctx, param, "secret-value",
			ssm.WithParameterType(types.ParameterTypeSecureString),
		)
		if err != nil {
			t.Fatalf("PutParameter SecureString failed: %v", err)
		}

		got, err := tc.client.GetParameter(ctx, param, ssm.WithDecryption(true))
		if err != nil {
			t.Fatalf("GetParameter failed: %v", err)
		}

		if got != "secret-value" {
			t.Errorf("value mismatch: got %s, want secret-value", got)
		}
	})
}

func TestGetParameter(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: create a parameter
	param := tc.prefix + "/get-test"
	value := "test-value-for-get"
	if err := tc.client.PutParameter(ctx, param, value); err != nil {
		t.Fatalf("setup PutParameter failed: %v", err)
	}

	t.Run("existing parameter", func(t *testing.T) {
		got, err := tc.client.GetParameter(ctx, param)
		if err != nil {
			t.Fatalf("GetParameter failed: %v", err)
		}

		if got != value {
			t.Errorf("value mismatch: got %s, want %s", got, value)
		}
	})

	t.Run("non-existent parameter", func(t *testing.T) {
		_, err := tc.client.GetParameter(ctx, tc.prefix+"/non-existent")
		if !errors.Is(err, ssm.ErrParameterNotFound) {
			t.Errorf("expected ErrParameterNotFound, got: %v", err)
		}
	})
}

func TestGetParameterInfo(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup
	param := tc.prefix + "/info-test"
	value := "info-test-value"
	if err := tc.client.PutParameter(ctx, param, value); err != nil {
		t.Fatalf("setup PutParameter failed: %v", err)
	}

	info, err := tc.client.GetParameterInfo(ctx, param)
	if err != nil {
		t.Fatalf("GetParameterInfo failed: %v", err)
	}

	if info.Name != param {
		t.Errorf("name mismatch: got %s, want %s", info.Name, param)
	}

	if info.Value != value {
		t.Errorf("value mismatch: got %s, want %s", info.Value, value)
	}

	if info.Version < 1 {
		t.Error("version should be at least 1")
	}

	if info.Type != string(types.ParameterTypeString) {
		t.Errorf("type mismatch: got %s, want String", info.Type)
	}
}

func TestGetParameters(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: create multiple parameters
	params := map[string]string{
		tc.prefix + "/multi/a": "value-a",
		tc.prefix + "/multi/b": "value-b",
		tc.prefix + "/multi/c": "value-c",
	}

	for name, value := range params {
		if err := tc.client.PutParameter(ctx, name, value); err != nil {
			t.Fatalf("setup PutParameter failed: %v", err)
		}
	}

	names := make([]string, 0, len(params))
	for name := range params {
		names = append(names, name)
	}

	// Add a non-existent parameter
	names = append(names, tc.prefix+"/multi/non-existent")

	values, invalid, err := tc.client.GetParameters(ctx, names)
	if err != nil {
		t.Fatalf("GetParameters failed: %v", err)
	}

	if len(values) != 3 {
		t.Errorf("values count mismatch: got %d, want 3", len(values))
	}

	if len(invalid) != 1 {
		t.Errorf("invalid count mismatch: got %d, want 1", len(invalid))
	}

	for name, expectedValue := range params {
		if got, ok := values[name]; !ok {
			t.Errorf("missing parameter %s", name)
		} else if got != expectedValue {
			t.Errorf("value mismatch for %s: got %s, want %s", name, got, expectedValue)
		}
	}
}

func TestListParametersByPath(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: create parameters under a path
	params := map[string]string{
		tc.prefix + "/list/param1":        "value1",
		tc.prefix + "/list/param2":        "value2",
		tc.prefix + "/list/nested/param3": "value3",
		tc.prefix + "/other/param4":       "value4",
	}

	for name, value := range params {
		if err := tc.client.PutParameter(ctx, name, value); err != nil {
			t.Fatalf("setup PutParameter failed: %v", err)
		}
	}

	t.Run("non-recursive", func(t *testing.T) {
		result, err := tc.client.ListParametersByPath(ctx, tc.prefix+"/list/")
		if err != nil {
			t.Fatalf("ListParametersByPath failed: %v", err)
		}

		// Should find param1 and param2, but not nested/param3
		if len(result) != 2 {
			t.Errorf("count mismatch: got %d, want 2", len(result))
		}
	})

	t.Run("recursive", func(t *testing.T) {
		result, err := tc.client.ListParametersByPath(ctx, tc.prefix+"/list/",
			ssm.WithRecursive(true),
		)
		if err != nil {
			t.Fatalf("ListParametersByPath recursive failed: %v", err)
		}

		// Should find all 3 under /list/
		if len(result) != 3 {
			t.Errorf("count mismatch: got %d, want 3", len(result))
		}
	})

	t.Run("with max results", func(t *testing.T) {
		result, err := tc.client.ListParametersByPath(ctx, tc.prefix+"/list/",
			ssm.WithRecursive(true),
			ssm.WithPathMaxResults(2),
		)
		if err != nil {
			t.Fatalf("ListParametersByPath with max results failed: %v", err)
		}

		if len(result) > 2 {
			t.Errorf("expected at most 2 results, got %d", len(result))
		}
	})
}

func TestListParameters(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: create some parameters
	for i := 0; i < 3; i++ {
		param := tc.prefix + "/listall/param" + string(rune('0'+i))
		if err := tc.client.PutParameter(ctx, param, "value"); err != nil {
			t.Fatalf("setup PutParameter failed: %v", err)
		}
	}

	names, err := tc.client.ListParameters(ctx)
	if err != nil {
		t.Fatalf("ListParameters failed: %v", err)
	}

	// Should find at least our 3 parameters
	found := 0
	for _, name := range names {
		if len(name) > len(tc.prefix) && name[:len(tc.prefix)] == tc.prefix {
			found++
		}
	}

	if found < 3 {
		t.Errorf("expected at least 3 parameters with prefix, found %d", found)
	}
}

func TestDeleteParameter(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: create a parameter
	param := tc.prefix + "/delete-test"
	if err := tc.client.PutParameter(ctx, param, "to-be-deleted"); err != nil {
		t.Fatalf("setup PutParameter failed: %v", err)
	}

	// Verify exists
	_, err := tc.client.GetParameter(ctx, param)
	if err != nil {
		t.Fatalf("parameter should exist: %v", err)
	}

	// Delete
	if err := tc.client.DeleteParameter(ctx, param); err != nil {
		t.Fatalf("DeleteParameter failed: %v", err)
	}

	// Verify deleted
	_, err = tc.client.GetParameter(ctx, param)
	if !errors.Is(err, ssm.ErrParameterNotFound) {
		t.Errorf("expected ErrParameterNotFound after delete, got: %v", err)
	}
}

func TestDeleteParameterNotFound(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	err := tc.client.DeleteParameter(ctx, tc.prefix+"/non-existent")
	if !errors.Is(err, ssm.ErrParameterNotFound) {
		t.Errorf("expected ErrParameterNotFound, got: %v", err)
	}
}

func TestValidationErrors(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	t.Run("empty name for get", func(t *testing.T) {
		_, err := tc.client.GetParameter(ctx, "")
		if err == nil {
			t.Error("expected error for empty parameter name")
		}
	})

	t.Run("empty name for put", func(t *testing.T) {
		err := tc.client.PutParameter(ctx, "", "value")
		if err == nil {
			t.Error("expected error for empty parameter name")
		}
	})

	t.Run("empty name for delete", func(t *testing.T) {
		err := tc.client.DeleteParameter(ctx, "")
		if err == nil {
			t.Error("expected error for empty parameter name")
		}
	})

	t.Run("empty path for list", func(t *testing.T) {
		_, err := tc.client.ListParametersByPath(ctx, "")
		if err == nil {
			t.Error("expected error for empty path")
		}
	})
}

func TestParameterVersioning(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	param := tc.prefix + "/versioning-test"

	// Create initial version
	if err := tc.client.PutParameter(ctx, param, "v1"); err != nil {
		t.Fatalf("initial PutParameter failed: %v", err)
	}

	info1, err := tc.client.GetParameterInfo(ctx, param)
	if err != nil {
		t.Fatalf("GetParameterInfo failed: %v", err)
	}

	if info1.Version != 1 {
		t.Errorf("initial version should be 1, got %d", info1.Version)
	}

	// Update to create new version
	if err := tc.client.PutParameter(ctx, param, "v2", ssm.WithOverwrite(true)); err != nil {
		t.Fatalf("update PutParameter failed: %v", err)
	}

	info2, err := tc.client.GetParameterInfo(ctx, param)
	if err != nil {
		t.Fatalf("GetParameterInfo after update failed: %v", err)
	}

	if info2.Version != 2 {
		t.Errorf("updated version should be 2, got %d", info2.Version)
	}

	if info2.Value != "v2" {
		t.Errorf("updated value should be v2, got %s", info2.Value)
	}
}
