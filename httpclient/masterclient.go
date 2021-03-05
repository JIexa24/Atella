package httpclient

import (
	"math/rand"

	"../../atella"
)

// Starting client in something sector.
func (c *Client) startMasterClient() {
	exit := false
	dontHaveMasters := make(chan struct{})
	go func() {
		select {
		case <-c.stopRequest:
			exit = true
		case <-dontHaveMasters:
			exit = true
		}
	}()
	c.masterElection()

	if !c.IAmIsMaster && c.masterIndex < 0 {
		c.logger.Errorf("Master servers not conigure!")
		close(dontHaveMasters)
	}

	for !exit {
		atella.Pause(10, &exit)
		c.sendToMaster()
	}
}

func (c *Client) sendToMaster() {
	if c.IAmIsMaster {
		v := c.Vector.GetVectorCopy()
		c.masterChannel <- v.List
		return
	}

	if err := c.setVectorAPIv1(c.hosts[c.masterIndex].Address,
		c.hosts[c.masterIndex].Port, false); err != nil {
		c.logger.Warning(err.Error())
	}
}

// Election new master host and set it index.
func (c *Client) masterElection() {
	if c.IAmIsMaster || c.mastersIndexes == nil {
		return
	}

	if len(c.mastersIndexes) <= 0 {
		c.masterIndex = -1
		return
	}

	if len(c.mastersIndexes) > 1 {
		c.logger.Warning("Only one master server supported")
		c.masterIndex = c.mastersIndexes[0]
	}

	if c.masterIndex < 0 {
		index := rand.Intn(len(c.mastersIndexes))
		c.masterIndex = int64(index)
		return
	}
}
