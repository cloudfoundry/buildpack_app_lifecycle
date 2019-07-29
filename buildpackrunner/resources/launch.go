package resources

import (
	"code.cloudfoundry.org/buildpackapplifecycle"
)

type LaunchData struct {
	Processes []Process `yaml:"processes"`
}

type Process struct {
	Type      string `yaml:"type" json:"type"`
	Command   string `yaml:"command" json:"command"`
	Platforms struct {
		Cloudfoundry struct {
			SidecarFor []string `yaml:"sidecar_for" json:"sidecar_for"`
		} `yaml:"cloudfoundry" json:"cloudfoundry"`
	} `yaml:"platforms" json:"platforms"`
}

func (p *Process) Replaceable(otherProc Process) bool {
	return p.Type == otherProc.Type
}

func ProcDataToProcesses(procData map[string]string) []Process {
	var result []Process
	for procType, procCommand := range procData {
		result = append(result, Process{
			Type:    procType,
			Command: procCommand,
		})
	}
	return result
}

func MergeProcesses(listA []Process, listB []Process) []Process {
	result := listA
	for _, procB := range listB {
		replaced := false
		for i, procA := range listA {
			if procA.Replaceable(procB) {
				result[i] = procB
				replaced = true
			}
		}
		if !replaced {
			result = append(result, procB)
		}
	}
	return result
}

func ConvertToResult(data LaunchData) buildpackapplifecycle.StagingResult {
	result := buildpackapplifecycle.StagingResult{}
	result.ProcessTypes = map[string]string{}
	for _, process := range data.Processes {
		result.ProcessList = append(result.ProcessList, buildpackapplifecycle.Process{
			Type:    process.Type,
			Command: process.Command,
		})

		result.ProcessTypes[process.Type] = process.Command

		sidecarTargets := process.Platforms.Cloudfoundry.SidecarFor
		if len(sidecarTargets) > 0 {
			result.Sidecars = append(result.Sidecars, buildpackapplifecycle.Sidecars{
				Name:         process.Type,
				ProcessTypes: sidecarTargets,
				Command:      process.Command,
			})
		}
	}
	return result
}
