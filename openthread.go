package main

import (
	"strings"

	"github.com/ThomasHabets/cmdg/cmdglib"
	"github.com/ThomasHabets/cmdg/ncwrap"
	gc "github.com/rthornton128/goncurses"
	gmail "google.golang.org/api/gmail/v1"
)

// Return true if cmdg should quit.
func openThreadMain(ts []*gmail.Thread, current int, marked map[string]bool, currentLabel string) bool {
	nc.Status("Opening thread")
	scroll := 0
	for {
		nc.ApplyMain(func(w *gc.Window) {
			openThreadPrint(w, ts, current, marked[ts[current].Id], currentLabel, scroll)
		})
		key := <-nc.Input
		switch key {
		case '?':
			helpWin(`q                 Quit
^P                Previous thread
^N                Next thread
Left, <, u        Exit thread
p, k              Previous message
n, j              Next message
Right, >, Enter   Open message
f                 Forward last
r                 Reply to last
a                 Reply all to last
e                 Archive
l                 Add label
L                 Remove label
x                 Mark thread (TODO)
Space             Page down
Backspace         Page up
`)
			nc.ApplyMain(func(w *gc.Window) { w.Clear() })
		case 'q':
			return true
		case gc.KEY_LEFT, '<', 'u':
			return false
		case gc.KEY_RIGHT, '>', '\n':
			// TODO
		case 'p', 'k':
			if current > 0 {
				current--
			}
		case 'n', 'j':
			if current < len(ts)-1 {
				current++
			}
		default:
			nc.Status("unknown key: %v", gc.KeyString(key))
		}
	}
	return false
}

func openThreadPrint(w *gc.Window, ts []*gmail.Thread, current int, marked bool, currentLabel string, scroll int) {
	t := ts[current]
	w.Move(0, 0)
	//height, width := w.MaxYX()
	tswidth := 7

	ncwrap.ColorPrint(w, "Thread: [bold]%s[unbold] (%d messages)\n", cmdglib.GetHeader(t.Messages[0], "Subject"), len(t.Messages))
	for n, m := range t.Messages {
		ncwrap.ColorPrint(w, "[green]%*.*s - %s\n", tswidth, tswidth, cmdglib.TimeString(m), cmdglib.GetHeader(m, "From"))

		if cmdglib.HasLabel(m.LabelIds, cmdglib.Unread) || n == len(t.Messages)-1 {
			bodyLines := breakLines(strings.Split(getBody(m), "\n"))
			body := strings.Join(bodyLines, "\n")
			ncwrap.ColorPrint(w, "%s\n", body)
		}
	}
}
