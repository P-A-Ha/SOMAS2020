package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/pkg/errors"
)

// runTurn runs a turn
func (s *SOMASServer) runTurn() error {
	s.logf("start runTurn")
	defer s.logf("finish runTurn")

	s.logf("TURN: %v, Season: %v", s.gameState.Turn, s.gameState.Season)

	s.startOfTurn()

	// run all orgs
	err := s.runOrgs()
	if err != nil {
		return errors.Errorf("Error running orgs: %v", err)
	}

	if err := s.endOfTurn(); err != nil {
		return errors.Errorf("Error running end of turn procedures: %v", err)
	}

	return nil
}

func (s *SOMASServer) startOfTurn() {
	s.logf("start startOfTurn")
	defer s.logf("finish startOfTurn")
	for _, clientID := range getNonDeadClientIDs(s.gameState.ClientInfos) {
		s.clientMap[clientID].StartOfTurn()
	}
}

// runOrgs runs all the orgs
func (s *SOMASServer) runOrgs() error {
	s.logf("start runOrgs")
	defer s.logf("finish runOrgs")

	if err := s.runIIGO(); err != nil {
		return errors.Errorf("IIGO error: %v", err)
	}

	if err := s.runIIFO(); err != nil {
		return errors.Errorf("IIFO error: %v", err)
	}

	if err := s.runIITO(); err != nil {
		return errors.Errorf("IITO error: %v", err)
	}

	return nil
}

// endOfTurn performs end of turn updates
// TODO: organise order of end of turn actions
func (s *SOMASServer) endOfTurn() error {
	s.logf("start endOfTurn")
	defer s.logf("finish endOfTurn")

	if err := s.runIIGOAllocations(); err != nil {
		return errors.Errorf("Failed to get common pool allocations at end of turn: %v", err)
	}

	// TODO : break foraging down into foraging investments and foraging returns
	if err := s.runForage(); err != nil {
		return errors.Errorf("Failed to run hunt at end of turn: %v", err)
	}

	if err := s.runIIFOEndOfTurn(); err != nil {
		return errors.Errorf("IIFO EndOfTurn error: %v", err)
	}

	// TODO: break IITO down into giving gifts and receiving gifts
	if err := s.runIITOEndOfTurn(); err != nil {
		return errors.Errorf("IITO EndOfTurn error: %v", err)
	}

	if err := s.runIIGOTax(); err != nil {
		return errors.Errorf("Failed to put taxes into common pool at end of turn: %v", err)
	}

	// probe for disaster
	updatedEnv, err := s.probeDisaster()
	if err != nil {
		return errors.Errorf("Failed to probe disaster: %v", err)
	}
	s.gameState.Environment = updatedEnv
	// increment turn & season if needed
	disasterHappened := updatedEnv.LastDisasterReport.Magnitude > 0
	s.incrementTurnAndSeason(disasterHappened)

	if disasterHappened {
		s.notifyClientsOfDisaster() // sends disaster report and effects to all non-dead clients
	}

	// deduct cost of living
	s.deductCostOfLiving(s.gameConfig.CostOfLiving)

	err = s.updateIslandLivingStatus()
	if err != nil {
		return errors.Errorf("Failed to update island living status: %v", err)
	}

	return nil
}

// incrementTurnAndSeason increments turn, and season if a disaster happened.
func (s *SOMASServer) incrementTurnAndSeason(disasterHappened bool) {
	s.logf("start incrementTurnAndSeason")
	defer s.logf("finish incrementTurnAndSeason")

	s.gameState.Turn++
	if disasterHappened {
		s.gameState.Season++
	}
}

func (s *SOMASServer) notifyClientsOfDisaster() {
	s.logf("start notifying clients of disaster")
	defer s.logf("finish notifying clients of disaster")

	nonDeadClients := getNonDeadClientIDs(s.gameState.ClientInfos)
	for _, id := range nonDeadClients {
		c := s.clientMap[id]
		effects := s.gameState.Environment.DisasterEffects() // gets effects of most recent disaster
		c.DisasterNotification(s.gameState.Environment.LastDisasterReport, effects)
	}
}

// deductCostOfLiving deducts CoL for all living islands, including critical ones
func (s *SOMASServer) deductCostOfLiving(costOfLiving shared.Resources) {
	s.logf("start deductCostOfLiving")
	defer s.logf("finish deductCostOfLiving")

	nonDeadClients := getNonDeadClientIDs(s.gameState.ClientInfos)
	for _, id := range nonDeadClients {
		ci := s.gameState.ClientInfos[id]
		if ci.Resources < costOfLiving {
			ci.Resources = 0
		} else {
			ci.Resources -= costOfLiving
		}
		s.gameState.ClientInfos[id] = ci
	}
}

// updateIslandLivingStatus changes the islands Alive and Critical state depending
// on the island's resource state.
// Dead islands are not resurrected.
func (s *SOMASServer) updateIslandLivingStatus() error {
	s.logf("start updateIslandLivingStatus")
	defer s.logf("finish updateIslandLivingStatus")

	nonDeadClients := getNonDeadClientIDs(s.gameState.ClientInfos)
	for _, id := range nonDeadClients {
		ci, err := updateIslandLivingStatusForClient(s.gameState.ClientInfos[id],
			s.gameConfig.MinimumResourceThreshold, s.gameConfig.MaxCriticalConsecutiveTurns)
		if err != nil {
			return errors.Errorf("Failed to update island living status for '%v': %v", id, err)
		}
		s.gameState.ClientInfos[id] = ci
	}

	return nil
}

func (s *SOMASServer) gameOver(maxTurns uint, maxSeasons uint) bool {
	st := s.gameState

	if !anyClientsAlive(st.ClientInfos) {
		s.logf("All clients are dead!")
		return true
	}

	// +1 due to 1-indexing
	if st.Turn >= maxTurns+1 {
		s.logf("Max turns '%v' reached or exceeded", maxTurns)
		return true
	}

	// +1 due to 1-indexing
	if st.Season >= maxSeasons+1 {
		s.logf("Max seasons '%v' reached or exceeded", maxSeasons)
		return true
	}

	return false
}
