// Copyright The Linux Foundation and each contributor to LFX.
// SPDX-License-Identifier: MIT

package model

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/gosimple/slug"
)

// Committee represents the core committee business entity
type Committee struct {
	CommitteeBase
	*CommitteeSettings
}

// Committee represents the core committee business entity
type CommitteeBase struct {
	UID              string    `json:"uid"`
	ProjectUID       string    `json:"project_uid"`
	ProjectName      string    `json:"project_name,omitempty"`
	Name             string    `json:"name"`
	Category         string    `json:"category"`
	Description      string    `json:"description,omitempty"`
	Website          *string   `json:"website,omitempty"`
	EnableVoting     bool      `json:"enable_voting"`
	SSOGroupEnabled  bool      `json:"sso_group_enabled"`
	SSOGroupName     string    `json:"sso_group_name,omitempty"`
	RequiresReview   bool      `json:"requires_review"`
	Public           bool      `json:"public"`
	Calendar         Calendar  `json:"calendar,omitempty"`
	DisplayName      string    `json:"display_name,omitempty"`
	ParentUID        *string   `json:"parent_uid,omitempty"`
	TotalMembers     int       `json:"total_members"`
	TotalVotingRepos int       `json:"total_voting_repos"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Calendar represents committee calendar settings
type Calendar struct {
	Public bool `json:"public"`
}

// SSOGroupNameBuild builds the SSO group name for the committee based on the project slug and committee name.
func (c *Committee) SSOGroupNameBuild(ctx context.Context, projectSlug string) error {

	baseName := slug.Make(fmt.Sprintf("%s-%s", projectSlug, c.Name))

	if c.SSOGroupName != "" {
		suffix := strings.TrimPrefix(c.SSOGroupName, baseName)

		if suffix == "" {
			suffix = "1"
		}
		suffix = strings.Trim(suffix, "-")

		// if the suffix is a number, we can increment it
		num, err := strconv.Atoi(suffix)
		if err != nil {
			slog.ErrorContext(ctx, "failed to parse SSO group name suffix as number",
				"error", err,
				"ssogroup_name", c.SSOGroupName,
			)
			return fmt.Errorf("failed to parse SSO group name suffix: %w", err)

		}

		baseName = fmt.Sprintf("%s-%d", baseName, num+1)

	}

	c.SSOGroupName = baseName

	return nil
}

// BuildIndexKey generates a SHA-256 hash for use as a NATS KV key.
// This is necessary because the original input may contain special characters,
// exceed length limits, or have inconsistent formatting, and we do not control its content.
// Using a hash ensures a safe, fixed-length, and deterministic key.
func (c *Committee) BuildIndexKey(ctx context.Context) string {
	// Combine project_uid and committee name with a delimiter
	data := fmt.Sprintf("%s|%s", c.ProjectUID, c.Name)

	hash := sha256.Sum256([]byte(data))

	key := hex.EncodeToString(hash[:])

	slog.DebugContext(ctx, "index key built",
		"project_uid", c.ProjectUID,
		"committee_name", c.Name,
		"key", key,
	)

	return key
}
