package ui

import (
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/app"
)

// prefixTree groups keys on "/" boundaries. Branch IDs end with "/", leaf IDs are the full key.
func prefixTree(sess *app.Session, onSelect func(key []byte)) fyne.CanvasObject {
	metas, err := sess.Keys()
	if err != nil {
		fyne.LogError("load keys for tree", err)
	}
	keys := make([]string, 0, len(metas))
	for _, m := range metas {
		keys = append(keys, m.Key)
	}
	children := buildPrefixTree(keys)

	t := widget.NewTree(
		func(uid widget.TreeNodeID) []widget.TreeNodeID {
			return children[uid]
		},
		func(uid widget.TreeNodeID) bool {
			return uid == "" || strings.HasSuffix(uid, "/")
		},
		func(branch bool) fyne.CanvasObject {
			l := widget.NewLabel("")
			l.TextStyle = fyne.TextStyle{Monospace: true}
			l.Truncation = fyne.TextTruncateEllipsis
			return l
		},
		func(uid widget.TreeNodeID, branch bool, c fyne.CanvasObject) {
			l := c.(*widget.Label)
			l.TextStyle = fyne.TextStyle{Monospace: true}
			l.SetText(treeLabel(uid))
		},
	)

	t.OnSelected = func(uid widget.TreeNodeID) {
		if strings.HasSuffix(uid, "/") || onSelect == nil {
			return
		}
		onSelect([]byte(uid))
	}

	return t
}

func buildPrefixTree(keys []string) map[string][]string {
	dedup := map[string]map[string]struct{}{"": {}}

	for _, k := range keys {
		parts := strings.Split(k, "/")
		parent := ""
		for i, p := range parts {
			_ = p
			isLast := i == len(parts)-1
			var id string
			if isLast {
				id = k
			} else {
				id = strings.Join(parts[:i+1], "/") + "/"
				if dedup[id] == nil {
					dedup[id] = map[string]struct{}{}
				}
			}
			if dedup[parent] == nil {
				dedup[parent] = map[string]struct{}{}
			}
			dedup[parent][id] = struct{}{}
			parent = id
		}
	}

	out := make(map[string][]string, len(dedup))
	for parent, set := range dedup {
		list := make([]string, 0, len(set))
		for id := range set {
			list = append(list, id)
		}
		sort.Strings(list)
		out[parent] = list
	}
	return out
}

func treeLabel(id string) string {
	if id == "" {
		return ""
	}
	if strings.HasSuffix(id, "/") {
		trimmed := strings.TrimSuffix(id, "/")
		if idx := strings.LastIndex(trimmed, "/"); idx >= 0 {
			return trimmed[idx+1:] + "/"
		}
		return trimmed + "/"
	}
	if idx := strings.LastIndex(id, "/"); idx >= 0 {
		return id[idx+1:]
	}
	return id
}
