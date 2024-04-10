package task

import "context"

func (s *service) DeleteData(ctx context.Context, userId string) error {
	return s.queries.DeleteUserTaskData(ctx, userId)
}
