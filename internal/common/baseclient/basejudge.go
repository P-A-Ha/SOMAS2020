package baseclient

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type BaseJudge struct {
}

// PayPresident pays the President a salary.
func (j *BaseJudge) PayPresident(presidentSalary shared.Resources) (shared.Resources, bool) {
	// TODO Implement opinion based salary payment.
	return presidentSalary, true
}

// inspectHistoryInternal is the base implementation of InspectHistory.
func (j *BaseJudge) InspectHistory(iigoHistory []shared.Accountability) (map[shared.ClientID]roles.EvaluationReturn, bool) {
	outputMap := map[shared.ClientID]roles.EvaluationReturn{}
	for _, entry := range iigoHistory {
		variablePairs := entry.Pairs
		clientID := entry.ClientID
		var rulesAffected []string
		for _, variable := range variablePairs {
			valuesToBeAdded, foundRules := rules.PickUpRulesByVariable(variable.VariableName, rules.RulesInPlay)
			if foundRules {
				rulesAffected = append(rulesAffected, valuesToBeAdded...)
			}
			updatedVariable := rules.UpdateVariable(variable.VariableName, variable)
			if !updatedVariable {
				return map[shared.ClientID]roles.EvaluationReturn{}, false
			}
		}
		if _, ok := outputMap[clientID]; !ok {
			outputMap[clientID] = roles.EvaluationReturn{
				Rules:       []rules.RuleMatrix{},
				Evaluations: []bool{},
			}
		}
		tempReturn := outputMap[clientID]
		for _, rule := range rulesAffected {
			evaluation, err := rules.BasicBooleanRuleEvaluator(rule)
			if err != nil {
				return outputMap, false
			}
			tempReturn.Rules = append(tempReturn.Rules, rules.RulesInPlay[rule])
			tempReturn.Evaluations = append(tempReturn.Evaluations, evaluation)
		}
		outputMap[clientID] = tempReturn
	}
	return outputMap, true
}

// CallPresidentElection is called by the judiciary to decide on power-transfer
func (j *BaseJudge) CallPresidentElection(turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.Plurality,
		IslandsToVote: allIslands,
		HoldElection:  true,
	}
	return electionsettings
}

// DecideNextPresident returns the ID of chosen next President
func (j *BaseJudge) DecideNextPresident(winner shared.ClientID) shared.ClientID {
	return winner
}
