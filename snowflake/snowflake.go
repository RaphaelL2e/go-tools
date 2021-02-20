package snowflake

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	//global var
	sequence = 0
	lastTime = -1

	//every segment bit
	workerIdBits     = 5
	datacenterIdBits = 5
	sequenceBits     = 12

	//every segment max number
	maxWorkerId     = -1 ^ (-1 << workerIdBits)
	maxDatacenterId = -1 ^ (-1 << datacenterIdBits)
	maxSequenceId   = -1 ^ (-1 << sequenceBits)

	//bit operation shift
	workerIdShift   = sequenceBits
	datacenterShift = workerIdBits + sequenceBits
	timestampShift  = datacenterIdBits + workerIdBits + sequenceBits
)

type Snowflake struct {
	datacenterId int
	workerId     int
	epoch        int
	mt           *sync.Mutex
}

func NewSnowflake(datacenterId, workerId, epoch int) (*Snowflake, error) {
	if datacenterId > maxDatacenterId || datacenterId < 0 {
		return nil, errors.New(fmt.Sprintf("datacenterId cant be greater than %d or less than 0", maxDatacenterId))
	}
	if workerId > maxWorkerId || workerId < 0 {
		return nil, errors.New(fmt.Sprintf("workerId cant be greater than %d or less than 0", maxWorkerId))
	}
	if epoch > getCurrentTime() {
		return nil, errors.New(fmt.Sprintf("epoch time cant be after now"))
	}
	sf := Snowflake{datacenterId, workerId, epoch, new(sync.Mutex)}
	return &sf, nil
}

func (sf *Snowflake) GetUniqueId() int {
	sf.mt.Lock()
	defer sf.mt.Unlock()
	//get current time
	currentTime := getCurrentTime()
	//compute sequence
	if currentTime < lastTime {
		currentTime = waitUntilNextTime(lastTime)
	} else if currentTime == lastTime {
		sequence = (sequence + 1) & maxSequenceId
		if sequence == 0 {
			currentTime = waitUntilNextTime(lastTime)
		}
	} else if currentTime > lastTime {
		sequence = 0
		lastTime = currentTime
	}
	//generate id
	return (currentTime-sf.epoch)<<timestampShift |
		sf.datacenterId<<datacenterShift |
		sf.workerId<<workerIdShift |
		sequence
}

func waitUntilNextTime(lasttime int) int {
	currentTime := getCurrentTime()
	for currentTime <= lasttime {
		time.Sleep(1 * time.Second / 1000) //sleep micro second
		currentTime = getCurrentTime()
	}
	return currentTime
}

func getCurrentTime() int {
	return int(time.Now().UnixNano() / 1e6) //micro second
}