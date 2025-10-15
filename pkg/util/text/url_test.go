package text

import "testing"

func TestSafeUnquote(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"test1",
			args{s: "name%3D%E5%BC%A0%E4%B8%89"},
			"name=张三",
		},
		{
			"test2",
			args{s: "email%3Dtest%40example.com"},
			"email=test@example.com",
		},
		{
			"test3",
			args{s: "message%3DHello%20World%21"},
			"message=Hello World!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SafeUnquote(tt.args.s); got != tt.want {
				t.Errorf("SafeUnquote() = %v, want %v", got, tt.want)
			}
		})
	}
}
