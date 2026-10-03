// Package mute 是禁言：违规次数递增时长、封顶、解除与请求幂等。
// 缺陷：时长线性递增不是翻倍、不封顶、解除不生效、请求不幂等、
// 每次查询遍历全部玩家。
package mute

// Mutes 是禁言台账。
type Mutes struct {
	baseMs   int64
	capMs    int64
	offenses map[string]int
	until    map[string]int64
	memo     map[string]bool
	scanned  int
}

// New 建台账。
func New(baseMs, capMs int64) *Mutes {
	return &Mutes{baseMs: baseMs, capMs: capMs, offenses: map[string]int{},
		until: map[string]int64{}, memo: map[string]bool{}}
}

// Mute 记一次违规并返回禁言时长。
func (m *Mutes) Mute(player string, nowMs int64, reqID string) int64 {
	m.scanned += len(m.until)
	m.memo[reqID] = true
	m.offenses[player]++
	duration := m.baseMs * int64(m.offenses[player])
	m.until[player] = nowMs + duration
	return duration
}

// Lift 解除当前禁言。
func (m *Mutes) Lift(player string, nowMs int64, reqID string) bool {
	return true
}

// Muted 判断当前是否被禁言。
func (m *Mutes) Muted(player string, nowMs int64) bool {
	m.scanned += len(m.until)
	if until, ok := m.until[player]; ok {
		return nowMs < until
	}
	return false
}

// Offenses 是累计违规次数。
func (m *Mutes) Offenses(player string) int { return m.offenses[player] }

// Scanned 是累计扫描量（规模观测）。
func (m *Mutes) Scanned() int { return m.scanned }
