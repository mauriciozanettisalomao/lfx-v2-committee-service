// Copyright The Linux Foundation and each contributor to LFX.
// SPDX-License-Identifier: MIT

package nats

import (
	"time"
)

// Config represents NATS configuration
type Config struct {
	// URL is the NATS server URL
	URL string `json:"url"`
	// Timeout is the request timeout duration
	Timeout time.Duration `json:"timeout"`
	// MaxReconnect is the maximum number of reconnection attempts
	MaxReconnect int `json:"max_reconnect"`
	// ReconnectWait is the time to wait between reconnection attempts
	ReconnectWait time.Duration `json:"reconnect_wait"`
}

// AccessCheckNATSRequest represents a NATS request for access checking
type AccessCheckNATSRequest struct {
	// Subject is the NATS subject for the request
	Subject string `json:"subject"`
	// Message is the serialized request data
	Message []byte `json:"message"`
	// Timeout is the request timeout duration
	Timeout time.Duration `json:"timeout"`
}

// AccessCheckNATSResponse represents a NATS response for access checking
type AccessCheckNATSResponse map[string]string
