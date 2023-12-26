package userfile

type File struct {
	bytes []byte
	parts []Part
	name  string
	len   int
}

func New(name string, bb []byte) File {
	return File{
		bytes: bb,
		parts: nil,
		name:  name,
		len:   len(bb),
	}
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Divide(n int) []Part {
	f.parts = make([]Part, n)
	fileParts := splitFile(f.bytes, n)
	for i := 0; i < len(fileParts); i++ {
		f.parts[i] = Part{
			b:    fileParts[i],
			meta: Meta{},
		}
	}
	f.bytes = nil

	return f.parts
}

func (f *File) Bytes() []byte {
	if len(f.bytes) > 0 {
		return f.bytes
	}

	var (
		l     = 0
		st    = 0
		fin   = 0
		bytes = make([]byte, f.len)
	)

	for i := 0; i < len(f.parts); i++ {
		l = len(f.parts[i].b)
		st = i * l
		fin = i*l + l
		copy(bytes[st:fin], f.parts[i].b)
	}

	return bytes
}

type Part struct {
	b    []byte
	meta Meta
}

func (p Part) Bytes() []byte {
	return p.b
}

func (p *Part) SetMeta(m Meta) {
	p.meta = m
}

type Meta struct {
	StorageServerHost string
}

func splitFile(data []byte, n int) [][]byte {
	var parts [][]byte
	partSize := len(data) / n

	for i := 0; i < n; i++ {
		start := i * partSize
		end := start + partSize
		if i == n-1 {
			end = len(data)
		}
		parts = append(parts, data[start:end])
	}

	return parts
}
