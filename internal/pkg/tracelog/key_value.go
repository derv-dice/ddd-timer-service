package tracelog

type KeyValue struct {
	vT int
	k  string
	v  any
}

func (kv *KeyValue) IsValid() bool {
	return kv.vT != vTNone
}

const (
	vTNone = iota
	vTString
	vTInt
	vTBool
)

func String(key string, val string) KeyValue {
	return KeyValue{
		vT: vTString,
		k:  key,
		v:  val,
	}
}

func Int(key string, val int) KeyValue {
	return KeyValue{
		vT: vTInt,
		k:  key,
		v:  val,
	}
}

func Bool(key string, val bool) KeyValue {
	return KeyValue{
		vT: vTBool,
		k:  key,
		v:  val,
	}
}
