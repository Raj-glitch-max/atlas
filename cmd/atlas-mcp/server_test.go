package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// mockAtlas stands in for cmd/atlas-server so the MCP layer can be tested in
// isolation.
func mockAtlas(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
	})
	mux.HandleFunc("/issue", func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"record": "eyJ...jws", "instance": "abc123", "scope": []string{"read:orders"}})
	})
	mux.HandleFunc("/verify", func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"decision": "accept", "accept": true, "causes": []string{}, "latencyMicros": 207, "trace": []any{}})
	})
	mux.HandleFunc("/revoke", func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"revoked": true, "instance": "abc123"})
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

// run feeds newline-delimited JSON-RPC requests through the server and returns
// the decoded responses.
func run(t *testing.T, srv *Server, reqs ...string) []rpcResponse {
	t.Helper()
	in := strings.NewReader(strings.Join(reqs, "\n") + "\n")
	var out bytes.Buffer
	if err := srv.Serve(in, &out); err != nil {
		t.Fatalf("serve: %v", err)
	}
	var resps []rpcResponse
	dec := json.NewDecoder(&out)
	for dec.More() {
		var r rpcResponse
		if err := dec.Decode(&r); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		resps = append(resps, r)
	}
	return resps
}

func TestInitializeAndToolsList(t *testing.T) {
	ts := mockAtlas(t)
	srv := NewServer(NewClient(ts.URL, ""))
	resps := run(t, srv,
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{}}}`,
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`,
	)
	// notification produces no response → expect exactly 2 responses.
	if len(resps) != 2 {
		t.Fatalf("want 2 responses, got %d", len(resps))
	}
	init, _ := json.Marshal(resps[0].Result)
	if !strings.Contains(string(init), "atlas-mcp") || !strings.Contains(string(init), "2024-11-05") {
		t.Fatalf("bad initialize result: %s", init)
	}
	list, _ := json.Marshal(resps[1].Result)
	for _, want := range []string{"atlas_issue", "atlas_verify", "atlas_revoke", "atlas_graph", "atlas_audit", "atlas_delegations"} {
		if !strings.Contains(string(list), want) {
			t.Fatalf("tools/list missing %q: %s", want, list)
		}
	}
}

func TestToolCallIssueAndVerify(t *testing.T) {
	ts := mockAtlas(t)
	srv := NewServer(NewClient(ts.URL, ""))

	resps := run(t, srv,
		`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"atlas_issue","arguments":{"principal":"spiffe://domain-a.test/p","delegate":"spiffe://domain-b.test/d","scope":["read:orders"]}}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"atlas_verify","arguments":{"record":"eyJ...jws"}}}`,
	)
	if len(resps) != 2 {
		t.Fatalf("want 2 responses, got %d", len(resps))
	}
	if text := resultText(t, resps[0]); !strings.Contains(text, "abc123") {
		t.Fatalf("issue result missing instance: %s", text)
	}
	if text := resultText(t, resps[1]); !strings.Contains(text, "accept") || !strings.Contains(text, "207") {
		t.Fatalf("verify result wrong: %s", text)
	}
}

func TestUnknownToolIsError(t *testing.T) {
	ts := mockAtlas(t)
	srv := NewServer(NewClient(ts.URL, ""))
	resps := run(t, srv, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"atlas_nope","arguments":{}}}`)
	if len(resps) != 1 || resps[0].Error == nil {
		t.Fatalf("unknown tool should return a JSON-RPC error, got %+v", resps)
	}
}

func TestUnreachableServerToolReportsError(t *testing.T) {
	srv := NewServer(NewClient("http://127.0.0.1:1", "")) // nothing listening
	resps := run(t, srv, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"atlas_delegations","arguments":{}}}`)
	if len(resps) != 1 {
		t.Fatalf("want 1 response, got %d", len(resps))
	}
	text := resultText(t, resps[0])
	if !strings.Contains(text, "failed") {
		t.Fatalf("expected an error result, got: %s", text)
	}
	// isError must be true so the agent knows the call failed.
	m, _ := json.Marshal(resps[0].Result)
	if !strings.Contains(string(m), "\"isError\":true") {
		t.Fatalf("expected isError:true, got %s", m)
	}
}

func TestMethodNotFound(t *testing.T) {
	ts := mockAtlas(t)
	srv := NewServer(NewClient(ts.URL, ""))
	resps := run(t, srv, `{"jsonrpc":"2.0","id":9,"method":"bogus/method"}`)
	if len(resps) != 1 || resps[0].Error == nil || resps[0].Error.Code != -32601 {
		t.Fatalf("want method-not-found error, got %+v", resps)
	}
}

func resultText(t *testing.T, r rpcResponse) string {
	t.Helper()
	m, _ := json.Marshal(r.Result)
	var res struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	_ = json.Unmarshal(m, &res)
	if len(res.Content) == 0 {
		t.Fatalf("no content in result: %s", m)
	}
	return res.Content[0].Text
}

// checkpoint: fix(stores): fix truststore backend
