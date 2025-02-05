// Copyright Contributors to the Open Cluster Management project

package database

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stolostron/search-indexer/pkg/model"
)

func Test_ResyncData(t *testing.T) {
	// Prepare a mock DAO instance.
	dao, mockPool := buildMockDAO(t)

	// Mock PosgreSQL apis
	mockPool.EXPECT().Exec(gomock.Any(), gomock.Eq(`DELETE FROM "search"."resources" WHERE (("cluster" = 'test-cluster') AND ("uid" != 'cluster__test-cluster'))`), gomock.Eq([]interface{}{})).Return(nil, nil)
	mockPool.EXPECT().Exec(gomock.Any(), gomock.Eq(`DELETE FROM "search"."edges" WHERE ("cluster" = 'test-cluster')`), gomock.Eq([]interface{}{})).Return(nil, nil)
	br := BatchResults{}
	mockPool.EXPECT().SendBatch(gomock.Any(), gomock.Any()).Return(br)

	// Prepare Request data.
	data, _ := os.Open("./mocks/simple.json")
	var syncEvent model.SyncEvent
	json.NewDecoder(data).Decode(&syncEvent) //nolint: errcheck

	// Supress console output to prevent log messages from polluting test output.
	defer SupressConsoleOutput()()

	// Execute function test.
	response := &model.SyncResponse{}
	dao.ResyncData(syncEvent, "test-cluster", response)
}

func Test_ResyncData_errors(t *testing.T) {
	// Prepare a mock DAO instance.
	dao, mockPool := buildMockDAO(t)

	// Mock PosgreSQL apis
	mockPool.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Eq([]interface{}{})).Return(nil, errors.New("Delete error")).Times(2)
	br := BatchResults{}
	mockPool.EXPECT().SendBatch(gomock.Any(), gomock.Any()).Return(br)

	// Prepare Request data.
	data, _ := os.Open("./mocks/simple.json")
	var syncEvent model.SyncEvent
	json.NewDecoder(data).Decode(&syncEvent) //nolint: errcheck

	// Supress console output to prevent log messages from polluting test output.
	defer SupressConsoleOutput()()

	// Execute function test.
	response := &model.SyncResponse{}
	dao.ResyncData(syncEvent, "test-cluster", response)
}
