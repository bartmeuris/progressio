package progressio

import "testing"

func Test_getUnit(t *testing.T) {

	// This is to more visually exercise the top edge case of the getUnit routine
	var distanceSize = SizeSystem{
		Name:       "distance",
		MultiPlier: 1000,
		Names:      []string{"metre", "kilometre"},
		Shorts:     []string{"m", "km"},
	}

	type args struct {
		ss   SizeSystem
		size int64
	}
	tests := []struct {
		name        string
		args        args
		wantDivider int64
		wantName    string
		wantShort   string
	}{
		{
			name:        "IEC 0",
			args:        args{ss: IEC, size: 0},
			wantDivider: Byte, wantName: "byte", wantShort: "B",
		},
		{
			name:        "IEC 1234",
			args:        args{ss: IEC, size: 1234},
			wantDivider: KibiByte, wantName: "kibibyte", wantShort: "KiB",
		},
		{
			name:        "IEC 1234567",
			args:        args{ss: IEC, size: 1234567},
			wantDivider: MebiByte, wantName: "mebibyte", wantShort: "MiB",
		},
		{
			name:        "JEDEC 1234",
			args:        args{ss: JEDEC, size: 1234},
			wantDivider: KibiByte, wantName: "kilobyte", wantShort: "KB",
		},
		{
			name:        "JEDEC 1234567",
			args:        args{ss: JEDEC, size: 1234567},
			wantDivider: MebiByte, wantName: "megabyte", wantShort: "MB",
		},
		{ // Uses top element
			name:        "IEC top edge case",
			args:        args{ss: IEC, size: PebiByte},
			wantDivider: PebiByte, wantName: "pebibyte", wantShort: "PiB",
		},
		{ // Exceeds top element
			name:        "IEC exceeds top edge case",
			args:        args{ss: IEC, size: IECMultiplier * PebiByte},
			wantDivider: PebiByte, wantName: "pebibyte", wantShort: "PiB",
		},
		{
			name:        "Distance 1",
			args:        args{ss: distanceSize, size: 1},
			wantDivider: 1, wantName: "metre", wantShort: "m",
		},
		{
			name:        "Distance 1000",
			args:        args{ss: distanceSize, size: 1000},
			wantDivider: 1000, wantName: "kilometre", wantShort: "km",
		},
		{ // There is no mega-metre defined, so should return highest mapped divisor - 1000 for km
			name:        "Distance 1000000",
			args:        args{ss: distanceSize, size: 1000000},
			wantDivider: 1000, wantName: "kilometre", wantShort: "km",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDivider, gotName, gotShort := getUnit(tt.args.ss, tt.args.size)
			if gotDivider != tt.wantDivider {
				t.Errorf("getUnit() gotDivider = %v, want %v", gotDivider, tt.wantDivider)
			}
			if gotName != tt.wantName {
				t.Errorf("getUnit() gotName = %v, want %v", gotName, tt.wantName)
			}
			if gotShort != tt.wantShort {
				t.Errorf("getUnit() gotShort = %v, want %v", gotShort, tt.wantShort)
			}
		})
	}
}
