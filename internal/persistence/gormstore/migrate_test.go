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

package gormstore

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/muixstudio/clio/internal/infra/database"
)

func TestMigration(t *testing.T) {
	//if os.Getenv("AUTH_MIGRATION_INTEGRATION_TEST") != "1" {
	//	t.Skip("set AUTH_MIGRATION_INTEGRATION_TEST=1 to run PostgreSQL migration integration test")
	//}

	db := database.MustNewDB(database.DBConfig{
		Host:                                     envOrDefault("AUTH_TEST_DB_HOST", "infra.codeclub.com.cn"),
		Port:                                     envIntOrDefault(t, "AUTH_TEST_DB_PORT", 5432),
		User:                                     envOrDefault("AUTH_TEST_DB_USER", "postgres"),
		Password:                                 envOrDefault("AUTH_TEST_DB_PASSWORD", "AX2F7nV!wNDrr26wKuqn"),
		DBName:                                   envOrDefault("AUTH_TEST_DB_NAME", "clio"),
		SSLMode:                                  envOrDefault("AUTH_TEST_DB_SSLMODE", "disable"),
		TimeZone:                                 envOrDefault("AUTH_TEST_DB_TIMEZONE", "Asia/Shanghai"),
		MaxIdleConns:                             1,
		MaxOpenConns:                             1,
		ConnMaxLifetime:                          time.Minute,
		ConnMaxIdleTime:                          time.Minute,
		DisableForeignKeyConstraintWhenMigrating: false,
	})

	if err := migrate(db); err != nil {
		t.Fatalf("migration failed: %v", err)
	}
}

func envOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func envIntOrDefault(t *testing.T, key string, defaultValue int) int {
	t.Helper()

	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		t.Fatalf("invalid %s: %v", key, err)
	}
	return parsed
}
