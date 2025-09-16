package stats_counter

import (
	"ddd-timer-service/models"
	"encoding/json"
	"fmt"
	"time"
)

type Stats struct {
	ok bool

	startTime time.Time
	endTime   time.Time

	now    time.Time
	left   time.Duration // Прошло наносекунд
	passed time.Duration // Осталось наносекунд
}

func NewStats(u *models.User, now time.Time) (*Stats, error) {
	if u.ServeFrom.IsZero() {
		return nil, fmt.Errorf("NewStats: u.startTime time is zero")
	}

	if u.ServeTo.IsZero() {
		return nil, fmt.Errorf("NewStats: u.endTime time is zero")
	}

	ss := &Stats{
		ok:        true,
		startTime: u.ServeFrom,
		endTime:   u.ServeTo,
		now:       now,
		left:      u.ServeTo.Sub(now),
		passed:    now.Sub(u.ServeFrom),
	}

	return ss, nil
}

func (s *Stats) LeftHours() float64 {
	return s.left.Hours()
}

func (s *Stats) PassedHours() float64 {
	return s.passed.Hours()
}

func (s *Stats) LeftDays() float64 {
	return s.LeftHours() / 24
}

func (s *Stats) PassedDays() float64 {
	return s.PassedHours() / 24
}

func (s *Stats) LeftWeeks() float64 {
	return s.LeftDays() / 7
}

func (s *Stats) PassedWeeks() float64 {
	return s.PassedDays() / 7
}

func (s *Stats) LeftPercents() float64 {
	if !s.ok {
		return 0.0
	}

	total := s.left + s.passed

	return float64(s.left) / float64(total) * 100
}

func (s *Stats) PassedPercents() float64 {
	if !s.ok {
		return 0.0
	}

	total := s.left + s.passed

	return float64(s.passed) / float64(total) * 100
}

func (s *Stats) MarshalJSON() ([]byte, error) {
	tmp := map[string]interface{}{
		"inDate":         s.startTime.String(),
		"outDate":        s.endTime.String(),
		"percentsLeft":   s.LeftPercents(),
		"percentsPassed": s.PassedPercents(),
	}

	return json.Marshal(tmp)
}

func (s *Stats) StringJSON() string {
	b, _ := json.MarshalIndent(s, "", "\t")
	return string(b)
}

func (s *Stats) PrettyShort() string {
	return fmt.Sprintf(tmplPrettyShort,
		s.PassedHours(), s.PassedDays(), s.PassedWeeks(), s.PassedPercents(),
		s.LeftHours(), s.LeftDays(), s.LeftWeeks(), s.LeftPercents())
}
