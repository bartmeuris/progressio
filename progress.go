package progressio

import (
	"time"
	"fmt"
)

const UpdateFreq = 100 * time.Millisecond
const Timeslots  = 5

type Progress struct {
	Transferred  int64   // in bytes
	TotalSize    int64   // in bytes
	Percent      float64
	SpeedAvg     int64   // Bytes/sec average
	Speed        int64   // Bytes/sec of last transfer
	Remaining    time.Duration // Estimated time remaining
	StartTime    time.Time
	StopTime     time.Time
}

type ioProgress struct {
	size      int64
	progress  int64
	ch        chan Progress
	closed    bool
	startTime time.Time
	lastSent  time.Time
	updatesW  []int64
	updatesT  []time.Time
	ts        int
}

func (p *Progress) String() string {
	timeS := fmt.Sprintf(" (Time: %s", FormatDuration(time.Since(p.StartTime)))
	// Build the Speed string
	speedS := ""
	if p.Speed > 0 {
		speedS = fmt.Sprintf(" (Speed: %s", FormatSize(IEC, p.Speed, true))+"/s"
	}
	if p.SpeedAvg > 0 {
		if len(speedS) > 0 {
			speedS += " / AVG: "
		} else {
			speedS = " (Speed AVG: "
		}
		speedS += FormatSize(IEC, p.SpeedAvg, true)+"/s"
	}
	if len(speedS) > 0 {
		speedS += ")"
	}

	if (p.TotalSize <= 0) {
		// No size was given, we can only show:
		// - Amount read/written
		// - average speed
		// - current speed
		return fmt.Sprintf("%s%s%s)",
			FormatSize(IEC, p.Transferred, true),
			speedS,
			timeS,
		)
	}
	// A size was given, we can add:
	// - Percentage
	// - Progress indicator
	// - Remaining time
	timeR := ""
	if p.Remaining >= time.Duration(0) {
		timeR = fmt.Sprintf(" / Remaining: %s", FormatDuration(p.Remaining))
	}


	return fmt.Sprintf("[%02.2f%%] (%s/%s)%s%s%s)",
		p.Percent,
		FormatSize(IEC, p.Transferred, true),
		FormatSize(IEC, p.TotalSize, true),
		speedS,
		timeS,
		timeR,
	)
}


func mkIoProgress(size int64) (*ioProgress) {
	return &ioProgress{
		size       : size,
		progress   : 0,
		ch         : make(chan Progress),
		closed     : false,
		startTime  : time.Time{},
		lastSent   : time.Time{},
		updatesW   : make([]int64, Timeslots),
		updatesT   : make([]time.Time, Timeslots),
		ts         : 0,
	}
}

func (p *ioProgress) updateProgress(written int64) {
	if p.closed && p.ch == nil {
		// Nothing to do
		return
	}
	if written > 0 {
		p.progress += written
	}
	// Throttle sending updated, limit to UpdateFreq - which should be 100ms
	if time.Since(p.lastSent) < UpdateFreq {
		return
	}
	if p.startTime.IsZero() {
		p.startTime = time.Now()
	}

	prog := Progress{
		StartTime: p.startTime,
		Transferred: p.progress,
		TotalSize: p.size,
	}

	// Calculate current speed based on the last `Timeslots` updates sent
	// Timeslots should be 20, which means this is minimum 20 100ms slots, being at least 2 seconds
	p.updatesW[p.ts % Timeslots] = p.progress
	p.updatesT[p.ts % Timeslots] = time.Now()
	p.ts++
	if !p.updatesT[p.ts % Timeslots].IsZero() {
		// Calculate the average speed of the last ~2 seconds
		prog.Speed = int64( (float64(p.progress - p.updatesW[p.ts % Timeslots]) / float64(time.Since(p.updatesT[p.ts % Timeslots]))) * float64(time.Second) )

		// Calculate the average speed since starting the transfer
		tp := time.Since(p.startTime)
		if tp > 0 {
			prog.SpeedAvg = int64( (float64(p.progress) / float64(tp)) * float64(time.Second) )
		} else {
			prog.SpeedAvg = -1
		}
		if p.size > 0 && prog.SpeedAvg > 0 {
			prog.Remaining = time.Duration((float64(p.size - p.progress) / float64(prog.SpeedAvg)) * float64(time.Second))
		} else {
			prog.Remaining = -1
		}
	} else {
		prog.Speed = -1
		prog.SpeedAvg = -1
		prog.Remaining = -1
	}
	
	// Calculate the percentage only if we have a size
	if (p.size > 0) {
		prog.Percent = float64(int64((float64(p.progress) / float64(p.size)) * 10000.0)) / 100.0
	}

	if p.closed || p.progress == p.size {
		// EOF or closed, we have to send this last message, and then close the chan
		// Prevent sending the last message multiple times
		if p.ch != nil {
			prog.StopTime = time.Now()
			p.ch <- prog
			p.cleanup()
		}
	} else {
		// Don't force send, only send when it would not block, the chan is non-buffered
		select {
		case p.ch <- prog:
			// update last sent values
			p.lastSent = time.Now()
		default:
		}
	}
}

func (p *ioProgress) cleanup() {
	p.closed = true
	if p.ch != nil {
		close(p.ch)
		p.ch = nil
	}
}

func (p *ioProgress) stopProgress() {
	p.closed = true
	p.updateProgress(-1)
}

func (p *ioProgress) GetChannel() (<- chan Progress) {
	return p.ch
}

