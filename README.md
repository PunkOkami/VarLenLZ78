# VarLenLZ78
A variation on lz78 that encoded text files using variable length codes. The number of bytes grows as it's needed to ensure most
efficient coding. When map holding codes for strings grows to size when a new byte is needed to encode that number, it is added to code.
Decryption code also extended buffer size when reading encoded file. Currently, the biggest allowed code is 8 bytes long which allows for
18 446 744 073 709 551 615 different codes. Enwik 9 requires only X bytes and uses Y of RAM.

## Installation


## Usage
```zsh
# Compresses with default output file name
vllz78 -i test.txt
# Decompresses with default output file name
vllz78 -d test.vl78
```

