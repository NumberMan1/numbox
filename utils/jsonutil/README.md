# jsonutil

基于 json-iterator，提供 ToJSONString、ToJSONBytes、UnmarshalFromJSON、GetIntSliceFromJSON、
GetInt32SliceFromJSON 和 GetMapFromJSON。编解码失败时 panic，保持原公开行为。

ToJSONString(nil) 返回空字符串；ToJSONBytes(nil) 返回 JSON 的 null。
动态类型序列化可能使用依赖库的反射机制；其他工具子包不连带依赖此包。
