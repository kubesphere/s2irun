// DO NOT EDIT.
//
// Generated by the Swift generator plugin for the protocol buffer compiler.
// Source: plugin.proto
//
// For information on using the generated types, please see the documenation:
//   https://github.com/apple/swift-protobuf/

// Copyright 2017 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// openapic (aka the OpenAPI Compiler) can be extended via plugins.  
// A plugin is just a program that reads a Request from stdin 
// and writes a Response to stdout.
//
// A plugin executable needs only to be placed somewhere in the path.  The
// plugin should be named "openapi_$NAME", and will then be used when the
// flag "--${NAME}_out" is passed to openapic.

import Foundation
import SwiftProtobuf

// If the compiler emits an error on this type, it is because this file
// was generated by a version of the `protoc` Swift plug-in that is
// incompatible with the version of SwiftProtobuf to which you are linking.
// Please ensure that your are building against the same version of the API
// that was used to generate this file.
fileprivate struct _GeneratedWithProtocGenSwiftVersion: SwiftProtobuf.ProtobufAPIVersionCheck {
  struct _2: SwiftProtobuf.ProtobufAPIVersion_2 {}
  typealias Version = _2
}

/// The version number of OpenAPI compiler.
struct Openapi_Plugin_V1_Version: SwiftProtobuf.Message {
  static let protoMessageName: String = _protobuf_package + ".Version"

  var major: Int32 = 0

  var minor: Int32 = 0

  var patch: Int32 = 0

  /// A suffix for alpha, beta or rc release, e.g., "alpha-1", "rc2". It should
  /// be empty for mainline stable releases.
  var suffix: String = String()

  var unknownFields = SwiftProtobuf.UnknownStorage()

  init() {}

  /// Used by the decoding initializers in the SwiftProtobuf library, not generally
  /// used directly. `init(serializedData:)`, `init(jsonUTF8Data:)`, and other decoding
  /// initializers are defined in the SwiftProtobuf library. See the Message and
  /// Message+*Additions` files.
  mutating func decodeMessage<D: SwiftProtobuf.Decoder>(decoder: inout D) throws {
    while let fieldNumber = try decoder.nextFieldNumber() {
      switch fieldNumber {
      case 1: try decoder.decodeSingularInt32Field(value: &self.major)
      case 2: try decoder.decodeSingularInt32Field(value: &self.minor)
      case 3: try decoder.decodeSingularInt32Field(value: &self.patch)
      case 4: try decoder.decodeSingularStringField(value: &self.suffix)
      default: break
      }
    }
  }

  /// Used by the encoding methods of the SwiftProtobuf library, not generally
  /// used directly. `Message.serializedData()`, `Message.jsonUTF8Data()`, and
  /// other serializer methods are defined in the SwiftProtobuf library. See the
  /// `Message` and `Message+*Additions` files.
  func traverse<V: SwiftProtobuf.Visitor>(visitor: inout V) throws {
    if self.major != 0 {
      try visitor.visitSingularInt32Field(value: self.major, fieldNumber: 1)
    }
    if self.minor != 0 {
      try visitor.visitSingularInt32Field(value: self.minor, fieldNumber: 2)
    }
    if self.patch != 0 {
      try visitor.visitSingularInt32Field(value: self.patch, fieldNumber: 3)
    }
    if !self.suffix.isEmpty {
      try visitor.visitSingularStringField(value: self.suffix, fieldNumber: 4)
    }
    try unknownFields.traverse(visitor: &visitor)
  }
}

/// A parameter passed to the plugin from (or through) the OpenAPI compiler.
struct Openapi_Plugin_V1_Parameter: SwiftProtobuf.Message {
  static let protoMessageName: String = _protobuf_package + ".Parameter"

  /// The name of the parameter as specified in the option string
  var name: String = String()

  /// The parameter value as specified in the option string
  var value: String = String()

  var unknownFields = SwiftProtobuf.UnknownStorage()

  init() {}

  /// Used by the decoding initializers in the SwiftProtobuf library, not generally
  /// used directly. `init(serializedData:)`, `init(jsonUTF8Data:)`, and other decoding
  /// initializers are defined in the SwiftProtobuf library. See the Message and
  /// Message+*Additions` files.
  mutating func decodeMessage<D: SwiftProtobuf.Decoder>(decoder: inout D) throws {
    while let fieldNumber = try decoder.nextFieldNumber() {
      switch fieldNumber {
      case 1: try decoder.decodeSingularStringField(value: &self.name)
      case 2: try decoder.decodeSingularStringField(value: &self.value)
      default: break
      }
    }
  }

  /// Used by the encoding methods of the SwiftProtobuf library, not generally
  /// used directly. `Message.serializedData()`, `Message.jsonUTF8Data()`, and
  /// other serializer methods are defined in the SwiftProtobuf library. See the
  /// `Message` and `Message+*Additions` files.
  func traverse<V: SwiftProtobuf.Visitor>(visitor: inout V) throws {
    if !self.name.isEmpty {
      try visitor.visitSingularStringField(value: self.name, fieldNumber: 1)
    }
    if !self.value.isEmpty {
      try visitor.visitSingularStringField(value: self.value, fieldNumber: 2)
    }
    try unknownFields.traverse(visitor: &visitor)
  }
}

/// An encoded Request is written to the plugin's stdin.
struct Openapi_Plugin_V1_Request: SwiftProtobuf.Message {
  static let protoMessageName: String = _protobuf_package + ".Request"

  /// A wrapped OpenAPI document to process.
  var wrapper: Openapi_Plugin_V1_Wrapper {
    get {return _storage._wrapper ?? Openapi_Plugin_V1_Wrapper()}
    set {_uniqueStorage()._wrapper = newValue}
  }
  /// Returns true if `wrapper` has been explicitly set.
  var hasWrapper: Bool {return _storage._wrapper != nil}
  /// Clears the value of `wrapper`. Subsequent reads from it will return its default value.
  mutating func clearWrapper() {_storage._wrapper = nil}

  /// Output path specified in the plugin invocation.
  var outputPath: String {
    get {return _storage._outputPath}
    set {_uniqueStorage()._outputPath = newValue}
  }

  /// Plugin parameters parsed from the invocation string.
  var parameters: [Openapi_Plugin_V1_Parameter] {
    get {return _storage._parameters}
    set {_uniqueStorage()._parameters = newValue}
  }

  /// The version number of openapi compiler.
  var compilerVersion: Openapi_Plugin_V1_Version {
    get {return _storage._compilerVersion ?? Openapi_Plugin_V1_Version()}
    set {_uniqueStorage()._compilerVersion = newValue}
  }
  /// Returns true if `compilerVersion` has been explicitly set.
  var hasCompilerVersion: Bool {return _storage._compilerVersion != nil}
  /// Clears the value of `compilerVersion`. Subsequent reads from it will return its default value.
  mutating func clearCompilerVersion() {_storage._compilerVersion = nil}

  var unknownFields = SwiftProtobuf.UnknownStorage()

  init() {}

  /// Used by the decoding initializers in the SwiftProtobuf library, not generally
  /// used directly. `init(serializedData:)`, `init(jsonUTF8Data:)`, and other decoding
  /// initializers are defined in the SwiftProtobuf library. See the Message and
  /// Message+*Additions` files.
  mutating func decodeMessage<D: SwiftProtobuf.Decoder>(decoder: inout D) throws {
    _ = _uniqueStorage()
    try withExtendedLifetime(_storage) { (_storage: _StorageClass) in
      while let fieldNumber = try decoder.nextFieldNumber() {
        switch fieldNumber {
        case 1: try decoder.decodeSingularMessageField(value: &_storage._wrapper)
        case 2: try decoder.decodeSingularStringField(value: &_storage._outputPath)
        case 3: try decoder.decodeRepeatedMessageField(value: &_storage._parameters)
        case 4: try decoder.decodeSingularMessageField(value: &_storage._compilerVersion)
        default: break
        }
      }
    }
  }

  /// Used by the encoding methods of the SwiftProtobuf library, not generally
  /// used directly. `Message.serializedData()`, `Message.jsonUTF8Data()`, and
  /// other serializer methods are defined in the SwiftProtobuf library. See the
  /// `Message` and `Message+*Additions` files.
  func traverse<V: SwiftProtobuf.Visitor>(visitor: inout V) throws {
    try withExtendedLifetime(_storage) { (_storage: _StorageClass) in
      if let v = _storage._wrapper {
        try visitor.visitSingularMessageField(value: v, fieldNumber: 1)
      }
      if !_storage._outputPath.isEmpty {
        try visitor.visitSingularStringField(value: _storage._outputPath, fieldNumber: 2)
      }
      if !_storage._parameters.isEmpty {
        try visitor.visitRepeatedMessageField(value: _storage._parameters, fieldNumber: 3)
      }
      if let v = _storage._compilerVersion {
        try visitor.visitSingularMessageField(value: v, fieldNumber: 4)
      }
    }
    try unknownFields.traverse(visitor: &visitor)
  }

  fileprivate var _storage = _StorageClass.defaultInstance
}

/// The plugin writes an encoded Response to stdout.
struct Openapi_Plugin_V1_Response: SwiftProtobuf.Message {
  static let protoMessageName: String = _protobuf_package + ".Response"

  /// Error message.  If non-empty, the plugin failed. 
  /// The plugin process should exit with status code zero 
  /// even if it reports an error in this way.
  ///
  /// This should be used to indicate errors which prevent the plugin from 
  /// operating as intended.  Errors which indicate a problem in openapic 
  /// itself -- such as the input Document being unparseable -- should be 
  /// reported by writing a message to stderr and exiting with a non-zero 
  /// status code.
  var errors: [String] = []

  /// file output, each file will be written by openapic to an appropriate location.
  var files: [Openapi_Plugin_V1_File] = []

  var unknownFields = SwiftProtobuf.UnknownStorage()

  init() {}

  /// Used by the decoding initializers in the SwiftProtobuf library, not generally
  /// used directly. `init(serializedData:)`, `init(jsonUTF8Data:)`, and other decoding
  /// initializers are defined in the SwiftProtobuf library. See the Message and
  /// Message+*Additions` files.
  mutating func decodeMessage<D: SwiftProtobuf.Decoder>(decoder: inout D) throws {
    while let fieldNumber = try decoder.nextFieldNumber() {
      switch fieldNumber {
      case 1: try decoder.decodeRepeatedStringField(value: &self.errors)
      case 2: try decoder.decodeRepeatedMessageField(value: &self.files)
      default: break
      }
    }
  }

  /// Used by the encoding methods of the SwiftProtobuf library, not generally
  /// used directly. `Message.serializedData()`, `Message.jsonUTF8Data()`, and
  /// other serializer methods are defined in the SwiftProtobuf library. See the
  /// `Message` and `Message+*Additions` files.
  func traverse<V: SwiftProtobuf.Visitor>(visitor: inout V) throws {
    if !self.errors.isEmpty {
      try visitor.visitRepeatedStringField(value: self.errors, fieldNumber: 1)
    }
    if !self.files.isEmpty {
      try visitor.visitRepeatedMessageField(value: self.files, fieldNumber: 2)
    }
    try unknownFields.traverse(visitor: &visitor)
  }
}

/// File describes a file generated by a plugin.
struct Openapi_Plugin_V1_File: SwiftProtobuf.Message {
  static let protoMessageName: String = _protobuf_package + ".File"

  /// name of the file
  var name: String = String()

  /// data to be written to the file
  var data: Data = SwiftProtobuf.Internal.emptyData

  var unknownFields = SwiftProtobuf.UnknownStorage()

  init() {}

  /// Used by the decoding initializers in the SwiftProtobuf library, not generally
  /// used directly. `init(serializedData:)`, `init(jsonUTF8Data:)`, and other decoding
  /// initializers are defined in the SwiftProtobuf library. See the Message and
  /// Message+*Additions` files.
  mutating func decodeMessage<D: SwiftProtobuf.Decoder>(decoder: inout D) throws {
    while let fieldNumber = try decoder.nextFieldNumber() {
      switch fieldNumber {
      case 1: try decoder.decodeSingularStringField(value: &self.name)
      case 2: try decoder.decodeSingularBytesField(value: &self.data)
      default: break
      }
    }
  }

  /// Used by the encoding methods of the SwiftProtobuf library, not generally
  /// used directly. `Message.serializedData()`, `Message.jsonUTF8Data()`, and
  /// other serializer methods are defined in the SwiftProtobuf library. See the
  /// `Message` and `Message+*Additions` files.
  func traverse<V: SwiftProtobuf.Visitor>(visitor: inout V) throws {
    if !self.name.isEmpty {
      try visitor.visitSingularStringField(value: self.name, fieldNumber: 1)
    }
    if !self.data.isEmpty {
      try visitor.visitSingularBytesField(value: self.data, fieldNumber: 2)
    }
    try unknownFields.traverse(visitor: &visitor)
  }
}

/// Wrapper wraps an OpenAPI document with its version.
struct Openapi_Plugin_V1_Wrapper: SwiftProtobuf.Message {
  static let protoMessageName: String = _protobuf_package + ".Wrapper"

  /// filename or URL of the wrapped document
  var name: String = String()

  /// version of the OpenAPI specification that is used by the wrapped document
  var version: String = String()

  /// valid serialized protocol buffer of the named OpenAPI specification version
  var value: Data = SwiftProtobuf.Internal.emptyData

  var unknownFields = SwiftProtobuf.UnknownStorage()

  init() {}

  /// Used by the decoding initializers in the SwiftProtobuf library, not generally
  /// used directly. `init(serializedData:)`, `init(jsonUTF8Data:)`, and other decoding
  /// initializers are defined in the SwiftProtobuf library. See the Message and
  /// Message+*Additions` files.
  mutating func decodeMessage<D: SwiftProtobuf.Decoder>(decoder: inout D) throws {
    while let fieldNumber = try decoder.nextFieldNumber() {
      switch fieldNumber {
      case 1: try decoder.decodeSingularStringField(value: &self.name)
      case 2: try decoder.decodeSingularStringField(value: &self.version)
      case 3: try decoder.decodeSingularBytesField(value: &self.value)
      default: break
      }
    }
  }

  /// Used by the encoding methods of the SwiftProtobuf library, not generally
  /// used directly. `Message.serializedData()`, `Message.jsonUTF8Data()`, and
  /// other serializer methods are defined in the SwiftProtobuf library. See the
  /// `Message` and `Message+*Additions` files.
  func traverse<V: SwiftProtobuf.Visitor>(visitor: inout V) throws {
    if !self.name.isEmpty {
      try visitor.visitSingularStringField(value: self.name, fieldNumber: 1)
    }
    if !self.version.isEmpty {
      try visitor.visitSingularStringField(value: self.version, fieldNumber: 2)
    }
    if !self.value.isEmpty {
      try visitor.visitSingularBytesField(value: self.value, fieldNumber: 3)
    }
    try unknownFields.traverse(visitor: &visitor)
  }
}

// MARK: - Code below here is support for the SwiftProtobuf runtime.

fileprivate let _protobuf_package = "openapi.plugin.v1"

extension Openapi_Plugin_V1_Version: SwiftProtobuf._MessageImplementationBase, SwiftProtobuf._ProtoNameProviding {
  static let _protobuf_nameMap: SwiftProtobuf._NameMap = [
    1: .same(proto: "major"),
    2: .same(proto: "minor"),
    3: .same(proto: "patch"),
    4: .same(proto: "suffix"),
  ]

  func _protobuf_generated_isEqualTo(other: Openapi_Plugin_V1_Version) -> Bool {
    if self.major != other.major {return false}
    if self.minor != other.minor {return false}
    if self.patch != other.patch {return false}
    if self.suffix != other.suffix {return false}
    if unknownFields != other.unknownFields {return false}
    return true
  }
}

extension Openapi_Plugin_V1_Parameter: SwiftProtobuf._MessageImplementationBase, SwiftProtobuf._ProtoNameProviding {
  static let _protobuf_nameMap: SwiftProtobuf._NameMap = [
    1: .same(proto: "name"),
    2: .same(proto: "value"),
  ]

  func _protobuf_generated_isEqualTo(other: Openapi_Plugin_V1_Parameter) -> Bool {
    if self.name != other.name {return false}
    if self.value != other.value {return false}
    if unknownFields != other.unknownFields {return false}
    return true
  }
}

extension Openapi_Plugin_V1_Request: SwiftProtobuf._MessageImplementationBase, SwiftProtobuf._ProtoNameProviding {
  static let _protobuf_nameMap: SwiftProtobuf._NameMap = [
    1: .same(proto: "wrapper"),
    2: .standard(proto: "output_path"),
    3: .same(proto: "parameters"),
    4: .standard(proto: "compiler_version"),
  ]

  fileprivate class _StorageClass {
    var _wrapper: Openapi_Plugin_V1_Wrapper? = nil
    var _outputPath: String = String()
    var _parameters: [Openapi_Plugin_V1_Parameter] = []
    var _compilerVersion: Openapi_Plugin_V1_Version? = nil

    static let defaultInstance = _StorageClass()

    private init() {}

    init(copying source: _StorageClass) {
      _wrapper = source._wrapper
      _outputPath = source._outputPath
      _parameters = source._parameters
      _compilerVersion = source._compilerVersion
    }
  }

  fileprivate mutating func _uniqueStorage() -> _StorageClass {
    if !isKnownUniquelyReferenced(&_storage) {
      _storage = _StorageClass(copying: _storage)
    }
    return _storage
  }

  func _protobuf_generated_isEqualTo(other: Openapi_Plugin_V1_Request) -> Bool {
    if _storage !== other._storage {
      let storagesAreEqual: Bool = withExtendedLifetime((_storage, other._storage)) { (_storage, other_storage) in
        if _storage._wrapper != other_storage._wrapper {return false}
        if _storage._outputPath != other_storage._outputPath {return false}
        if _storage._parameters != other_storage._parameters {return false}
        if _storage._compilerVersion != other_storage._compilerVersion {return false}
        return true
      }
      if !storagesAreEqual {return false}
    }
    if unknownFields != other.unknownFields {return false}
    return true
  }
}

extension Openapi_Plugin_V1_Response: SwiftProtobuf._MessageImplementationBase, SwiftProtobuf._ProtoNameProviding {
  static let _protobuf_nameMap: SwiftProtobuf._NameMap = [
    1: .same(proto: "errors"),
    2: .same(proto: "files"),
  ]

  func _protobuf_generated_isEqualTo(other: Openapi_Plugin_V1_Response) -> Bool {
    if self.errors != other.errors {return false}
    if self.files != other.files {return false}
    if unknownFields != other.unknownFields {return false}
    return true
  }
}

extension Openapi_Plugin_V1_File: SwiftProtobuf._MessageImplementationBase, SwiftProtobuf._ProtoNameProviding {
  static let _protobuf_nameMap: SwiftProtobuf._NameMap = [
    1: .same(proto: "name"),
    2: .same(proto: "data"),
  ]

  func _protobuf_generated_isEqualTo(other: Openapi_Plugin_V1_File) -> Bool {
    if self.name != other.name {return false}
    if self.data != other.data {return false}
    if unknownFields != other.unknownFields {return false}
    return true
  }
}

extension Openapi_Plugin_V1_Wrapper: SwiftProtobuf._MessageImplementationBase, SwiftProtobuf._ProtoNameProviding {
  static let _protobuf_nameMap: SwiftProtobuf._NameMap = [
    1: .same(proto: "name"),
    2: .same(proto: "version"),
    3: .same(proto: "value"),
  ]

  func _protobuf_generated_isEqualTo(other: Openapi_Plugin_V1_Wrapper) -> Bool {
    if self.name != other.name {return false}
    if self.version != other.version {return false}
    if self.value != other.value {return false}
    if unknownFields != other.unknownFields {return false}
    return true
  }
}
