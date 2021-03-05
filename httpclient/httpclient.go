package httpclient

import (
	"sync"

	"../../atella"
)

// universalResponse describe tmplate for api responses in json format
type universalResponse struct {
	ResponseBody interface{} `json:"response"`
}

// Client is client config structure.
type Client struct {
	waitGroup         *sync.WaitGroup
	sectorMapper      map[string]atella.Sector
	sectorMapperMutex sync.RWMutex
	hosts             []atella.Host
	logger            atella.Logger
	code              string
	hostname          string
	connectivity      int64
	IAmIsMaster       bool
	stopRequest       chan struct{}
	interval          int
	masterInterval    int
	masterChannel     chan map[string]atella.HostVector
	mastersIndexes    []int64
	masterIndex       int64
	Vector            *atella.Vector
}

// NewHTTPClient return new HTTP client descriptor.
func NewHTTPClient(hostname, code string, hosts []atella.Host, parsedHosts *atella.ParsedHosts,
	connectivity int64, logger atella.Logger) (*Client,
	chan map[string]atella.HostVector) {
	client := &Client{
		waitGroup:         &sync.WaitGroup{},
		sectorMapper:      parsedHosts.SectorMapper,
		sectorMapperMutex: sync.RWMutex{},
		hosts:             hosts,
		logger:            logger,
		hostname:          hostname,
		code:              code,
		connectivity:      connectivity,
		IAmIsMaster: atella.SubsetInt64(parsedHosts.MasterIndexes,
			parsedHosts.SelfIndexes),
		Vector:         atella.NewVector(),
		masterInterval: 10,
		stopRequest:    make(chan struct{}),
		masterChannel:  make(chan map[string]atella.HostVector),
		mastersIndexes: parsedHosts.MasterIndexes,
		masterIndex:    -1,
	}
	return client, client.masterChannel
}

// Start client.
func (c *Client) Start() {
	for sectorName, sector := range c.sectorMapper {
		if len(sector.SelfIndexes) > 0 {
			c.logger.Infof("Start client for sector %s", sectorName)
			c.waitGroup.Add(1)
			go c.startSectorClient(sectorName)
		}
	}
	go c.startMasterClient()
}

// Starting client in something sector.
func (c *Client) startSectorClient(sectorName string) {
	var (
		exit bool = false
		// startIndex int64 = c.sectorMapper[sectorName].SelfIndexes[0]
		currentConnectivity int64   = 0
		indexesToCheck      []int64 = make([]int64, 0)
	)

	// find hosts who will be checked.
	c.sectorMapperMutex.RLock()
	hostIndexes := c.sectorMapper[sectorName].HostsIndexes
	selfIndexes := c.sectorMapper[sectorName].SelfIndexes
	c.sectorMapperMutex.RUnlock()
	for index := 0; index < len(hostIndexes) &&
		currentConnectivity < c.connectivity; index++ {
		if !atella.ElExistsInt64(selfIndexes,
			hostIndexes[index]) {
			currentConnectivity++
			indexesToCheck = append(indexesToCheck, hostIndexes[index])
		}
	}

	c.logger.Infof("Sector '%s', check indexes %v", sectorName, indexesToCheck)
	if len(indexesToCheck) <= 0 {
		c.logger.Warningf("Sector '%s', nothing to check", sectorName)
		return
	}
	go func() {
		<-c.stopRequest
		exit = true
	}()
	for !exit {
		for _, index := range indexesToCheck {
			h, code, err := c.getHostnameAPIv1(c.hosts[index].Address,
				c.hosts[index].Port,
				false)
			vector, empty := c.Vector.GetElement(c.hosts[index].Address, c.hosts[index].Port)
			if empty {
				vector.Address = c.hosts[index].Address
				vector.Port = c.hosts[index].Port
				vector.Hostname = c.hosts[index].Hostname
			}

			if err != nil {
				c.logger.Warning(err.Error())
				vector.Status = false
				c.logger.Infof("%s:%s,response code %d",
					c.hosts[index].Address,
					c.hosts[index].Port, code)
			} else if h == vector.Hostname {
				vector.Status = true
			}

			c.Vector.SetElement(vector.Address, vector.Port, vector)
		}
		atella.Pause(10, &exit)
	}
}

// Stop client.
func (c *Client) Stop() {
	close(c.stopRequest)
}
