// Package mute 是禁言：违规次数递增时长、封顶、解除与请求幂等。
package mute

// Mutes 是禁言台账。
type Mutes struct {
	baseMs   int64
	capMs    int64
	offenses map[string]int
	until    map[string]int64
	memo     map[string]int64
	scanned  int
}

// New 建台账。
func New(baseMs, capMs int64) *Mutes {
	return &Mutes{baseMs: baseMs, capMs: capMs, offenses: map[string]int{},
		until: map[string]int64{}, memo: map[string]int64{}}
}

// Mute 记一次违规并返回禁言时长；同一个 reqID 只记一次。
func (m *Mutes) Mute(player string, nowMs int64, reqID string) int64 {
	if duration, ok := m.memo[reqID]; ok {
		return duration
	}
	m.offenses[player]++
	duration := m.baseMs
	for n := 1; n < m.offenses[player] && duration < m.capMs; n++ {
		if duration > m.capMs/2 {
			duration = m.capMs
			break
		}
		duration *= 2
	}
	if duration > m.capMs {
		duration = m.capMs
	}
	m.until[player] = nowMs + duration
	m.memo[reqID] = duration
	return duration
}

// Lift 解除当前禁言但保留违规次数；没有禁言时返回 false。
func (m *Mutes) Lift(player string, nowMs int64, reqID string) bool {
	until, ok := m.until[player]
	if ok && nowMs < until {
		delete(m.until, player)
		return true
	}
	return false
}

// Muted 判断当前是否被禁言。
func (m *Mutes) Muted(player string, nowMs int64) bool {
	until, ok := m.until[player]
	return ok && nowMs < until
}

// Offenses 是累计违规次数。
func (m *Mutes) Offenses(player string) int { return m.offenses[player] }

// Scanned 是累计扫描量（规模观测）。
func (m *Mutes) Scanned() int { return m.scanned }
