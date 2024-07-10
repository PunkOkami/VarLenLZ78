# VarLenLZ78
A variation on lz78 that encoded text files using variable length codes. The number of bytes grows as it's needed to ensure most
efficient coding. When map holding codes for strings grows to size when a new byte is needed to encode that number, it is added to code.
Decryption code also extended buffer size when reading encoded file. Currently, the biggest allowed code is 8 bytes long which allows for
18 446 744 073 709 551 615 different codes. 

## Enwik9
Performance was tested using [enwik9](https://mattmahoney.net/dc/text.html) compression benchmark. With result file being XXX bytes long, 
it placed VarLenLZ78 on YYY position with 38.86% compression rate. It was not added to the list yet, but I contacted the author asking about it.

## Installation
Clone code repo...
```zsh
git clone https://github.com/PunkOkami/VarLenLZ78.git
```
...and build the binary
```zsh
go build -o vllz78 vllz78.go
```
You are good too go!

## Usage
```zsh
# Compresses with default output file name
./vllz78 -i test.txt
# Decompresses with default output file name
./vllz78 -d test.vl78
```
## Options
- i/input_file STRING    relative path to input file
- o/output_file STRING   relative path to output file (default: {input_file}_compressed.vl78)
- d/decompress           decompress input file
- h/help                 help
- v/verbose              verbose mode