package metric

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/onuryartasi/scaler/pkg/request"
	p2 "github.com/onuryartasi/scaler/pkg/types"
)

/*func INITMetric(){
	projects := v1.GetProjects()
	for _,project := range *projects {
		go ProjectMetric(project)
	}

}*/

func ProjectMetric(project p2.Project) {
	//var sum int64

	channels := make([]chan p2.Metric, len(project.Containers))
	for i, id := range project.Containers {
		channels[i] = make(chan p2.Metric)
		go ContainerMetric(channels[i], id)
	}
	for {
		//time.Sleep(time.Second * 3)
		var sum float64 = 0.0
		for _, value := range channels {
			met := <-value

			CpuPercent := calculateCPUPercentUnix(project.Cpu, met.PrecpuStats.ThrottlingData.ThrottledTime, met.CPUStats.ThrottlingData.ThrottledTime, met.CPUStats.CPUUsage.TotalUsage, met.CPUStats.SystemCPUUsage, met.PrecpuStats.CPUUsage.TotalUsage, met.PrecpuStats.SystemCPUUsage, met.CPUStats.OnlineCpus)
			log.Println(CpuPercent)
			sum += CpuPercent
		}

	}
}

func ContainerMetric(stream chan p2.Metric, id string) {
	var metric = p2.Metric{}
	client := request.NewClient()

	response, err := client.Get(fmt.Sprintf("http://v1.37/containers/%s/stats", id))
	if err != nil {
		panic(err)
	}
	for {
		json.NewDecoder(response.Body).Decode(&metric)
		stream <- metric
	}

}

func calculateCPUPercentUnix(cpuLimit float32, PreThrottledTime int, ThrottledTime int, TotalCpu int64, TotalSystem int64, previousCPU, previousSystem int64, cpu int) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(TotalCpu) - float64(previousCPU)
		// calculate the change for the entire system between readings
		systemDelta = float64(TotalSystem) - float64(previousSystem)
	)
	var Quota float64 = float64(cpuLimit) * 1000000000
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		//log.Println("->", cpuDelta, "->", ThrottledTimeDelta, "Divedi->", (cpuDelta/ThrottledTimeDelta)*float64(cpu)*100, "-", cpu)
		//cpuPercent = (cpuDelta / systemDelta) * float64(cpu) * 100
		//cpuPercent = (cpuDelta / ThrottledTimeDelta) * float64(cpu) * 100
		cpuPercent = (cpuDelta / Quota) * 100.0
	}
	return cpuPercent
}
