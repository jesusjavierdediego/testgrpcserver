package grpc

import (
	pb "xqledger/grpcserver/protobuf"

	"golang.org/x/net/context"
)


type RecordHistoryService struct {
	query *pb.Query
}

func NewRecordHistoryService(query *pb.Query) *RecordHistoryService {
	return &RecordHistoryService{query: query}
}


func (s *RecordHistoryService) GetRecordHistory(ctx context.Context, query *pb.Query) (*pb.RecordHistory, error) {
	var result pb.RecordHistory
	var commit pb.Commit
	commit.AuthorEmail = "TestOrchetstrator@gmail.com"
	commit.AuthorName = "TestOrchetstrator"
	commit.AuthorWhen = "01/12/1987"
	commit.CommitterEmail = ""
	commit.CommitterName = "TestOrchetstrator"
	commit.CommitterWhen = "TestOrchetstrator@gmail.com"
	commit.Hash = "123456789"
	commit.Message = "1234"
	commit.Parents = 1
	var list []*pb.Commit
	list = append(list, &commit)
	result.Commits = list

	return &result, nil
}

func (s *RecordHistoryService) GetContentInCommit(ctx context.Context, query *pb.Query) (*pb.CommitContent, error) {
	var result pb.CommitContent
	result.Content = "hello"
	return &result, nil
}

func (s *RecordHistoryService) GetDiffTwoCommitsInFile(ctx context.Context, query *pb.Query) (*pb.DiffHtml, error) {
	var result pb.DiffHtml
	result.Html = "<head>hello</head>"
	return &result, nil
}
