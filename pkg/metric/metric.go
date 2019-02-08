package metric

import (
	"time"
	"fmt"
	"encoding/json"
	project "github.com/onuryartasi/scaler/pkg/types"
	"github.com/onuryartasi/scaler/pkg/request"
	"log"
)
type Metric struct {
	Read      time.Time `json:"read"`
	Preread   time.Time `json:"preread"`
	PidsStats struct {
		Current int `json:"current"`
	} `json:"pids_stats"`
	BlkioStats struct {
		IoServiceBytesRecursive []struct {
			Major int    `json:"major"`
			Minor int    `json:"minor"`
			Op    string `json:"op"`
			Value int    `json:"value"`
		} `json:"io_service_bytes_recursive"`
		IoServicedRecursive []struct {
			Major int    `json:"major"`
			Minor int    `json:"minor"`
			Op    string `json:"op"`
			Value int    `json:"value"`
		} `json:"io_serviced_recursive"`
		IoQueueRecursive       []interface{} `json:"io_queue_recursive"`
		IoServiceTimeRecursive []interface{} `json:"io_service_time_recursive"`
		IoWaitTimeRecursive    []interface{} `json:"io_wait_time_recursive"`
		IoMergedRecursive      []interface{} `json:"io_merged_recursive"`
		IoTimeRecursive        []interface{} `json:"io_time_recursive"`
		SectorsRecursive       []interface{} `json:"sectors_recursive"`
	} `json:"blkio_stats"`
	NumProcs     int `json:"num_procs"`
	StorageStats struct {
	} `json:"storage_stats"`
	CPUStats struct {
		CPUUsage struct {
			TotalUsage        int64   `json:"total_usage"`
			PercpuUsage       []int64 `json:"percpu_usage"`
			UsageInKernelmode int64   `json:"usage_in_kernelmode"`
			UsageInUsermode   int64   `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		SystemCPUUsage int64 `json:"system_cpu_usage"`
		OnlineCpus     int   `json:"online_cpus"`
		ThrottlingData struct {
			Periods          int `json:"periods"`
			ThrottledPeriods int `json:"throttled_periods"`
			ThrottledTime    int `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"cpu_stats"`
	PrecpuStats struct {
		CPUUsage struct {
			TotalUsage        int64   `json:"total_usage"`
			PercpuUsage       []int64 `json:"percpu_usage"`
			UsageInKernelmode int64   `json:"usage_in_kernelmode"`
			UsageInUsermode   int64   `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		SystemCPUUsage int64 `json:"system_cpu_usage"`
		OnlineCpus     int   `json:"online_cpus"`
		ThrottlingData struct {
			Periods          int `json:"periods"`
			ThrottledPeriods int `json:"throttled_periods"`
			ThrottledTime    int `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"precpu_stats"`
	MemoryStats struct {
		Usage    int `json:"usage"`
		MaxUsage int `json:"max_usage"`
		Stats    struct {
			ActiveAnon              int   `json:"active_anon"`
			ActiveFile              int   `json:"active_file"`
			Cache                   int   `json:"cache"`
			Dirty                   int   `json:"dirty"`
			HierarchicalMemoryLimit int64 `json:"hierarchical_memory_limit"`
			HierarchicalMemswLimit  int   `json:"hierarchical_memsw_limit"`
			InactiveAnon            int   `json:"inactive_anon"`
			InactiveFile            int   `json:"inactive_file"`
			MappedFile              int   `json:"mapped_file"`
			Pgfault                 int   `json:"pgfault"`
			Pgmajfault              int   `json:"pgmajfault"`
			Pgpgin                  int   `json:"pgpgin"`
			Pgpgout                 int   `json:"pgpgout"`
			Rss                     int   `json:"rss"`
			RssHuge                 int   `json:"rss_huge"`
			TotalActiveAnon         int   `json:"total_active_anon"`
			TotalActiveFile         int   `json:"total_active_file"`
			TotalCache              int   `json:"total_cache"`
			TotalDirty              int   `json:"total_dirty"`
			TotalInactiveAnon       int   `json:"total_inactive_anon"`
			TotalInactiveFile       int   `json:"total_inactive_file"`
			TotalMappedFile         int   `json:"total_mapped_file"`
			TotalPgfault            int   `json:"total_pgfault"`
			TotalPgmajfault         int   `json:"total_pgmajfault"`
			TotalPgpgin             int   `json:"total_pgpgin"`
			TotalPgpgout            int   `json:"total_pgpgout"`
			TotalRss                int   `json:"total_rss"`
			TotalRssHuge            int   `json:"total_rss_huge"`
			TotalUnevictable        int   `json:"total_unevictable"`
			TotalWriteback          int   `json:"total_writeback"`
			Unevictable             int   `json:"unevictable"`
			Writeback               int   `json:"writeback"`
		} `json:"stats"`
		Limit int64 `json:"limit"`
	} `json:"memory_stats"`
	Name     string `json:"name"`
	ID       string `json:"id"`
	Networks struct {
		Eth0 struct {
			RxBytes   int `json:"rx_bytes"`
			RxPackets int `json:"rx_packets"`
			RxErrors  int `json:"rx_errors"`
			RxDropped int `json:"rx_dropped"`
			TxBytes   int `json:"tx_bytes"`
			TxPackets int `json:"tx_packets"`
			TxErrors  int `json:"tx_errors"`
			TxDropped int `json:"tx_dropped"`
		} `json:"eth0"`
	} `json:"networks"`
}


/*func INITMetric(){
	projects := v1.GetProjects()
	for _,project := range *projects {
		go ProjectMetric(project)
	}

}*/

func ProjectMetric(project project.Project) (){
	//var sum int64

	channels := make([]chan Metric,len(project.Containers))
	for i,id := range project.Containers {
		channels[i] = make(chan Metric)
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

func ContainerMetric(stream chan Metric,id string){
	var metric = Metric{}
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