// Autogenerated by Thrift Compiler (1.0.0-dev)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package bar

import (
	"bytes"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = bytes.Equal

var GoUnusedProtection__ int

//SHA3-256
type ID []byte

func IDPtr(v ID) *ID { return &v }

type IDs [][]byte

func IDsPtr(v IDs) *IDs { return &v }

type Manifests []*Manifest

func ManifestsPtr(v Manifests) *Manifests { return &v }

// Bard server info
//
// Attributes:
//  - HttpEndpoint: HTTP endpoint (http://bard.served:3000/v1)
//  - RpcEndpoints: Thrift rpc endpoints (tcp://bard.served:3000)
//  - ChunkSize: Preferred chunk size.
//  - MaxConn: Preferred max connections for client
//  - BufferSize: Thrift client buffer size
type ServerInfo struct {
	HttpEndpoint string   `thrift:"httpEndpoint,1" json:"httpEndpoint"`
	RpcEndpoints []string `thrift:"rpcEndpoints,2" json:"rpcEndpoints"`
	ChunkSize    int64    `thrift:"chunkSize,3" json:"chunkSize"`
	MaxConn      int32    `thrift:"maxConn,4" json:"maxConn"`
	BufferSize   int32    `thrift:"bufferSize,5" json:"bufferSize"`
}

func NewServerInfo() *ServerInfo {
	return &ServerInfo{}
}

func (p *ServerInfo) GetHttpEndpoint() string {
	return p.HttpEndpoint
}

func (p *ServerInfo) GetRpcEndpoints() []string {
	return p.RpcEndpoints
}

func (p *ServerInfo) GetChunkSize() int64 {
	return p.ChunkSize
}

func (p *ServerInfo) GetMaxConn() int32 {
	return p.MaxConn
}

func (p *ServerInfo) GetBufferSize() int32 {
	return p.BufferSize
}
func (p *ServerInfo) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.readField2(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.readField3(iprot); err != nil {
				return err
			}
		case 4:
			if err := p.readField4(iprot); err != nil {
				return err
			}
		case 5:
			if err := p.readField5(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *ServerInfo) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.HttpEndpoint = v
	}
	return nil
}

func (p *ServerInfo) readField2(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]string, 0, size)
	p.RpcEndpoints = tSlice
	for i := 0; i < size; i++ {
		var _elem0 string
		if v, err := iprot.ReadString(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			_elem0 = v
		}
		p.RpcEndpoints = append(p.RpcEndpoints, _elem0)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *ServerInfo) readField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.ChunkSize = v
	}
	return nil
}

func (p *ServerInfo) readField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 4: ", err)
	} else {
		p.MaxConn = v
	}
	return nil
}

func (p *ServerInfo) readField5(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 5: ", err)
	} else {
		p.BufferSize = v
	}
	return nil
}

func (p *ServerInfo) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("ServerInfo"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := p.writeField4(oprot); err != nil {
		return err
	}
	if err := p.writeField5(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *ServerInfo) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("httpEndpoint", thrift.STRING, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:httpEndpoint: ", p), err)
	}
	if err := oprot.WriteString(string(p.HttpEndpoint)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.httpEndpoint (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:httpEndpoint: ", p), err)
	}
	return err
}

func (p *ServerInfo) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("rpcEndpoints", thrift.LIST, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:rpcEndpoints: ", p), err)
	}
	if err := oprot.WriteListBegin(thrift.STRING, len(p.RpcEndpoints)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range p.RpcEndpoints {
		if err := oprot.WriteString(string(v)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
	}
	if err := oprot.WriteListEnd(); err != nil {
		return thrift.PrependError("error writing list end: ", err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:rpcEndpoints: ", p), err)
	}
	return err
}

func (p *ServerInfo) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("chunkSize", thrift.I64, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:chunkSize: ", p), err)
	}
	if err := oprot.WriteI64(int64(p.ChunkSize)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.chunkSize (3) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:chunkSize: ", p), err)
	}
	return err
}

func (p *ServerInfo) writeField4(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("maxConn", thrift.I32, 4); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:maxConn: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.MaxConn)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.maxConn (4) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 4:maxConn: ", p), err)
	}
	return err
}

func (p *ServerInfo) writeField5(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("bufferSize", thrift.I32, 5); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:bufferSize: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.BufferSize)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.bufferSize (5) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 5:bufferSize: ", p), err)
	}
	return err
}

func (p *ServerInfo) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ServerInfo(%+v)", *p)
}

// Info about data entity
//
//
// Attributes:
//  - Id
//  - Size: data size
type DataInfo struct {
	Id   ID    `thrift:"id,1" json:"id"`
	Size int64 `thrift:"size,2" json:"size"`
}

func NewDataInfo() *DataInfo {
	return &DataInfo{}
}

func (p *DataInfo) GetId() ID {
	return p.Id
}

func (p *DataInfo) GetSize() int64 {
	return p.Size
}
func (p *DataInfo) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.readField2(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *DataInfo) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBinary(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		temp := ID(v)
		p.Id = temp
	}
	return nil
}

func (p *DataInfo) readField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Size = v
	}
	return nil
}

func (p *DataInfo) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("DataInfo"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *DataInfo) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("id", thrift.STRING, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:id: ", p), err)
	}
	if err := oprot.WriteBinary(p.Id); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.id (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:id: ", p), err)
	}
	return err
}

func (p *DataInfo) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("size", thrift.I64, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:size: ", p), err)
	}
	if err := oprot.WriteI64(int64(p.Size)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.size (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:size: ", p), err)
	}
	return err
}

func (p *DataInfo) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("DataInfo(%+v)", *p)
}

// Blob manifest.
//
//
// Attributes:
//  - Info
//  - Chunks
type Manifest struct {
	Info   *DataInfo `thrift:"info,1" json:"info"`
	Chunks []*Chunk  `thrift:"chunks,2" json:"chunks"`
}

func NewManifest() *Manifest {
	return &Manifest{}
}

var Manifest_Info_DEFAULT *DataInfo

func (p *Manifest) GetInfo() *DataInfo {
	if !p.IsSetInfo() {
		return Manifest_Info_DEFAULT
	}
	return p.Info
}

func (p *Manifest) GetChunks() []*Chunk {
	return p.Chunks
}
func (p *Manifest) IsSetInfo() bool {
	return p.Info != nil
}

func (p *Manifest) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.readField2(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *Manifest) readField1(iprot thrift.TProtocol) error {
	p.Info = &DataInfo{}
	if err := p.Info.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Info), err)
	}
	return nil
}

func (p *Manifest) readField2(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*Chunk, 0, size)
	p.Chunks = tSlice
	for i := 0; i < size; i++ {
		_elem1 := &Chunk{}
		if err := _elem1.Read(iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem1), err)
		}
		p.Chunks = append(p.Chunks, _elem1)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *Manifest) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Manifest"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *Manifest) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("info", thrift.STRUCT, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:info: ", p), err)
	}
	if err := p.Info.Write(oprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Info), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:info: ", p), err)
	}
	return err
}

func (p *Manifest) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("chunks", thrift.LIST, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:chunks: ", p), err)
	}
	if err := oprot.WriteListBegin(thrift.STRUCT, len(p.Chunks)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range p.Chunks {
		if err := v.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
		}
	}
	if err := oprot.WriteListEnd(); err != nil {
		return thrift.PrependError("error writing list end: ", err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:chunks: ", p), err)
	}
	return err
}

func (p *Manifest) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Manifest(%+v)", *p)
}

// Chunk info
//
//
// Attributes:
//  - Info
//  - Offset
type Chunk struct {
	Info *DataInfo `thrift:"info,1" json:"info"`
	// unused field # 2
	Offset int64 `thrift:"offset,3" json:"offset"`
}

func NewChunk() *Chunk {
	return &Chunk{}
}

var Chunk_Info_DEFAULT *DataInfo

func (p *Chunk) GetInfo() *DataInfo {
	if !p.IsSetInfo() {
		return Chunk_Info_DEFAULT
	}
	return p.Info
}

func (p *Chunk) GetOffset() int64 {
	return p.Offset
}
func (p *Chunk) IsSetInfo() bool {
	return p.Info != nil
}

func (p *Chunk) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.readField3(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *Chunk) readField1(iprot thrift.TProtocol) error {
	p.Info = &DataInfo{}
	if err := p.Info.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Info), err)
	}
	return nil
}

func (p *Chunk) readField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.Offset = v
	}
	return nil
}

func (p *Chunk) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Chunk"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *Chunk) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("info", thrift.STRUCT, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:info: ", p), err)
	}
	if err := p.Info.Write(oprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Info), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:info: ", p), err)
	}
	return err
}

func (p *Chunk) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("offset", thrift.I64, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:offset: ", p), err)
	}
	if err := oprot.WriteI64(int64(p.Offset)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.offset (3) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:offset: ", p), err)
	}
	return err
}

func (p *Chunk) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Chunk(%+v)", *p)
}

// Attributes:
//  - Id
//  - Timestamp
//  - Blobs
//  - Removes
type Spec struct {
	Id        ID            `thrift:"id,1" json:"id"`
	Timestamp int64         `thrift:"timestamp,2" json:"timestamp"`
	Blobs     map[string]ID `thrift:"blobs,3" json:"blobs"`
	Removes   []string      `thrift:"removes,4" json:"removes"`
}

func NewSpec() *Spec {
	return &Spec{}
}

func (p *Spec) GetId() ID {
	return p.Id
}

func (p *Spec) GetTimestamp() int64 {
	return p.Timestamp
}

func (p *Spec) GetBlobs() map[string]ID {
	return p.Blobs
}

func (p *Spec) GetRemoves() []string {
	return p.Removes
}
func (p *Spec) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.readField2(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.readField3(iprot); err != nil {
				return err
			}
		case 4:
			if err := p.readField4(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *Spec) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBinary(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		temp := ID(v)
		p.Id = temp
	}
	return nil
}

func (p *Spec) readField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Timestamp = v
	}
	return nil
}

func (p *Spec) readField3(iprot thrift.TProtocol) error {
	_, _, size, err := iprot.ReadMapBegin()
	if err != nil {
		return thrift.PrependError("error reading map begin: ", err)
	}
	tMap := make(map[string]ID, size)
	p.Blobs = tMap
	for i := 0; i < size; i++ {
		var _key2 string
		if v, err := iprot.ReadString(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			_key2 = v
		}
		var _val3 ID
		if v, err := iprot.ReadBinary(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			temp := ID(v)
			_val3 = temp
		}
		p.Blobs[_key2] = _val3
	}
	if err := iprot.ReadMapEnd(); err != nil {
		return thrift.PrependError("error reading map end: ", err)
	}
	return nil
}

func (p *Spec) readField4(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]string, 0, size)
	p.Removes = tSlice
	for i := 0; i < size; i++ {
		var _elem4 string
		if v, err := iprot.ReadString(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			_elem4 = v
		}
		p.Removes = append(p.Removes, _elem4)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *Spec) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Spec"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := p.writeField4(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *Spec) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("id", thrift.STRING, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:id: ", p), err)
	}
	if err := oprot.WriteBinary(p.Id); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.id (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:id: ", p), err)
	}
	return err
}

func (p *Spec) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("timestamp", thrift.I64, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:timestamp: ", p), err)
	}
	if err := oprot.WriteI64(int64(p.Timestamp)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.timestamp (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:timestamp: ", p), err)
	}
	return err
}

func (p *Spec) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("blobs", thrift.MAP, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:blobs: ", p), err)
	}
	if err := oprot.WriteMapBegin(thrift.STRING, thrift.STRING, len(p.Blobs)); err != nil {
		return thrift.PrependError("error writing map begin: ", err)
	}
	for k, v := range p.Blobs {
		if err := oprot.WriteString(string(k)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
		if err := oprot.WriteBinary(v); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
	}
	if err := oprot.WriteMapEnd(); err != nil {
		return thrift.PrependError("error writing map end: ", err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:blobs: ", p), err)
	}
	return err
}

func (p *Spec) writeField4(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("removes", thrift.LIST, 4); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:removes: ", p), err)
	}
	if err := oprot.WriteListBegin(thrift.STRING, len(p.Removes)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range p.Removes {
		if err := oprot.WriteString(string(v)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
	}
	if err := oprot.WriteListEnd(); err != nil {
		return thrift.PrependError("error writing list end: ", err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 4:removes: ", p), err)
	}
	return err
}

func (p *Spec) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Spec(%+v)", *p)
}
