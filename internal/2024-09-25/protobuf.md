# protobuf

phase 1:

~~~go
func (Message) AddVarint(Number, Varint)
func (Message) AddFixed64(Number, Fixed64)
func (Message) AddFixed32(Number, Fixed32)
func (Message) AddBytes(Number, Bytes)
func (Message) Add(Number, Message)
func (Message) AddFunc(Number, func(Message))

func (Message) GetVarint(Number) func() (Varint, bool)
func (Message) GetFixed64(Number) func() (Fixed64, bool)
func (Message) GetFixed32(Number) func() (Fixed32, bool)
func (Message) GetBytes(Number) func() (Bytes, bool)
func (Message) Get(Number) func() (Message, bool)
func (Message) GetFunc(Number, func(Message) bool)
~~~

possible phase 2:

~~~go
func (Message) Add(Number, Message)
func (Message) AddBytes(Number, Bytes)
func (Message) AddFixed32(Number, Fixed32)
func (Message) AddFixed64(Number, Fixed64)
func (Message) AddFunc(Number, func(Message))
func (Message) AddVarint(Number, Varint)

func (Message) Get(Number) (Message, bool)
func (Message) GetBytes(Number) (Bytes, bool)
func (Message) GetFixed32(Number) (Fixed32, bool)
func (Message) GetFixed64(Number) (Fixed64, bool)
func (Message) GetFunc(Number, func(Message))
func (Message) GetVarint(Number) (Varint, bool)

func (Message) GetAll(Number) func() (Message, bool)
func (Message) GetAllBytes(Number) func() (Bytes, bool)
func (Message) GetAllFixed32(Number) func() (Fixed32, bool)
func (Message) GetAllFixed64(Number) func() (Fixed64, bool)
func (Message) GetAllFunc(Number, func(Message) bool)
func (Message) GetAllVarint(Number) func() (Varint, bool)
~~~
