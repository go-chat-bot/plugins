package cachet

import (
	"encoding/json"
	"fmt"
	"github.com/go-chat-bot/bot"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	statusFailed = 4
)

var (
	cachetAPI               = os.Getenv("CACHET_API")
	configFilePath          = os.Getenv("CACHET_ALERT_CONFIG")
	outageReportConfig      []ChannelConfig
	pastOutageNotifications map[string]time.Time
	pastOutageMutex         = sync.RWMutex{}
)

// cachetComponents is Go representation of https://docs.cachethq.io/reference#get-components
type cachetComponents struct {
	Meta struct {
		Pagination struct {
			Total       int `json:"total"`
			Count       int `json:"count"`
			PerPage     int `json:"per_page"`
			CurrentPage int `json:"current_page"`
			TotalPages  int `json:"total_pages"`
			Links       struct {
				NextPage     string `json:"next_page"`
				PreviousPage string `json:"previous_page"`
			} `json:"links"`
		} `json:"pagination"`
	} `json:"meta"`
	Data []struct {
		ID          int           `json:"id"`
		Name        string        `json:"name"`
		Description string        `json:"description"`
		Link        string        `json:"link"`
		Status      int           `json:"status"`
		Order       int           `json:"order"`
		GroupID     int           `json:"group_id"`
		Enabled     bool          `json:"enabled"`
		Meta        interface{}   `json:"meta"`
		CreatedAt   string        `json:"created_at"`
		UpdatedAt   string        `json:"updated_at"`
		DeletedAt   interface{}   `json:"deleted_at"`
		StatusName  string        `json:"status_name"`
		Tags        []interface{} `json:"tags"`
	} `json:"data"`
}

// ChannelConfig is representation of alert configuration for single channel
type ChannelConfig struct {
	Channel   string   `json:"channel"`
	Services  []string `json:"services"`
	RepeatGap int      `json:"repeatGap"`
}

func cachetGetComponentsFromURL(url string) (components cachetComponents, err error) {
	err = nil
	log.Printf("Getting components from Cachet URL %s", url)
	resp, err := http.Get(url)

	if err != nil {
		log.Printf("Cachet API call failed: %v", err)
		return
	}

	if resp.StatusCode != 200 {
		log.Printf("Cachet API call failed with: %d", resp.StatusCode)
		err = fmt.Errorf("Cachet API call failed with code: %d", resp.StatusCode)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed reading cachet response body: %v", err)
		return
	}
	err = json.Unmarshal(body, &components)
	if err != nil {
		log.Printf("Failed to unmarshal JSON response: %v", err)
		return
	}
	return
}

func cachetGetComponentNames(params string) (names []string, err error) {
	url := fmt.Sprintf("%s/v1/components?%s", cachetAPI, params)
	var components cachetComponents
	for {
		components, err = cachetGetComponentsFromURL(url)
		if err != nil {
			return
		}

		for _, component := range components.Data {
			names = append(names, component.Name)
		}

		url = components.Meta.Pagination.Links.NextPage
		if url == "" {
			// end of paging
			break
		}
	}

	return
}

func getChannelNamesForServiceNotification(service string) (ret []string) {
	for _, channelConfig := range outageReportConfig {
		for _, serviceName := range channelConfig.Services {
			if service == serviceName {
				ret = append(ret, channelConfig.Channel)
			}
		}
	}
	return
}

func getChannelConfig(channel string) (ret *ChannelConfig) {
	for i := range outageReportConfig {
		if outageReportConfig[i].Channel == channel {
			return &outageReportConfig[i]
		}
	}
	return nil
}

func recordOutage(channel string, service string) {
	cc := getChannelConfig(channel)
	if cc == nil {
		log.Printf("Could not find channel config for %s", channel)
		return
	}
	until := time.Now().UTC().Add(time.Duration(cc.RepeatGap) * time.Minute)
	key := fmt.Sprintf("%s-%s", channel, service)
	pastOutageMutex.Lock()
	pastOutageNotifications[key] = until
	pastOutageMutex.Unlock()

	go func() {
		time.Sleep(time.Duration(cc.RepeatGap) * time.Minute)
		pastOutageMutex.Lock()
		delete(pastOutageNotifications, key)
		pastOutageMutex.Unlock()
	}()
}

func checkCachet() (ret []bot.CmdResult, err error) {
	failedNames, err := cachetGetComponentNames("status=4")
	if err != nil {
		log.Printf("Failure while getting failed components: %v", err)
		return
	}
	anyChannels := getChannelNamesForServiceNotification("any")

	for _, failedService := range failedNames {
		notifyChannels := []string{}
		notifyChannels = append(notifyChannels, anyChannels...)
		notifyChannels = append(notifyChannels,
			getChannelNamesForServiceNotification(failedService)...)
		log.Printf("Reporting alerts for %s to %s", failedService, notifyChannels)
		for _, notifyChannel := range notifyChannels {
			key := fmt.Sprintf("%s-%s", notifyChannel, failedService)
			pastOutageMutex.RLock()
			until, found := pastOutageNotifications[key]
			pastOutageMutex.RUnlock()
			if found {
				log.Printf("Skipping notification for %s in %s (until %v)",
					failedService, notifyChannel, until)
				continue
			}
			recordOutage(notifyChannel, failedService)
			log.Printf("Alerting about %s outage in %s", failedService, notifyChannel)
			ret = append(ret, bot.CmdResult{
				Message: fmt.Sprintf("Service '%s' is in outage as per %s",
					failedService, cachetAPI),
				Channel: notifyChannel,
			})
		}
	}
	return
}

func reloadConfig() {
	configFile, err := os.Open(configFilePath)
	if err != nil {
		log.Printf("Failed to open config file: %v", err)
		return
	}
	defer configFile.Close()
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&outageReportConfig)
	if err != nil {
		log.Printf("Failed to parse config file: %v", err)
		return
	}
	log.Printf("Loaded config: %v", outageReportConfig)
}

func saveConfig() {
	log.Printf("Config before save: %v", outageReportConfig)
	configFile, err := os.Create(configFilePath)
	if err != nil {
		log.Printf("Failed to open/write config file: %v", err)
		return
	}
	defer configFile.Close()
	encoder := json.NewEncoder(configFile)
	err = encoder.Encode(&outageReportConfig)
	if err != nil {
		log.Printf("Failed to encode config file: %v", err)
		return
	}
}

func getChannelKey(cmd *bot.Cmd) string {
	if cmd.ChannelData.IsPrivate {
		return cmd.User.Nick
	}
	return cmd.Channel
}

func listComponents(cmd *bot.Cmd) (bot.CmdResultV3, error) {
	componentNames, err := cachetGetComponentNames("")
	log.Printf("Listing services in %s", getChannelKey(cmd))
	result := bot.CmdResultV3{
		Channel: getChannelKey(cmd),
		Message: make(chan string),
		Done:    make(chan bool, 1)}
	if err != nil {
		log.Printf("Failed getting components from cachet: %v", err)
		result.Message <- fmt.Sprintf("Failed getting components from cachet: %v", err)
		result.Done <- true
		return result, err
	}
	go func() {
		result.Message <- "Services known in cachet:"
		curMsgLen := 0
		curComponents := []string{}
		for _, componentName := range componentNames {
			if curMsgLen > 80 {
				log.Printf("Returning partial list of components: %v", curComponents)
				result.Message <- strings.Join(curComponents, ", ")
				curMsgLen = 0
				curComponents = []string{}
			}
			curMsgLen = curMsgLen + len(componentName)
			curComponents = append(curComponents, componentName)
		}
		log.Printf("Returning last part of components: %v", curComponents)
		result.Message <- strings.Join(curComponents, ", ")
		result.Done <- true
	}()
	return result, err
}

func listSubscriptions(cmd *bot.Cmd) (string, error) {
	channelKey := getChannelKey(cmd)
	channelConfig := getChannelConfig(channelKey)
	if channelConfig != nil && channelConfig.Channel == channelKey {
		return fmt.Sprintf("This channel is subscribed to notifications for: %v",
			channelConfig.Services), nil
	}
	return "This channel has no subscriptions", nil
}

func subscribeChannel(cmd *bot.Cmd) (ret string, err error) {
	if len(cmd.Args) != 1 {
		return "Expecting 1 argument: <name of service>", nil
	}
	channelKey := getChannelKey(cmd)
	newService := cmd.Args[0]
	channelConfig := getChannelConfig(channelKey)
	ret = fmt.Sprintf("Succesfully subscribed channel %s to outage notifications for '%s'",
		channelKey, newService)
	defer saveConfig()
	if channelConfig == nil {
		log.Printf("Channel %s has no config yet. Adding new one", channelKey)
		outageReportConfig = append(outageReportConfig, ChannelConfig{
			Channel:   channelKey,
			Services:  []string{newService},
			RepeatGap: 5,
		})
		return
	}

	for _, service := range channelConfig.Services {
		if service == newService {
			return fmt.Sprintf(
				"This channel is already subscribed to '%s' outage notifications",
				service), nil
		}
	}
	log.Printf("Channel already has a config. Appending new service notification")
	channelConfig.Services = append(channelConfig.Services, newService)
	log.Printf("New notifications: %s", channelConfig.Services)
	return
}

func unsubscribeChannel(cmd *bot.Cmd) (string, error) {
	if len(cmd.Args) != 1 {
		return "Expecting 1 argument: <name of service>", nil
	}
	channelKey := getChannelKey(cmd)
	channelConfig := getChannelConfig(channelKey)
	newService := cmd.Args[0]
	if channelConfig == nil {
		return "Channel is not subscribed to anything", nil
	}

	newServices := []string{}
	for _, service := range channelConfig.Services {
		if service == newService {
			continue
		}
		newServices = append(newServices, service)
	}

	log.Printf("Channel already has a config. Appending new service notification")
	channelConfig.Services = newServices
	log.Printf("New notifications: %s", channelConfig.Services)
	saveConfig()
	return fmt.Sprintf(
		"Succesfully unsubscribed channel %s from outage notifications for %s",
		channelKey, newService), nil
}

func outageRepeatGap(cmd *bot.Cmd) (ret string, err error) {
	if len(cmd.Args) != 1 {
		return "Expecting 1 argument: <number of minutes>", nil
	}
	min, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		return "Argument must be exactly 1 number (of minutes between notifications)", nil
	}
	channelKey := getChannelKey(cmd)
	channelConfig := getChannelConfig(channelKey)
	defer saveConfig()
	ret = fmt.Sprintf("Succesfully configured notification gap to be %d minutes", min)
	if channelConfig == nil {
		log.Printf("Channel has no config yet. Adding new one")
		outageReportConfig = append(outageReportConfig, ChannelConfig{
			Channel:   cmd.Channel,
			Services:  []string{},
			RepeatGap: min,
		})
		return
	}
	channelConfig.RepeatGap = min
	return
}

func init() {
	pastOutageNotifications = make(map[string]time.Time)
	reloadConfig()

	bot.RegisterPeriodicCommandV2(
		"systemStatusCheck",
		bot.PeriodicConfig{
			CronSpec:  "@every 1m",
			CmdFuncV2: checkCachet,
		})
	bot.RegisterCommandV3(
		"services",
		"List services available for subscriptions",
		"",
		listComponents)
	bot.RegisterCommand(
		"subscriptions",
		"Lists active outage subscriptions",
		"",
		listSubscriptions)
	bot.RegisterCommand(
		"subscribe",
		"Subscribes this channel to outage notifications of specific service (or 'any' for all outages)",
		"<service>",
		subscribeChannel)
	bot.RegisterCommand(
		"unsubscribe",
		"Unsubscribes this channel from outage notifications of specific service",
		"<service>",
		unsubscribeChannel)
	bot.RegisterCommand(
		"repeatgap",
		"Sets number of minutes between notification of specific service outage",
		"60",
		outageRepeatGap)
}
