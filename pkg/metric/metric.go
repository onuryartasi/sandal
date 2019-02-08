package metric

import (
	"fmt"
	"encoding/json"
	p2 "github.com/onuryartasi/scaler/pkg/types"
	"github.com/onuryartasi/scaler/pkg/request"
	"log"
)


/*func INITMetric(){
	projects := v1.GetProjects()
	for _,project := range *projects {
		go ProjectMetric(project)
	}

}*/

func ProjectMetric(project p2.Project) (){
	//var sum int64

	channels := make([]chan p2.Metric,len(project.Containers))
	for i,id := range project.Containers {
		channels[i] = make(chan p2.Metric)
		go ContainerMetric(channels[i],id)
	}
	for {
		//time.Sleep(time.Second * 3)
		var sum float64 = 0.0
		for i, value := range channels {
			met := <-value
			CpuPercent := calculateCPUPercentUnix(met.CPUStats.CPUUsage.TotalUsage,met.CPUStats.SystemCPUUsage,met.PrecpuStats.CPUUsage.TotalUsage,met.PrecpuStats.SystemCPUUsage,len(met.CPUStats.CPUUsage.PercpuUsage))
			sum+= CpuPercent
			log.Printf("%v.Container Cpu : %v",i,CpuPercent)
			}
		log.Printf("Sum: %v",sum/float64(len(channels)))

	}
}

func ContainerMetric(stream chan p2.Metric,id string){
	var metric = p2.Metric{}
	client := request.NewClient()

		response, err := client.Get(fmt.Sprintf("http://v1.28/containers/%s/stats",id))
		if err != nil {
			panic(err)
		}
	for{
		json.NewDecoder(response.Body).Decode(&metric)
		stream <- metric
	}

}

func calculateCPUPercentUnix(TotalCpu int64,TotalSystem int64,previousCPU, previousSystem int64,cpucount int) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(TotalCpu) - float64(previousCPU)
		// calculate the change for the entire system between readings
		systemDelta = float64(TotalSystem) - float64(previousSystem)
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(cpucount) * 100.0
	}
	return cpuPercent
}