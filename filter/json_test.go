package filter

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMarshalJSON(t *testing.T) {
	type args struct {
		f Filter
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name:    "",
			args:    args{f: Eq{"f": "abc"}},
			want:    []byte(`[["f","==","abc"]]`),
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: Eq{"f": 3.14}},
			want:    []byte(`[["f","==",3.14]]`),
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: Eq{"f": 100}},
			want:    []byte(`[["f","==",100]]`),
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: Eq{"f": true}},
			want:    []byte(`[["f","==",true]]`),
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: Eq{"f": nil}},
			want:    []byte(`[["f","==",null]]`),
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: Eq{"f": time.Unix(0, 0)}},
			want:    []byte(`[["f","==","1970-01-01T03:00:00+03:00"]]`),
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: Eq{"f": time.Unix(0, 0).UTC()}},
			want:    []byte(`[["f","==","1970-01-01T00:00:00Z"]]`),
			wantErr: nil,
		},
		{
			name:    "",
			args:    args{f: Eq{"f": uuid.UUID{}}},
			want:    []byte(`[["f","==","00000000-0000-0000-0000-000000000000"]]`),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalJSON(tt.args.f)
			require.ErrorIs(t, tt.wantErr, err)
			require.Equal(t, tt.want, got)
		})
	}
}
