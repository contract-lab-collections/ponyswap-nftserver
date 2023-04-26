package bctools

import (
	"testing"
)

func TestEip712V4Sign(t *testing.T) {
	const privKey = "dd0927dbac09fb670433cd3adc302788763166609e5807c68fd37b83542097f2"
	const domain = "0x5bd49e0db6cf393076fe141d41e36a8eaf3013b2264eff33410abb11c6b7d2a3"

	type args struct {
		privKey         string
		domainSeparator string
		structHash      string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"1) local signature", args{
			privKey, domain,
			"0x9c0bdcc036090dd65b5ab4241dbf17f2001d1f6933c3ad2a7875f8eae2091d17"},
			"75c4483354cab17e3e5eb64b5fd7119d319d0d244a38afd5a1febf48be67b8015f4c169beec517d0e595bf8d40085806ce4e15ca60ce812bfd1d02a732a3398600",
			false},
		{"2) local signature", args{
			privKey, domain,
			"0x5a9bb6922c4d635e8d5106cd1a028ac851bb26ac7117600b4a34eacefdb4c63a"},
			"9435efaa92bdaa2821e22c17488de08eff1271005e8215e410efe3731dc68e910a51a27fa46efdb2dcab1477975065f165ea6aaa79b6e7e2ab3266ccf193f87100",
			// 591a94ff989b5cf22d523fbebb501b1b9799935144e5f2ed0eb136556c7bb7e15cd47fffdd760b0bf0e809a5d0eb165c1b2f15c2040e5877fb12c4523e1a3c441c
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Eip712V4Sign(tt.args.privKey, tt.args.domainSeparator, tt.args.structHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("Eip712V4Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Eip712V4Sign() = %v, want %v", got, tt.want)
			}
		})
	}
}
