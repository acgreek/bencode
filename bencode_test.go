package bencode

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestParseToMap(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Item
		wantErr error
	}{
		{"string", args{[]byte("3:foo")}, Item{Str: "foo"}, nil},
		{"integer", args{[]byte("i400e")}, Item{Num: 400}, nil},
		{"negative integer", args{[]byte("i-99934e")}, Item{Num: -99934}, nil},
		{"list", args{[]byte("li-99934e3:fooe")}, Item{List: []Item{{Num: -99934}, {Str: "foo"}}}, nil},
		{"dictionary", args{[]byte("d3:bari-99934e4:rats3:fooe")}, Item{Dict: map[string]Item{"bar": {Num: -99934}, "rats": {Str: "foo"}}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := ParseToMap(tt.args.in); !reflect.DeepEqual(got, tt.want) || err != tt.wantErr {
				t.Errorf("ParseToMap(%s) = %v, want %v; error is %v, wantError %v", tt.name, got, tt.want, err, tt.wantErr)
			}
		})
	}
}

func TestParseToMapReadSample(t *testing.T) {
	torrent, err := ioutil.ReadFile("samples/ubuntu-22.04.1-desktop-amd64.iso.torrent")
	if err != nil {
		t.Errorf("failed to open sample file: %v", err)
	}
	item, err := ParseToMap(torrent)
	if err != nil {
		t.Errorf("failed to parse torrent: %v", err)
	}

	if item.Dict == nil {
		t.Errorf("item.Dict  should not be nil")
	}
}
