package types

import (
	"fmt"
	"math"
	"path/filepath"
	"sync"

	"github.com/ismilent/nuclei/v2/pkg/catalog/config"
	"github.com/rs/xid"
)

// Default resume file
const DefaultResumeFileName = "resume-%s.cfg"

func DefaultResumeFilePath() string {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return fmt.Sprintf(DefaultResumeFileName, xid.New().String())
	}
	resumeFile := filepath.Join(configDir, fmt.Sprintf(DefaultResumeFileName, xid.New().String()))
	return resumeFile
}

// ResumeCfg contains the scan progression
type ResumeCfg struct {
	sync.RWMutex
	ResumeFrom map[string]*ResumeInfo `json:"resumeFrom"`
	Current    map[string]*ResumeInfo `json:"-"`
}

type ResumeInfo struct {
	sync.RWMutex
	Completed bool                `json:"completed"`
	InFlight  map[uint32]struct{} `json:"inFlight"`
	SkipUnder uint32              `json:"-"`
	Repeat    map[uint32]struct{} `json:"-"`
	DoAbove   uint32              `json:"-"`
}

// NewResumeCfg creates a new scan progression structure
func NewResumeCfg() *ResumeCfg {
	return &ResumeCfg{
		ResumeFrom: make(map[string]*ResumeInfo),
		Current:    make(map[string]*ResumeInfo),
	}
}

// Clone the resume structure
func (resumeCfg *ResumeCfg) Clone() ResumeCfg {
	return ResumeCfg{
		ResumeFrom: resumeCfg.ResumeFrom,
		Current:    resumeCfg.Current,
	}
}

// Clone the resume structure
func (resumeCfg *ResumeCfg) Compile() {
	resumeCfg.Lock()
	defer resumeCfg.Unlock()

	for _, resumeInfo := range resumeCfg.ResumeFrom {
		if resumeInfo.Completed && len(resumeInfo.InFlight) > 0 {
			resumeInfo.InFlight = make(map[uint32]struct{})
		}
		min := uint32(math.MaxUint32)
		max := uint32(0)
		for index := range resumeInfo.InFlight {
			if index < min {
				min = index
			}
			if index > max {
				max = index
			}
		}
		// maybe redundant but ensures we track the indexes to be repeated
		resumeInfo.Repeat = map[uint32]struct{}{}
		for index := range resumeInfo.InFlight {
			resumeInfo.Repeat[index] = struct{}{}
		}
		resumeInfo.SkipUnder = min
		resumeInfo.DoAbove = max
	}
}
