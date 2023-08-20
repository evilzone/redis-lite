package core

import (
	"errors"
	"testing"
)

func TestProtocol(t *testing.T) {

	t.Run("isParamValid", func(t *testing.T) {
		t.Run("should return true", func(t *testing.T) {
			cmd := CMDGet
			params := []string{"key"}

			if !cmd.isParamValid(params) {
				t.Errorf("Expected true, got false")
			}
		})

		t.Run("should return false", func(t *testing.T) {
			cmd := CMDGet
			params := []string{}

			if cmd.isParamValid(params) {
				t.Errorf("Expected false, got true")
			}
		})
	})

	t.Run("parseCommand", func(t *testing.T) {
		t.Run("should return CMDGet", func(t *testing.T) {
			command, err := parseCommand("GET")

			if command != CMDGet {
				t.Errorf("Expected CMDGet, got %v", command)
			}

			if err != nil {
				t.Errorf("Expected nil, got not nil")
			}
		})

		t.Run("should return CMDSet", func(t *testing.T) {
			command, err := parseCommand("SET")

			if command != CMDSet {
				t.Errorf("Expected CMDSet, got %v", command)
			}

			if err != nil {
				t.Errorf("Expected nil, got not nil")
			}
		})

		t.Run("should return CMDDel", func(t *testing.T) {
			command, err := parseCommand("DEL")

			if command != CMDDel {
				t.Errorf("Expected CMDDel, got %v", command)
			}

			if err != nil {
				t.Errorf("Expected nil, got not nil")
			}
		})

		t.Run("should return CMDPing", func(t *testing.T) {
			command, err := parseCommand("PING")

			if command != CMDPing {
				t.Errorf("Expected CMDPing, got %v", command)
			}

			if err != nil {
				t.Errorf("Expected nil, got not nil")
			}
		})

		t.Run("should return empty command", func(t *testing.T) {
			command, err := parseCommand("INVALID_COMMAND")

			if command != (Command{}) {
				t.Errorf("Expected Command{}, got %v", command)
			}

			if !errors.Is(err, ErrorInvalidCommand) {
				t.Errorf("Expected ErrorInvalidCommand, got not %v", err)
			}
		})
	})

	t.Run("parseProtocol", func(t *testing.T) {
		t.Run("should return error", func(t *testing.T) {

			_, err := parseProtocol("")

			if err == nil {
				t.Errorf("Expected error, got nil")
			}

			if !errors.Is(err, ErrorInvalidCommand) {
				t.Errorf("Expected ErrorInvalidCommand, got %v", err)
			}
		})

		t.Run("parses valid GET command", func(t *testing.T) {

			req, err := parseProtocol("GET key")

			if err != nil {
				t.Errorf("Expected nil, got error %v", err)
			}

			if req.Command != CMDGet {
				t.Errorf("Expected CMDGet, got command %s", req.Command.Cmd)
			}

			if len(req.Params) != 1 {
				t.Errorf("Expected len 1, got command %d", len(req.Params))
			}

			if req.Params[0] != "key" {
				t.Errorf("Expected 'key', got  %s", req.Params[0])
			}
		})

		t.Run("parses valid SET command", func(t *testing.T) {

			req, err := parseProtocol("SET key value 1000")

			if err != nil {
				t.Errorf("Expected nil, got error %v", err)
			}

			if req.Command != CMDSet {
				t.Errorf("Expected CMDSet, got command %s", req.Command.Cmd)
			}

			if len(req.Params) != 3 {
				t.Errorf("Expected len 1, got command %d", len(req.Params))
			}

			if req.Params[0] != "key" {
				t.Errorf("Expected 'key', got  %s", req.Params[0])
			}

			if req.Params[1] != "value" {
				t.Errorf("Expected 'value', got  %s", req.Params[1])
			}

			if req.Params[2] != "1000" {
				t.Errorf("Expected 'key', got  %s", req.Params[2])
			}
		})
	})
}
