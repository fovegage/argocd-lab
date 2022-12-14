// Copyright 2021 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package blockwriter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/pingcap/tidb-operator/tests/pkg/util"
)

const (
	defaultTableNum  = 16
	defaultRecordNum = 256
	defaultBatchNum  = 100
	defaultBatchSize = 1024

	createTableExpr = `
CREATE TABLE IF NOT EXISTS %s (
    id BIGINT NOT NULL AUTO_INCREMENT,
    raw_bytes BLOB NOT NULL,
    PRIMARY KEY (id)
);
`

	insertExpr = `
INSERT INTO %s (raw_bytes) VALUES %s;
`
)

type option func(bw *blockWriter)

func WithTableNum(tableNum int) option {
	return func(bw *blockWriter) {
		bw.tableNum = tableNum
	}
}

func WithRecordNum(recordNum int) option {
	return func(bw *blockWriter) {
		bw.recordNum = recordNum
	}
}

func WithGenTableName(genTableName func(nr int) string) option {
	return func(bw *blockWriter) {
		bw.genTableName = genTableName
	}
}

func WithBatchNum(batchNum int) option {
	return func(bw *blockWriter) {
		bw.batchNum = batchNum
	}
}

func WithBatchSize(batchSize int) option {
	return func(bw *blockWriter) {
		bw.batchSize = batchSize
	}
}

// BlockWriter write test data to a database.
type BlockWriter interface {
	Write(ctx context.Context, dsn string) error
}

type blockWriter struct {
	tableNum  int
	recordNum int
	batchNum  int
	batchSize int

	genTableName func(nr int) string
}

func New(opts ...option) BlockWriter {
	bw := &blockWriter{
		tableNum:     defaultTableNum,
		recordNum:    defaultRecordNum,
		batchNum:     defaultBatchNum,
		batchSize:    defaultBatchSize,
		genTableName: genTableNameDefault,
	}

	for _, opt := range opts {
		opt(bw)
	}

	return bw
}

func (bw *blockWriter) Write(ctx context.Context, dsn string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	for i := 0; i < bw.tableNum; i++ {
		expr := fmt.Sprintf(createTableExpr, bw.genTableName(i))
		if _, err := db.Exec(expr); err != nil {
			return err
		}
	}

	errors := []error{}
	wg := sync.WaitGroup{}
	for i := 0; i < bw.tableNum; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			conn, err := db.Conn(ctx)
			if err != nil {
				errors = append(errors, err)
				return
			}
			defer conn.Close()

			for k := 0; k < bw.recordNum; k++ {
				values := make([]string, bw.batchNum)
				for k := 0; k < bw.batchNum; k++ {
					blockData := util.RandString(bw.batchSize)
					values[k] = fmt.Sprintf("('%s')", blockData)
				}
				expr := fmt.Sprintf(insertExpr, bw.genTableName(index), strings.Join(values, ","))

				if _, err := conn.ExecContext(ctx, expr); err != nil {
					errors = append(errors, err)
					break
				}
			}
		}(i)
	}
	wg.Wait()
	if len(errors) != 0 {
		return fmt.Errorf("write errors: %v", errors)
	}
	return nil
}

func genTableNameDefault(nr int) string {
	return fmt.Sprintf("block_writer%d", nr)
}
