package datastore

import (
	context "context"

	grpcds "github.com/chernser/authtogo/grpcdatastore"
	grpc "google.golang.org/grpc"
)

// GRPCDataStore - structure defining grpc datastore
type GRPCDataStore struct {
	client grpcds.VolatileDataStoreClient
}

// Get returns fields for requested row
func (store *GRPCDataStore) Get(id string, fields []string) (map[string]string, bool) {
	respose, err := store.client.GetRow(context.Background(), &grpcds.GetRowRequest{Id: id, Fields: fields})

	if err != nil || respose.IsError {
		return nil, false
	}

	return respose.Row, true
}

// Insert records new record to store
func (store *GRPCDataStore) Insert(id string, values map[string]string) bool {
	respose, err := store.client.InsertRow(context.Background(), &grpcds.InsertRowRequest{Id: id, Row: values})
	if err != nil || respose.IsError {
		return false
	}

	return true
}

// Update records new values
func (store *GRPCDataStore) Update(id string, values map[string]string) bool {
	respose, err := store.client.PutRow(context.Background(), &grpcds.PutRowRequest{Id: id, Row: values})
	if err != nil || respose.IsError {
		return false
	}

	return true
}

// Delete removes record from datastore
func (store *GRPCDataStore) Delete(id string) bool {
	respose, err := store.client.DeleteRow(context.Background(), &grpcds.DeleteRowRequest{Id: id})
	if err != nil || respose.IsError {
		return false
	}

	return true
}

// CreateGRPCStorage creates an grpc storage connected to specified server
func CreateGRPCStorage(server string) (*GRPCDataStore, error) {
	connection, err := grpc.Dial(server)
	if err != nil {
		return nil, err
	}

	datastore := &GRPCDataStore{
		client: grpcds.NewVolatileDataStoreClient(connection),
	}

	return datastore, nil
}
