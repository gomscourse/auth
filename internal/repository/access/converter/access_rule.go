package converter

import (
	"github.com/gomscourse/auth/internal/model"
	modelRepo "github.com/gomscourse/auth/internal/repository/access/model"
)

func ToAccessRuleFromRepo(rule *modelRepo.AccessRule) *model.AccessRule {
	return &model.AccessRule{
		Role:     rule.Role,
		Endpoint: rule.Endpoint,
	}
}
