package main

import (
	"time"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/pkg/gitinfo"
)

type runInfo struct {
	TimeStart       time.Time
	TimeEnd         time.Time
	DurationSeconds float64
	Version         string
	GOOS            string
	GOARCH          string
}

// output represents what is output into the output.json file
type output struct {
	Config     config.Config
	GitInfo    gitinfo.GitInfo
	RunInfo    runInfo
	GameStates []gamestate.GameState
}
