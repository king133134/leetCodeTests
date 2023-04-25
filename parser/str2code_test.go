package parser

import (
	"reflect"
	"testing"
)

func TestCreateTests(t *testing.T) {
	type args struct {
		code string
		con  *string
	}
	tests := []struct {
		name string
		args args
		want *Tests
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateTests(tt.args.code, tt.args.con); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateTests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toListNode(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"test1", args{str: "[1,2,3]"}, []byte("&ListNode{Val:1,Next:&ListNode{Val:2,Next:&ListNode{Val:3,Next:nil}}}")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toListNode(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toListNode() = %v, want %v", got, tt.want)
			}
		})
	}
}
