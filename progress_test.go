package progressio

import (
	"testing"
	"time"
)

func TestPrintNoSize(t *testing.T) {
	var s, expect string
	p := Progress{}
	// Full scenario
	p.Speed = 100 * KibiByte    // 100KiB/sec
	p.SpeedAvg = 100 * KibiByte // 100KiB/sec
	p.Remaining = -1
	p.Transferred = MebiByte * 10
	p.TotalSize = 0
	p.Percent = 0
	p.StartTime = time.Now().Add(time.Second * -5)

	expect = "10.00MiB (Speed: 100.00KiB/s / AVG: 100.00KiB/s) (Time: 5 seconds)"
	s = p.String()
	if s != expect {
		t.Log("TestPrintNoSize: full failed:")
		t.Logf("   Got     : '%s'\n", s)
		t.Logf("   Expected: '%s'\n", expect)
		t.Fail()
	}
}
func TestPrintSize(t *testing.T) {
	var s, expect string
	p := Progress{}

	// Test full scenario
	p.Speed = 100 * KibiByte    // 100KiB/sec
	p.SpeedAvg = 100 * KibiByte // 100KiB/sec
	p.Remaining = time.Second * 10
	p.Transferred = MebiByte * 10
	p.TotalSize = MebiByte * 20
	p.Percent = 50.0
	p.StartTime = time.Now().Add(time.Second * -5)
	expect = "[50.00%] (10.00MiB/20.00MiB) (Speed: 100.00KiB/s / AVG: 100.00KiB/s) (Time: 5 seconds / Remaining: 10 seconds)"
	s = p.String()
	if s != expect {
		t.Log("TestPrintSize: full failed:")
		t.Logf("   Got     : '%s'\n", s)
		t.Logf("   Expected: '%s'\n", expect)
		t.Fail()
	}
	// Test without p.SpeedAvg
	p.SpeedAvg = 0
	p.StartTime = time.Now().Add(time.Second * -5)
	expect = "[50.00%] (10.00MiB/20.00MiB) (Speed: 100.00KiB/s) (Time: 5 seconds / Remaining: 10 seconds)"
	s = p.String()
	if s != expect {
		t.Log("TestPrintSize: without p.SpeedAvg failed:")
		t.Logf("   Got     : '%s'\n", s)
		t.Logf("   Expected: '%s'\n", expect)
		t.Fail()
	}
	// Test without p.Remaining
	p.SpeedAvg = 100 * KibiByte
	p.Remaining = -1
	p.StartTime = time.Now().Add(time.Second * -5)
	expect = "[50.00%] (10.00MiB/20.00MiB) (Speed: 100.00KiB/s / AVG: 100.00KiB/s) (Time: 5 seconds)"
	s = p.String()
	if s != expect {
		t.Log("TestPrintSize: without p.Remaining failed:")
		t.Logf("   Got     : '%s'\n", s)
		t.Logf("   Expected: '%s'\n", expect)
		t.Fail()
	}
	// Test p.Remaining == 0
	p.Remaining = 0
	p.StartTime = time.Now().Add(time.Second * -5)
	expect = "[50.00%] (10.00MiB/20.00MiB) (Speed: 100.00KiB/s / AVG: 100.00KiB/s) (Time: 5 seconds / Remaining: 0 seconds)"
	s = p.String()
	if s != expect {
		t.Log("TestPrintSize: with p.Remaining == 0 failed:")
		t.Logf("   Got     : '%s'\n", s)
		t.Logf("   Expected: '%s'\n", expect)
		t.Fail()
	}
}

func TestIOProgress(t *testing.T) {
	iop := mkIoProgress(100 * MebiByte)
	iop.progress = 50 * MebiByte
	iop.updatesW[1] = 40 * MebiByte
	iop.updatesT[1] = time.Now().Add(time.Second * -1)
	iop.startTime = time.Now().Add(time.Second * -10)
	go iop.updateProgress(0)
	p := <-iop.ch
	t.Logf("P: %p\n", &p)
	t.Logf("P: %s\n", p.String())
	//t.Fail()
}
