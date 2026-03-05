package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer("md2gdoc", "1.0.0",
		server.WithToolCapabilities(true),
	)

	// Tool: convert markdown to plain text (structured for Google Docs)
	s.AddTool(
		mcp.NewTool("md_to_text",
			mcp.WithDescription("Convert Markdown content to structured plain text suitable for Google Docs import"),
			mcp.WithString("markdown", mcp.Required(), mcp.Description("The Markdown content to convert")),
		),
		handleMdToText,
	)

	// Tool: parse markdown headings
	s.AddTool(
		mcp.NewTool("parse_headings",
			mcp.WithDescription("Extract all headings from a Markdown document with their levels"),
			mcp.WithString("markdown", mcp.Required(), mcp.Description("The Markdown content to parse")),
		),
		handleParseHeadings,
	)

	// Tool: health/ping
	s.AddTool(
		mcp.NewTool("ping",
			mcp.WithDescription("Health check - returns pong"),
		),
		handlePing,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting md2gdoc MCP server on %s", addr)

	httpServer := server.NewStreamableHTTPServer(s)
	if err := httpServer.Start(addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func handleMdToText(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	markdown, err := req.RequireString("markdown")
	if err != nil {
		return mcp.NewToolResultError("markdown parameter is required"), nil
	}

	// Simple markdown to structured text conversion
	lines := strings.Split(markdown, "\n")
	var result []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Convert headings
		if strings.HasPrefix(trimmed, "### ") {
			result = append(result, strings.ToUpper(strings.TrimPrefix(trimmed, "### ")))
		} else if strings.HasPrefix(trimmed, "## ") {
			result = append(result, "")
			result = append(result, strings.ToUpper(strings.TrimPrefix(trimmed, "## ")))
			result = append(result, strings.Repeat("-", len(strings.TrimPrefix(trimmed, "## "))))
		} else if strings.HasPrefix(trimmed, "# ") {
			result = append(result, "")
			result = append(result, strings.ToUpper(strings.TrimPrefix(trimmed, "# ")))
			result = append(result, strings.Repeat("=", len(strings.TrimPrefix(trimmed, "# "))))
		} else if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			// Convert list items
			result = append(result, "  • "+trimmed[2:])
		} else if strings.HasPrefix(trimmed, "```") {
			// Skip code fences
			continue
		} else if strings.HasPrefix(trimmed, "**") && strings.HasSuffix(trimmed, "**") {
			// Bold text
			result = append(result, strings.Trim(trimmed, "*"))
		} else {
			result = append(result, line)
		}
	}

	return mcp.NewToolResultText(strings.Join(result, "\n")), nil
}

func handleParseHeadings(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	markdown, err := req.RequireString("markdown")
	if err != nil {
		return mcp.NewToolResultError("markdown parameter is required"), nil
	}

	lines := strings.Split(markdown, "\n")
	var headings []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			level := 0
			for _, c := range trimmed {
				if c == '#' {
					level++
				} else {
					break
				}
			}
			title := strings.TrimSpace(strings.TrimLeft(trimmed, "#"))
			headings = append(headings, fmt.Sprintf("H%d: %s", level, title))
		}
	}

	if len(headings) == 0 {
		return mcp.NewToolResultText("No headings found."), nil
	}

	return mcp.NewToolResultText(strings.Join(headings, "\n")), nil
}

func handlePing(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText("pong"), nil
}
