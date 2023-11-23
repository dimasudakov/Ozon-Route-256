package main

import (
	"reflect"
	"testing"
)

func Test_putApples(t *testing.T) {
	type args struct {
		apples      []Apple
		boxCapacity int
	}
	tests := []struct {
		name string
		args args
		want []Car
	}{
		{
			name: "Test Case 1",
			args: args{
				apples: []Apple{
					{Id: 1}, {Id: 2}, {Id: 3}, {Id: 4},
					{Id: 5}, {Id: 6}, {Id: 7}, {Id: 8},
					{Id: 9}, {Id: 10}, {Id: 11}, {Id: 12},
					{Id: 13}, {Id: 14}, {Id: 15}, {Id: 16},
					{Id: 17}, {Id: 18}, {Id: 19}, {Id: 20},
				},
				boxCapacity: 4,
			},
			want: []Car{
				{
					Id: 1,
					Boxes: []Box{
						{Id: 1, Apples: []Apple{{Id: 20}, {Id: 19}, {Id: 18}, {Id: 17}}},
						{Id: 3, Apples: []Apple{{Id: 12}, {Id: 11}, {Id: 10}, {Id: 9}}},
						{Id: 5, Apples: []Apple{{Id: 4}, {Id: 3}, {Id: 2}, {Id: 1}}},
					},
				},
				{
					Id: 2,
					Boxes: []Box{
						{Id: 2, Apples: []Apple{{Id: 16}, {Id: 15}, {Id: 14}, {Id: 13}}},
						{Id: 4, Apples: []Apple{{Id: 8}, {Id: 7}, {Id: 6}, {Id: 5}}},
					},
				},
			},
		},
		{
			name: "Test Case 2",
			args: args{
				apples: []Apple{
					{Id: 1}, {Id: 2}, {Id: 3}, {Id: 4}, {Id: 5},
				},
				boxCapacity: 100,
			},
			want: []Car{
				{
					Id: 1,
					Boxes: []Box{
						{Id: 1, Apples: []Apple{{Id: 5}, {Id: 4}, {Id: 3}, {Id: 2}, {Id: 1}}},
					},
				},
				{
					Id:    2,
					Boxes: nil,
				},
			},
		},
		{
			name: "Test case 3",
			args: args{
				apples: []Apple{
					{1}, {2}, {3}, {4}, {5}, {6}, {7},
				},
				boxCapacity: 1,
			},
			want: []Car{
				{
					Id: 1,
					Boxes: []Box{
						{Id: 1, Apples: []Apple{{7}}},
						{Id: 3, Apples: []Apple{{5}}},
						{Id: 5, Apples: []Apple{{3}}},
						{Id: 7, Apples: []Apple{{1}}},
					},
				},
				{
					Id: 2,
					Boxes: []Box{
						{Id: 2, Apples: []Apple{{6}}},
						{Id: 4, Apples: []Apple{{4}}},
						{Id: 6, Apples: []Apple{{2}}},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := putApples(tt.args.apples, tt.args.boxCapacity); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("\nresult:\t %v,\n want:\t %v", got, tt.want)
				}
			},
		)
	}
}
