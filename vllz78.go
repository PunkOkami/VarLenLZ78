package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"io"
	"math"
	"math/big"
	"math/bits"
	"os"
	"runtime"
	"strings"
)

func print_mem_usage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("\nAlloc = %v GiB", m.Alloc/1024/1024/1024)
	fmt.Printf("\tTotalAlloc = %v GiB", m.TotalAlloc/1024/1024/1024)
	fmt.Printf("\tSys = %v GiB", m.Sys/1024/1024/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func lz78_encoder(input_file_path string, encoded_file_path string) error {
	
	var x = 0.0
	var total float64
	
	input_file, err := os.Open(input_file_path)
	if err != nil {
		return err
	}
	stats, err := input_file.Stat()
	if err != nil {
		return err
	}
	total = float64(stats.Size())
	output_file, err := os.Create(encoded_file_path)
	if err != nil {
		return err
	}
	defer input_file.Close()
	defer output_file.Close()
	
	reader := bufio.NewReader(input_file)
	writer := bufio.NewWriter(output_file)
	
	var lz78_map = make(map[string]uint64)
	var i uint64 = 1
	var text string
	
	char, _, err := reader.ReadRune()
	for ; err == nil; char, _, err = reader.ReadRune() {
		x++
		_, ok := lz78_map[text+string(char)]
		if ok {
			text += string(char)
		} else {
			buf := new(bytes.Buffer)
			err_b := binary.Write(buf, binary.BigEndian, lz78_map[text])
			if err_b != nil {
				return err_b
			}
			byte_code := buf.Bytes()
			byte_code = byte_code[bits.LeadingZeros64(i)>>3:]
			
			_, err_r := writer.Write(byte_code)
			
			if err_r != nil {
				return err_r
			}
			_, err_r = writer.WriteRune(char)
			if err_r != nil {
				return err_r
			}
			
			lz78_map[text+string(char)] = i
			i++
			
			text = ""
		}
		perc := math.Round(x/total*10000.0) / 100.0
		fmt.Printf("Compressed %v percent of input file. Number of bytes needed: %v\r", perc, math.Log2(float64(i+1))/8)
	}
	
	if err != io.EOF {
		return err
	} else if text != "" {
		buf := new(bytes.Buffer)
		err_b := binary.Write(buf, binary.BigEndian, lz78_map[text])
		if err_b != nil {
			return err_b
		}
		byte_code := buf.Bytes()
		byte_code = byte_code[bits.LeadingZeros64(i)>>3:]
		
		_, err_r := writer.Write(byte_code)
		if err_r != nil {
			return err_r
		}
		_, err_r = writer.WriteRune(char)
		if err_r != nil {
			return err_r
		}
		i++
	}
	fmt.Printf("\n")
	fmt.Printf("Needed at most %v bytes to compress\n", int(math.Ceil(math.Log2(float64(i+1))/8)))
	
	err = writer.Flush()
	
	if err != nil {
		return err
	}
	
	return nil
}

func lz78_decoder(encoded_file_path string, output_file_path string) error {
	encded_file, err := os.Open(encoded_file_path)
	if err != nil {
		return err
	}
	defer encded_file.Close()
	
	output_file, err := os.Create(output_file_path)
	if err != nil {
		return nil
	}
	defer output_file.Close()
	
	reader := bufio.NewReader(encded_file)
	writer := bufio.NewWriter(output_file)
	
	var i uint64 = 1
	var byte_num uint8 = 1
	var text string
	lz78_map := make(map[uint64]string)
	
	for {
		if math.Log2(float64(i+1))/8 > float64(byte_num) {
			byte_num++
		}
		text = ""
		
		buf := make([]byte, byte_num)
		_, err := io.ReadFull(reader, buf)
		code := big.NewInt(0).SetBytes(buf).Uint64()
		
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		
		if code != 0 {
			text += lz78_map[code]
		}
		
		char, _, err := reader.ReadRune()
		
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		
		if char != 0 {
			text += string(char)
		}
		
		_, err_w := writer.WriteString(text)
		if err_w != nil {
			return err_w
		}
		
		lz78_map[i] = text
		i++
	}
	
	err = writer.Flush()
	
	if err != nil {
		return err
	}
	
	return nil
}

func main() {
	flag_set := ff.NewFlagSet("VarLenLZ78")
	var (
		input_file  = flag_set.String('i', "input_file", "", "relative path to input file")
		output_file = flag_set.String('o', "output_file", "{input_file}_compressed.vl78", "relative path to output file")
		decompress  = flag_set.Bool('d', "decompress", "decompress input file")
		help        = flag_set.Bool('h', "help", "help")
		verbose     = flag_set.Bool('v', "verbose", "verbose mode")
	)
	
	err := ff.Parse(flag_set, os.Args[1:])
	if err != nil {
		fmt.Println("Paring arguments failed: ", err)
	}
	
	if *help {
		fmt.Printf("%s\n", ffhelp.Flags(flag_set))
		os.Exit(0)
	}
	
	if *input_file == "" {
		fmt.Println("No input file specified!\n")
		fmt.Printf("%s\n", ffhelp.Flags(flag_set))
		os.Exit(1)
	}
	if *output_file == "{input_file}_compressed.vl78" {
		if *decompress {
			*output_file = fmt.Sprintf("%s_decompressed.file", strings.Split(*input_file, ".")[0])
		} else {
			*output_file = fmt.Sprintf("%s_compressed.vl78", strings.Split(*input_file, ".")[0])
		}
	}
	
	if *decompress {
		fmt.Println("Decompressing input file...", *input_file)
		err = lz78_decoder(*input_file, *output_file)
		if err != nil {
			fmt.Printf("Error decoding file %v : %v\n", *input_file, err.Error())
		}
		if *verbose {
			print_mem_usage()
		}
	} else {
		fmt.Println("Compressing input file...", *input_file)
		err = lz78_encoder(*input_file, *output_file)
		if err != nil {
			fmt.Printf("Error encoding file %v : %v\n", *input_file, err.Error())
		}
		if *verbose {
			print_mem_usage()
		}
	}
}
