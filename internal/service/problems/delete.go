package service_problems

import "context"

func (s *ProblemsService) Delete(ctx context.Context, id string) error {
	return s.problemsRepository.Delete(ctx, id)
}
