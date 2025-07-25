// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: product/v1/category.proto

package productv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ListByTreeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListByTreeRequest) Reset() {
	*x = ListByTreeRequest{}
	mi := &file_product_v1_category_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListByTreeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListByTreeRequest) ProtoMessage() {}

func (x *ListByTreeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_product_v1_category_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListByTreeRequest.ProtoReflect.Descriptor instead.
func (*ListByTreeRequest) Descriptor() ([]byte, []int) {
	return file_product_v1_category_proto_rawDescGZIP(), []int{0}
}

type ListByTreeResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CategoryTree  []*CategoryResponse    `protobuf:"bytes,1,rep,name=category_tree,json=categoryTree,proto3" json:"category_tree,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListByTreeResponse) Reset() {
	*x = ListByTreeResponse{}
	mi := &file_product_v1_category_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListByTreeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListByTreeResponse) ProtoMessage() {}

func (x *ListByTreeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_product_v1_category_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListByTreeResponse.ProtoReflect.Descriptor instead.
func (*ListByTreeResponse) Descriptor() ([]byte, []int) {
	return file_product_v1_category_proto_rawDescGZIP(), []int{1}
}

func (x *ListByTreeResponse) GetCategoryTree() []*CategoryResponse {
	if x != nil {
		return x.CategoryTree
	}
	return nil
}

type CategoryResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt     int64                  `protobuf:"varint,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     int64                  `protobuf:"varint,3,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	Name          string                 `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	ParentCid     int64                  `protobuf:"varint,5,opt,name=parent_cid,json=parentCid,proto3" json:"parent_cid,omitempty"`
	CatLevel      int32                  `protobuf:"varint,6,opt,name=cat_level,json=catLevel,proto3" json:"cat_level,omitempty"`
	ShowStatus    int32                  `protobuf:"varint,7,opt,name=show_status,json=showStatus,proto3" json:"show_status,omitempty"`
	Sort          int32                  `protobuf:"varint,8,opt,name=sort,proto3" json:"sort,omitempty"`
	Icon          string                 `protobuf:"bytes,9,opt,name=icon,proto3" json:"icon,omitempty"`
	ProductUint   string                 `protobuf:"bytes,10,opt,name=product_uint,json=productUint,proto3" json:"product_uint,omitempty"`
	ProductCount  int32                  `protobuf:"varint,11,opt,name=product_count,json=productCount,proto3" json:"product_count,omitempty"`
	Children      []*CategoryResponse    `protobuf:"bytes,12,rep,name=children,proto3" json:"children,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CategoryResponse) Reset() {
	*x = CategoryResponse{}
	mi := &file_product_v1_category_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CategoryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CategoryResponse) ProtoMessage() {}

func (x *CategoryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_product_v1_category_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CategoryResponse.ProtoReflect.Descriptor instead.
func (*CategoryResponse) Descriptor() ([]byte, []int) {
	return file_product_v1_category_proto_rawDescGZIP(), []int{2}
}

func (x *CategoryResponse) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *CategoryResponse) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *CategoryResponse) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *CategoryResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CategoryResponse) GetParentCid() int64 {
	if x != nil {
		return x.ParentCid
	}
	return 0
}

func (x *CategoryResponse) GetCatLevel() int32 {
	if x != nil {
		return x.CatLevel
	}
	return 0
}

func (x *CategoryResponse) GetShowStatus() int32 {
	if x != nil {
		return x.ShowStatus
	}
	return 0
}

func (x *CategoryResponse) GetSort() int32 {
	if x != nil {
		return x.Sort
	}
	return 0
}

func (x *CategoryResponse) GetIcon() string {
	if x != nil {
		return x.Icon
	}
	return ""
}

func (x *CategoryResponse) GetProductUint() string {
	if x != nil {
		return x.ProductUint
	}
	return ""
}

func (x *CategoryResponse) GetProductCount() int32 {
	if x != nil {
		return x.ProductCount
	}
	return 0
}

func (x *CategoryResponse) GetChildren() []*CategoryResponse {
	if x != nil {
		return x.Children
	}
	return nil
}

var File_product_v1_category_proto protoreflect.FileDescriptor

const file_product_v1_category_proto_rawDesc = "" +
	"\n" +
	"\x19product/v1/category.proto\x12\n" +
	"product.v1\"\x13\n" +
	"\x11ListByTreeRequest\"W\n" +
	"\x12ListByTreeResponse\x12A\n" +
	"\rcategory_tree\x18\x01 \x03(\v2\x1c.product.v1.CategoryResponseR\fcategoryTree\"\xfb\x02\n" +
	"\x10CategoryResponse\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x1d\n" +
	"\n" +
	"created_at\x18\x02 \x01(\x03R\tcreatedAt\x12\x1d\n" +
	"\n" +
	"updated_at\x18\x03 \x01(\x03R\tupdatedAt\x12\x12\n" +
	"\x04name\x18\x04 \x01(\tR\x04name\x12\x1d\n" +
	"\n" +
	"parent_cid\x18\x05 \x01(\x03R\tparentCid\x12\x1b\n" +
	"\tcat_level\x18\x06 \x01(\x05R\bcatLevel\x12\x1f\n" +
	"\vshow_status\x18\a \x01(\x05R\n" +
	"showStatus\x12\x12\n" +
	"\x04sort\x18\b \x01(\x05R\x04sort\x12\x12\n" +
	"\x04icon\x18\t \x01(\tR\x04icon\x12!\n" +
	"\fproduct_uint\x18\n" +
	" \x01(\tR\vproductUint\x12#\n" +
	"\rproduct_count\x18\v \x01(\x05R\fproductCount\x128\n" +
	"\bchildren\x18\f \x03(\v2\x1c.product.v1.CategoryResponseR\bchildren2^\n" +
	"\x0fCategoryService\x12K\n" +
	"\n" +
	"ListByTree\x12\x1d.product.v1.ListByTreeRequest\x1a\x1e.product.v1.ListByTreeResponseB\xa6\x01\n" +
	"\x0ecom.product.v1B\rCategoryProtoP\x01Z<github.com/WeiXinao/daily_fresh/api/gen/product/v1;productv1\xa2\x02\x03PXX\xaa\x02\n" +
	"Product.V1\xca\x02\n" +
	"Product\\V1\xe2\x02\x16Product\\V1\\GPBMetadata\xea\x02\vProduct::V1b\x06proto3"

var (
	file_product_v1_category_proto_rawDescOnce sync.Once
	file_product_v1_category_proto_rawDescData []byte
)

func file_product_v1_category_proto_rawDescGZIP() []byte {
	file_product_v1_category_proto_rawDescOnce.Do(func() {
		file_product_v1_category_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_product_v1_category_proto_rawDesc), len(file_product_v1_category_proto_rawDesc)))
	})
	return file_product_v1_category_proto_rawDescData
}

var file_product_v1_category_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_product_v1_category_proto_goTypes = []any{
	(*ListByTreeRequest)(nil),  // 0: product.v1.ListByTreeRequest
	(*ListByTreeResponse)(nil), // 1: product.v1.ListByTreeResponse
	(*CategoryResponse)(nil),   // 2: product.v1.CategoryResponse
}
var file_product_v1_category_proto_depIdxs = []int32{
	2, // 0: product.v1.ListByTreeResponse.category_tree:type_name -> product.v1.CategoryResponse
	2, // 1: product.v1.CategoryResponse.children:type_name -> product.v1.CategoryResponse
	0, // 2: product.v1.CategoryService.ListByTree:input_type -> product.v1.ListByTreeRequest
	1, // 3: product.v1.CategoryService.ListByTree:output_type -> product.v1.ListByTreeResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_product_v1_category_proto_init() }
func file_product_v1_category_proto_init() {
	if File_product_v1_category_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_product_v1_category_proto_rawDesc), len(file_product_v1_category_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_product_v1_category_proto_goTypes,
		DependencyIndexes: file_product_v1_category_proto_depIdxs,
		MessageInfos:      file_product_v1_category_proto_msgTypes,
	}.Build()
	File_product_v1_category_proto = out.File
	file_product_v1_category_proto_goTypes = nil
	file_product_v1_category_proto_depIdxs = nil
}
