package owl

import (
	"context"
	"errors"
)

func (s *Service) DeleteConsumerGroup(ctx context.Context, groupID string) (string, error) {
	commitResponse, err := s.kafkaSvc.DeleteGroupsRequest(ctx, groupID)
	if err != nil {
		return "nil", err
	}
	if len(commitResponse.Groups) > 0 {
		return commitResponse.Groups[0].Group, nil
	}
	return "", errors.New("delete group empty...")
}
