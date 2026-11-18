package imapwire_test

import (
	"bufio"
	"strings"
	"testing"

	"github.com/emersion/go-imap/v2/internal/imapwire"
)

func newTestDecoder(s string) *imapwire.Decoder {
	return imapwire.NewDecoder(bufio.NewReader(strings.NewReader(s)), imapwire.ConnSideServer)
}

func TestExpectMailboxInvalidUTF8(t *testing.T) {
	for _, s := range []string{"\"bad\xffname\"", "\"a\x7fb\"", "{3}\r\na\x01b"} {
		dec := newTestDecoder(s + "\r\n")
		dec.MailboxUTF8 = true

		var name string
		if dec.ExpectMailbox(&name) {
			t.Errorf("ExpectMailbox(%q) accepted an invalid name", s)
		}
	}
}

func TestExpectMailbox(t *testing.T) {
	tests := []struct {
		s           string
		mailboxUTF8 bool
		want        string
	}{
		{s: "R&-D", want: "R&D"},
		{s: "Caf&AOg-", want: "Cafè"},
		{s: "R&D", mailboxUTF8: true, want: "R&D"},
		{s: `"Cafè"`, mailboxUTF8: true, want: "Cafè"},
		{s: "R&-D", mailboxUTF8: true, want: "R&-D"},
		{s: "inbox", mailboxUTF8: true, want: "INBOX"},
	}
	for _, tc := range tests {
		dec := newTestDecoder(tc.s + "\r\n")
		dec.MailboxUTF8 = tc.mailboxUTF8

		var name string
		if !dec.ExpectMailbox(&name) {
			t.Errorf("ExpectMailbox(%q) = %v", tc.s, dec.Err())
		} else if name != tc.want {
			t.Errorf("ExpectMailbox(%q) = %q, want %q", tc.s, name, tc.want)
		}
	}
}
