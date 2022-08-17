package protorand

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

var (
	// Chars is the set of characters used to generate random strings.
	Chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

// ProtoRand is a source of random values for protobuf fields.
type ProtoRand struct {
	rand *rand.Rand
}

// New creates a new ProtoRand.
func New() *ProtoRand {
	return &ProtoRand{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Seed sets the seed of the random generator.
func (p *ProtoRand) Seed(seed int64) {
	p.rand = rand.New(rand.NewSource(seed))
}

// Gen generates a new proto.Message having randoms value in its fields.
// The input is used to specify the type of the generated message.
// The input itself is immutable.
func (p *ProtoRand) Gen(in proto.Message) (proto.Message, error) {
	mds := in.ProtoReflect().Descriptor()
	dm, err := p.NewDynamicProtoRand(mds)
	if err != nil {
		return nil, err
	}

	out := reflect.New(reflect.ValueOf(in).Elem().Type()).Interface().(proto.Message)
	proto.Merge(out, dm)
	return out, nil
}

// NewDynamicProtoRand creates dynamicpb with assiging random value to proto.
func (p *ProtoRand) NewDynamicProtoRand(mds protoreflect.MessageDescriptor) (*dynamicpb.Message, error) {
	fmt.Printf("NewDynamicProtoRand(%v)\n", mds.Name())

	getRandValue := func(fd protoreflect.FieldDescriptor) (protoreflect.Value, error) {
		switch fd.Kind() {
		case protoreflect.Int32Kind:
			return protoreflect.ValueOfInt32(p.randInt32()), nil
		case protoreflect.Int64Kind:
			return protoreflect.ValueOfInt64(p.randInt64()), nil
		case protoreflect.Sint32Kind:
			return protoreflect.ValueOfInt32(p.randInt32()), nil
		case protoreflect.Sint64Kind:
			return protoreflect.ValueOfInt64(p.randInt64()), nil
		case protoreflect.Uint32Kind:
			return protoreflect.ValueOfUint32(p.randUint32()), nil
		case protoreflect.Uint64Kind:
			return protoreflect.ValueOfUint64(p.randUint64()), nil
		case protoreflect.Fixed32Kind:
			return protoreflect.ValueOfUint32(p.randUint32()), nil
		case protoreflect.Fixed64Kind:
			return protoreflect.ValueOfUint64(p.randUint64()), nil
		case protoreflect.Sfixed32Kind:
			return protoreflect.ValueOfInt32(p.randInt32()), nil
		case protoreflect.Sfixed64Kind:
			return protoreflect.ValueOfInt64(p.randInt64()), nil
		case protoreflect.FloatKind:
			return protoreflect.ValueOfFloat32(p.randFloat32()), nil
		case protoreflect.DoubleKind:
			return protoreflect.ValueOfFloat64(p.randFloat64()), nil
		case protoreflect.StringKind:
			return protoreflect.ValueOfString(p.randString()), nil
		case protoreflect.BoolKind:
			return protoreflect.ValueOfBool(p.randBool()), nil
		case protoreflect.EnumKind:
			return protoreflect.ValueOfEnum(p.chooseEnumValueRandomly(fd.Enum().Values())), nil
		case protoreflect.BytesKind:
			return protoreflect.ValueOfBytes(p.randBytes()), nil
		case protoreflect.MessageKind:
			// process recursively
			rm, err := p.NewDynamicProtoRand(fd.Message())
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
		populatedOneOfField[oneOf.Name()] = p.chooseOneOfFieldRandomly(oneOf).Number()
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

		fmt.Printf("calling getRandValue on field %v\n", fd.Name())
		value, err := getRandValue(fd)
		if err != nil {
			return nil, err
		}
		dm.Set(fd, value)
	}

	return dm, nil
}

func (p *ProtoRand) randInt32() int32 {
	return p.rand.Int31()
}

func (p *ProtoRand) randInt64() int64 {
	return p.rand.Int63()
}

func (p *ProtoRand) randUint32() uint32 {
	return p.rand.Uint32()
}

func (p *ProtoRand) randUint64() uint64 {
	return p.rand.Uint64()
}

func (p *ProtoRand) randFloat32() float32 {
	return p.rand.Float32()
}

func (p *ProtoRand) randFloat64() float64 {
	return p.rand.Float64()
}

func (p *ProtoRand) randBytes() []byte {
	return []byte(p.randString())
}

func (p *ProtoRand) randString() string {
	b := make([]rune, 10) // TODO: make the length randomly or use a predefined length?
	for i := range b {
		b[i] = Chars[p.rand.Intn(len(Chars))]
	}
	return string(b)
}

func (p *ProtoRand) randBool() bool {
	return p.rand.Int31()%2 == 0
}

func (p *ProtoRand) chooseEnumValueRandomly(values protoreflect.EnumValueDescriptors) protoreflect.EnumNumber {
	ln := values.Len()
	if ln <= 1 {
		return 0
	}

	value := values.Get(p.rand.Intn(ln - 1))
	return value.Number()
}

func (p *ProtoRand) chooseOneOfFieldRandomly(oneOf protoreflect.OneofDescriptor) protoreflect.FieldDescriptor {
	ln := oneOf.Fields().Len()
	if ln <= 1 {
		return oneOf.Fields().Get(0)
	}
	index := p.rand.Intn(ln - 1)
	return oneOf.Fields().Get(index)
}
