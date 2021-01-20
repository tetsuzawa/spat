package spat

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	BitLenShort   = 16
	BitLenFloat   = 32
	BitLenDouble  = 64
	ByteLenShort  = 2
	ByteLenFloat  = 4
	ByteLenDouble = 8
)

var (
	ErrUnknownDataType = errors.New("unknown data type")
)

// DataType is type of
// DataType behaves as enum.
type DataType int

const (
	DSA DataType = iota + 1
	DFA
	DDA
	DSB
	DFB
	DDB
)

// String returns data type name as string.
func (dt DataType) String() string {
	switch dt {
	case DSA:
		return "DSA"
	case DFA:
		return "DFA"
	case DDA:
		return "DDA"
	case DSB:
		return "DSB"
	case DFB:
		return "DFB"
	case DDB:
		return "DDB"
	default:
		return "unknown data type" // unreachable code
	}
}

// StringToDataType determines data type from specified string.
// If the specified string is invalid, this func returns error.
func StringToDataType(s string) (DataType, error) {
	switch s {
	case "DSA":
		return DSA, nil
	case "DFA":
		return DFA, nil
	case "DDA":
		return DDA, nil
	case "DSB":
		return DSB, nil
	case "DFB":
		return DFB, nil
	case "DDB":
		return DDB, nil
	default:
		return 0, ErrUnknownDataType
	}
}

// BitLen returns the bit length of data type.
func (dt DataType) BitLen() int {
	switch dt {
	case DSA:
		return BitLenShort
	case DFA:
		return BitLenFloat
	case DDA:
		return BitLenDouble
	case DSB:
		return BitLenShort
	case DFB:
		return BitLenFloat
	case DDB:
		return BitLenDouble
	default:
		return -1 // unreachable code
	}
}

// ByteLen returns the byte length of data type.
func (dt DataType) ByteLen() int {
	switch dt {
	case DSA:
		return ByteLenShort
	case DFA:
		return ByteLenFloat
	case DDA:
		return ByteLenDouble
	case DSB:
		return ByteLenShort
	case DFB:
		return ByteLenFloat
	case DDB:
		return ByteLenDouble
	default:
		return -1 // unreachable code
	}
}

// LenDXXFile reads sample size of .DXX file.
// This func determines the data type from the filename extension and reads that data.
func LenDXXFile(filename string) (int, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	dt, err := StringToDataType(ext(filename))
	if err != nil {
		return 0, err
	}
	return int(fi.Size()) / dt.ByteLen(), nil
}

// Read reads data from reader as specified data type.
// The return type is []float64 to make the data easier to handle.
func Read(r io.Reader, dt DataType, length int) ([]float64, error) {
	switch dt {
	case DSA:
		i16s, err := readDSA(r, length)
		if err != nil {
			return nil, err
		}
		return Int16sToFloat64s(i16s), nil
	case DFA:
		f32s, err := readDFA(r, length)
		if err != nil {
			return nil, err
		}
		return Float32sToFloat64s(f32s), nil
	case DDA:
		f64s, err := readDDA(r, length)
		if err != nil {
			return nil, err
		}
		return f64s, nil
	case DSB:
		i16s, err := readDSB(r, length)
		if err != nil {
			return nil, err
		}
		return Int16sToFloat64s(i16s), nil
	case DFB:
		f32s, err := readDFB(r, length)
		if err != nil {
			return nil, err
		}
		return Float32sToFloat64s(f32s), nil
	case DDB:
		f64s, err := readDDB(r, length)
		if err != nil {
			return nil, err
		}
		return f64s, nil
	default:
		return nil, ErrUnknownDataType
	}
}

// ReadDXXFile reads .DXX file.
// This func determines the data type from the filename extension and reads that data.
// The return type is []float64 to make the data easier to handle.
func ReadDXXFile(filename string) ([]float64, error) {
	dt, err := StringToDataType(ext(filename))
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(b)
	length := r.Len() / dt.ByteLen()
	return Read(r, dt, length)
}

func readDSA(r io.Reader, length int) ([]int16, error) {
	sc := bufio.NewScanner(r)
	data := make([]int16, length)
	for i := range data {
		sc.Scan()
		v, err := strconv.ParseInt(sc.Text(), 10, 16)
		if err != nil {
			return nil, err
		}
		data[i] = int16(v)
	}
	return data, nil
}

func readDFA(r io.Reader, length int) ([]float32, error) {
	sc := bufio.NewScanner(r)
	data := make([]float32, length)
	for i := range data {
		sc.Scan()
		v, err := strconv.ParseFloat(sc.Text(), 32)
		if err != nil {
			return nil, err
		}
		data[i] = float32(v)
	}
	return data, nil
}

func readDDA(r io.Reader, length int) ([]float64, error) {
	sc := bufio.NewScanner(r)
	data := make([]float64, length)
	for i := range data {
		sc.Scan()
		v, err := strconv.ParseFloat(sc.Text(), 64)
		if err != nil {
			return nil, err
		}
		data[i] = v
	}
	return data, nil
}

func readDSB(r io.Reader, length int) ([]int16, error) {
	buf := make([]byte, ByteLenShort)
	data := make([]int16, length)
	for i := range data {
		_, err := io.ReadFull(r, buf)
		if err != nil {
			if err == io.EOF {
				return data, nil
			}
			return data, err
		}
		v, err := BytesToInt16(buf)
		if err != nil {
			return data, err
		}
		data[i] = v
	}
	return data, nil
}

func readDFB(r io.Reader, length int) ([]float32, error) {
	buf := make([]byte, ByteLenFloat)
	data := make([]float32, length)
	for i := range data {
		_, err := io.ReadFull(r, buf)
		if err != nil {
			if err == io.EOF {
				return data, nil
			}
			return data, err
		}
		v, err := BytesToFloat32(buf)
		if err != nil {
			return data, err
		}
		data[i] = v
	}
	return data, nil
}

func readDDB(r io.Reader, length int) ([]float64, error) {
	buf := make([]byte, ByteLenDouble)
	data := make([]float64, length)
	for i := range data {
		_, err := io.ReadFull(r, buf)
		if err != nil {
			if err == io.EOF {
				return data, nil
			}
			return data, err
		}
		v, err := BytesToFloat64(buf)
		if err != nil {
			return data, err
		}
		data[i] = v
	}
	return data, nil
}

// Writes writes data to writer as specified data type.
// The return type is []float64 to make the data easier to handle.
func Write(w io.Writer, dt DataType, data []float64) error {
	buf := &bytes.Buffer{}
	var err error
	switch dt {
	case DSA:
		err = writeDSA(buf, Float64sToInt16s(data))
	case DFA:
		err = writeDFA(buf, Float64sToFloat32s(data))
	case DDA:
		err = writeDDA(buf, data)
	case DSB:
		err = writeDSB(buf, Float64sToInt16s(data))
	case DFB:
		err = writeDFB(buf, Float64sToFloat32s(data))
	case DDB:
		err = writeDDB(buf, data)
	default:
		err = ErrUnknownDataType
	}
	if err != nil {
		return err
	}
	_, err = io.Copy(w, buf)
	return err
}

// WriteDXXFile writes data to .DXX file.
// This func determines the data type from the filename extension and writes the data to the file.
// The return type is []float64 to make the data easier to handle.
func WriteDXXFile(filename string, data []float64) error {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		return err
	}

	dt, err := StringToDataType(ext(filename))
	if err != nil {
		return err
	}
	return Write(f, dt, data)
}

func writeDSA(w io.Writer, data []int16) error {
	for _, v := range data {
		if _, err := fmt.Fprintf(w, "%d\n", v); err != nil {
			return err
		}
	}
	return nil
}

func writeDFA(w io.Writer, data []float32) error {
	for _, v := range data {
		if _, err := fmt.Fprintf(w, "%e\n", v); err != nil {
			return err
		}
	}
	return nil
}

func writeDDA(w io.Writer, data []float64) error {
	for _, v := range data {
		if _, err := fmt.Fprintf(w, "%e\n", v); err != nil {
			return err
		}
	}
	return nil
}

func writeDSB(w io.Writer, data []int16) error {
	for _, v := range data {
		buf, err := Int16ToBytes(v)
		if err != nil {
			return err
		}
		_, err = w.Write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeDFB(w io.Writer, data []float32) error {
	for _, v := range data {
		buf, err := Float32ToBytes(v)
		if err != nil {
			return err
		}
		_, err = w.Write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeDDB(w io.Writer, data []float64) error {
	for _, v := range data {
		buf, err := Float64ToBytes(v)
		if err != nil {
			return err
		}
		_, err = w.Write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

// ext returns the path of extension *without* dot.
// eg: ext(/path/to/file.aaa) -> aaa
func ext(path string) string {
	return strings.TrimPrefix(filepath.Ext(path), ".")
}
