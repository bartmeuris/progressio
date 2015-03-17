package progressio

import "fmt"

// Various constants related to the units
const (
	Byte int64 = 1	// Byte is the representation of a single byte

	MetricMultiplier = 1000	// Metric uses 1 10^3 multiplier
	KiloByte         = Byte * MetricMultiplier	// Metric unit "KiloByte" constant
	MegaByte         = KiloByte * MetricMultiplier	// Metric unit MegaByte constant
	GigaByte         = MegaByte * MetricMultiplier	// Metric unit GigaByte constant
	TeraByte         = GigaByte * MetricMultiplier	// Metric unit TerraByte constant
	PetaByte         = TeraByte * MetricMultiplier	// Metric unit PetaByte constant

	IECMultiplier = 1024	// IEC Standard multiplier, 1024 based
	KibiByte      = Byte * IECMultiplier		// IEC standard unit KibiByte constant
	MebiByte      = KibiByte * IECMultiplier	// IEC standard unit MebiByte constant
	GibiByte      = MebiByte * IECMultiplier	// IEC standard unit GibiByte constant
	TebiByte      = GibiByte * IECMultiplier	// IEC standard unit TebiByte constant
	PebiByte      = TebiByte * IECMultiplier	// IEC standard unit PebiByte constant

	JEDECKiloByte = KibiByte // JEDEC uses IEC multipliers, but Metric names, JEDEC KiloByte constant
	JEDECMegaByte = MebiByte // JEDEC uses IEC multipliers, but Metric names, JEDEC MegaByte constant
	JEDECGigaByte = GibiByte // JEDEC uses IEC multipliers, but Metric names, JEDEC GigaByte constant
)


// IECNames is an array containing the unit names for the IEC standards
var IECNames = []string{
	"byte",
	"kibibyte",
	"mebibyte",
	"gibibyte",
	"tebibyte",
	"pebibyte",
}
// IECShorts is an array containing the shortened unit names for the IEC standard
var IECShorts = []string{
	"B",
	"KiB",
	"MiB",
	"GiB",
	"TiB",
	"PiB",
}

// JEDECNames is an array containing the unit names for the JEDEC standard
var JEDECNames = []string{
	"byte",
	"kilobyte",
	"megabyte",
	"gigabyte",
}
// JEDECShorts is an array containing the shortened unit names for the JEDEC standard
var JEDECShorts = []string{
	"B",
	"KB",
	"MB",
	"GB",
}


// MetricNames is an array containing the unit names for the metric units
var MetricNames = []string{
	"byte",
	"kilobyte",
	"megabyte",
	"gigabyte",
	"terabyte",
	"petabyte",
}
// MetricShorts is an array containing the shortened unit names for the metric units
var MetricShorts = []string{
	"B",
	"kB",
	"MB",
	"GB",
	"TB",
	"PB",
}


// SizeSystem is a structure representing a unit standard
type SizeSystem struct {
	Name       string	// The name of the unit standard
	MultiPlier int64	// The multiplier used by the unit standard
	Names      []string	// The names used by the unit standard
	Shorts     []string	// The shortened names used by the unit standard
}

// Metric is a SizeSystem instance representing the metric system
var Metric = SizeSystem{
	Name:       "metric",
	MultiPlier: MetricMultiplier,
	Names:      MetricNames,
	Shorts:     MetricShorts,
}
// IEC is a SizeSystem instance representing the IEC standard
var IEC = SizeSystem{
	Name:       "IEC",
	MultiPlier: IECMultiplier,
	Names:      IECNames,
	Shorts:     IECShorts,
}
// JEDEC is a SizeSystem instance representing the JEDEC standard
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

// FormatSize formats a number of bytes using the given unit standard system.
// If the 'short' flag is set to true, it uses the shortened names.
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
