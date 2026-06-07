package ui

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	fynetheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/mohsenm4/kv-explorer/internal/i18n"
)

// Read-only tree view; falls back to textBody on invalid JSON.
func jsonTreeBody(v fyne.ThemeVariant, value []byte) (fyne.CanvasObject, func() ([]byte, error), func()) {
	dec := json.NewDecoder(strings.NewReader(string(value)))
	dec.UseNumber()
	var root any
	if err := dec.Decode(&root); err != nil {
		return textBody(value)
	}

	// Decoded manually because map[string]any loses field order and the tree would shuffle between renders.
	type orderedObj struct {
		keys []string
		vals map[string]any
	}

	var decodeOrdered func(d *json.Decoder) (any, error)
	decodeOrdered = func(d *json.Decoder) (any, error) {
		tok, err := d.Token()
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case json.Delim:
			switch t {
			case '{':
				obj := orderedObj{vals: map[string]any{}}
				for d.More() {
					kt, err := d.Token()
					if err != nil {
						return nil, err
					}
					k := kt.(string)
					val, err := decodeOrdered(d)
					if err != nil {
						return nil, err
					}
					obj.keys = append(obj.keys, k)
					obj.vals[k] = val
				}
				if _, err := d.Token(); err != nil {
					return nil, err
				}
				return obj, nil
			case '[':
				var arr []any
				for d.More() {
					val, err := decodeOrdered(d)
					if err != nil {
						return nil, err
					}
					arr = append(arr, val)
				}
				if _, err := d.Token(); err != nil { // consume ']'
					return nil, err
				}
				return arr, nil
			}
		}
		return tok, nil
	}

	dec2 := json.NewDecoder(strings.NewReader(string(value)))
	dec2.UseNumber()
	ordered, err := decodeOrdered(dec2)
	if err != nil {
		return textBody(value)
	}

	type node struct {
		label    string
		children []widget.TreeNodeID
		branch   bool
	}
	nodes := map[widget.TreeNodeID]node{}

	var walk func(id widget.TreeNodeID, key string, val any)
	walk = func(id widget.TreeNodeID, key string, val any) {
		prefix := ""
		if key != "" {
			prefix = key + ": "
		}
		switch x := val.(type) {
		case orderedObj:
			var children []widget.TreeNodeID
			for _, k := range x.keys {
				cid := id + "/" + k
				children = append(children, cid)
				walk(cid, k, x.vals[k])
			}
			nodes[id] = node{
				label:    prefix + fmt.Sprintf("{%d}", len(x.keys)),
				children: children,
				branch:   true,
			}
		case []any:
			var children []widget.TreeNodeID
			for i, item := range x {
				cid := id + "/" + strconv.Itoa(i)
				children = append(children, cid)
				walk(cid, "["+strconv.Itoa(i)+"]", item)
			}
			nodes[id] = node{
				label:    prefix + fmt.Sprintf("[%d]", len(x)),
				children: children,
				branch:   true,
			}
		default:
			nodes[id] = node{label: prefix + formatLeaf(x), branch: false}
		}
	}
	walk("", "", ordered)

	muted := themeColor(v, fynetheme.ColorNamePlaceHolder)

	tree := widget.NewTree(
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			return nodes[id].children
		},
		func(id widget.TreeNodeID) bool {
			return nodes[id].branch
		},
		func(branch bool) fyne.CanvasObject {
			l := widget.NewLabel("")
			l.TextStyle = fyne.TextStyle{Monospace: true}
			return l
		},
		func(id widget.TreeNodeID, branch bool, c fyne.CanvasObject) {
			c.(*widget.Label).SetText(nodes[id].label)
		},
	)
	tree.OpenBranch("")

	hint := canvas.NewText(i18n.T("editor.jsonHint"), muted)
	hint.TextSize = 11

	body := container.NewBorder(nil, container.NewPadded(hint), nil, nil, tree)

	current := func() ([]byte, error) { return value, nil }
	reset := func() {}
	return body, current, reset
}

func formatLeaf(v any) string {
	switch x := v.(type) {
	case nil:
		return "null"
	case string:
		return strconv.Quote(x)
	case bool:
		if x {
			return "true"
		}
		return "false"
	case json.Number:
		return x.String()
	default:
		return fmt.Sprintf("%v", x)
	}
}
