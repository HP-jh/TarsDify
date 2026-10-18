package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HP-jh/TarsDify/pkg/config"
	"github.com/HP-jh/TarsDify/pkg/upstream"
)

var Version = "1.0.1"
const RepoName = "TarsDify"

func main() {
	version := flag.Bool("version", false, "print version and exit")
	flag.Parse()
	if *version { fmt.Printf("%s v%s\n", RepoName, Version); return }
	cfg, err := config.Load("")
	if err != nil { log.Fatalf("加载配置失败：%v", err) }
	client := upstream.NewClient(cfg.Endpoint, cfg.APIKey, cfg.Timeout)
	log.Printf("[%s] starting v%s endpoint=%s", RepoName, Version, cfg.Endpoint)
	mode := os.Getenv("TARS_MCP_MODE")
	if mode == "" { mode = "stdio" }
	srv := &http.Server{Addr: ":8080", ReadTimeout: time.Duration(cfg.Timeout) * time.Second, WriteTimeout: time.Duration(cfg.Timeout) * time.Second}
	http.HandleFunc("/health", withClient(client, healthHandler))
	http.HandleFunc("/tools", withClient(client, toolsHandler))
	http.HandleFunc("/invoke", withClient(client, invokeHandler))
	http.HandleFunc("/models", withClient(client, modelsHandler))
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Printf("[%s] shutting down", RepoName)
		_ = srv.Close()
	}()
	if mode == "stdio" { runStdioServer(client); return }
	log.Printf("[%s] HTTP listening on %s", RepoName, srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func withClient(client *upstream.Client, fn func(http.ResponseWriter, *http.Request, *upstream.Client)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions { w.WriteHeader(http.StatusNoContent); return }
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		fn(w, r, client)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request, client *upstream.Client) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	up := client.Health(ctx)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok", "version": Version, "name": RepoName, "upstream": up,
		"time": time.Now().UTC().Format(time.RFC3339),
	})
}

func toolsHandler(w http.ResponseWriter, r *http.Request, client *upstream.Client) {
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"tools": []map[string]interface{}{
		{
			"name":        "dify_chat",
			"description": "Send a chat message to a Dify app",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{"conversation_id": map[string]interface{}{"type": "string", "description": "Parameter conversation_id"}, "query": map[string]interface{}{"type": "string", "description": "Parameter query"}, "user": map[string]interface{}{"type": "string", "description": "Parameter user"}},
				"required":   []interface{}{"query", "user"},
			},
		},
		{
			"name":        "dify_workflow",
			"description": "Run a Dify workflow",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{"inputs": map[string]interface{}{"type": "object", "description": "Parameter inputs"}, "user": map[string]interface{}{"type": "string", "description": "Parameter user"}},
				"required":   []interface{}{"inputs", "user"},
			},
		},
		{
			"name":        "dify_knowledge_retrieve",
			"description": "Retrieve from Dify knowledge base",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{"dataset_id": map[string]interface{}{"type": "string", "description": "Parameter dataset_id"}, "query": map[string]interface{}{"type": "string", "description": "Parameter query"}},
				"required":   []interface{}{"dataset_id", "query"},
			},
		},
		{
			"name":        "dify_files",
			"description": "Upload a file to Dify",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{"file": map[string]interface{}{"type": "string", "description": "Parameter file"}},
				"required":   []interface{}{"file"},
			},
		},
		{
			"name":        "dify_conversations",
			"description": "List conversations",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{"limit": map[string]interface{}{"type": "number", "description": "Parameter limit"}, "user": map[string]interface{}{"type": "string", "description": "Parameter user"}},
				"required":   []interface{}{"user"},
			},
		},
		{
			"name":        "dify_health",
			"description": "Check Dify instance health",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
				"required":   []interface{}{},
			},
		},
		{
			"name":        "dify_status",
			"description": "Integration status",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
				"required":   []interface{}{},
			},
		},
		},
	})
}

func modelsHandler(w http.ResponseWriter, r *http.Request, client *upstream.Client) {
	if r.Method != http.MethodGet { w.WriteHeader(http.StatusMethodNotAllowed); _ = json.NewEncoder(w).Encode(map[string]interface{}{"error": "只接受 GET"}); return }
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	models, err := client.ListModels(ctx)
	if err != nil { w.WriteHeader(http.StatusBadGateway); _ = json.NewEncoder(w).Encode(map[string]interface{}{"error": fmt.Sprintf("获取失败：%v", err)}); return }
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"models": models})
}

func invokeHandler(w http.ResponseWriter, r *http.Request, client *upstream.Client) {
	if r.Method != http.MethodPost { w.WriteHeader(http.StatusMethodNotAllowed); _ = json.NewEncoder(w).Encode(map[string]interface{}{"error": "只接受 POST"}); return }
	var req struct { Tool string `json:"tool"`; Params map[string]interface{} `json:"params"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { w.WriteHeader(http.StatusBadRequest); _ = json.NewEncoder(w).Encode(map[string]interface{}{"error": fmt.Sprintf("解析失败：%v", err)}); return }
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()
	switch req.Tool {
	case "dify_chat":
		if req.Params["query"] == nil { w.WriteHeader(http.StatusBadRequest); _ = json.NewEncoder(w).Encode(map[string]interface{}{"error": "缺少参数：query"}); return }
		if req.Params["user"] == nil { w.WriteHeader(http.StatusBadRequest); _ = json.NewEncoder(w).Encode(map[string]interface{}{"error": "缺少参数：user"}); return }
		body := map[string]interface{}{"conversation_id": req.Params["conversation_id"], "query": req.Params["query"], "user": req.Params["user"]}
		result, err := client.InvokeRaw(ctx, "POST", "/v1/chat-messages", body)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": fmt.Sprintf("调用失败：%v", err)})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"result": result})
	case "dify_workflow":
		if req.Params["inputs"] == nil { w.WriteHeader(http.StatusBadRequest); _ = json.NewEncoder(w).Encode(map[string]interface{}{"error": "缺少参数：inputs"}); return }
		if req.Params["user"] == nil { w.WriteHeader(http.StatusBadRequest); _ = json.NewEncoder(w).Encode(map[string]interface{}{"error": "缺少参数：user"}); return }
		body := map[string]interface{}{"inputs": req.Params["inputs"], "user": req.Params["user"]}
		result, err := client.InvokeRaw(ctx, "POST", "/v1/workflows/run", body)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": fmt.Sprintf("调用失败：%v", err)})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"result": result})
	case "dify_knowledge_retrieve":
		if req.Params["dataset_id"] == nil { w.WriteHeader(http.StatusBadRequest); _ = json.NewEncoder(w).Encode(map[string]interface{}{"error": "缺少参数：dataset_id"}); return }
		if req.Params["query"] == nil { w.WriteHeader(http.StatusBadRequest); _ = json.NewEncoder(w).Encode(map[string]interface{}{"error": "缺少参数：query"}); return }
		body := map[string]interface{}{"dataset_id": req.Params["dataset_id"], "query": req.Params["query"]}
		result, err := client.InvokeRaw(ctx, "POST", "/v1/datasets/{dataset_id}/retrieve", body)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": fmt.Sprintf("调用失败：%v", err)})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"result": result})
	case "dify_files":
		if req.Params["file"] == nil { w.WriteHeader(http.StatusBadRequest); _ = json.NewEncoder(w).Encode(map[string]interface{}{"error": "缺少参数：file"}); return }
		body := map[string]interface{}{"file": req.Params["file"]}
		result, err := client.InvokeRaw(ctx, "POST", "/v1/files", body)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": fmt.Sprintf("调用失败：%v", err)})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"result": result})
	case "dify_conversations":
		if req.Params["user"] == nil { w.WriteHeader(http.StatusBadRequest); _ = json.NewEncoder(w).Encode(map[string]interface{}{"error": "缺少参数：user"}); return }
		limitF, _ := req.Params["limit"].(float64); limit := int(limitF); if limit <= 0 { limit = 10 }
		path := fmt.Sprintf("/v1/conversations?limit=%d", limit)
		result, err := client.InvokeRaw(ctx, "GET", path, nil)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": fmt.Sprintf("调用失败：%v", err)})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"result": result})
	case "dify_health":
		result, err := client.InvokeRaw(ctx, "GET", "/v1/health", nil)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": fmt.Sprintf("调用失败：%v", err)})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"result": result})
	case "dify_status":
		result, err := client.InvokeRaw(ctx, "GET", "/status", nil)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": fmt.Sprintf("调用失败：%v", err)})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"result": result})
	default:
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": fmt.Sprintf("不认识工具：%s", req.Tool)})
	}
}

func runStdioServer(client *upstream.Client) {
	log.Printf("[%s] stdio mode", RepoName)
	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	for {
		var msg map[string]interface{}
		if err := dec.Decode(&msg); err != nil { return }
		id := msg["id"]; method, _ := msg["method"].(string)
		var result map[string]interface{}
		switch method {
		case "tools/list":
			result = map[string]interface{}{"tools": []map[string]interface{}{{"name": "dify_chat", "description": "Send a chat message to a Dify app"}, {"name": "dify_workflow", "description": "Run a Dify workflow"}, {"name": "dify_knowledge_retrieve", "description": "Retrieve from Dify knowledge base"}, {"name": "dify_files", "description": "Upload a file to Dify"}, {"name": "dify_conversations", "description": "List conversations"}, {"name": "dify_health", "description": "Check Dify instance health"}, {"name": "dify_status", "description": "Integration status"}}}
		case "tools/call":
			params, _ := msg["params"].(map[string]interface{})
			toolName, _ := params["name"].(string)
			args, _ := params["arguments"].(map[string]interface{})
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			result = handleStdioTool(ctx, client, toolName, args)
			cancel()
		default:
			result = map[string]interface{}{"message": fmt.Sprintf("%s v%s ready", RepoName, Version)}
		}
		_ = enc.Encode(map[string]interface{}{"jsonrpc": "2.0", "id": id, "result": result})
	}
}

func handleStdioTool(ctx context.Context, client *upstream.Client, toolName string, args map[string]interface{}) map[string]interface{} {
	switch toolName {
	case "dify_chat":
		if args["query"] == nil { return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": "缺少参数：query"}}} }
		if args["user"] == nil { return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": "缺少参数：user"}}} }
		result, err := client.InvokeRaw(ctx, "POST", "/v1/chat-messages", map[string]interface{}{"conversation_id": args["conversation_id"], "query": args["query"], "user": args["user"]})
		if err != nil {
			return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": fmt.Sprintf("调用失败：%v", err)}}}
		}
		data, _ := json.Marshal(result)
		return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": string(data)}}}
	case "dify_workflow":
		if args["inputs"] == nil { return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": "缺少参数：inputs"}}} }
		if args["user"] == nil { return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": "缺少参数：user"}}} }
		result, err := client.InvokeRaw(ctx, "POST", "/v1/workflows/run", map[string]interface{}{"inputs": args["inputs"], "user": args["user"]})
		if err != nil {
			return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": fmt.Sprintf("调用失败：%v", err)}}}
		}
		data, _ := json.Marshal(result)
		return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": string(data)}}}
	case "dify_knowledge_retrieve":
		if args["dataset_id"] == nil { return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": "缺少参数：dataset_id"}}} }
		if args["query"] == nil { return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": "缺少参数：query"}}} }
		result, err := client.InvokeRaw(ctx, "POST", "/v1/datasets/{dataset_id}/retrieve", map[string]interface{}{"dataset_id": args["dataset_id"], "query": args["query"]})
		if err != nil {
			return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": fmt.Sprintf("调用失败：%v", err)}}}
		}
		data, _ := json.Marshal(result)
		return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": string(data)}}}
	case "dify_files":
		if args["file"] == nil { return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": "缺少参数：file"}}} }
		result, err := client.InvokeRaw(ctx, "POST", "/v1/files", map[string]interface{}{"file": args["file"]})
		if err != nil {
			return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": fmt.Sprintf("调用失败：%v", err)}}}
		}
		data, _ := json.Marshal(result)
		return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": string(data)}}}
	case "dify_conversations":
		if args["user"] == nil { return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": "缺少参数：user"}}} }
		result, err := client.InvokeRaw(ctx, "GET", "/v1/conversations", nil)
		if err != nil {
			return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": fmt.Sprintf("调用失败：%v", err)}}}
		}
		data, _ := json.Marshal(result)
		return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": string(data)}}}
	case "dify_health":
		result, err := client.InvokeRaw(ctx, "GET", "/v1/health", nil)
		if err != nil {
			return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": fmt.Sprintf("调用失败：%v", err)}}}
		}
		data, _ := json.Marshal(result)
		return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": string(data)}}}
	case "dify_status":
		result, err := client.InvokeRaw(ctx, "GET", "/status", nil)
		if err != nil {
			return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": fmt.Sprintf("调用失败：%v", err)}}}
		}
		data, _ := json.Marshal(result)
		return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": string(data)}}}
	default:
		return map[string]interface{}{"content": []map[string]string{{"type": "text", "text": fmt.Sprintf("不认识工具：%s", toolName)}}}
	}
}
