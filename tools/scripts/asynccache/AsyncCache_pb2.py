# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: AsyncCache.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
from google.protobuf import descriptor_pb2
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='AsyncCache.proto',
  package='',
  serialized_pb=_b('\n\x10\x41syncCache.proto\"\x95\x01\n\x11\x41syncCacheRequest\x12\x10\n\x08\x62lock_id\x18\x01 \x01(\x03\x12\x13\n\x0bsource_host\x18\x02 \x01(\t\x12\x13\n\x0bsource_port\x18\x03 \x01(\x05\x12\x34\n\x16open_ufs_block_options\x18\x04 \x01(\x0b\x32\x14.OpenUfsBlockOptions\x12\x0e\n\x06length\x18\x05 \x01(\x03\"\xa3\x01\n\x13OpenUfsBlockOptions\x12\x10\n\x08ufs_path\x18\x01 \x01(\t\x12\x16\n\x0eoffset_in_file\x18\x02 \x01(\x03\x12\x12\n\nblock_size\x18\x03 \x01(\x03\x12\x1d\n\x15maxUfsReadConcurrency\x18\x04 \x01(\x05\x12\x0f\n\x07mountId\x18\x05 \x01(\x03\x12\x10\n\x08no_cache\x18\x06 \x01(\x08\x12\x0c\n\x04user\x18\x07 \x01(\t\"&\n\x16LocalBlockOpenResponse\x12\x0c\n\x04path\x18\x01 \x01(\t')
)
_sym_db.RegisterFileDescriptor(DESCRIPTOR)




_ASYNCCACHEREQUEST = _descriptor.Descriptor(
  name='AsyncCacheRequest',
  full_name='AsyncCacheRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='block_id', full_name='AsyncCacheRequest.block_id', index=0,
      number=1, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='source_host', full_name='AsyncCacheRequest.source_host', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='source_port', full_name='AsyncCacheRequest.source_port', index=2,
      number=3, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='open_ufs_block_options', full_name='AsyncCacheRequest.open_ufs_block_options', index=3,
      number=4, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='length', full_name='AsyncCacheRequest.length', index=4,
      number=5, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=21,
  serialized_end=170,
)


_OPENUFSBLOCKOPTIONS = _descriptor.Descriptor(
  name='OpenUfsBlockOptions',
  full_name='OpenUfsBlockOptions',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='ufs_path', full_name='OpenUfsBlockOptions.ufs_path', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='offset_in_file', full_name='OpenUfsBlockOptions.offset_in_file', index=1,
      number=2, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='block_size', full_name='OpenUfsBlockOptions.block_size', index=2,
      number=3, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='maxUfsReadConcurrency', full_name='OpenUfsBlockOptions.maxUfsReadConcurrency', index=3,
      number=4, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='mountId', full_name='OpenUfsBlockOptions.mountId', index=4,
      number=5, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='no_cache', full_name='OpenUfsBlockOptions.no_cache', index=5,
      number=6, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='user', full_name='OpenUfsBlockOptions.user', index=6,
      number=7, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=173,
  serialized_end=336,
)


_LOCALBLOCKOPENRESPONSE = _descriptor.Descriptor(
  name='LocalBlockOpenResponse',
  full_name='LocalBlockOpenResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='path', full_name='LocalBlockOpenResponse.path', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=338,
  serialized_end=376,
)

_ASYNCCACHEREQUEST.fields_by_name['open_ufs_block_options'].message_type = _OPENUFSBLOCKOPTIONS
DESCRIPTOR.message_types_by_name['AsyncCacheRequest'] = _ASYNCCACHEREQUEST
DESCRIPTOR.message_types_by_name['OpenUfsBlockOptions'] = _OPENUFSBLOCKOPTIONS
DESCRIPTOR.message_types_by_name['LocalBlockOpenResponse'] = _LOCALBLOCKOPENRESPONSE

AsyncCacheRequest = _reflection.GeneratedProtocolMessageType('AsyncCacheRequest', (_message.Message,), dict(
  DESCRIPTOR = _ASYNCCACHEREQUEST,
  __module__ = 'AsyncCache_pb2'
  # @@protoc_insertion_point(class_scope:AsyncCacheRequest)
  ))
_sym_db.RegisterMessage(AsyncCacheRequest)

OpenUfsBlockOptions = _reflection.GeneratedProtocolMessageType('OpenUfsBlockOptions', (_message.Message,), dict(
  DESCRIPTOR = _OPENUFSBLOCKOPTIONS,
  __module__ = 'AsyncCache_pb2'
  # @@protoc_insertion_point(class_scope:OpenUfsBlockOptions)
  ))
_sym_db.RegisterMessage(OpenUfsBlockOptions)

LocalBlockOpenResponse = _reflection.GeneratedProtocolMessageType('LocalBlockOpenResponse', (_message.Message,), dict(
  DESCRIPTOR = _LOCALBLOCKOPENRESPONSE,
  __module__ = 'AsyncCache_pb2'
  # @@protoc_insertion_point(class_scope:LocalBlockOpenResponse)
  ))
_sym_db.RegisterMessage(LocalBlockOpenResponse)


# @@protoc_insertion_point(module_scope)
