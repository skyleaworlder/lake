/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package helper

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/incubator-devlake/plugins/core"
	"gorm.io/datatypes"
)

// Table structure for raw data storage
type RawData struct {
	ID        uint64 `gorm:"primaryKey"`
	Params    string `gorm:"type:varchar(255);index"`
	Data      []byte
	Url       string
	Input     datatypes.JSON
	CreatedAt time.Time
}

type RawDataSubTaskArgs struct {
	Ctx core.SubTaskContext

	//	Table store raw data
	Table string `comment:"Raw data table name"`

	//	This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
	//	set of data to be process, for example, we process JiraIssues by Board
	Params interface{} `comment:"To identify a set of records with same UrlTemplate, i.e. {ConnectionId, BoardId} for jira entities"`
}

// Common features for raw data sub tasks
type RawDataSubTask struct {
	args   *RawDataSubTaskArgs
	table  string
	params string
}

func newRawDataSubTask(args RawDataSubTaskArgs) (*RawDataSubTask, error) {
	if args.Ctx == nil {
		return nil, fmt.Errorf("Ctx is required for RawDataSubTask")
	}
	if args.Table == "" {
		return nil, fmt.Errorf("Table is required for RawDataSubTask")
	}
	paramsString := ""
	if args.Params == nil {
		args.Ctx.GetLogger().Warn("Missing `Params` for raw data subtask %s", args.Ctx.GetName())
	} else {
		// TODO: maybe sort it to make it consisitence
		paramsBytes, err := json.Marshal(args.Params)
		if err != nil {
			return nil, err
		}
		paramsString = string(paramsBytes)
	}
	return &RawDataSubTask{
		args:   &args,
		table:  fmt.Sprintf("_raw_%s", args.Table),
		params: paramsString,
	}, nil
}
