package progressio

import "fmt"

const (
	Byte int64 = 1

	MetricMultiplier = 1000
	KiloByte         = Byte * MetricMultiplier
	MegaByte         = KiloByte * MetricMultiplier
	GigaByte         = MegaByte * MetricMultiplier
	TeraByte         = GigaByte * MetricMultiplier
	PetaByte         = TeraByte * MetricMultiplier

	IECMultiplier = 1024
	KibiByte      = Byte * IECMultiplier
	MebiByte      = KibiByte * IECMultiplier
	GibiByte      = MebiByte * IECMultiplier
	TebiByte      = GibiByte * IECMultiplier
	PebiByte      = TebiByte * IECMultiplier

	JEDECKiloByte = KibiByte
	JEDECMegaByte = MebiByte
	JEDECGigaByte = GibiByte
)

var IECNames = []string{
	"byte",
	"kibibyte",
	"mebibyte",
	"gibibyte",
	"tebibyte",
	"pebibyte",
}
var IECShorts = []string{
	"B",
	"KiB",
	"MiB",
	"GiB",
	"TiB",
	"PiB",
}

var JEDECShorts = []string{
	"B",
	"KB",
	"MB",
	"GB",
}
var JEDECNames = []string{
	"byte",
	"kilobyte",
	"megabyte",
	"gigabyte",
}

var MetricShorts = []string{
	"B",
	"kB",
	"MB",
	"GB",
	"TB",
	"PB",
}

var MetricNames = []string{
	"byte",
	"kilobyte",
	"megabyte",
	"gigabyte",
	"terabyte",
	"petabyte",
}

type SizeSystem struct {
	Name       string
	MultiPlier int64
	Names      []string
	Shorts     []string
}

var Metric = SizeSystem{
	Name:       "metric",
	MultiPlier: MetricMultiplier,
	Names:      MetricNames,
	Shorts:     MetricShorts,
}
var IEC = SizeSystem{
	Name:       "IEC",
	MultiPlier: IECMultiplier,
	Names:      IECNames,
	Shorts:     IECShorts,
}

var JEDEC = SizeSystem{
	Name:       "JEDEC",
	MultiPlier: IECMultiplier,
	Names:      JEDECNames,
	Shorts:     JEDECShorts,
}

func getUnit(ss SizeSystem, size int64) (divider int64, name, short string) {
	if size < 0 {
		size = -size
	}
	if size == 0 {
		return 1, ss.Names[0], ss.Shorts[0]
	}
	div := Byte
	for i := 0; i < len(ss.Names); i++ {
		//fmt.Printf("TEST[%d]: %d / DIV: %d | %s / %s\n", i, size, div, ss.Names[i], ss.Shorts[i])
		if div <= size {
			div *= ss.MultiPlier
			continue
		}
		return div / ss.MultiPlier, ss.Names[i-1], ss.Shorts[i-1]
	}
	return div, ss.Names[len(ss.Names)-1], ss.Shorts[len(ss.Shorts)-1]
}

func FormatSize(ss SizeSystem, size int64, short bool) string {
	div, name, shortnm := getUnit(ss, size)
	ds := float64(size) / float64(div)
	numfm := "%.2f"
	if div == 1 {
		numfm = "%.0f"
	}
	if short {
		return fmt.Sprintf(numfm+"%s", ds, shortnm)
	}
	return fmt.Sprintf(numfm+" %s", ds, name)
}

func testUnit(ss SizeSystem, sz int64) []interface{} {
	div, name, short := getUnit(ss, sz)
	return []interface{}{sz, div, name, short, FormatSize(ss, sz, true), FormatSize(ss, sz, false)}
}

func testSize(ss SizeSystem, sz int64) {
	fmt.Printf("---- %s: %d ----\n", ss.Name, sz)
	fmt.Printf("  -1: %d: M: %d | '%s', '%s' -- %s / %s\n", testUnit(ss, sz-1)...)
	fmt.Printf("  +0: %d: M: %d | '%s', '%s' -- %s / %s\n", testUnit(ss, sz)...)
	fmt.Printf("  +1: %d: M: %d | '%s', '%s' -- %s / %s\n", testUnit(ss, sz+1)...)
}

func test() {
	testSize(Metric, -1)
	testSize(Metric, 0)
	testSize(Metric, 1)
	testSize(Metric, KiloByte)
	testSize(Metric, MegaByte)
	testSize(Metric, GigaByte)
	testSize(Metric, TeraByte)
	testSize(Metric, PetaByte)

	testSize(IEC, 0)
	testSize(IEC, 1)
	testSize(IEC, KibiByte)
	testSize(IEC, MebiByte)
	testSize(IEC, GibiByte)
	testSize(IEC, TebiByte)
	testSize(IEC, PebiByte)

	testSize(JEDEC, 0)
	testSize(JEDEC, 1)
	testSize(JEDEC, KibiByte)
	testSize(JEDEC, MebiByte)
	testSize(JEDEC, GibiByte)
	testSize(JEDEC, TebiByte)
	testSize(JEDEC, PebiByte)
}
