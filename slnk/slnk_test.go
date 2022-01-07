// package slnk
package slnk

import (
	"errors"
	"os"
	"reflect"
	"testing"
)

func TestNewSlist(t *testing.T) {
	tests := []struct {
		name string
		want *slist
	}{
		// test cases.
		{
			name: "NewSlist passing",
			want: &slist{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSlist(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSlist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_slist_Error(t *testing.T) {
	type fields struct {
		head *node
		tail *node
		curr *node
		cnt  int
		cmp  func(interface{}, interface{}) int
		err  error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// test cases
		{
			name:    "Error passing",
			fields:  fields{},
			wantErr: false,
		},
		{
			name:    "Error failing",
			fields:  fields{err: errors.New("some error")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &slist{
				head: tt.fields.head,
				tail: tt.fields.tail,
				curr: tt.fields.curr,
				cnt:  tt.fields.cnt,
				cmp:  tt.fields.cmp,
				err:  tt.fields.err,
			}
			if err := s.Error(); (err != nil) != tt.wantErr {
				t.Errorf("slist.Error() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_slist_SetCompareFunc(t *testing.T) {
	fu := func(interface{}, interface{}) int { return 0 }
	s := &slist{cmp: fu}

	type args struct {
		f func(interface{}, interface{}) int
	}
	tests := []struct {
		name string
		args args
		want *slist
	}{
		// test cases
		{
			name: "SetCompareFunc passing",
			args: args{f: fu},
			want: s,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.SetCompareFunc(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("slist.SetCompareFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_slist_Len(t *testing.T) {
	tests := []struct {
		name string
		s    *slist
		want int
	}{
		// test cases
		{
			name: "Len passing",
			s:    &slist{cnt: 0},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Len(); got != tt.want {
				t.Errorf("slist.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_slist_Add_One(t *testing.T) {
	s := NewSlist()
	s.Add(1)
	if s == nil {
		t.Errorf("slist.Add() returned nil")
	}
	if err := s.Error(); err != nil {
		t.Errorf("slist.Add() returned error %v", err)
	}
	if len := s.Len(); len != 1 {
		t.Errorf("slist.Add() returned len %d", len)
	}
}

func Test_slist_Add_Five(t *testing.T) {
	s := NewSlist()
	for i, v := range "12345" { // do 5 times
		s.Add(v)
		if s == nil {
			t.Errorf("slist.Add() returned nil")
		}
		if err := s.Error(); err != nil {
			t.Errorf("slist.Add() returned error: %v", err)
		}
		if len := s.Len(); len != i+1 {
			t.Errorf("slist.Add() returned len: %d wanted: %d", len, i)
		}
	}
	s.Debug(true).Print(os.Stdout)
}

func Test_slist_Size(t *testing.T) {
	type fields struct {
		head  *node
		tail  *node
		curr  *node
		cnt   int
		cmp   func(interface{}, interface{}) int
		err   error
		debug bool
		start int
		end   int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// test cases
		{
			name:   "Size passing",
			fields: fields{cnt: 1, start: 0, end: 1, cmp: nil, debug: false, err: nil},
			want:   80,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &slist{
				head:  tt.fields.head,
				tail:  tt.fields.tail,
				curr:  tt.fields.curr,
				cnt:   tt.fields.cnt,
				cmp:   tt.fields.cmp,
				err:   tt.fields.err,
				debug: tt.fields.debug,
				start: tt.fields.start,
				end:   tt.fields.end,
			}
			if got := s.Size(); got != tt.want {
				t.Errorf("slist.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
