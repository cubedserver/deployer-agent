package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"

	"github.com/shirou/gopsutil/mem"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var actionCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Collect server metrics",
	Long:  `Collects the latest server monitoring metrics and sends them to your Deployer account.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Read in the existing accounts and get ready for adding another user
		viper.ReadInConfig()

		var accounts []Account
		configErr := viper.UnmarshalKey("accounts", &accounts)

		if configErr != nil {
			color.Red("There was a problem setting up the user account.\nPlease try again or contact Deployer for assistance.")
			os.Exit(1)
		}

		// Get the organisation id
		var orgId string
		viper.UnmarshalKey("orgId", &orgId)

		if len(orgId) <= 0 {
			color.Red("The organisation id is missing from your Deployer configuration.\nPlease check you've correctly configured Deployer on this server and try again.")
			os.Exit(1)
		}

		// Get the server api key
		var serverId string
		viper.UnmarshalKey("serverId", &serverId)

		if len(serverId) <= 0 {
			color.Red("The server API key is missing.\nPlease check you've correctly configured Deployer on this server and try again.")
			os.Exit(1)
		}

		// Get the team api key
		var teamAPIKey string
		viper.UnmarshalKey("teamAPIKey", &teamAPIKey)

		if len(teamAPIKey) <= 0 {
			color.Red("The team API key is missing.\nPlease check you've correctly configured Deployer on this server and try again.")
			os.Exit(1)
		}

		// Get the base domain, which can optionally be overridden
		var baseDomain string
		viper.UnmarshalKey("baseDomain", &baseDomain)

		if len(baseDomain) <= 0 {
			// No overridden base domain, fall back to the default
			baseDomain = "https://deployer.codions.com/api/"
		}

		// Build the base url
		baseURL := baseDomain + "monitoring"

		// Get current system stats
		memory, _ := mem.VirtualMemory()
		loadAvg, _ := load.Avg()
		cpuHw, _ := cpu.Info()
		//percent, _ := cpu.Percent(time.Second, true)
		uptime, _ := host.Uptime()
		misc, _ := load.Misc()
		platform, family, version, _ := host.PlatformInformation()
		diskStat, err := disk.Usage("/")
		currentTime := time.Now()
		timeZone, timeOffset := currentTime.Zone()

		// Create http client
		httpClient := http.Client{
			Timeout: time.Second * 10, // Maximum of 10 secs
		}

		form := url.Values{}

		// Mem
		form.Add("mem[total]", fmt.Sprint(memory.Total))
		form.Add("mem[free]", fmt.Sprint(memory.Free))
		form.Add("mem[used]", fmt.Sprint(memory.Used))
		form.Add("mem[used_percent]", fmt.Sprint(memory.UsedPercent))

		// Loadavg
		form.Add("load[1]", fmt.Sprint(loadAvg.Load1))
		form.Add("load[5]", fmt.Sprint(loadAvg.Load5))
		form.Add("load[15]", fmt.Sprint(loadAvg.Load15))

		// Processes
		form.Add("procs[running]", fmt.Sprint(misc.ProcsRunning))
		form.Add("procs[blocked]", fmt.Sprint(misc.ProcsBlocked))
		form.Add("procs[total]", fmt.Sprint(misc.ProcsTotal))

		// Generic CPU HW
		form.Add("cpu", fmt.Sprint(cpuHw))

		// Uptime
		form.Add("uptime_seconds", fmt.Sprint(uptime))

		// Platform info
		form.Add("platform[name]", fmt.Sprint(platform))
		form.Add("platform[family]", fmt.Sprint(family))
		form.Add("platform[version]", fmt.Sprint(version))

		// Disk info
		form.Add("disk[total]", strconv.FormatUint(diskStat.Total, 10))
		form.Add("disk[used]", strconv.FormatUint(diskStat.Used, 10))
		form.Add("disk[free]", strconv.FormatUint(diskStat.Free, 10))
		form.Add("disk[percent_free]", strconv.FormatFloat(diskStat.UsedPercent, 'f', 2, 64))
		form.Add("disk[inodes_total]", strconv.FormatUint(diskStat.InodesTotal, 10))
		form.Add("disk[inodes_used]", strconv.FormatUint(diskStat.InodesUsed, 10))
		form.Add("disk[inodes_free]", strconv.FormatUint(diskStat.InodesFree, 10))

		// Time info
		form.Add("time[zone]", fmt.Sprint(timeZone))
		form.Add("time[offset]", fmt.Sprint(timeOffset))
		form.Add("time[now]", fmt.Sprint(currentTime.Unix()))

		// Create a request
		req, err := http.NewRequest(http.MethodPost, baseURL, strings.NewReader(form.Encode()))
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("User-Agent", "DeployerAgent-v1.0.0;"+runtime.GOOS)
		req.Header.Set("TeamApiKey", teamAPIKey)
		req.Header.Set("ServerApiKey", serverId)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		httpClient.Do(req)
	},
}

func init() {
	rootCmd.AddCommand(actionCmd)
}
