// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type AIChatMessage struct {
	Type   string      `json:"Type"`
	Result *RowsResult `json:"Result,omitempty"`
	Text   string      `json:"Text"`
}

type ChatInput struct {
	PreviousConversation string  `json:"PreviousConversation"`
	Query                string  `json:"Query"`
	Model                string  `json:"Model"`
	Token                *string `json:"Token,omitempty"`
}

type Column struct {
	Type string `json:"Type"`
	Name string `json:"Name"`
}

type GraphUnit struct {
	Unit      *StorageUnit             `json:"Unit"`
	Relations []*GraphUnitRelationship `json:"Relations"`
}

type GraphUnitRelationship struct {
	Name         string                    `json:"Name"`
	Relationship GraphUnitRelationshipType `json:"Relationship"`
}

type LoginCredentials struct {
	ID       *string        `json:"Id,omitempty"`
	Type     string         `json:"Type"`
	Hostname string         `json:"Hostname"`
	Username string         `json:"Username"`
	Password string         `json:"Password"`
	Database string         `json:"Database"`
	Advanced []*RecordInput `json:"Advanced,omitempty"`
}

type LoginProfile struct {
	Alias    *string      `json:"Alias,omitempty"`
	ID       string       `json:"Id"`
	Type     DatabaseType `json:"Type"`
	Database *string      `json:"Database,omitempty"`
	Source   string       `json:"Source"`
}

type LoginProfileInput struct {
	ID       string       `json:"Id"`
	Type     DatabaseType `json:"Type"`
	Database *string      `json:"Database,omitempty"`
}

type Mutation struct {
}

type Query struct {
}

type Record struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type RecordInput struct {
	Key   string         `json:"Key"`
	Value string         `json:"Value"`
	Extra []*RecordInput `json:"Extra,omitempty"`
}

type RowsResult struct {
	Columns       []*Column  `json:"Columns"`
	Rows          [][]string `json:"Rows"`
	DisableUpdate bool       `json:"DisableUpdate"`
}

type SettingsConfig struct {
	MetricsEnabled *bool `json:"MetricsEnabled,omitempty"`
}

type SettingsConfigInput struct {
	MetricsEnabled *string `json:"MetricsEnabled,omitempty"`
}

type StatusResponse struct {
	Status bool `json:"Status"`
}

type StorageUnit struct {
	Name       string    `json:"Name"`
	Attributes []*Record `json:"Attributes"`
}

type DatabaseType string

const (
	DatabaseTypePostgres      DatabaseType = "Postgres"
	DatabaseTypeMySQL         DatabaseType = "MySQL"
	DatabaseTypeSqlite3       DatabaseType = "Sqlite3"
	DatabaseTypeMongoDb       DatabaseType = "MongoDB"
	DatabaseTypeRedis         DatabaseType = "Redis"
	DatabaseTypeElasticSearch DatabaseType = "ElasticSearch"
	DatabaseTypeMariaDb       DatabaseType = "MariaDB"
	DatabaseTypeClickHouse    DatabaseType = "ClickHouse"
)

var AllDatabaseType = []DatabaseType{
	DatabaseTypePostgres,
	DatabaseTypeMySQL,
	DatabaseTypeSqlite3,
	DatabaseTypeMongoDb,
	DatabaseTypeRedis,
	DatabaseTypeElasticSearch,
	DatabaseTypeMariaDb,
	DatabaseTypeClickHouse,
}

func (e DatabaseType) IsValid() bool {
	switch e {
	case DatabaseTypePostgres, DatabaseTypeMySQL, DatabaseTypeSqlite3, DatabaseTypeMongoDb, DatabaseTypeRedis, DatabaseTypeElasticSearch, DatabaseTypeMariaDb, DatabaseTypeClickHouse:
		return true
	}
	return false
}

func (e DatabaseType) String() string {
	return string(e)
}

func (e *DatabaseType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DatabaseType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DatabaseType", str)
	}
	return nil
}

func (e DatabaseType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type GraphUnitRelationshipType string

const (
	GraphUnitRelationshipTypeOneToOne   GraphUnitRelationshipType = "OneToOne"
	GraphUnitRelationshipTypeOneToMany  GraphUnitRelationshipType = "OneToMany"
	GraphUnitRelationshipTypeManyToOne  GraphUnitRelationshipType = "ManyToOne"
	GraphUnitRelationshipTypeManyToMany GraphUnitRelationshipType = "ManyToMany"
	GraphUnitRelationshipTypeUnknown    GraphUnitRelationshipType = "Unknown"
)

var AllGraphUnitRelationshipType = []GraphUnitRelationshipType{
	GraphUnitRelationshipTypeOneToOne,
	GraphUnitRelationshipTypeOneToMany,
	GraphUnitRelationshipTypeManyToOne,
	GraphUnitRelationshipTypeManyToMany,
	GraphUnitRelationshipTypeUnknown,
}

func (e GraphUnitRelationshipType) IsValid() bool {
	switch e {
	case GraphUnitRelationshipTypeOneToOne, GraphUnitRelationshipTypeOneToMany, GraphUnitRelationshipTypeManyToOne, GraphUnitRelationshipTypeManyToMany, GraphUnitRelationshipTypeUnknown:
		return true
	}
	return false
}

func (e GraphUnitRelationshipType) String() string {
	return string(e)
}

func (e *GraphUnitRelationshipType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = GraphUnitRelationshipType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid GraphUnitRelationshipType", str)
	}
	return nil
}

func (e GraphUnitRelationshipType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
