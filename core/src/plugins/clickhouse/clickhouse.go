package clickhouse

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/clidey/whodb/core/src/engine"
	"github.com/clidey/whodb/core/src/plugins/common"
	"gorm.io/gorm"
)

type ClickHousePlugin struct{}

func (p *ClickHousePlugin) IsAvailable(config *engine.PluginConfig) bool {
	db, err := DB(config)
	if err != nil {
		return false
	}
	sqlDb, err := db.DB()
	if err != nil {
		return false
	}
	sqlDb.Close()
	return true
}

func (p *ClickHousePlugin) GetDatabases() ([]string, error) {
	return nil, errors.New("unsupported")
}

func (p *ClickHousePlugin) GetSchema(config *engine.PluginConfig) ([]string, error) {
	db, err := DB(config)
	if err != nil {
		return nil, err
	}
	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer sqlDb.Close()
	var schemas []struct {
		SchemaName string `gorm:"column:database"`
	}
	if err := db.Raw("SELECT name AS database FROM system.databases").Scan(&schemas).Error; err != nil {
		return nil, err
	}
	schemaNames := []string{}
	for _, schema := range schemas {
		schemaNames = append(schemaNames, schema.SchemaName)
	}
	return schemaNames, nil
}

func (p *ClickHousePlugin) GetStorageUnits(config *engine.PluginConfig, schema string) ([]engine.StorageUnit, error) {
	db, err := DB(config)
	if err != nil {
		return nil, err
	}
	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer sqlDb.Close()
	storageUnits := []engine.StorageUnit{}
	rows, err := db.Raw(fmt.Sprintf(`
		SELECT
			table AS table_name,
			database AS table_schema,
			total_bytes AS total_size,
			data_bytes AS data_size,
			rows AS row_count
		FROM system.parts
		WHERE database = '%v'
		GROUP BY table, database
	`, schema)).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	allTablesWithColumns, err := getTableSchema(db, schema)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tableName, tableSchema string
		var totalSize, dataSize int64
		var rowCount int64
		if err := rows.Scan(&tableName, &tableSchema, &totalSize, &dataSize, &rowCount); err != nil {
			log.Fatal(err)
		}

		rowCountRecordValue := "unknown"
		if rowCount >= 0 {
			rowCountRecordValue = fmt.Sprintf("%d", rowCount)
		}

		attributes := []engine.Record{
			{Key: "Table Type", Value: "Table"},
			{Key: "Total Size", Value: fmt.Sprintf("%d", totalSize)},
			{Key: "Data Size", Value: fmt.Sprintf("%d", dataSize)},
			{Key: "Count", Value: rowCountRecordValue},
		}

		attributes = append(attributes, allTablesWithColumns[tableName]...)

		storageUnits = append(storageUnits, engine.StorageUnit{
			Name:       tableName,
			Attributes: attributes,
		})
	}
	return storageUnits, nil
}

func getTableSchema(db *gorm.DB, schema string) (map[string][]engine.Record, error) {
	var result []struct {
		TableName  string `gorm:"column:table"`
		ColumnName string `gorm:"column:name"`
		DataType   string `gorm:"column:type"`
	}

	query := fmt.Sprintf(`
		SELECT table, name, type
		FROM system.columns
		WHERE database = '%v'
		ORDER BY table, position
	`, schema)

	if err := db.Raw(query).Scan(&result).Error; err != nil {
		return nil, err
	}

	tableColumnsMap := make(map[string][]engine.Record)
	for _, row := range result {
		tableColumnsMap[row.TableName] = append(tableColumnsMap[row.TableName], engine.Record{Key: row.ColumnName, Value: row.DataType})
	}

	return tableColumnsMap, nil
}

func (p *ClickHousePlugin) GetRows(config *engine.PluginConfig, schema string, storageUnit string, where string, pageSize int, pageOffset int) (*engine.GetRowsResult, error) {
	if !common.IsValidSQLTableName(storageUnit) {
		return nil, errors.New("invalid table name")
	}

	query := fmt.Sprintf("SELECT * FROM `%v`.`%s`", schema, storageUnit)
	if len(where) > 0 {
		query = fmt.Sprintf("%v WHERE %v", query, where)
	}
	query = fmt.Sprintf("%v LIMIT %d OFFSET %d", query, pageSize, pageOffset)
	return p.executeRawSQL(config, query)
}

func (p *ClickHousePlugin) executeRawSQL(config *engine.PluginConfig, query string, params ...interface{}) (*engine.GetRowsResult, error) {
	db, err := DB(config)
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer sqlDb.Close()

	// Execute the query
	rows, err := db.Raw(query, params...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get column information
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	// Prepare result structure
	result := &engine.GetRowsResult{}

	// Populate column information
	for _, col := range columns {
		for _, colType := range columnTypes {
			if col == colType.Name() {
				result.Columns = append(result.Columns, engine.Column{Name: col, Type: colType.DatabaseTypeName()})
				break
			}
		}
	}

	// Process each row
	for rows.Next() {
		// Prepare pointers to scan into
		columnPointers := make([]interface{}, len(columns))
		row := make([]string, len(columns))

		for i := range columns {
			columnPointers[i] = new(sql.NullString)
		}

		// Scan the row into columnPointers
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		// Convert scanned values to strings
		for i, colPtr := range columnPointers {
			val := colPtr.(*sql.NullString)
			if val.Valid {
				row[i] = val.String
			} else {
				row[i] = "" // Handle NULL values gracefully
			}
		}

		// Append the row to result.Rows
		result.Rows = append(result.Rows, row)
	}

	return result, nil
}

func (p *ClickHousePlugin) RawExecute(config *engine.PluginConfig, query string) (*engine.GetRowsResult, error) {
	return p.executeRawSQL(config, query)
}

func NewClickHousePlugin() *engine.Plugin {
	return &engine.Plugin{
		Type:            engine.DatabaseType_ClickHouse,
		PluginFunctions: &ClickHousePlugin{},
	}
}
