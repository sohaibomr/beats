package decoder

// import (
// 	"fmt"
// 	"strconv"

// 	"github.com/elastic/beats/libbeat/common"
// 	"github.com/elastic/beats/libbeat/common/streambuf"
// 	"github.com/elastic/beats/libbeat/logp"
// 	"github.com/elastic/beats/packetbeat/protos/tls"
// )

// type st struct {
// 	// Buf provides the buffering with parsing support
// 	Buf streambuf.Buffer
// }

// type parserResult int8
// type handshakeType uint8

// const (
// 	resultOK parserResult = iota
// 	resultFailed
// 	resultMore
// 	resultEncrypted
// )

// const (
// 	maxTLSRecordLength  = (1 << 14) + 2048
// 	maxHandshakeSize    = 1 << 16
// 	recordHeaderSize    = 5
// 	recordTypeHandshake = 22
// 	handshakeHeaderSize = 4
// 	helloHeaderLength   = 7
// 	randomDataLength    = 28
// 	clientHello         = 1
// )

// type recordType uint8

// type tlsVersion struct {
// 	major, minor uint8
// }

// type recordHeader struct {
// 	recordType recordType
// 	version    tlsVersion
// 	length     uint16
// }
// type handshakeHeader struct {
// 	handshakeType handshakeType
// 	length        int
// }

// type helloMsg struct {
// 	extensions tls.Extensions
// }

// type bufferView struct {
// 	buf         *streambuf.Buffer
// 	base, limit int
// }

// type parserHello struct {
// 	// Buffer to accumulate records until a full handshake message
// 	// is received
// 	handshakeBuf streambuf.Buffer
// 	hello        *helloMsg
// }

// func (header *recordHeader) String() string {
// 	return fmt.Sprintf("recordHeader type[%v] version[%v] length[%d]",
// 		header.recordType, header.version, header.length)
// }

// func readRecordHeader(buf *streambuf.Buffer) (*recordHeader, error) {
// 	var (
// 		header recordHeader
// 		err    error
// 		record uint8
// 	)
// 	if record, err = buf.ReadNetUint8At(0); err != nil {
// 		return nil, err
// 	}
// 	header.recordType = recordType(record)
// 	if header.version.major, err = buf.ReadNetUint8At(1); err != nil {
// 		return nil, err
// 	}
// 	if header.version.minor, err = buf.ReadNetUint8At(2); err != nil {
// 		return nil, err
// 	}
// 	if header.length, err = buf.ReadNetUint16At(3); err != nil {
// 		return nil, err
// 	}
// 	return &header, nil
// }

// func (header *recordHeader) isValid() bool {
// 	return header.version.major == 3 && header.length <= maxTLSRecordLength
// }

// func readHandshakeHeader(buf *streambuf.Buffer) (*handshakeHeader, error) {
// 	var err error
// 	var len8, typ uint8
// 	var len16 uint16
// 	if typ, err = buf.ReadNetUint8At(0); err != nil {
// 		return nil, err
// 	}
// 	if len8, err = buf.ReadNetUint8At(1); err != nil {
// 		return nil, err
// 	}
// 	if len16, err = buf.ReadNetUint16At(2); err != nil {
// 		return nil, err
// 	}
// 	return &handshakeHeader{handshakeType(typ),
// 		int(len16) | (int(len8) << 16)}, nil
// }

// func (parser *parserHello) bufferHandshake(buf *streambuf.Buffer, length int) error {
// 	// TODO: parse in-place if message in received buffer is complete
// 	if err := parser.handshakeBuf.Append(buf.Bytes()[recordHeaderSize : recordHeaderSize+length]); err != nil {
// 		logp.Warn("failed appending to buffer: %v", err)
// 		// Discard buffer
// 		parser.handshakeBuf.Init(nil, false)
// 		return err
// 	}
// 	for parser.handshakeBuf.Avail(handshakeHeaderSize) {
// 		// type
// 		header, err := readHandshakeHeader(&parser.handshakeBuf)
// 		if err != nil {
// 			logp.Warn("read failed: %v", err)
// 			parser.handshakeBuf.Init(nil, false)
// 			return err
// 		}
// 		if header.length > maxHandshakeSize {
// 			// Discard buffer
// 			parser.handshakeBuf.Init(nil, false)
// 			return fmt.Errorf("message too large (%d bytes)", header.length)
// 		}
// 		limit := handshakeHeaderSize + header.length
// 		if limit > parser.handshakeBuf.Len() {
// 			break
// 		}
// 		if !parser.parseHandshake(header.handshakeType,
// 			bufferView{&parser.handshakeBuf, handshakeHeaderSize, limit}) {
// 			parser.handshakeBuf.Advance(limit)
// 			return fmt.Errorf("bad handshake %+v", header)
// 		}
// 		parser.handshakeBuf.Advance(limit)
// 	}
// 	if parser.handshakeBuf.Len() == 0 {
// 		parser.handshakeBuf.Reset()
// 	}
// 	return nil
// }

// func (parser *parserHello) parseHandshake(handshakeType handshakeType, buffer bufferView) bool {

// 	switch handshakeType {

// 	case clientHello:
// 		if parser.hello = parseClientHello(buffer); parser.hello == nil {
// 			return false
// 		}
// 		return true

// 	}
// 	return true
// }

// func parseHello(buf *streambuf.Buffer) parserResult {

// 	fmt.Println("In Decoder Parser")

// 	for buf.Avail(recordHeaderSize) {
// 		header, err := readRecordHeader(buf)
// 		if header.recordType == 22 {
// 			fmt.Println(header)
// 		}
// 		if err != nil || !header.isValid() {
// 			if err != nil {
// 				logp.Warn("internal buffer error: %v", err)
// 			}
// 			return resultFailed
// 		}
// 		limit := recordHeaderSize + int(header.length)
// 		if !buf.Avail(limit) {
// 			// wait for complete record
// 			return resultMore
// 		}

// 		switch header.recordType {

// 		case recordTypeHandshake:
// 			fmt.Println(header)

// 			if err = bufferHandshake(buf, int(header.length)); err != nil {
// 				logp.Warn("Error parsing handshake message: %v", err)
// 				return resultFailed
// 			}
// 		}

// 		buf.Advance(limit)
// 	}

// 	if buf.Len() == 0 {
// 		return resultOK
// 	}

// 	return resultMore
// }

// func (hello *helloMsg) parseExtensions(buffer bufferView) {
// 	hello.extensions = ParseExtensions(buffer)
// }

// func parseClientHello(buffer bufferView) *helloMs {
// 	var result helloMsg
// 	pos, ok := parseCommonHello(buffer, &result)
// 	if !ok {
// 		return nil
// 	}

// 	var cipherSuitesLength uint16
// 	if !buffer.read16Net(pos, &cipherSuitesLength) {
// 		logp.Warn("failed parsing client hello cipher suite length")
// 		return nil
// 	}

// 	for base := pos + 2; base < pos+2+int(cipherSuitesLength); base += 2 {
// 		var cipher uint16
// 		if !buffer.read16Net(base, &cipher) {
// 			logp.Warn("failed parsing client hello cipher suite")
// 			return nil
// 		}
// 		if !isGreaseValue(cipher) {
// 			result.supported.cipherSuites = append(result.supported.cipherSuites, cipherSuite(cipher))
// 		}
// 	}

// 	pos += 2 + int(cipherSuitesLength)
// 	var compMethodsLength uint8
// 	if !buffer.read8(pos, &compMethodsLength) {
// 		logp.Warn("failed parsing client hello compression methods length")
// 		return nil
// 	}
// 	limit := pos + 1 + int(compMethodsLength)
// 	for base := pos + 1; base < limit; base++ {
// 		var method uint8
// 		if !buffer.read8(base, &method) {
// 			logp.Warn("failed parsing client hello compression methods")
// 			return nil
// 		}
// 		result.supported.compression = append(result.supported.compression, compressionMethod(method))
// 	}

// 	result.parseExtensions(buffer.subview(limit, buffer.limit-limit))
// 	// fmt.Println("client hosts", result.extensions.Parsed["server_name_indication"])
// 	return &result
// }

// func (r bufferView) subview(start, length int) bufferView {
// 	if 0 <= start && 0 <= length && start+length <= r.limit {
// 		return bufferView{
// 			base:  r.base + start,
// 			limit: r.base + start + length,
// 			buf:   r.buf,
// 		}
// 	}

// 	panic(fmt.Sprintf("invalid buffer view requested start:%d len:%d limit:%d",
// 		start, length, r.limit))
// }

// type extensionParser func(reader bufferView) interface{}

// type extension struct {
// 	label   string
// 	parser  extensionParser
// 	saveRaw bool
// }

// var extensionMap = map[uint16]extension{
// 	0: {"server_name_indication", parseSni, false},
// }

// // ParseExtensions returns an tls.Extensions object parsed from the supplied buffer
// func ParseExtensions(buffer bufferView) tls.Extensions {

// 	var extensionsLength uint16
// 	if !buffer.read16Net(0, &extensionsLength) || extensionsLength == 0 {
// 		// No extensions
// 		return tls.Extensions{}
// 	}

// 	limit := 2 + int(extensionsLength)
// 	result := tls.Extensions{
// 		Parsed: common.MapStr{},
// 		Raw:    make(map[ExtensionID][]byte),
// 	}

// 	var unknown []string
// 	for base := 2; base < limit; {
// 		var code, length uint16
// 		if !buffer.read16Net(base, &code) || !buffer.read16Net(base+2, &length) {
// 			logp.Warn("failed parsing extensions")
// 			return tls.Extensions{}
// 		}

// 		extBuffer := buffer.subview(base+4, int(length))
// 		base += 4 + int(length)

// 		// Skip GREASE extensions
// 		if isGreaseValue(code) {
// 			continue
// 		}

// 		result.InOrder = append(result.InOrder, ExtensionID(code))
// 		label, parsed, saveRaw := parseExtension(code, extBuffer)
// 		if parsed != nil {
// 			// find here k sni kha aata hai
// 			result.Parsed[label] = parsed
// 		} else {
// 			unknown = append(unknown, label)
// 		}
// 		if saveRaw {
// 			result.Raw[ExtensionID(code)] = extBuffer.readBytes(0, extBuffer.length())
// 		}
// 	}
// 	if len(unknown) != 0 {
// 		result.Parsed["_unparsed_"] = unknown
// 	}
// 	return result
// }

// func parseExtension(code uint16, buffer bufferView) (string, interface{}, bool) {
// 	if ext, ok := extensionMap[code]; ok {
// 		parsed := ext.parser(buffer)
// 		return ext.label, parsed, ext.saveRaw
// 	}
// 	return strconv.Itoa(int(code)), nil, false
// }

// func parseSni(buffer bufferView) interface{} {
// 	var listLength uint16
// 	if !buffer.read16Net(0, &listLength) {
// 		return nil
// 	}
// 	var hosts []string
// 	// fmt.Println("######################## In parse SNI ########################")
// 	for pos, limit := 2, 2+int(listLength); pos+3 <= limit; {
// 		// fmt.Println("In parse sni", buffer.buf.)
// 		var nameType uint8
// 		var nameLen uint16
// 		var host string
// 		if !buffer.read8(pos, &nameType) || !buffer.read16Net(pos+1, &nameLen) ||
// 			limit < pos+3+int(nameLen) || !buffer.readString(pos+3, int(nameLen), &host) {
// 			logp.Warn("SNI hostname list truncated")
// 			break
// 		}
// 		if nameType == 0 {
// 			hosts = append(hosts, host)
// 		}
// 		pos += 3 + int(nameLen)
// 	}
// 	// fmt.Println("############################### Out of parse SNI ############################################")
// 	return hosts
// }

// func newBufferView(buf *streambuf.Buffer, pos int, length int) *bufferView {
// 	return &bufferView{buf, pos, pos + length}
// }

// func (r bufferView) read8(pos int, dest *uint8) bool {
// 	offset := pos + r.base
// 	if offset+1 > r.limit {
// 		return false
// 	}
// 	val, err := r.buf.ReadNetUint8At(offset)
// 	*dest = val
// 	return err == nil
// }

// func (r bufferView) read16Net(pos int, dest *uint16) bool {
// 	offset := pos + r.base
// 	if offset+2 > r.limit {
// 		return false
// 	}
// 	val, err := r.buf.ReadNetUint16At(offset)
// 	*dest = val
// 	return err == nil
// }

// func (r bufferView) read24Net(pos int, dest *uint32) bool {
// 	offset := pos + r.base
// 	if offset+3 > r.limit {
// 		return false
// 	}
// 	val8, err := r.buf.ReadNetUint8At(offset)
// 	if err != nil {
// 		return false
// 	}
// 	val16, err := r.buf.ReadNetUint16At(offset + 1)
// 	if err != nil {
// 		return false
// 	}
// 	*dest = uint32(val16) | (uint32(val8) << 16)
// 	return true
// }

// func (r bufferView) read32Net(pos int, dest *uint32) bool {
// 	offset := pos + r.base
// 	if offset+4 > r.limit {
// 		return false
// 	}
// 	val, err := r.buf.ReadNetUint32At(offset)
// 	*dest = val
// 	return err == nil
// }

// func (r bufferView) readString(pos int, len int, dest *string) bool {
// 	if slice := r.readBytes(pos, len); slice != nil {
// 		*dest = string(slice)
// 		return true
// 	}
// 	return false
// }

// func (r bufferView) readBytes(pos int, length int) []byte {
// 	offset := pos + r.base
// 	lim := offset + length
// 	if lim > r.limit {
// 		return nil
// 	}
// 	bytes := r.buf.Bytes()
// 	if lim <= len(bytes) {
// 		return r.buf.Bytes()[offset:lim]
// 	}
// 	return nil
// }

// func (r bufferView) length() int {
// 	return r.limit - r.base
// }
