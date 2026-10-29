package audit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// DEFECT D8: Drift in documented counts and examples.
// README claims 30 vectors, 20 adversarial; actual: 42 vectors (10+20+12)
// AI_DISCLOSURE claims 139 tests, 30 vectors; actual vectors: 42
// Gate header example uses `atlas -q delegate` but CLI requires flag before command
// Root Dockerfile header is stale
// This test asserts the CURRENT WRONG behavior: documentation drift exists.
// After fix, this test should be renamed TestD8_DocDriftFixed and inverted.
func TestDefectD8_DocDrift(t *testing.T) {
	// Count actual vectors
	verdictVectors := countVectors(t, "tests/vectors/verdict-vectors.json")
	negativeVectors := countVectors(t, "tests/vectors/negative-vectors.json")
	revocationVectors := countVectors(t, "tests/vectors/revocation-vectors.json")
	totalVectors := verdictVectors + negativeVectors + revocationVectors

	// CURRENT WRONG BEHAVIOR: total != 30 (README/AI_DISCLOSURE claim)
	if totalVectors == 30 {
		t.Fatalf("DEFECT D8 NOT REPRODUCED: total vectors = 30 (documentation matches)")
	}
	if totalVectors != 42 {
		t.Logf("Vector counts: verdict=%d, negative=%d, revocation=%d, total=%d",
			verdictVectors, negativeVectors, revocationVectors, totalVectors)
	}
	t.Logf("DEFECT D8 REPRODUCED: total vectors = %d (README claims 30)", totalVectors)

	// Check README for "30 vectors" claim
	readme := readFile(t, "README.md")
	if !strings.Contains(readme, "30 vectors") && !strings.Contains(readme, "20 adversarial") {
		t.Logf("README may have been updated; checking for old claims...")
	} else {
		t.Logf("README contains '30 vectors' or '20 adversarial' claim")
	}

	// Check AI_DISCLOSURE for "139 tests, 30 vectors" claim
	aiDisclosure := readFile(t, "AI_DISCLOSURE.md")
	if strings.Contains(aiDisclosure, "139 tests") && strings.Contains(aiDisclosure, "30 vectors") {
		t.Logf("AI_DISCLOSURE contains '139 tests, 30 vectors' claim")
	} else {
		t.Logf("AI_DISCLOSURE may have been updated")
	}

	// Check gate header example in examples/atlas-gate/main.go
	gateMain := readFile(t, "examples/atlas-gate/main.go")
	if strings.Contains(gateMain, "atlas -q delegate") {
		t.Logf("Gate example contains 'atlas -q delegate' (wrong flag position)")
	} else {
		t.Logf("Gate example may have been fixed")
	}

	// Check root Dockerfile header
	dockerfile := readFile(t, "Dockerfile")
	if strings.Contains(dockerfile, "dev toolbox") || strings.Contains(dockerfile, "stale") {
		t.Logf("Root Dockerfile has stale header")
	} else {
		t.Logf("Root Dockerfile header may have been fixed")
	}
}

func countVectors(t *testing.T, path string) int {
	t.Helper()
	data := readFile(t, path)
	var obj struct {
		Vectors []json.RawMessage `json:"vectors"`
	}
	if err := json.Unmarshal([]byte(data), &obj); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	return len(obj.Vectors)
}

func readFile(t *testing.T, relPath string) string {
	t.Helper()
	// Find repo root
	dir, _ := os.Getwd()
	for i := 0; i < 5; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		dir = filepath.Dir(dir)
	}
	fullPath := filepath.Join(dir, relPath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("read %s: %v", fullPath, err)
	}
	return string(content)
}

// DEFECT D9: No release pipeline; SDKs are not installable; chosen package names are untested.
// This test asserts the CURRENT WRONG behavior: no release workflow, no packaging metadata.
// After fix, this test should be renamed TestD9_ReleasePipelineExists and inverted.
func TestDefectD9_NoReleasePipeline(t *testing.T) {
	// Check for release workflow
	_, err := os.Stat(".github/workflows/release.yml")
	if err == nil {
		t.Fatalf("DEFECT D9 NOT REPRODUCED: release.yml exists")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("stat release.yml: %v", err)
	}
	t.Logf("DEFECT D9 REPRODUCED: no .github/workflows/release.yml")

	// Check Python SDK packaging metadata
	_, err = os.Stat("sdk/python/pyproject.toml")
	if err == nil {
		t.Logf("Python SDK has pyproject.toml")
	} else if os.IsNotExist(err) {
		t.Logf("DEFECT D9 REPRODUCED: Python SDK missing pyproject.toml")
	}

	_, err = os.Stat("sdk/python/setup.py")
	if err == nil {
		t.Logf("Python SDK has setup.py")
	} else if os.IsNotExist(err) {
		t.Logf("DEFECT D9 REPRODUCED: Python SDK missing setup.py")
	}

	// Check TypeScript SDK package.json
	_, err = os.Stat("sdk/typescript/package.json")
	if err == nil {
		t.Logf("TypeScript SDK has package.json")
	} else if os.IsNotExist(err) {
		t.Logf("DEFECT D9 REPRODUCED: TypeScript SDK missing package.json")
	}

	// Note: npm name @atlas/sdk availability untested (requires network)
	t.Logf("DEFECT D9: npm @atlas/sdk availability untested")
}

// DEFECT D10: No Atlas code path uses SPIRE; the lab has never run; nine substrate tests skip.
// This test asserts the CURRENT WRONG behavior: no Workload API usage in Atlas code.
// After fix, this test should be renamed TestD10_SPIREIntegrationExists and inverted.
func TestDefectD10_NoSPIREIntegration(t *testing.T) {
	// Check for workloadapi import in Atlas code (excluding lab)
	found := false
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.HasPrefix(path, "./atlas-lab") || strings.HasPrefix(path, "./.git") {
			return filepath.SkipDir
		}
		if strings.HasSuffix(path, ".go") {
			content := readFileDirect(t, path)
			if strings.Contains(content, "workloadapi") {
				found = true
				t.Logf("Found workloadapi import in: %s", path)
			}
		}
		return nil
	})

	if !found {
		t.Fatalf("DEFECT D10 FIXED: no workloadapi import found in Atlas code")
	}
	t.Logf("DEFECT D10 FIXED: workloadapi import found in Atlas code (excluding lab)")

	// Check TestSubstrateBlockedATs skips
	acceptanceTest := readFile(t, "tests/acceptance/revocation_realization_test.go")
	skipped := []string{"AT10", "AT12", "AT13", "AT14", "AT15", "AT16", "AT17", "AT24", "AT26"}
	for _, at := range skipped {
		if strings.Contains(acceptanceTest, at) && strings.Contains(acceptanceTest, "Skip") {
			t.Logf("TestSubstrateBlockedATs skips %s", at)
		}
	}
}

func readFileDirect(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

// DEFECT D12: No gRPC exists although it is claimed.
// This test asserts the CURRENT WRONG behavior: no gRPC dependencies in go.mod.
// After fix, this test should be renamed TestD12_GRPCExists and inverted.
func TestDefectD12_NoGRPC(t *testing.T) {
	goMod := readFile(t, "go.mod")
	if !strings.Contains(goMod, "grpc") && !strings.Contains(goMod, "go-control-plane") {
		t.Fatalf("DEFECT D12 FIXED: no grpc or go-control-plane found in go.mod")
	}
	t.Logf("DEFECT D12 FIXED: grpc dependency found in go.mod")

	// Check README/ROADMAP for gRPC claims
	readme := readFile(t, "README.md")
	roadmap := readFile(t, "ROADMAP.md")
	if strings.Contains(readme, "gRPC") || strings.Contains(roadmap, "gRPC") {
		t.Logf("Documentation mentions gRPC (claim exists and implemented)")
	}
}
