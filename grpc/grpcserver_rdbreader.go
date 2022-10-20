package grpc

import (
	pb "xqledger/grpcserver/protobuf"
	"golang.org/x/net/context"
	"encoding/json"
)


type RecordQueryService struct {
	query *pb.RDBQuery
}

func NewRecordQueryService(query *pb.RDBQuery) *RecordQueryService {
	return &RecordQueryService{query: query}
}


type Sample struct {
	Field1   string `json:"field1"`
	Field2   string `json:"field2"`
}

func getArraySamples() []string{
	var result []string
	var a1 Sample
	a1.Field1 = "Ranjan"
	a1.Field2 = "Vahit"
	a, _ := json.Marshal(a1)
	result = append(result, string(a))

	var a2 Sample
	a2.Field1 = "Fred"
	a2.Field2 = "Vargas"
	b, _ := json.Marshal(a2)
	result = append(result, string(b))

	var a3 Sample
	a3.Field1 = "Roco"
	a3.Field2 = "Sifredi"
	c, _ := json.Marshal(a3)
	result = append(result, string(c))

	return result

}


func (s *RecordQueryService) GetRDBRecords(ctx context.Context, query *pb.RDBQuery) (*pb.RecordSet, error) {
	var result pb.RecordSet
	result.Records = getArraySamples()

	return &result, nil
}


func (s *RecordQueryService) GetNumberRecordsFromColl(ctx context.Context, query *pb.RDBQuery) (*pb.RDCColCount, error) {
	var result pb.RDCColCount
	result.Count = 60

	return &result, nil
}