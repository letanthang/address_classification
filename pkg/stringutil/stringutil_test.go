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
