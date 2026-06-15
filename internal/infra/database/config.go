/*
 * Copyright 2026 MuixStudio
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package database

import (
	"time"
)

// DBConfig database configuration
type DBConfig struct {
	// Connection information
	Host     string `mapstructure:"host"`      // Database host
	Port     int    `mapstructure:"port"`      // Database port
	User     string `mapstructure:"user"`      // Database user
	Password string `mapstructure:"password"`  // Database password
	DBName   string `mapstructure:"db_name"`   // Database name
	SSLMode  string `mapstructure:"ssl_mode"`  // SSL mode: disable, require, verify-ca, verify-full
	TimeZone string `mapstructure:"time_zone"` // Time zone

	// Connection pool configuration
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`     // Maximum idle connections
	MaxOpenConns    int           `mapstructure:"max_open_conns"`     // Maximum open connections
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`  // Maximum connection lifetime
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"` // Maximum connection idle time

	// GORM configuration
	ParameterizedQueries                     bool `mapstructure:"parameterized_queries"`                         // Use parameterized queries
	PrepareStmt                              bool `mapstructure:"prepare_stmt"`                                  // Cache prepared statements
	DisableForeignKeyConstraintWhenMigrating bool `mapstructure:"disable_foreign_key_constraint_when_migrating"` // Disable foreign key constraints when migrating

	// Log config
	Log           LogConfig     `mapstructure:"log"`            // Log configuration
	SlowThreshold time.Duration `mapstructure:"slow_threshold"` // Slow query threshold

	// Table name configuration
	TablePrefix   string `mapstructure:"table_prefix"`   // Table name prefix
	SingularTable bool   `mapstructure:"singular_table"` // Use singular table names
}

type LogConfig struct {
	Level      string // debug, info, warn, error
	Format     string // json, console
	OutputPath string // 输出文件路径，空则只输出到控制台
	MaxSize    int    // 单个日志文件最大 MB，默认 100
	MaxBackups int    // 保留旧日志文件数，默认 3
	MaxAge     int    // 保留天数，默认 7
	Compress   bool   // 是否压缩旧日志

}
