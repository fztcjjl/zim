package websocket

import (
	"bufio"
	"bytes"
	"github.com/panjf2000/gnet"
	"io"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
)

const (
	crlf          = "\r\n"
	colonAndSpace = ": "
	commaAndSpace = ", "
)

const (
	textHeadUpgrade = "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\n"
)

var (
	textHeadBadRequest          = statusText(http.StatusBadRequest)
	textHeadInternalServerError = statusText(http.StatusInternalServerError)
	textHeadUpgradeRequired     = statusText(http.StatusUpgradeRequired)

	textTailErrHandshakeBadProtocol   = errorText(ErrHandshakeBadProtocol)
	textTailErrHandshakeBadMethod     = errorText(ErrHandshakeBadMethod)
	textTailErrHandshakeBadHost       = errorText(ErrHandshakeBadHost)
	textTailErrHandshakeBadUpgrade    = errorText(ErrHandshakeBadUpgrade)
	textTailErrHandshakeBadConnection = errorText(ErrHandshakeBadConnection)
	textTailErrHandshakeBadSecAccept  = errorText(ErrHandshakeBadSecAccept)
	textTailErrHandshakeBadSecKey     = errorText(ErrHandshakeBadSecKey)
	textTailErrHandshakeBadSecVersion = errorText(ErrHandshakeBadSecVersion)
	textTailErrUpgradeRequired        = errorText(ErrHandshakeUpgradeRequired)
)

var (
	headerHost          = "Host"
	headerUpgrade       = "Upgrade"
	headerConnection    = "Connection"
	headerSecVersion    = "Sec-WebSocket-Version"
	headerSecProtocol   = "Sec-WebSocket-Protocol"
	headerSecExtensions = "Sec-WebSocket-Extensions"
	headerSecKey        = "Sec-WebSocket-Key"
	headerSecAccept     = "Sec-WebSocket-Accept"

	headerHostCanonical          = textproto.CanonicalMIMEHeaderKey(headerHost)
	headerUpgradeCanonical       = textproto.CanonicalMIMEHeaderKey(headerUpgrade)
	headerConnectionCanonical    = textproto.CanonicalMIMEHeaderKey(headerConnection)
	headerSecVersionCanonical    = textproto.CanonicalMIMEHeaderKey(headerSecVersion)
	headerSecProtocolCanonical   = textproto.CanonicalMIMEHeaderKey(headerSecProtocol)
	headerSecExtensionsCanonical = textproto.CanonicalMIMEHeaderKey(headerSecExtensions)
	headerSecKeyCanonical        = textproto.CanonicalMIMEHeaderKey(headerSecKey)
	headerSecAcceptCanonical     = textproto.CanonicalMIMEHeaderKey(headerSecAccept)
)

var (
	specHeaderValueUpgrade         = "websocket"
	specHeaderValueConnection      = "Upgrade"
	specHeaderValueConnectionLower = "upgrade"
	specHeaderValueSecVersion      = "13"
)

func httpWriteHeaderKey(bw *bufio.Writer, key string) {
	bw.WriteString(key)
	bw.WriteString(colonAndSpace)
}

func httpWriteResponseError(bw *bufio.Writer, err error, code int, header HandshakeHeaderFunc) {
	switch code {
	case http.StatusBadRequest:
		bw.WriteString(textHeadBadRequest)
	case http.StatusInternalServerError:
		bw.WriteString(textHeadInternalServerError)
	case http.StatusUpgradeRequired:
		bw.WriteString(textHeadUpgradeRequired)
	default:
		writeStatusText(bw, code)
	}

	// Write custom headers.
	if header != nil {
		header(bw)
	}

	switch err {
	case ErrHandshakeBadProtocol:
		bw.WriteString(textTailErrHandshakeBadProtocol)
	case ErrHandshakeBadMethod:
		bw.WriteString(textTailErrHandshakeBadMethod)
	case ErrHandshakeBadHost:
		bw.WriteString(textTailErrHandshakeBadHost)
	case ErrHandshakeBadUpgrade:
		bw.WriteString(textTailErrHandshakeBadUpgrade)
	case ErrHandshakeBadConnection:
		bw.WriteString(textTailErrHandshakeBadConnection)
	case ErrHandshakeBadSecAccept:
		bw.WriteString(textTailErrHandshakeBadSecAccept)
	case ErrHandshakeBadSecKey:
		bw.WriteString(textTailErrHandshakeBadSecKey)
	case ErrHandshakeBadSecVersion:
		bw.WriteString(textTailErrHandshakeBadSecVersion)
	case ErrHandshakeUpgradeRequired:
		bw.WriteString(textTailErrUpgradeRequired)
	case nil:
		bw.WriteString(crlf)
	default:
		writeErrorText(bw, err)
	}
}

func writeStatusText(bw *bufio.Writer, code int) {
	bw.WriteString("HTTP/1.1 ")
	bw.WriteString(strconv.Itoa(code))
	bw.WriteByte(' ')
	bw.WriteString(http.StatusText(code))
	bw.WriteString(crlf)
	bw.WriteString("Content-Type: text/plain; charset=utf-8")
	bw.WriteString(crlf)
}

func writeErrorText(bw *bufio.Writer, err error) {
	body := err.Error()
	bw.WriteString("Content-Length: ")
	bw.WriteString(strconv.Itoa(len(body)))
	bw.WriteString(crlf)
	bw.WriteString(crlf)
	bw.WriteString(body)
}

// httpError is like the http.Error with WebSocket context exception.
func httpError(w http.ResponseWriter, body string, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(code)
	w.Write([]byte(body))
}

// statusText is a non-performant status text generator.
// NOTE: Used only to generate constants.
func statusText(code int) string {
	var buf bytes.Buffer
	bw := bufio.NewWriter(&buf)
	writeStatusText(bw, code)
	bw.Flush()
	return buf.String()
}

// errorText is a non-performant error text generator.
// NOTE: Used only to generate constants.
func errorText(err error) string {
	var buf bytes.Buffer
	bw := bufio.NewWriter(&buf)
	writeErrorText(bw, err)
	bw.Flush()
	return buf.String()
}

// HandshakeHeader is the interface that writes both upgrade request or
// response headers into a given io.Writer.
type HandshakeHeader interface {
	io.WriterTo
}

// HandshakeHeaderString is an adapter to allow the use of headers represented
// by ordinary string as HandshakeHeader.
type HandshakeHeaderString string

// WriteTo implements HandshakeHeader (and io.WriterTo) interface.
func (s HandshakeHeaderString) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(s))
	return int64(n), err
}

// HandshakeHeaderBytes is an adapter to allow the use of headers represented
// by ordinary slice of bytes as HandshakeHeader.
type HandshakeHeaderBytes []byte

// WriteTo implements HandshakeHeader (and io.WriterTo) interface.
func (b HandshakeHeaderBytes) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b)
	return int64(n), err
}

// HandshakeHeaderFunc is an adapter to allow the use of headers represented by
// ordinary function as HandshakeHeader.
type HandshakeHeaderFunc func(io.Writer) (int64, error)

// WriteTo implements HandshakeHeader (and io.WriterTo) interface.
func (f HandshakeHeaderFunc) WriteTo(w io.Writer) (int64, error) {
	return f(w)
}

// HandshakeHeaderHTTP is an adapter to allow the use of http.Header as
// HandshakeHeader.
type HandshakeHeaderHTTP http.Header

// WriteTo implements HandshakeHeader (and io.WriterTo) interface.
func (h HandshakeHeaderHTTP) WriteTo(w io.Writer) (int64, error) {
	wr := writer{w: w}
	err := http.Header(h).Write(&wr)
	return wr.n, err
}

type writer struct {
	n int64
	w io.Writer
}

func (w *writer) WriteString(s string) (int, error) {
	n, err := io.WriteString(w.w, s)
	w.n += int64(n)
	return n, err
}

func (w *writer) Write(p []byte) (int, error) {
	n, err := w.w.Write(p)
	w.n += int64(n)
	return n, err
}

// ========================

type Request struct {
	Method     string
	RequestURI string
	Proto      string // "HTTP/1.0"
	ProtoMajor int    // 1
	ProtoMinor int    // 0

	Header http.Header
}

func ReadRequest(conn gnet.Conn) (req *Request, out []byte, err error) {
	buf := conn.Read()

	var index int
	if index = bytes.Index(buf, []byte("\r\n\r\n")); index == -1 {
		err = ErrShortPackaet
		return
	}
	lines := bytes.Split(buf[:index], []byte("\r\n"))
	if len(lines) == 0 {
		err = ErrShortPackaet
		return
	}

	req = new(Request)
	var ok bool
	if req.Method, req.RequestURI, req.Proto, ok = parseRequestLine(string(lines[0])); !ok {
		err = ErrMalformedRequest
	}

	if req.ProtoMajor, req.ProtoMinor, ok = parseHTTPVersion(req.Proto); !ok {
		err = ErrHandshakeBadProtocol
	}

	if len(lines) > 1 {
		req.Header, err = parseMIMEHeader(lines[1:])
	}

	conn.ShiftN(index + 4)

	if err != nil {
		var code int
		if rej, ok := err.(*rejectConnectionError); ok {
			code = rej.code
		}
		if code == 0 {
			code = http.StatusInternalServerError
		}
		var buf bytes.Buffer
		bw := bufio.NewWriter(&buf)
		httpWriteResponseError(bw, err, code, nil)
		bw.Flush()
		out = buf.Bytes()
	}

	return
}

func parseRequestLine(line string) (method, requestURI, proto string, ok bool) {
	s1 := strings.Index(line, " ")
	s2 := strings.Index(line[s1+1:], " ")
	if s1 < 0 || s2 < 0 {
		return
	}
	s2 += s1 + 1
	return line[:s1], line[s1+1 : s2], line[s2+1:], true
}

func parseHTTPVersion(vers string) (major, minor int, ok bool) {
	const Big = 1000000 // arbitrary upper bound
	switch vers {
	case "HTTP/1.1":
		return 1, 1, true
	case "HTTP/1.0":
		return 1, 0, true
	}
	if !strings.HasPrefix(vers, "HTTP/") {
		return 0, 0, false
	}
	dot := strings.Index(vers, ".")
	if dot < 0 {
		return 0, 0, false
	}
	major, err := strconv.Atoi(vers[5:dot])
	if err != nil || major < 0 || major > Big {
		return 0, 0, false
	}
	minor, err = strconv.Atoi(vers[dot+1:])
	if err != nil || minor < 0 || minor > Big {
		return 0, 0, false
	}
	return major, minor, true
}

func parseHeaderLine(line []byte) (k, v []byte, ok bool) {
	colon := bytes.IndexByte(line, ':')
	if colon == -1 {
		return
	}

	k = btrim(line[:colon])
	// TODO(gobwas): maybe use just lower here?
	canonicalizeHeaderKey(k)

	v = btrim(line[colon+1:])

	return k, v, true
}

func parseMIMEHeader(lines [][]byte) (header http.Header, err error) {
	header = make(http.Header)

	for i := 0; i < len(lines); i++ {
		if len(lines[i]) == 0 {
			break
		}

		k, v, ok := parseHeaderLine(lines[i])
		if !ok {
			err = ErrMalformedRequest
			break
		}

		header.Add(string(k), string(v))
	}

	return
}
