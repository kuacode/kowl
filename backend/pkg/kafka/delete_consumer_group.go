package kafka

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kmsg"
)

func (s *Service) DeleteGroupsRequest(ctx context.Context, groupID string) (*kmsg.DeleteGroupsResponse, error) {
	req := kmsg.NewDeleteGroupsRequest()
	req.Groups = []string{groupID}

	res, err := req.RequestWith(ctx, s.KafkaClient)
	if err != nil {
		return nil, fmt.Errorf("failed to delete group -> '%s': %w", groupID, err)
	}

	return res, nil
}
