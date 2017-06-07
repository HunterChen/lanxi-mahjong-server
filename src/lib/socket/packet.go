/**********************************************************
 * Author        : Michael
 * Email         : dolotech@163.com
 * Last modified : 2016-01-23 10:24
 * Filename      : packet.go
 * Description   : 通讯协议处理，主要处理封包和解包的过程
 * *******************************************************/

package socket

import "errors"

var HEADER_ERR = errors.New("header error")

var HeaderLen uint32 = 1 //　包头长度
var HANDDLen uint32 = 9

const (
	PROTOLen uint32 = 4
	DataLen  uint32 = 4 //包信息数据长度占位长度
)

//Big Endian
func DecodeUint32(data []byte) uint32 {
	return (uint32(data[0]) << 24) | (uint32(data[1]) << 16) | (uint32(data[2]) << 8) | uint32(data[3])
}

//Big Endian
func EncodeUint32(n uint32) []byte {
	b := make([]byte, 4)
	b[3] = byte(n & 0xFF)
	b[2] = byte((n >> 8) & 0xFF)
	b[1] = byte((n >> 16) & 0xFF)
	b[0] = byte((n >> 24) & 0xFF)
	return b
}

//封包
func Pack(proto uint32, message []byte, count uint32) []byte {
	buff := make([]byte, int(HANDDLen)+len(message))
	msglen := uint32(len(message))
	buff[0] = byte(count)
	copy(buff[HeaderLen:HeaderLen+PROTOLen], EncodeUint32(proto))
	copy(buff[HeaderLen+PROTOLen:HeaderLen+PROTOLen+DataLen], EncodeUint32(msglen))
	copy(buff[HANDDLen:HANDDLen+msglen], message)
	return buff
}

//解包
func Unpack(buffer []byte, readerChannel chan *Packet) ([]byte, error) {
	length := uint32(len(buffer))
	var i uint32
	for i = 0; i < length; {
		// 包头都不足
		if length < i+HANDDLen {
			break
		}
		buf := buffer[i:1]
		count := uint32(buf[0])
		// 读取信息数据长度
		messageLength := DecodeUint32(buffer[i+HeaderLen+PROTOLen : i+HANDDLen])

		// 只有包头，数据不足一包
		if length < i+HANDDLen+messageLength {
			break
		}
		// 读取整包信息数据
		data := buffer[i+HANDDLen : i+HANDDLen+messageLength]

		p := &Packet{
			proto:   DecodeUint32(buffer[i+HeaderLen : i+HANDDLen]),
			content: data,
			count:   count,
		}
		readerChannel <- p
		i += HANDDLen + messageLength
	}
	// 刚好整包，包括2,3,4...个包
	if i == length {
		return make([]byte, 0, HeaderLen+1), nil
	}
	// 不足一包，或者残包余数
	return buffer[i:], nil
}
