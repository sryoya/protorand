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

	// These aim to enable to inject the random value to be fixed in testing
	randomInt32      = genRandInt32
	randomFloat      = genRandFloat
	randomString     = genRandString
	randomBool       = genRandBool
	randIntForEnum   = rand.Intn
	randIndexForEnum = rand.Intn
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// EmbedValues embeds randoms value to fields in the provided proto message
func EmbedValues(msg proto.Message) error {
	mds := msg.ProtoReflect().Descriptor()
	dm, err := NewDynamicProtoRand(mds)
	if err != nil {
		return err
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
		case protoreflect.EnumKind:
			return protoreflect.ValueOfEnum(chooseEnumValueRandomly(fd.Enum().Values())), nil
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

	// decide which fields in each OneOf will be populated in advance
	populatedOneOfField := map[protoreflect.Name]protoreflect.FieldNumber{}
	oneOfs := mds.Oneofs()
	for i := 0; i < oneOfs.Len(); i++ {
		oneOf := oneOfs.Get(i)
		populatedOneOfField[oneOf.Name()] = chooseOneOfFieldRandomly(oneOf).Number()
	}

	dm := dynamicpb.NewMessage(mds)
	fds := mds.Fields()
	for k := 0; k < fds.Len(); k++ {
		fd := fds.Get(k)

		// If a field is in OneOf, check if the field should be populated
		if oneOf := fd.ContainingOneof(); oneOf != nil {
			populatedFieldNum := populatedOneOfField[oneOf.Name()]
			if populatedFieldNum != fd.Number() {
				continue
			}
		}

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

func chooseEnumValueRandomly(values protoreflect.EnumValueDescriptors) protoreflect.EnumNumber {
	ln := values.Len()
	if ln <= 1 {
		return 0
	}

	value := values.Get(randIntForEnum(ln - 1))

	return value.Number()
}

func chooseOneOfFieldRandomly(oneOf protoreflect.OneofDescriptor) protoreflect.FieldDescriptor {
	index := randIndexForEnum(oneOf.Fields().Len() - 1)
	return oneOf.Fields().Get(index)
}
