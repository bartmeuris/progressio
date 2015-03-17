/*
Package progressio contains io.Reader and io.Writer wrappers to easily get
progress feedback, including speed/sec, average speed, %, time remaining,
size, transferred size, ...  over a channel in a progressio.Progress object.

Important note is that the returned object implements the io.Closer interface
and you have to close the progressio.ProgressReader and
progressio.ProgressWriter objects in order to clean everything up.

Usage is pretty simple:


	preader, pchan := progressio.NewProgressReader(myreader, -1)
	defer preader.Close()
	go func() {
		for p := range pchan {
			fmt.Printf("Progress: %s\n", p.String())
		}
	}
	// read from your new reader object
	io.Copy(mywriter, preader)


A helper function is available that opens a file, determines it's size, and
wraps it's os.File io.Reader object:


	if pr, pc, err := progressio.NewProgressFileReader(myfile); err != nil {
		return err
	} else {
		defer pr.Close()
		go func() {
			for p := range pc{
				fmt.Printf("Progress: %s\n", p.String())
			}
		}
		// read from your new reader object
		io.Copy(mywriter, pr)
	}


A wrapper for an io.WriterCloser is available too, but no helper function 
is available to write to an os.File since the target size is not known.
Usually, wrapping the io.Writer is more accurate, since writing potentially
takes up more time and happens last. Useage is similar to wrapping the 
io.Reader:


▸   pwriter, pchan := progressio.NewProgressWriter(mywriter, -1)
▸   defer pwriter.Close()
▸   go func() {
▸   ▸   for p := range pchan {
▸   ▸   ▸   fmt.Printf("Progress: %s\n", p.String())
▸   ▸   }
▸   }
▸   // write to your new writer object
▸   io.Copy(pwriter, myreader)


Note that you can also implement your own formatting. See the String() function
implementation or consult the Progress struct layout and documentation

*/
package progressio

import (
	"fmt"
	"time"
)

// Frequency of the updates over the channels
const UpdateFreq = 100 * time.Millisecond
const timeSlots = 5

// Progress is the object sent back over the progress channel.
type Progress struct {
	Transferred int64         // Transferred data in bytes
	TotalSize   int64         // Total size of the transfer in bytes. <= 0 if size is unknown.
	Percent     float64       // If the size is known, the progress of the transfer in %
	SpeedAvg    int64         // Bytes/sec average over the entire transfer
	Speed       int64         // Bytes/sec of the last few reads/writes
	Remaining   time.Duration // Estimated time remaining, only available if the size is known.
	StartTime   time.Time     // When the transfer was started
	StopTime    time.Time     // only specified when the transfer is completed: when the transfer was stopped
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

// String returns a string representation of the progress. It takes into account
// if the size was known, and only tries to display relevant data.
func (p *Progress) String() string {
	timeS := fmt.Sprintf(" (Time: %s", FormatDuration(time.Since(p.StartTime)))
	// Build the Speed string
	speedS := ""
	if p.Speed > 0 {
		speedS = fmt.Sprintf(" (Speed: %s", FormatSize(IEC, p.Speed, true)) + "/s"
	}
	if p.SpeedAvg > 0 {
		if len(speedS) > 0 {
			speedS += " / AVG: "
		} else {
			speedS = " (Speed AVG: "
		}
		speedS += FormatSize(IEC, p.SpeedAvg, true) + "/s"
	}
	if len(speedS) > 0 {
		speedS += ")"
	}

	if p.TotalSize <= 0 {
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

func mkIoProgress(size int64) *ioProgress {
	return &ioProgress{
		size:      size,
		progress:  0,
		ch:        make(chan Progress),
		closed:    false,
		startTime: time.Time{},
		lastSent:  time.Time{},
		updatesW:  make([]int64, timeSlots),
		updatesT:  make([]time.Time, timeSlots),
		ts:        0,
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
	// Always send when finished
	if (time.Since(p.lastSent) < UpdateFreq) && ((p.size > 0) && (p.progress != p.size)) {
		return
	}
	if p.startTime.IsZero() {
		p.startTime = time.Now()
	}

	prog := Progress{
		StartTime:   p.startTime,
		Transferred: p.progress,
		TotalSize:   p.size,
	}

	// Calculate current speed based on the last `timeSlots` updates sent
	p.updatesW[p.ts%timeSlots] = p.progress
	p.updatesT[p.ts%timeSlots] = time.Now()
	p.ts++
	if !p.updatesT[p.ts%timeSlots].IsZero() {
		// Calculate the average speed of the last ~2 seconds
		prog.Speed = int64((float64(p.progress-p.updatesW[p.ts%timeSlots]) / float64(time.Since(p.updatesT[p.ts%timeSlots]))) * float64(time.Second))

		// Calculate the average speed since starting the transfer
		tp := time.Since(p.startTime)
		if tp > 0 {
			prog.SpeedAvg = int64((float64(p.progress) / float64(tp)) * float64(time.Second))
		} else {
			prog.SpeedAvg = -1
		}
		if p.size > 0 && prog.SpeedAvg > 0 {
			prog.Remaining = time.Duration((float64(p.size-p.progress) / float64(prog.SpeedAvg)) * float64(time.Second))
		} else {
			prog.Remaining = -1
		}
	} else {
		prog.Speed = -1
		prog.SpeedAvg = -1
		prog.Remaining = -1
	}

	// Calculate the percentage only if we have a size
	if p.size > 0 {
		prog.Percent = float64(int64((float64(p.progress)/float64(p.size))*10000.0)) / 100.0
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
