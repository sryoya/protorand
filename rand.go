package protorand

import (
	"fmt"
	"math/rand"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

var (
	chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	randomInt32  = genRandInt32
	randomFloat  = genRandFloat
	randomString = genRandString
	randomBool   = genRandBool
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// EmbedValues embeds randoms value to fields in the provided proto message
func EmbedValues(msg proto.Message) error {
	mds := msg.ProtoReflect().Descriptor()
	dm, err := NewDynamicProtoRand(mds)
	if err != nil {
		return nil
	}

	proto.Merge(msg, dm)
	return nil
}

// NewDynamicProtoRand created dynamicpb with assiging random value to proto
func NewDynamicProtoRand(mds protoreflect.MessageDescriptor) (*dynamicpb.Message, error) {
	getRandValue := func(fd protoreflect.FieldDescriptor) (protoreflect.Value, error) {
		switch fd.Kind() {
		case protoreflect.Int32Kind:
			return protoreflect.ValueOfInt32(randomInt32()), nil
		case protoreflect.FloatKind:
			return protoreflect.ValueOfFloat32(randomFloat()), nil
		case protoreflect.StringKind:
			return protoreflect.ValueOfString(randomString(10)), nil
		case protoreflect.BoolKind:
			return protoreflect.ValueOfBool(randomBool()), nil
		case protoreflect.MessageKind:
			// process recursively
			rm, err := NewDynamicProtoRand(fd.Message())
			if err != nil {
				return protoreflect.Value{}, err
			}
			return protoreflect.ValueOfMessage(rm), nil
		default:
			return protoreflect.Value{}, fmt.Errorf("unexpected type: %v", fd.Kind())
		}
	}

	dm := dynamicpb.NewMessage(mds)
	fds := mds.Fields()
	for k := 0; k < fds.Len(); k++ {
		fd := fds.Get(k)

		if fd.IsList() {
			list := dm.Mutable(fd).List()
			// TODO: decide the number of elements randomly
			value, err := getRandValue(fd)
			if err != nil {
				return nil, err
			}
			list.Append(value)
			dm.Set(fd, protoreflect.ValueOfList(list))
			continue
		}
		if fd.IsMap() {
			mp := dm.Mutable(fd).Map()
			// TODO: make the number of elements randomly
			key, err := getRandValue(fd.MapKey())
			if err != nil {
				return nil, err
			}
			value, err := getRandValue(fd.MapValue())
			if err != nil {
				return nil, err
			}
			mp.Set(protoreflect.MapKey(key), protoreflect.Value(value))
			dm.Set(fd, protoreflect.ValueOfMap(mp))
			continue
		}

		value, err := getRandValue(fd)
		if err != nil {
			return nil, err
		}
		dm.Set(fd, value)
	}

	return dm, nil
}

func genRandInt32() int32 {
	return rand.Int31()
}

func genRandFloat() float32 {
	return rand.Float32()
}

func genRandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func genRandBool() bool {
	return rand.Int31()%2 == 0
}
