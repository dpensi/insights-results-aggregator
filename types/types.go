// Copyright 2020, 2021, 2022 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package types contains declaration of various data types (usually structures)
// used elsewhere in the aggregator code.
package types

import (
	"time"

	types "github.com/RedHatInsights/insights-results-types"
)

// OrgID represents organization ID
type OrgID = types.OrgID

// Account represent the account number
type Account = uint32

// ClusterName represents name of cluster in format c8590f31-e97e-4b85-b506-c45ce1911a12
type ClusterName = types.ClusterName

// UserID represents type for user id
type UserID = types.UserID

// ClusterReport represents cluster report
type ClusterReport = types.ClusterReport

// Timestamp represents any timestamp in a form gathered from database
// TODO: need to be improved
type Timestamp = types.Timestamp

// UserVote is a type for user's vote
type UserVote = types.UserVote

// InfoKey represents type for info key
type InfoKey string

const (
	// UserVoteDislike shows user's dislike
	UserVoteDislike = types.UserVoteDislike
	// UserVoteNone shows no vote from user
	UserVoteNone = types.UserVoteNone
	// UserVoteLike shows user's like
	UserVoteLike = types.UserVoteLike
)

type (
	// RequestID is used to store the request ID supplied in input Kafka records as
	// a unique identifier of payloads. Empty string represents a missing request ID.
	RequestID = types.RequestID
	// RuleOnReport represents a single (hit) rule of the string encoded report
	RuleOnReport = types.RuleOnReport
	// ReportRules is a helper struct for easy JSON unmarshalling of string encoded report
	ReportRules = types.ReportRules
	// ReportResponse represents the response of /report endpoint
	ReportResponse = types.ReportResponse
	// ReportResponseMeta contains metadata about the report
	ReportResponseMeta = types.ReportResponseMeta
	// DisabledRuleResponse represents a single disabled rule displaying only identifying information
	DisabledRuleResponse = types.DisabledRuleResponse
	// RuleID represents type for rule id
	RuleID = types.RuleID
	// RuleSelector represents type for rule id + error key in the rule_id|error_key format
	RuleSelector = types.RuleSelector
	// ErrorKey represents type for error key
	ErrorKey = types.ErrorKey
	// Rule represents the content of rule table
	Rule = types.Rule
	// RuleErrorKey represents the content of rule_error_key table
	RuleErrorKey = types.RuleErrorKey
	// RuleWithContent represents a rule with content, basically the mix of rule and rule_error_key tables' content
	RuleWithContent = types.RuleWithContent
	// Version represents the version of the cluster
	Version = types.Version
	// KafkaOffset type for kafka offset
	KafkaOffset = types.KafkaOffset
	// DBDriver type for db driver enum
	DBDriver = types.DBDriver
	// Identity contains internal user info
	Identity = types.Identity
	// Token is x-rh-identity struct
	Token = types.Token
	// User is single-user identifying struct
	User = types.User
	// RuleFQDN is rule module, older notation
	RuleFQDN = types.RuleFQDN
)

const (
	// DBDriverSQLite3 shows that db driver is sqlite
	DBDriverSQLite3 = types.DBDriverSQLite3
	// DBDriverPostgres shows that db driver is postgres
	DBDriverPostgres = types.DBDriverPostgres
	// DBDriverGeneral general sql(used for mock now)
	DBDriverGeneral = types.DBDriverGeneral
)

const (
	// ContextKeyUser is a constant for user authentication token in request
	ContextKeyUser = types.ContextKeyUser
)

// FeedbackRequest contains message of user feedback
type FeedbackRequest struct {
	Message string `json:"message"`
}

// ReportItem represents a single (hit) rule of the string encoded report
type ReportItem = types.ReportItem

// InfoItem represents a single (hit) info rule of the string encoded report
type InfoItem struct {
	InfoID  RuleID            `json:"info_id"`
	InfoKey InfoKey           `json:"key"`
	Details map[string]string `json:"details"`
}

// ClusterReports is a data structure containing list of clusters, list of
// errors and dictionary with results per cluster.
type ClusterReports = types.ClusterReports

// SchemaVersion represents the current version of data schema
type SchemaVersion = types.SchemaVersion

// RuleRating represents a rule rating element sent by the user
type RuleRating = types.RuleRating

// Metadata represents the metadata field in the Kafka topic messages
type Metadata struct {
	GatheredAt time.Time `json:"gathering_time"`
}
