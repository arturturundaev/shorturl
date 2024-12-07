package service

import "testing"

type MockRepository struct {
}

func (r MockRepository) GetUrlsCount() int32 {
	return 1
}
func (r MockRepository) GetUsersCount() int32 {
	return 1
}
func TestStatService_GetUrlsAndUsersStat(t *testing.T) {
	rep := new(MockRepository)

	tests := []struct {
		name  string
		want  int32
		want1 int32
	}{
		{
			name:  "success",
			want:  1,
			want1: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStatService(rep)
			got, got1 := s.GetUrlsAndUsersStat()
			if got != tt.want {
				t.Errorf("GetUrlsAndUsersStat() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetUrlsAndUsersStat() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
