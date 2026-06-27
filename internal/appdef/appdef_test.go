package appdef

import (
	"encoding/json"
	"net/http"
	"testing"
)

// 프리셋 문자열(기존)과 세부 권한 객체(신규)가 모두 같은 결과로 파싱되는지.
func TestParsePresetAndGranular(t *testing.T) {
	raw := `{"name":"t","collections":{
		"resp":"submissions",
		"board":"public",
		"secret":"private",
		"notice":{"read":true,"write":false,"edit":false},
		"full":{"read":true,"write":true,"edit":true}
	}}`
	var d Def
	if err := json.Unmarshal([]byte(raw), &d); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := d.validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
	want := map[string]CollectionPerm{
		"resp":   {Write: true},                          // submissions == write-only
		"board":  {Read: true, Write: true, Edit: true},  // public == all
		"secret": {},                                      // private == none
		"notice": {Read: true},                            // granular read-only
		"full":   {Read: true, Write: true, Edit: true},
	}
	for name, w := range want {
		if got := d.Collections[name]; got != w {
			t.Errorf("%s: got %+v want %+v", name, got, w)
		}
	}
}

// verb 집행: read=GET, write=POST, edit=PATCH/DELETE.
func TestPublicAllows(t *testing.T) {
	d := &Def{Collections: map[string]CollectionPerm{
		"sub":    {Write: true},
		"notice": {Read: true},
		"pub":    {Read: true, Write: true, Edit: true},
		"none":   {},
	}}
	cases := []struct {
		col, method string
		want        bool
	}{
		{"sub", http.MethodPost, true},
		{"sub", http.MethodGet, false},   // 제출만 — 읽기 금지
		{"sub", http.MethodPatch, false}, // 편집 금지
		{"notice", http.MethodGet, true},
		{"notice", http.MethodPost, false}, // 읽기 전용 — 쓰기 금지
		{"pub", http.MethodDelete, true},
		{"none", http.MethodGet, false},
		{"none", http.MethodPost, false},
	}
	for _, c := range cases {
		got, declared := d.PublicAllows(c.col, c.method)
		if !declared {
			t.Errorf("%s 미선언으로 나옴", c.col)
		}
		if got != c.want {
			t.Errorf("PublicAllows(%s,%s)=%v want %v", c.col, c.method, got, c.want)
		}
	}
	// 미선언 컬렉션
	if _, declared := d.PublicAllows("ghost", http.MethodGet); declared {
		t.Error("미선언 컬렉션이 declared=true")
	}
}

// 잘못된 프리셋 문자열은 거부.
func TestBadPreset(t *testing.T) {
	var d Def
	err := json.Unmarshal([]byte(`{"collections":{"x":"open"}}`), &d)
	if err == nil {
		t.Fatal("잘못된 프리셋이 통과됨")
	}
}

// Label: 프리셋 일치 시 친화명, 아니면 조합.
func TestLabel(t *testing.T) {
	cases := map[CollectionPerm]string{
		{Write: true}:                         "제출 받기",
		{Read: true, Write: true, Edit: true}: "공개",
		{}:                                     "나만(차단)",
		{Read: true}:                           "읽기",
		{Read: true, Edit: true}:               "읽기·편집",
	}
	for p, want := range cases {
		if got := p.Label(); got != want {
			t.Errorf("Label(%+v)=%q want %q", p, got, want)
		}
	}
}
