package stringutil

import "testing"

func Test_IsInteger(t *testing.T) {
	cases := []struct {
		args string
		want bool
	}{
		{
			args: "123",
			want: true,
		},
		{
			args: "123.1",
			want: false,
		},
		{
			args: "abc",
			want: false,
		},
		{
			args: "123abc",
			want: false,
		},
	}

	for _, c := range cases {
		got := IsInteger(c.args)
		if got != c.want {
			t.Errorf("IsInteger(%q) == %t, want %t", c.args, got, c.want)
		}
	}
}

func Test_Reverse(t *testing.T) {
	cases := []struct {
		args string
		want string
	}{
		{
			args: "123",
			want: "321",
		},
		{
			args: "123.1",
			want: "1.321",
		},
		{
			args: "abc",
			want: "cba",
		},
		{
			args: "123abc",
			want: "cba321",
		},
	}

	for _, c := range cases {
		got := Reverse(c.args)
		if got != c.want {
			t.Errorf("Reverse (%q) == %s, want %s", c.args, got, c.want)
		}
	}
}
