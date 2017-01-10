# protoc-go-inject-field
[![Build Status](https://travis-ci.org/favadi/protoc-go-inject-field.svg?branch=master)](https://travis-ci.org/favadi/protoc-go-inject-field)
[![Go Report Card](https://goreportcard.com/badge/github.com/favadi/protoc-go-inject-field)](https://goreportcard.com/report/github.com/favadi/protoc-go-inject-field)

## Why?

Sometimes it is useful to have custom unexported fields in generated golang codes, but this use case is 
[unsupported](https://github.com/golang/protobuf/issues/38). This tool injects custom fields along with its 
getter/setter methods to generated .pb.go files.  

## Install

`go get github.com/favadi/protoc-go-inject-field`

## Usage

Add one or more comments with syntax `// @inject_field: field_name field_type` before messages.

Example:

```
// file: test.proto
syntax = "proto3";

package pb;

// @inject_field: age int
// @inject_field: spouse *Human
// @inject_field: IgnoreMe
message Human {
    string name = 1;
}

// @inject_field: model string
message Robot {
    string name = 1;
}

message Alien {
    string name = 1;
}
```

Generate with protoc command.

```
protoc --go_out=. test.proto
```

Run `protoc-go-inject-field` with generated file.

```
protoc-go-inject-field -input=./test.pb.go
```

The custom fields will be injected along with its getter/setter methods.

```
diff --git a/playground.pb.go.orig b/playground.pb.go
index 679c243..5b4ec2f 100644
--- a/playground.pb.go.orig
+++ b/playground.pb.go
@@ -38,6 +38,24 @@ const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package
 // @inject_field: IgnoreMe
 type Human struct {
     Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
+
+    // custom fields
+    age int
+    spouse *Human
+}
+
+// custom fields getter/setter
+func (m *Human) Age() int {
+    return m.age
+}
+func (m *Human) SetAge(in int){
+    m.age = in
+}
+func (m *Human) Spouse() *Human {
+    return m.spouse
+}
+func (m *Human) SetSpouse(in *Human){
+    m.spouse = in
 }
 
 func (m *Human) Reset()                    { *m = Human{} }
@@ -48,6 +66,17 @@ func (*Human) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }
 // @inject_field: model string
 type Robot struct {
     Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
+
+    // custom fields
+    model string
+}
+
+// custom fields getter/setter
+func (m *Robot) Model() string {
+    return m.model
+}
+func (m *Robot) SetModel(in string){
+    m.model = in
 }
 
 func (m *Robot) Reset()                    { *m = Robot{} }
 ```
 
