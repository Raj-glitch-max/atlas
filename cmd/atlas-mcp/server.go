package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Minimal Model Context Protocol (MCP) server: JSON-RPC 2.0 over newline-
// delimited stdio, hand-rolled (no third-party dependency). It implements the
// subset a tool server needs: initialize, tools/list, tools/call, ping.
//
// The tools bridge to the Atlas Server so an agent can issue, verify, revoke,
// and inspect delegations on the same trust state the human UI shows.

const protocolVersion = "2024-11-05"

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"` // absent => notification
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  any             `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// tool is one MCP tool: a JSON Schema for its arguments and a handler that
// bridges to the Atlas API and returns human/agent-readable text.
type tool struct {
	Name        string
	Description string
	Schema      json.RawMessage
	Handler     func(ctx context.Context, c *Client, args json.RawMessage) (string, bool)
}

// Server holds the tool registry and the Atlas client.
type Server struct {
	client *Client
	tools  []tool
}

func NewServer(client *Client) *Server {
	return &Server{client: client, tools: atlasTools()}
}

// Serve runs the stdio JSON-RPC loop until stdin closes.
func (s *Server) Serve(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Buffer(make([]byte, 0, 64*1024), 8*1024*1024) // allow large messages
	enc := json.NewEncoder(out)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var req rpcRequest
		if err := json.Unmarshal(line, &req); err != nil {
			continue // malformed frame; MCP hosts don't expect a reply to garbage
		}
		resp, isNotification := s.dispatch(&req)
		if isNotification {
			continue
		}
		if err := enc.Encode(resp); err != nil {
			return err
		}
	}
	return sc.Err()
}

func (s *Server) dispatch(req *rpcRequest) (*rpcResponse, bool) {
	isNotification := len(req.ID) == 0
	reply := func(result any, e *rpcError) (*rpcResponse, bool) {
		if isNotification {
			return nil, true
		}
		return &rpcResponse{JSONRPC: "2.0", ID: req.ID, Result: result, Error: e}, false
	}

	switch req.Method {
	case "initialize":
		return reply(map[string]any{
			"protocolVersion": protocolVersion,
			"capabilities":    map[string]any{"tools": map[string]any{}},
			"serverInfo":      map[string]any{"name": "atlas-mcp", "version": mcpVersion},
			"instructions":    "Atlas issues and verifies offline-verifiable, attenuable delegation tokens bound to SPIFFE workload identity. Use atlas_issue to create a delegation, atlas_verify to check one, atlas_revoke to revoke by instance id, and atlas_graph/atlas_delegations/atlas_audit to inspect trust state.",
		}, nil)
	case "notifications/initialized", "notifications/cancelled":
		return nil, true
	case "ping":
		return reply(map[string]any{}, nil)
	case "tools/list":
		return reply(map[string]any{"tools": s.toolList()}, nil)
	case "tools/call":
		return s.callTool(req, reply)
	default:
		return reply(nil, &rpcError{Code: -32601, Message: "method not found: " + req.Method})
	}
}

func (s *Server) toolList() []map[string]any {
	out := make([]map[string]any, 0, len(s.tools))
	for _, t := range s.tools {
		out = append(out, map[string]any{
			"name":        t.Name,
			"description": t.Description,
			"inputSchema": t.Schema,
		})
	}
	return out
}

func (s *Server) callTool(req *rpcRequest, reply func(any, *rpcError) (*rpcResponse, bool)) (*rpcResponse, bool) {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return reply(nil, &rpcError{Code: -32602, Message: "invalid params: " + err.Error()})
	}
	for _, t := range s.tools {
		if t.Name == params.Name {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			text, isErr := t.Handler(ctx, s.client, params.Arguments)
			return reply(toolResult(text, isErr), nil)
		}
	}
	return reply(nil, &rpcError{Code: -32602, Message: "unknown tool: " + params.Name})
}

// toolResult builds the MCP tools/call result envelope.
func toolResult(text string, isErr bool) map[string]any {
	return map[string]any{
		"content": []map[string]any{{"type": "text", "text": text}},
		"isError": isErr,
	}
}

// pretty re-indents raw JSON for readable tool output.
func pretty(raw json.RawMessage) string {
	var buf []byte
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return string(raw)
	}
	buf, _ = json.MarshalIndent(v, "", "  ")
	return string(buf)
}

func errText(action string, err error) (string, bool) {
	return fmt.Sprintf("atlas_%s failed: %s", action, err.Error()), true
}
