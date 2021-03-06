package txscript

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"

	chainhash "git.parallelcoin.io/dev/9/pkg/chain/hash"
	"git.parallelcoin.io/dev/9/pkg/chain/wire"
	ec "git.parallelcoin.io/dev/9/pkg/util/elliptic"
	"golang.org/x/crypto/ripemd160"
)

// An opcode defines the information related to a txscript opcode.  opfunc, if present, is the function to call to perform the opcode on the script.  The current script is passed in as a slice with the first member being the opcode itself.

type opcode struct {
	value  byte
	name   string
	length int
	opfunc func(
		*parsedOpcode, *Engine) error
}

// These constants are the values of the official opcodes used on the btc wiki, in bitcoin core and in
// most if not all other references and software related to handling BTC scripts.
// Since anyway these opcodes are not used directly in their mnemonic form
// and we aren't living in the 70s, these are all re-named in standard Go case format
// Because `golint` is awesome.
const (
	OpZero                = 0x00 // 0
	OpFalse               = 0x00 // 0 - AKA OpZero
	OpData1               = 0x01 // 1
	OpData2               = 0x02 // 2
	OpData3               = 0x03 // 3
	OpData4               = 0x04 // 4
	OpData5               = 0x05 // 5
	OpData6               = 0x06 // 6
	OpData7               = 0x07 // 7
	OpData8               = 0x08 // 8
	OpData9               = 0x09 // 9
	OpData10              = 0x0a // 10
	OpData11              = 0x0b // 11
	OpData12              = 0x0c // 12
	OpData13              = 0x0d // 13
	OpData14              = 0x0e // 14
	OpData15              = 0x0f // 15
	OpData16              = 0x10 // 16
	OpData17              = 0x11 // 17
	OpData18              = 0x12 // 18
	OpData19              = 0x13 // 19
	OpData20              = 0x14 // 20
	OpData21              = 0x15 // 21
	OpData22              = 0x16 // 22
	OpData23              = 0x17 // 23
	OpData24              = 0x18 // 24
	OpData25              = 0x19 // 25
	OpData26              = 0x1a // 26
	OpData27              = 0x1b // 27
	OpData28              = 0x1c // 28
	OpData29              = 0x1d // 29
	OpData30              = 0x1e // 30
	OpData31              = 0x1f // 31
	OpData32              = 0x20 // 32
	OpData33              = 0x21 // 33
	OpData34              = 0x22 // 34
	OpData35              = 0x23 // 35
	OpData36              = 0x24 // 36
	OpData37              = 0x25 // 37
	OpData38              = 0x26 // 38
	OpData39              = 0x27 // 39
	OpData40              = 0x28 // 40
	OpData41              = 0x29 // 41
	OpData42              = 0x2a // 42
	OpData43              = 0x2b // 43
	OpData44              = 0x2c // 44
	OpData45              = 0x2d // 45
	OpData46              = 0x2e // 46
	OpData47              = 0x2f // 47
	OpData48              = 0x30 // 48
	OpData49              = 0x31 // 49
	OpData50              = 0x32 // 50
	OpData51              = 0x33 // 51
	OpData52              = 0x34 // 52
	OpData53              = 0x35 // 53
	OpData54              = 0x36 // 54
	OpData55              = 0x37 // 55
	OpData56              = 0x38 // 56
	OpData57              = 0x39 // 57
	OpData58              = 0x3a // 58
	OpData59              = 0x3b // 59
	OpData60              = 0x3c // 60
	OpData61              = 0x3d // 61
	OpData62              = 0x3e // 62
	OpData63              = 0x3f // 63
	OpData64              = 0x40 // 64
	OpData65              = 0x41 // 65
	OpData66              = 0x42 // 66
	OpData67              = 0x43 // 67
	OpData68              = 0x44 // 68
	OpData69              = 0x45 // 69
	OpData70              = 0x46 // 70
	OpData71              = 0x47 // 71
	OpData72              = 0x48 // 72
	OpData73              = 0x49 // 73
	OpData74              = 0x4a // 74
	OpData75              = 0x4b // 75
	OpPushData1           = 0x4c // 76
	OpPushData2           = 0x4d // 77
	OpPushData4           = 0x4e // 78
	Op1Negate             = 0x4f // 79
	OpReserved            = 0x50 // 80
	Op1                   = 0x51 // 81 - AKA OpTrue
	OpTrue                = 0x51 // 81
	Op2                   = 0x52 // 82
	Op3                   = 0x53 // 83
	Op4                   = 0x54 // 84
	Op5                   = 0x55 // 85
	Op6                   = 0x56 // 86
	Op7                   = 0x57 // 87
	Op8                   = 0x58 // 88
	Op9                   = 0x59 // 89
	Op10                  = 0x5a // 90
	Op11                  = 0x5b // 91
	Op12                  = 0x5c // 92
	Op13                  = 0x5d // 93
	Op14                  = 0x5e // 94
	Op15                  = 0x5f // 95
	Op16                  = 0x60 // 96
	OpNoOp                = 0x61 // 97
	OpVer                 = 0x62 // 98
	OpIf                  = 0x63 // 99
	OpIfNot               = 0x64 // 100
	OpVerIf               = 0x65 // 101
	OpVerIfNot            = 0x66 // 102
	OpElse                = 0x67 // 103
	OpEndIf               = 0x68 // 104
	OpVerify              = 0x69 // 105
	OpReturn              = 0x6a // 106
	OpToAltStack          = 0x6b // 107
	OpFromAltStack        = 0x6c // 108
	Op2Drop               = 0x6d // 109
	Op2Dup                = 0x6e // 110
	Op3Dup                = 0x6f // 111
	Op2Over               = 0x70 // 112
	Op2Rot                = 0x71 // 113
	Op2Swap               = 0x72 // 114
	OpIfDup               = 0x73 // 115
	OpDepth               = 0x74 // 116
	OpDrop                = 0x75 // 117
	OpDup                 = 0x76 // 118
	OpNip                 = 0x77 // 119
	OpOver                = 0x78 // 120
	OpPick                = 0x79 // 121
	OpRoll                = 0x7a // 122
	OpRot                 = 0x7b // 123
	OpSwap                = 0x7c // 124
	OpTuck                = 0x7d // 125
	OpCat                 = 0x7e // 126
	OpSubstr              = 0x7f // 127
	OpLeft                = 0x80 // 128
	OpRight               = 0x81 // 129
	OpSize                = 0x82 // 130
	OpInvert              = 0x83 // 131
	OpAnd                 = 0x84 // 132
	OpOr                  = 0x85 // 133
	OpXor                 = 0x86 // 134
	OpEqual               = 0x87 // 135
	OpEqualVerify         = 0x88 // 136
	OpReserved1           = 0x89 // 137
	OpReserved2           = 0x8a // 138
	Op1Add                = 0x8b // 139
	Op1Sub                = 0x8c // 140
	Op1Mul                = 0x8d // 141
	Op2Div                = 0x8e // 142
	OpNegate              = 0x8f // 143
	OpAbs                 = 0x90 // 144
	OpNot                 = 0x91 // 145
	Op0NotEqual           = 0x92 // 146
	OpAdd                 = 0x93 // 147
	OpSub                 = 0x94 // 148
	OpMul                 = 0x95 // 149
	OpDiv                 = 0x96 // 150
	OpMod                 = 0x97 // 151
	OpLShift              = 0x98 // 152
	OpRShift              = 0x99 // 153
	OpBoolAnd             = 0x9a // 154
	OpBoolOr              = 0x9b // 155
	OpNumEqual            = 0x9c // 156
	OpNumEqualVerify      = 0x9d // 157
	OpNumNotEqual         = 0x9e // 158
	OpLessThan            = 0x9f // 159
	OpGreaterThan         = 0xa0 // 160
	OpLessThanOrEqual     = 0xa1 // 161
	OpGreaterThanOrEqual  = 0xa2 // 162
	OpMin                 = 0xa3 // 163
	OpMax                 = 0xa4 // 164
	OpWithin              = 0xa5 // 165
	OpRipeMD160           = 0xa6 // 166
	OpSHA1                = 0xa7 // 167
	OpSHA256              = 0xa8 // 168
	OpHash160             = 0xa9 // 169
	OpHash256             = 0xaa // 170
	OpCodeSeparator       = 0xab // 171
	OpCheckSig            = 0xac // 172
	OpCheckSigVerify      = 0xad // 173
	OpCheckMultiSig       = 0xae // 174
	OpCheckMultiSigVerify = 0xaf // 175
	OpNoOp1               = 0xb0 // 176
	OpNoOp2               = 0xb1 // 177
	OpCheckLockTimeVerify = 0xb1 // 177 - AKA OpNoOp2
	OpNoOp3               = 0xb2 // 178
	OpCheckSequenceVerify = 0xb2 // 178 - AKA OpNoOp3
	OpNoOp4               = 0xb3 // 179
	OpNoOp5               = 0xb4 // 180
	OpNoOp6               = 0xb5 // 181
	OpNoOp7               = 0xb6 // 182
	OpNoOp8               = 0xb7 // 183
	OpNoOp9               = 0xb8 // 184
	OpNoOp10              = 0xb9 // 185
	OpUnknown186          = 0xba // 186
	OpUnknown187          = 0xbb // 187
	OpUnknown188          = 0xbc // 188
	OpUnknown189          = 0xbd // 189
	OpUnknown190          = 0xbe // 190
	OpUnknown191          = 0xbf // 191
	OpUnknown192          = 0xc0 // 192
	OpUnknown193          = 0xc1 // 193
	OpUnknown194          = 0xc2 // 194
	OpUnknown195          = 0xc3 // 195
	OpUnknown196          = 0xc4 // 196
	OpUnknown197          = 0xc5 // 197
	OpUnknown198          = 0xc6 // 198
	OpUnknown199          = 0xc7 // 199
	OpUnknown200          = 0xc8 // 200
	OpUnknown201          = 0xc9 // 201
	OpUnknown202          = 0xca // 202
	OpUnknown203          = 0xcb // 203
	OpUnknown204          = 0xcc // 204
	OpUnknown205          = 0xcd // 205
	OpUnknown206          = 0xce // 206
	OpUnknown207          = 0xcf // 207
	OpUnknown208          = 0xd0 // 208
	OpUnknown209          = 0xd1 // 209
	OpUnknown210          = 0xd2 // 210
	OpUnknown211          = 0xd3 // 211
	OpUnknown212          = 0xd4 // 212
	OpUnknown213          = 0xd5 // 213
	OpUnknown214          = 0xd6 // 214
	OpUnknown215          = 0xd7 // 215
	OpUnknown216          = 0xd8 // 216
	OpUnknown217          = 0xd9 // 217
	OpUnknown218          = 0xda // 218
	OpUnknown219          = 0xdb // 219
	OpUnknown220          = 0xdc // 220
	OpUnknown221          = 0xdd // 221
	OpUnknown222          = 0xde // 222
	OpUnknown223          = 0xdf // 223
	OpUnknown224          = 0xe0 // 224
	OpUnknown225          = 0xe1 // 225
	OpUnknown226          = 0xe2 // 226
	OpUnknown227          = 0xe3 // 227
	OpUnknown228          = 0xe4 // 228
	OpUnknown229          = 0xe5 // 229
	OpUnknown230          = 0xe6 // 230
	OpUnknown231          = 0xe7 // 231
	OpUnknown232          = 0xe8 // 232
	OpUnknown233          = 0xe9 // 233
	OpUnknown234          = 0xea // 234
	OpUnknown235          = 0xeb // 235
	OpUnknown236          = 0xec // 236
	OpUnknown237          = 0xed // 237
	OpUnknown238          = 0xee // 238
	OpUnknown239          = 0xef // 239
	OpUnknown240          = 0xf0 // 240
	OpUnknown241          = 0xf1 // 241
	OpUnknown242          = 0xf2 // 242
	OpUnknown243          = 0xf3 // 243
	OpUnknown244          = 0xf4 // 244
	OpUnknown245          = 0xf5 // 245
	OpUnknown246          = 0xf6 // 246
	OpUnknown247          = 0xf7 // 247
	OpUnknown248          = 0xf8 // 248
	OpUnknown249          = 0xf9 // 249
	OpSmallInteger        = 0xfa // 250 - bitcoin core internal
	OpPubKeys             = 0xfb // 251 - bitcoin core internal
	OpUnknown252          = 0xfc // 252
	OpPubKeyHash          = 0xfd // 253 - bitcoin core internal
	OpPubKey              = 0xfe // 254 - bitcoin core internal
	OpInvalidOpCode       = 0xff // 255 - bitcoin core internal
)

// Conditional execution constants.
const (
	OpCondFalse = 0
	OpCondTrue  = 1
	OpCondSkip  = 2
)

// opcodeArray holds details about all possible opcodes such as how many bytes the opcode and any associated data should take, its human-readable name, and the handler function.
var opcodeArray = [256]opcode{

	// Data push opcodes.
	OpFalse:     {OpFalse, "OpZero", 1, opcodeFalse},
	OpData1:     {OpData1, "OpData1", 2, opcodePushData},
	OpData2:     {OpData2, "OpData2", 3, opcodePushData},
	OpData3:     {OpData3, "OpData3", 4, opcodePushData},
	OpData4:     {OpData4, "OpData4", 5, opcodePushData},
	OpData5:     {OpData5, "OpData5", 6, opcodePushData},
	OpData6:     {OpData6, "OpData6", 7, opcodePushData},
	OpData7:     {OpData7, "OpData7", 8, opcodePushData},
	OpData8:     {OpData8, "OpData8", 9, opcodePushData},
	OpData9:     {OpData9, "OpData9", 10, opcodePushData},
	OpData10:    {OpData10, "OpData10", 11, opcodePushData},
	OpData11:    {OpData11, "OpData11", 12, opcodePushData},
	OpData12:    {OpData12, "OpData12", 13, opcodePushData},
	OpData13:    {OpData13, "OpData13", 14, opcodePushData},
	OpData14:    {OpData14, "OpData14", 15, opcodePushData},
	OpData15:    {OpData15, "OpData15", 16, opcodePushData},
	OpData16:    {OpData16, "OpData16", 17, opcodePushData},
	OpData17:    {OpData17, "OpData17", 18, opcodePushData},
	OpData18:    {OpData18, "OpData18", 19, opcodePushData},
	OpData19:    {OpData19, "OpData19", 20, opcodePushData},
	OpData20:    {OpData20, "OpData20", 21, opcodePushData},
	OpData21:    {OpData21, "OpData21", 22, opcodePushData},
	OpData22:    {OpData22, "OpData22", 23, opcodePushData},
	OpData23:    {OpData23, "OpData23", 24, opcodePushData},
	OpData24:    {OpData24, "OpData24", 25, opcodePushData},
	OpData25:    {OpData25, "OpData25", 26, opcodePushData},
	OpData26:    {OpData26, "OpData26", 27, opcodePushData},
	OpData27:    {OpData27, "OpData27", 28, opcodePushData},
	OpData28:    {OpData28, "OpData28", 29, opcodePushData},
	OpData29:    {OpData29, "OpData29", 30, opcodePushData},
	OpData30:    {OpData30, "OpData30", 31, opcodePushData},
	OpData31:    {OpData31, "OpData31", 32, opcodePushData},
	OpData32:    {OpData32, "OpData32", 33, opcodePushData},
	OpData33:    {OpData33, "OpData33", 34, opcodePushData},
	OpData34:    {OpData34, "OpData34", 35, opcodePushData},
	OpData35:    {OpData35, "OpData35", 36, opcodePushData},
	OpData36:    {OpData36, "OpData36", 37, opcodePushData},
	OpData37:    {OpData37, "OpData37", 38, opcodePushData},
	OpData38:    {OpData38, "OpData38", 39, opcodePushData},
	OpData39:    {OpData39, "OpData39", 40, opcodePushData},
	OpData40:    {OpData40, "OpData40", 41, opcodePushData},
	OpData41:    {OpData41, "OpData41", 42, opcodePushData},
	OpData42:    {OpData42, "OpData42", 43, opcodePushData},
	OpData43:    {OpData43, "OpData43", 44, opcodePushData},
	OpData44:    {OpData44, "OpData44", 45, opcodePushData},
	OpData45:    {OpData45, "OpData45", 46, opcodePushData},
	OpData46:    {OpData46, "OpData46", 47, opcodePushData},
	OpData47:    {OpData47, "OpData47", 48, opcodePushData},
	OpData48:    {OpData48, "OpData48", 49, opcodePushData},
	OpData49:    {OpData49, "OpData49", 50, opcodePushData},
	OpData50:    {OpData50, "OpData50", 51, opcodePushData},
	OpData51:    {OpData51, "OpData51", 52, opcodePushData},
	OpData52:    {OpData52, "OpData52", 53, opcodePushData},
	OpData53:    {OpData53, "OpData53", 54, opcodePushData},
	OpData54:    {OpData54, "OpData54", 55, opcodePushData},
	OpData55:    {OpData55, "OpData55", 56, opcodePushData},
	OpData56:    {OpData56, "OpData56", 57, opcodePushData},
	OpData57:    {OpData57, "OpData57", 58, opcodePushData},
	OpData58:    {OpData58, "OpData58", 59, opcodePushData},
	OpData59:    {OpData59, "OpData59", 60, opcodePushData},
	OpData60:    {OpData60, "OpData60", 61, opcodePushData},
	OpData61:    {OpData61, "OpData61", 62, opcodePushData},
	OpData62:    {OpData62, "OpData62", 63, opcodePushData},
	OpData63:    {OpData63, "OpData63", 64, opcodePushData},
	OpData64:    {OpData64, "OpData64", 65, opcodePushData},
	OpData65:    {OpData65, "OpData65", 66, opcodePushData},
	OpData66:    {OpData66, "OpData66", 67, opcodePushData},
	OpData67:    {OpData67, "OpData67", 68, opcodePushData},
	OpData68:    {OpData68, "OpData68", 69, opcodePushData},
	OpData69:    {OpData69, "OpData69", 70, opcodePushData},
	OpData70:    {OpData70, "OpData70", 71, opcodePushData},
	OpData71:    {OpData71, "OpData71", 72, opcodePushData},
	OpData72:    {OpData72, "OpData72", 73, opcodePushData},
	OpData73:    {OpData73, "OpData73", 74, opcodePushData},
	OpData74:    {OpData74, "OpData74", 75, opcodePushData},
	OpData75:    {OpData75, "OpData75", 76, opcodePushData},
	OpPushData1: {OpPushData1, "OpPushData1", -1, opcodePushData},
	OpPushData2: {OpPushData2, "OpPushData2", -2, opcodePushData},
	OpPushData4: {OpPushData4, "OpPushData4", -4, opcodePushData},
	Op1Negate:   {Op1Negate, "Op1Negate", 1, opcode1Negate},
	OpReserved:  {OpReserved, "OpReserved", 1, opcodeReserved},
	OpTrue:      {OpTrue, "Op1", 1, opcodeN},
	Op2:         {Op2, "Op2", 1, opcodeN},
	Op3:         {Op3, "Op3", 1, opcodeN},
	Op4:         {Op4, "Op4", 1, opcodeN},
	Op5:         {Op5, "Op5", 1, opcodeN},
	Op6:         {Op6, "Op6", 1, opcodeN},
	Op7:         {Op7, "Op7", 1, opcodeN},
	Op8:         {Op8, "Op8", 1, opcodeN},
	Op9:         {Op9, "Op9", 1, opcodeN},
	Op10:        {Op10, "Op10", 1, opcodeN},
	Op11:        {Op11, "Op11", 1, opcodeN},
	Op12:        {Op12, "Op12", 1, opcodeN},
	Op13:        {Op13, "Op13", 1, opcodeN},
	Op14:        {Op14, "Op14", 1, opcodeN},
	Op15:        {Op15, "Op15", 1, opcodeN},
	Op16:        {Op16, "Op16", 1, opcodeN},

	// Control opcodes.
	OpNoOp:                {OpNoOp, "OpNoOp", 1, opcodeNop},
	OpVer:                 {OpVer, "OpVer", 1, opcodeReserved},
	OpIf:                  {OpIf, "OpIf", 1, opcodeIf},
	OpIfNot:               {OpIfNot, "OpIfNot", 1, opcodeNotIf},
	OpVerIf:               {OpVerIf, "OpVerIf", 1, opcodeReserved},
	OpVerIfNot:            {OpVerIfNot, "OpVerIfNot", 1, opcodeReserved},
	OpElse:                {OpElse, "OpElse", 1, opcodeElse},
	OpEndIf:               {OpEndIf, "OpEndIf", 1, opcodeEndif},
	OpVerify:              {OpVerify, "OpVerify", 1, opcodeVerify},
	OpReturn:              {OpReturn, "OpReturn", 1, opcodeReturn},
	OpCheckLockTimeVerify: {OpCheckLockTimeVerify, "OpCheckLockTimeVerify", 1, opcodeCheckLockTimeVerify},
	OpCheckSequenceVerify: {OpCheckSequenceVerify, "OpCheckSequenceVerify", 1, opcodeCheckSequenceVerify},

	// Stack opcodes.
	OpToAltStack:   {OpToAltStack, "OpToAltStack", 1, opcodeToAltStack},
	OpFromAltStack: {OpFromAltStack, "OpFromAltStack", 1, opcodeFromAltStack},
	Op2Drop:        {Op2Drop, "Op2Drop", 1, opcode2Drop},
	Op2Dup:         {Op2Dup, "Op2Dup", 1, opcode2Dup},
	Op3Dup:         {Op3Dup, "Op3Dup", 1, opcode3Dup},
	Op2Over:        {Op2Over, "Op2Over", 1, opcode2Over},
	Op2Rot:         {Op2Rot, "Op2Rot", 1, opcode2Rot},
	Op2Swap:        {Op2Swap, "Op2Swap", 1, opcode2Swap},
	OpIfDup:        {OpIfDup, "OpIfDup", 1, opcodeIfDup},
	OpDepth:        {OpDepth, "OpDepth", 1, opcodeDepth},
	OpDrop:         {OpDrop, "OpDrop", 1, opcodeDrop},
	OpDup:          {OpDup, "OpDup", 1, opcodeDup},
	OpNip:          {OpNip, "OpNip", 1, opcodeNip},
	OpOver:         {OpOver, "OpOver", 1, opcodeOver},
	OpPick:         {OpPick, "OpPick", 1, opcodePick},
	OpRoll:         {OpRoll, "OpRoll", 1, opcodeRoll},
	OpRot:          {OpRot, "OpRot", 1, opcodeRot},
	OpSwap:         {OpSwap, "OpSwap", 1, opcodeSwap},
	OpTuck:         {OpTuck, "OpTuck", 1, opcodeTuck},

	// Splice opcodes.
	OpCat:    {OpCat, "OpCat", 1, opcodeDisabled},
	OpSubstr: {OpSubstr, "OpSubstr", 1, opcodeDisabled},
	OpLeft:   {OpLeft, "OpLeft", 1, opcodeDisabled},
	OpRight:  {OpRight, "OpRight", 1, opcodeDisabled},
	OpSize:   {OpSize, "OpSize", 1, opcodeSize},

	// Bitwise logic opcodes.
	OpInvert:      {OpInvert, "OpInvert", 1, opcodeDisabled},
	OpAnd:         {OpAnd, "OpAnd", 1, opcodeDisabled},
	OpOr:          {OpOr, "OpOr", 1, opcodeDisabled},
	OpXor:         {OpXor, "OpXor", 1, opcodeDisabled},
	OpEqual:       {OpEqual, "OpEqual", 1, opcodeEqual},
	OpEqualVerify: {OpEqualVerify, "OpEqualVerify", 1, opcodeEqualVerify},
	OpReserved1:   {OpReserved1, "OpReserved1", 1, opcodeReserved},
	OpReserved2:   {OpReserved2, "OpReserved2", 1, opcodeReserved},

	// Numeric related opcodes.
	Op1Add:               {Op1Add, "Op1Add", 1, opcode1Add},
	Op1Sub:               {Op1Sub, "Op1Sub", 1, opcode1Sub},
	Op1Mul:               {Op1Mul, "Op1Mul", 1, opcodeDisabled},
	Op2Div:               {Op2Div, "Op2Div", 1, opcodeDisabled},
	OpNegate:             {OpNegate, "OpNegate", 1, opcodeNegate},
	OpAbs:                {OpAbs, "OpAbs", 1, opcodeAbs},
	OpNot:                {OpNot, "OpNot", 1, opcodeNot},
	Op0NotEqual:          {Op0NotEqual, "Op0NotEqual", 1, opcode0NotEqual},
	OpAdd:                {OpAdd, "OpAdd", 1, opcodeAdd},
	OpSub:                {OpSub, "OpSub", 1, opcodeSub},
	OpMul:                {OpMul, "OpMul", 1, opcodeDisabled},
	OpDiv:                {OpDiv, "OpDiv", 1, opcodeDisabled},
	OpMod:                {OpMod, "OpMod", 1, opcodeDisabled},
	OpLShift:             {OpLShift, "OpLShift", 1, opcodeDisabled},
	OpRShift:             {OpRShift, "OpRShift", 1, opcodeDisabled},
	OpBoolAnd:            {OpBoolAnd, "OpBoolAnd", 1, opcodeBoolAnd},
	OpBoolOr:             {OpBoolOr, "OpBoolOr", 1, opcodeBoolOr},
	OpNumEqual:           {OpNumEqual, "OpNumEqual", 1, opcodeNumEqual},
	OpNumEqualVerify:     {OpNumEqualVerify, "OpNumEqualVerify", 1, opcodeNumEqualVerify},
	OpNumNotEqual:        {OpNumNotEqual, "OpNumNotEqual", 1, opcodeNumNotEqual},
	OpLessThan:           {OpLessThan, "OpLessThan", 1, opcodeLessThan},
	OpGreaterThan:        {OpGreaterThan, "OpGreaterThan", 1, opcodeGreaterThan},
	OpLessThanOrEqual:    {OpLessThanOrEqual, "OpLessThanOrEqual", 1, opcodeLessThanOrEqual},
	OpGreaterThanOrEqual: {OpGreaterThanOrEqual, "OpGreaterThanOrEqual", 1, opcodeGreaterThanOrEqual},
	OpMin:                {OpMin, "OpMin", 1, opcodeMin},
	OpMax:                {OpMax, "OpMax", 1, opcodeMax},
	OpWithin:             {OpWithin, "OpWithin", 1, opcodeWithin},

	// Crypto opcodes.
	OpRipeMD160:           {OpRipeMD160, "OpRipeMD160", 1, opcodeRipemd160},
	OpSHA1:                {OpSHA1, "OpSHA1", 1, opcodeSha1},
	OpSHA256:              {OpSHA256, "OpSHA256", 1, opcodeSha256},
	OpHash160:             {OpHash160, "OpHash160", 1, opcodeHash160},
	OpHash256:             {OpHash256, "OpHash256", 1, opcodeHash256},
	OpCodeSeparator:       {OpCodeSeparator, "OpCodeSeparator", 1, opcodeCodeSeparator},
	OpCheckSig:            {OpCheckSig, "OpCheckSig", 1, opcodeCheckSig},
	OpCheckSigVerify:      {OpCheckSigVerify, "OpCheckSigVerify", 1, opcodeCheckSigVerify},
	OpCheckMultiSig:       {OpCheckMultiSig, "OpCheckMultiSig", 1, opcodeCheckMultiSig},
	OpCheckMultiSigVerify: {OpCheckMultiSigVerify, "OpCheckMultiSigVerify", 1, opcodeCheckMultiSigVerify},

	// Reserved opcodes.
	OpNoOp1:  {OpNoOp1, "OpNoOp1", 1, opcodeNop},
	OpNoOp4:  {OpNoOp4, "OpNoOp4", 1, opcodeNop},
	OpNoOp5:  {OpNoOp5, "OpNoOp5", 1, opcodeNop},
	OpNoOp6:  {OpNoOp6, "OpNoOp6", 1, opcodeNop},
	OpNoOp7:  {OpNoOp7, "OpNoOp7", 1, opcodeNop},
	OpNoOp8:  {OpNoOp8, "OpNoOp8", 1, opcodeNop},
	OpNoOp9:  {OpNoOp9, "OpNoOp9", 1, opcodeNop},
	OpNoOp10: {OpNoOp10, "OpNoOp10", 1, opcodeNop},

	// Undefined opcodes.
	OpUnknown186: {OpUnknown186, "OpUnknown186", 1, opcodeInvalid},
	OpUnknown187: {OpUnknown187, "OpUnknown187", 1, opcodeInvalid},
	OpUnknown188: {OpUnknown188, "OpUnknown188", 1, opcodeInvalid},
	OpUnknown189: {OpUnknown189, "OpUnknown189", 1, opcodeInvalid},
	OpUnknown190: {OpUnknown190, "OpUnknown190", 1, opcodeInvalid},
	OpUnknown191: {OpUnknown191, "OpUnknown191", 1, opcodeInvalid},
	OpUnknown192: {OpUnknown192, "OpUnknown192", 1, opcodeInvalid},
	OpUnknown193: {OpUnknown193, "OpUnknown193", 1, opcodeInvalid},
	OpUnknown194: {OpUnknown194, "OpUnknown194", 1, opcodeInvalid},
	OpUnknown195: {OpUnknown195, "OpUnknown195", 1, opcodeInvalid},
	OpUnknown196: {OpUnknown196, "OpUnknown196", 1, opcodeInvalid},
	OpUnknown197: {OpUnknown197, "OpUnknown197", 1, opcodeInvalid},
	OpUnknown198: {OpUnknown198, "OpUnknown198", 1, opcodeInvalid},
	OpUnknown199: {OpUnknown199, "OpUnknown199", 1, opcodeInvalid},
	OpUnknown200: {OpUnknown200, "OpUnknown200", 1, opcodeInvalid},
	OpUnknown201: {OpUnknown201, "OpUnknown201", 1, opcodeInvalid},
	OpUnknown202: {OpUnknown202, "OpUnknown202", 1, opcodeInvalid},
	OpUnknown203: {OpUnknown203, "OpUnknown203", 1, opcodeInvalid},
	OpUnknown204: {OpUnknown204, "OpUnknown204", 1, opcodeInvalid},
	OpUnknown205: {OpUnknown205, "OpUnknown205", 1, opcodeInvalid},
	OpUnknown206: {OpUnknown206, "OpUnknown206", 1, opcodeInvalid},
	OpUnknown207: {OpUnknown207, "OpUnknown207", 1, opcodeInvalid},
	OpUnknown208: {OpUnknown208, "OpUnknown208", 1, opcodeInvalid},
	OpUnknown209: {OpUnknown209, "OpUnknown209", 1, opcodeInvalid},
	OpUnknown210: {OpUnknown210, "OpUnknown210", 1, opcodeInvalid},
	OpUnknown211: {OpUnknown211, "OpUnknown211", 1, opcodeInvalid},
	OpUnknown212: {OpUnknown212, "OpUnknown212", 1, opcodeInvalid},
	OpUnknown213: {OpUnknown213, "OpUnknown213", 1, opcodeInvalid},
	OpUnknown214: {OpUnknown214, "OpUnknown214", 1, opcodeInvalid},
	OpUnknown215: {OpUnknown215, "OpUnknown215", 1, opcodeInvalid},
	OpUnknown216: {OpUnknown216, "OpUnknown216", 1, opcodeInvalid},
	OpUnknown217: {OpUnknown217, "OpUnknown217", 1, opcodeInvalid},
	OpUnknown218: {OpUnknown218, "OpUnknown218", 1, opcodeInvalid},
	OpUnknown219: {OpUnknown219, "OpUnknown219", 1, opcodeInvalid},
	OpUnknown220: {OpUnknown220, "OpUnknown220", 1, opcodeInvalid},
	OpUnknown221: {OpUnknown221, "OpUnknown221", 1, opcodeInvalid},
	OpUnknown222: {OpUnknown222, "OpUnknown222", 1, opcodeInvalid},
	OpUnknown223: {OpUnknown223, "OpUnknown223", 1, opcodeInvalid},
	OpUnknown224: {OpUnknown224, "OpUnknown224", 1, opcodeInvalid},
	OpUnknown225: {OpUnknown225, "OpUnknown225", 1, opcodeInvalid},
	OpUnknown226: {OpUnknown226, "OpUnknown226", 1, opcodeInvalid},
	OpUnknown227: {OpUnknown227, "OpUnknown227", 1, opcodeInvalid},
	OpUnknown228: {OpUnknown228, "OpUnknown228", 1, opcodeInvalid},
	OpUnknown229: {OpUnknown229, "OpUnknown229", 1, opcodeInvalid},
	OpUnknown230: {OpUnknown230, "OpUnknown230", 1, opcodeInvalid},
	OpUnknown231: {OpUnknown231, "OpUnknown231", 1, opcodeInvalid},
	OpUnknown232: {OpUnknown232, "OpUnknown232", 1, opcodeInvalid},
	OpUnknown233: {OpUnknown233, "OpUnknown233", 1, opcodeInvalid},
	OpUnknown234: {OpUnknown234, "OpUnknown234", 1, opcodeInvalid},
	OpUnknown235: {OpUnknown235, "OpUnknown235", 1, opcodeInvalid},
	OpUnknown236: {OpUnknown236, "OpUnknown236", 1, opcodeInvalid},
	OpUnknown237: {OpUnknown237, "OpUnknown237", 1, opcodeInvalid},
	OpUnknown238: {OpUnknown238, "OpUnknown238", 1, opcodeInvalid},
	OpUnknown239: {OpUnknown239, "OpUnknown239", 1, opcodeInvalid},
	OpUnknown240: {OpUnknown240, "OpUnknown240", 1, opcodeInvalid},
	OpUnknown241: {OpUnknown241, "OpUnknown241", 1, opcodeInvalid},
	OpUnknown242: {OpUnknown242, "OpUnknown242", 1, opcodeInvalid},
	OpUnknown243: {OpUnknown243, "OpUnknown243", 1, opcodeInvalid},
	OpUnknown244: {OpUnknown244, "OpUnknown244", 1, opcodeInvalid},
	OpUnknown245: {OpUnknown245, "OpUnknown245", 1, opcodeInvalid},
	OpUnknown246: {OpUnknown246, "OpUnknown246", 1, opcodeInvalid},
	OpUnknown247: {OpUnknown247, "OpUnknown247", 1, opcodeInvalid},
	OpUnknown248: {OpUnknown248, "OpUnknown248", 1, opcodeInvalid},
	OpUnknown249: {OpUnknown249, "OpUnknown249", 1, opcodeInvalid},

	// Bitcoin Core internal use opcode.  Defined here for completeness.
	OpSmallInteger:  {OpSmallInteger, "OpSmallInteger", 1, opcodeInvalid},
	OpPubKeys:       {OpPubKeys, "OpPubKeys", 1, opcodeInvalid},
	OpUnknown252:    {OpUnknown252, "OpUnknown252", 1, opcodeInvalid},
	OpPubKeyHash:    {OpPubKeyHash, "OpPubKeyHash", 1, opcodeInvalid},
	OpPubKey:        {OpPubKey, "OpPubKey", 1, opcodeInvalid},
	OpInvalidOpCode: {OpInvalidOpCode, "OpInvalidOpCode", 1, opcodeInvalid},
}

// opcodeOnelineRepls defines opcode names which are replaced when doing a one-line disassembly.  This is done to match the output of the reference implementation while not changing the opcode names in the nicer full disassembly.
var opcodeOnelineRepls = map[string]string{
	"Op1Negate": "-1",
	"OpZero":    "0",
	"Op1":       "1",
	"Op2":       "2",
	"Op3":       "3",
	"Op4":       "4",
	"Op5":       "5",
	"Op6":       "6",
	"Op7":       "7",
	"Op8":       "8",
	"Op9":       "9",
	"Op10":      "10",
	"Op11":      "11",
	"Op12":      "12",
	"Op13":      "13",
	"Op14":      "14",
	"Op15":      "15",
	"Op16":      "16",
}

// parsedOpcode represents an opcode that has been parsed and includes any potential data associated with it.

type parsedOpcode struct {
	opcode *opcode
	data   []byte
}

// isDisabled returns whether or not the opcode is disabled and thus is always bad to see in the instruction stream (even if turned off by a conditional).
func (pop *parsedOpcode) isDisabled() bool {

	switch pop.opcode.value {

	case OpCat:
		return true
	case OpSubstr:
		return true
	case OpLeft:
		return true
	case OpRight:
		return true
	case OpInvert:
		return true
	case OpAnd:
		return true
	case OpOr:
		return true
	case OpXor:
		return true
	case Op1Mul:
		return true
	case Op2Div:
		return true
	case OpMul:
		return true
	case OpDiv:
		return true
	case OpMod:
		return true
	case OpLShift:
		return true
	case OpRShift:
		return true
	default:
		return false
	}
}

// alwaysIllegal returns whether or not the opcode is always illegal when passed over by the program counter even if in a non-executed branch (it isn't a coincidence that they are conditionals).
func (pop *parsedOpcode) alwaysIllegal() bool {

	switch pop.opcode.value {

	case OpVerIf:
		return true
	case OpVerIfNot:
		return true
	default:
		return false
	}
}

// isConditional returns whether or not the opcode is a conditional opcode which changes the conditional execution stack when executed.
func (pop *parsedOpcode) isConditional() bool {

	switch pop.opcode.value {

	case OpIf:
		return true
	case OpIfNot:
		return true
	case OpElse:
		return true
	case OpEndIf:
		return true
	default:
		return false
	}
}

// checkMinimalDataPush returns whether or not the current data push uses the smallest possible opcode to represent it.  For example, the value 15 could be pushed with OpData1 15 (among other variations); however, Op15 is a single opcode that represents the same value and is only a single byte versus two bytes.
func (pop *parsedOpcode) checkMinimalDataPush() error {

	data := pop.data
	dataLen := len(data)
	opcode := pop.opcode.value
	if dataLen == 0 && opcode != OpZero {

		str := fmt.Sprintf("zero length data push is encoded with "+
			"opcode %s instead of OpZero", pop.opcode.name)
		return scriptError(ErrMinimalData, str)
	} else if dataLen == 1 && data[0] >= 1 && data[0] <= 16 {

		if opcode != Op1+data[0]-1 {

			// Should have used Op1 .. Op16
			str := fmt.Sprintf("data push of the value %d encoded "+
				"with opcode %s instead of OP_%d", data[0],
				pop.opcode.name, data[0])
			return scriptError(ErrMinimalData, str)
		}
	} else if dataLen == 1 && data[0] == 0x81 {

		if opcode != Op1Negate {

			str := fmt.Sprintf("data push of the value -1 encoded "+
				"with opcode %s instead of Op1Negate",
				pop.opcode.name)
			return scriptError(ErrMinimalData, str)
		}
	} else if dataLen <= 75 {

		if int(opcode) != dataLen {

			// Should have used a direct push
			str := fmt.Sprintf("data push of %d bytes encoded "+
				"with opcode %s instead of OpData%d", dataLen,
				pop.opcode.name, dataLen)
			return scriptError(ErrMinimalData, str)
		}
	} else if dataLen <= 255 {

		if opcode != OpPushData1 {

			str := fmt.Sprintf("data push of %d bytes encoded "+
				"with opcode %s instead of OpPushData1",
				dataLen, pop.opcode.name)
			return scriptError(ErrMinimalData, str)
		}
	} else if dataLen <= 65535 {

		if opcode != OpPushData2 {

			str := fmt.Sprintf("data push of %d bytes encoded "+
				"with opcode %s instead of OpPushData2",
				dataLen, pop.opcode.name)
			return scriptError(ErrMinimalData, str)
		}
	}
	return nil
}

// print returns a human-readable string representation of the opcode for use in script disassembly.
func (pop *parsedOpcode) print(oneline bool) string {

	// The reference implementation one-line disassembly replaces opcodes which represent values (e.g. OpZero through Op16 and Op1Negate) with the raw value.  However, when not doing a one-line dissassembly, we prefer to show the actual opcode names.  Thus, only replace the opcodes in question when the oneline flag is set.
	opcodeName := pop.opcode.name
	if oneline {

		if replName, ok := opcodeOnelineRepls[opcodeName]; ok {

			opcodeName = replName
		}
		// Nothing more to do for non-data push opcodes.

		if pop.opcode.length == 1 {

			return opcodeName
		}
		return fmt.Sprintf("%x", pop.data)
	}

	// Nothing more to do for non-data push opcodes.
	if pop.opcode.length == 1 {

		return opcodeName
	}

	// Add length for the OP_PUSHDATA# opcodes.
	retString := opcodeName

	switch pop.opcode.length {

	case -1:
		retString += fmt.Sprintf(" 0x%02x", len(pop.data))
	case -2:
		retString += fmt.Sprintf(" 0x%04x", len(pop.data))
	case -4:
		retString += fmt.Sprintf(" 0x%08x", len(pop.data))
	}
	return fmt.Sprintf("%s 0x%02x", retString, pop.data)
}

// bytes returns any data associated with the opcode encoded as it would be in a script.  This is used for unparsing scripts from parsed opcodes.
func (pop *parsedOpcode) bytes() ([]byte, error) {

	var retbytes []byte
	if pop.opcode.length > 0 {

		retbytes = make([]byte, 1, pop.opcode.length)
	} else {

		retbytes = make([]byte, 1, 1+len(pop.data)-
			pop.opcode.length)
	}
	retbytes[0] = pop.opcode.value
	if pop.opcode.length == 1 {

		if len(pop.data) != 0 {

			str := fmt.Sprintf("internal consistency error - "+
				"parsed opcode %s has data length %d when %d "+
				"was expected", pop.opcode.name, len(pop.data),
				0)
			return nil, scriptError(ErrInternal, str)
		}
		return retbytes, nil
	}
	nbytes := pop.opcode.length
	if pop.opcode.length < 0 {

		l := len(pop.data)
		// tempting just to hardcode to avoid the complexity here.

		switch pop.opcode.length {

		case -1:
			retbytes = append(retbytes, byte(l))
			nbytes = int(retbytes[1]) + len(retbytes)
		case -2:
			retbytes = append(retbytes, byte(l&0xff),
				byte(l>>8&0xff))
			nbytes = int(binary.LittleEndian.Uint16(retbytes[1:])) +
				len(retbytes)
		case -4:
			retbytes = append(retbytes, byte(l&0xff),
				byte((l>>8)&0xff), byte((l>>16)&0xff),
				byte((l>>24)&0xff))
			nbytes = int(binary.LittleEndian.Uint32(retbytes[1:])) +
				len(retbytes)
		}
	}
	retbytes = append(retbytes, pop.data...)
	if len(retbytes) != nbytes {

		str := fmt.Sprintf("internal consistency error - "+
			"parsed opcode %s has data length %d when %d was "+
			"expected", pop.opcode.name, len(retbytes), nbytes)
		return nil, scriptError(ErrInternal, str)
	}
	return retbytes, nil
}

// Opcode implementation functions start here.

// opcodeDisabled is a common handler for disabled opcodes.  It returns an appropriate error indicating the opcode is disabled.  While it would ordinarily make more sense to detect if the script contains any disabled opcodes before executing in an initial parse step, the consensus rules dictate the script doesn't fail until the program counter passes over a disabled opcode (even when they appear in a branch that is not executed).
func opcodeDisabled(
	op *parsedOpcode, vm *Engine) error {

	str := fmt.Sprintf("attempt to execute disabled opcode %s",
		op.opcode.name)
	return scriptError(ErrDisabledOpcode, str)
}

// opcodeReserved is a common handler for all reserved opcodes.  It returns an appropriate error indicating the opcode is reserved.
func opcodeReserved(
	op *parsedOpcode, vm *Engine) error {

	str := fmt.Sprintf("attempt to execute reserved opcode %s",
		op.opcode.name)
	return scriptError(ErrReservedOpcode, str)
}

// opcodeInvalid is a common handler for all invalid opcodes.  It returns an appropriate error indicating the opcode is invalid.
func opcodeInvalid(
	op *parsedOpcode, vm *Engine) error {

	str := fmt.Sprintf("attempt to execute invalid opcode %s",
		op.opcode.name)
	return scriptError(ErrReservedOpcode, str)
}

// opcodeFalse pushes an empty array to the data stack to represent false.  Note that 0, when encoded as a number according to the numeric encoding consensus rules, is an empty array.
func opcodeFalse(
	op *parsedOpcode, vm *Engine) error {

	vm.dstack.PushByteArray(nil)
	return nil
}

// opcodePushData is a common handler for the vast majority of opcodes that push raw data (bytes) to the data stack.
func opcodePushData(
	op *parsedOpcode, vm *Engine) error {

	vm.dstack.PushByteArray(op.data)
	return nil
}

// opcode1Negate pushes -1, encoded as a number, to the data stack.
func opcode1Negate(
	op *parsedOpcode, vm *Engine) error {

	vm.dstack.PushInt(scriptNum(-1))
	return nil
}

// opcodeN is a common handler for the small integer data push opcodes.  It pushes the numeric value the opcode represents (which will be from 1 to 16) onto the data stack.
func opcodeN(
	op *parsedOpcode, vm *Engine) error {

	// The opcodes are all defined consecutively, so the numeric value is the difference.
	vm.dstack.PushInt(scriptNum((op.opcode.value - (Op1 - 1))))
	return nil
}

// opcodeNop is a common handler for the NOP family of opcodes.  As the name implies it generally does nothing, however, it will return an error when the flag to discourage use of NOPs is set for select opcodes.
func opcodeNop(
	op *parsedOpcode, vm *Engine) error {

	switch op.opcode.value {

	case OpNoOp1, OpNoOp4, OpNoOp5,
		OpNoOp6, OpNoOp7, OpNoOp8, OpNoOp9, OpNoOp10:

		if vm.hasFlag(ScriptDiscourageUpgradableNops) {

			str := fmt.Sprintf("OpNoOp%d reserved for soft-fork "+
				"upgrades", op.opcode.value-(OpNoOp1-1))
			return scriptError(ErrDiscourageUpgradableNOPs, str)
		}
	}
	return nil
}

// popIfBool enforces the "minimal if" policy during script execution if the particular flag is set.  If so, in order to eliminate an additional source of nuisance malleability, post-segwit for version 0 witness programs, we now require the following: for OpIf and OP_NOT_IF, the top stack item MUST either be an empty byte slice, or [0x01]. Otherwise, the item at the top of the stack will be popped and interpreted as a boolean.
func popIfBool(
	vm *Engine) (bool, error) {

	// When not in witness execution mode, not executing a v0 witness program, or the minimal if flag isn't set pop the top stack item as a normal bool.
	if !vm.isWitnessVersionActive(0) || !vm.hasFlag(ScriptVerifyMinimalIf) {

		return vm.dstack.PopBool()
	}

	// At this point, a v0 witness program is being executed and the minimal if flag is set, so enforce additional constraints on the top stack item.
	so, err := vm.dstack.PopByteArray()
	if err != nil {

		return false, err
	}

	// The top element MUST have a length of at least one.
	if len(so) > 1 {

		str := fmt.Sprintf("minimal if is active, top element MUST "+
			"have a length of at least, instead length is %v",
			len(so))
		return false, scriptError(ErrMinimalIf, str)
	}

	// Additionally, if the length is one, then the value MUST be 0x01.
	if len(so) == 1 && so[0] != 0x01 {

		str := fmt.Sprintf("minimal if is active, top stack item MUST "+
			"be an empty byte array or 0x01, is instead: %v",
			so[0])
		return false, scriptError(ErrMinimalIf, str)
	}
	return asBool(so), nil
}

// opcodeIf treats the top item on the data stack as a boolean and removes it. An appropriate entry is added to the conditional stack depending on whether the boolean is true and whether this if is on an executing branch in order to allow proper execution of further opcodes depending on the conditional logic.  When the boolean is true, the first branch will be executed (unless this opcode is nested in a non-executed branch). <expression> if [statements] [else [statements]] endif Note that, unlike for all non-conditional opcodes, this is executed even when it is on a non-executing branch so proper nesting is maintained.
// Data stack transformation: [... bool] -> [...]
// Conditional stack transformation: [...] -> [... OpCondValue]
func opcodeIf(
	op *parsedOpcode, vm *Engine) error {

	condVal := OpCondFalse
	if vm.isBranchExecuting() {

		ok, err := popIfBool(vm)

		if err != nil {

			return err
		}

		if ok {

			condVal = OpCondTrue
		}
	} else {

		condVal = OpCondSkip
	}
	vm.condStack = append(vm.condStack, condVal)
	return nil
}

// opcodeNotIf treats the top item on the data stack as a boolean and removes it.
// An appropriate entry is added to the conditional stack depending on whether the boolean is true and whether this if is on an executing branch in order to allow proper execution of further opcodes depending on the conditional logic.  When the boolean is false, the first branch will be executed (unless this opcode is nested in a non-executed branch). <expression> notif [statements] [else [statements]] endif Note that, unlike for all non-conditional opcodes, this is executed even when it is on a non-executing branch so proper nesting is maintained.
// Data stack transformation: [... bool] -> [...]
// Conditional stack transformation: [...] -> [... OpCondValue]
func opcodeNotIf(
	op *parsedOpcode, vm *Engine) error {

	condVal := OpCondFalse
	if vm.isBranchExecuting() {

		ok, err := popIfBool(vm)

		if err != nil {

			return err
		}

		if !ok {

			condVal = OpCondTrue
		}
	} else {

		condVal = OpCondSkip
	}
	vm.condStack = append(vm.condStack, condVal)
	return nil
}

// opcodeElse inverts conditional execution for other half of if/else/endif.
// An error is returned if there has not already been a matching OpIf. Conditional stack transformation: [... OpCondValue] -> [... !OpCondValue]
func opcodeElse(
	op *parsedOpcode, vm *Engine) error {

	if len(vm.condStack) == 0 {

		str := fmt.Sprintf("encountered opcode %s with no matching "+
			"opcode to begin conditional execution", op.opcode.name)
		return scriptError(ErrUnbalancedConditional, str)
	}
	conditionalIdx := len(vm.condStack) - 1

	switch vm.condStack[conditionalIdx] {

	case OpCondTrue:
		vm.condStack[conditionalIdx] = OpCondFalse
	case OpCondFalse:
		vm.condStack[conditionalIdx] = OpCondTrue
	case OpCondSkip:
		// Value doesn't change in skip since it indicates this opcode is nested in a non-executed branch.
	}
	return nil
}

// opcodeEndif terminates a conditional block, removing the value from the conditional execution stack.
// An error is returned if there has not already been a matching OpIf.
// Conditional stack transformation: [... OpCondValue] -> [...]
func opcodeEndif(
	op *parsedOpcode, vm *Engine) error {

	if len(vm.condStack) == 0 {

		str := fmt.Sprintf("encountered opcode %s with no matching "+
			"opcode to begin conditional execution", op.opcode.name)
		return scriptError(ErrUnbalancedConditional, str)
	}
	vm.condStack = vm.condStack[:len(vm.condStack)-1]
	return nil
}

// abstractVerify examines the top item on the data stack as a boolean value and verifies it evaluates to true.
// An error is returned either when there is no item on the stack or when that item evaluates to false.  In the latter case where the verification fails specifically due to the top item evaluating to false, the returned error will use the passed error code.
func abstractVerify(
	op *parsedOpcode, vm *Engine, c ErrorCode) error {

	verified, err := vm.dstack.PopBool()
	if err != nil {

		return err
	}
	if !verified {

		str := fmt.Sprintf("%s failed", op.opcode.name)
		return scriptError(c, str)
	}
	return nil
}

// opcodeVerify examines the top item on the data stack as a boolean value and verifies it evaluates to true.  An error is returned if it does not.
func opcodeVerify(
	op *parsedOpcode, vm *Engine) error {

	return abstractVerify(op, vm, ErrVerify)
}

// opcodeReturn returns an appropriate error since it is always an error to return early from a script.
func opcodeReturn(
	op *parsedOpcode, vm *Engine) error {

	return scriptError(ErrEarlyReturn, "script returned early")
}

// verifyLockTime is a helper function used to validate locktimes.
func verifyLockTime(
	txLockTime, threshold, lockTime int64) error {

	// The lockTimes in both the script and transaction must be of the same type.
	if !((txLockTime < threshold && lockTime < threshold) ||
		(txLockTime >= threshold && lockTime >= threshold)) {

		str := fmt.Sprintf("mismatched locktime types -- tx locktime "+
			"%d, stack locktime %d", txLockTime, lockTime)
		return scriptError(ErrUnsatisfiedLockTime, str)
	}
	if lockTime > txLockTime {

		str := fmt.Sprintf("locktime requirement not satisfied -- "+
			"locktime is greater than the transaction locktime: "+
			"%d > %d", lockTime, txLockTime)
		return scriptError(ErrUnsatisfiedLockTime, str)
	}
	return nil
}

// opcodeCheckLockTimeVerify compares the top item on the data stack to the LockTime field of the transaction containing the script signature validating if the transaction outputs are spendable yet.  If flag ScriptVerifyCheckLockTimeVerify is not set, the code continues as if OpNoOp2 were executed.
func opcodeCheckLockTimeVerify(
	op *parsedOpcode, vm *Engine) error {

	// If the ScriptVerifyCheckLockTimeVerify script flag is not set, treat opcode as OpNoOp2 instead.
	if !vm.hasFlag(ScriptVerifyCheckLockTimeVerify) {

		if vm.hasFlag(ScriptDiscourageUpgradableNops) {

			return scriptError(ErrDiscourageUpgradableNOPs,
				"OpNoOp2 reserved for soft-fork upgrades")
		}
		return nil
	}

	// The current transaction locktime is a uint32 resulting in a maximum locktime of 2^32-1 (the year 2106).  However, scriptNums are signed and therefore a standard 4-byte scriptNum would only support up to a maximum of 2^31-1 (the year 2038).  Thus, a 5-byte scriptNum is used here since it will support up to 2^39-1 which allows dates beyond the current locktime limit.

	// PeekByteArray is used here instead of PeekInt because we do not want to be limited to a 4-byte integer for reasons specified above.
	so, err := vm.dstack.PeekByteArray(0)
	if err != nil {

		return err
	}
	lockTime, err := makeScriptNum(so, vm.dstack.verifyMinimalData, 5)
	if err != nil {

		return err
	}

	// In the rare event that the argument needs to be < 0 due to some arithmetic being done first, you can always use 0 OpMax OpCheckLockTimeVerify.
	if lockTime < 0 {

		str := fmt.Sprintf("negative lock time: %d", lockTime)
		return scriptError(ErrNegativeLockTime, str)
	}

	// The lock time field of a transaction is either a block height at which the transaction is finalized or a timestamp depending on if the value is before the txscript.LockTimeThreshold.  When it is under the threshold it is a block height.
	err = verifyLockTime(int64(vm.tx.LockTime), LockTimeThreshold,
		int64(lockTime))
	if err != nil {

		return err
	}

	// The lock time feature can also be disabled, thereby bypassing OpCheckLockTimeVerify, if every transaction input has been finalized by setting its sequence to the maximum value (wire.MaxTxInSequenceNum).  This condition would result in the transaction being allowed into the blockchain making the opcode ineffective.

	// This condition is prevented by enforcing that the input being used by the opcode is unlocked (its sequence number is less than the max value).  This is sufficient to prove correctness without having to check every input.

	// NOTE: This implies that even if the transaction is not finalized due to another input being unlocked, the opcode execution will still fail when the input being used by the opcode is locked.
	if vm.tx.TxIn[vm.txIdx].Sequence == wire.MaxTxInSequenceNum {

		return scriptError(ErrUnsatisfiedLockTime,
			"transaction input is finalized")
	}
	return nil
}

// opcodeCheckSequenceVerify compares the top item on the data stack to the LockTime field of the transaction containing the script signature validating if the transaction outputs are spendable yet.  If flag ScriptVerifyCheckSequenceVerify is not set, the code continues as if OpNoOp3 were executed.
func opcodeCheckSequenceVerify(
	op *parsedOpcode, vm *Engine) error {

	// If the ScriptVerifyCheckSequenceVerify script flag is not set, treat opcode as OpNoOp3 instead.
	if !vm.hasFlag(ScriptVerifyCheckSequenceVerify) {

		if vm.hasFlag(ScriptDiscourageUpgradableNops) {

			return scriptError(ErrDiscourageUpgradableNOPs,
				"OpNoOp3 reserved for soft-fork upgrades")
		}
		return nil
	}

	// The current transaction sequence is a uint32 resulting in a maximum sequence of 2^32-1.  However, scriptNums are signed and therefore a standard 4-byte scriptNum would only support up to a maximum of 2^31-1.  Thus, a 5-byte scriptNum is used here since it will support up to 2^39-1 which allows sequences beyond the current sequence limit.

	// PeekByteArray is used here instead of PeekInt because we do not want to be limited to a 4-byte integer for reasons specified above.
	so, err := vm.dstack.PeekByteArray(0)
	if err != nil {

		return err
	}
	stackSequence, err := makeScriptNum(so, vm.dstack.verifyMinimalData, 5)
	if err != nil {

		return err
	}

	// In the rare event that the argument needs to be < 0 due to some arithmetic being done first, you can always use 0 OpMax OpCheckSequenceVerify.
	if stackSequence < 0 {

		str := fmt.Sprintf("negative sequence: %d", stackSequence)
		return scriptError(ErrNegativeLockTime, str)
	}
	sequence := int64(stackSequence)

	// To provide for future soft-fork extensibility, if the operand has the disabled lock-time flag set, CHECKSEQUENCEVERIFY behaves as a NOP.
	if sequence&int64(wire.SequenceLockTimeDisabled) != 0 {

		return nil
	}

	// Transaction version numbers not high enough to trigger CSV rules must fail.
	if vm.tx.Version < 2 {

		str := fmt.Sprintf("invalid transaction version: %d",
			vm.tx.Version)
		return scriptError(ErrUnsatisfiedLockTime, str)
	}

	// Sequence numbers with their most significant bit set are not consensus constrained. Testing that the transaction's sequence number does not have this bit set prevents using this property to get around a CHECKSEQUENCEVERIFY check.
	txSequence := int64(vm.tx.TxIn[vm.txIdx].Sequence)
	if txSequence&int64(wire.SequenceLockTimeDisabled) != 0 {

		str := fmt.Sprintf("transaction sequence has sequence "+
			"locktime disabled bit set: 0x%x", txSequence)
		return scriptError(ErrUnsatisfiedLockTime, str)
	}

	// Mask off non-consensus bits before doing comparisons.
	lockTimeMask := int64(wire.SequenceLockTimeIsSeconds |
		wire.SequenceLockTimeMask)
	return verifyLockTime(txSequence&lockTimeMask,
		wire.SequenceLockTimeIsSeconds, sequence&lockTimeMask)
}

// opcodeToAltStack removes the top item from the main data stack and pushes it onto the alternate data stack.
// Main data stack transformation: [... x1 x2 x3] -> [... x1 x2]
// Alt data stack transformation:  [... y1 y2 y3] -> [... y1 y2 y3 x3]
func opcodeToAltStack(
	op *parsedOpcode, vm *Engine) error {

	so, err := vm.dstack.PopByteArray()
	if err != nil {

		return err
	}
	vm.astack.PushByteArray(so)
	return nil
}

// opcodeFromAltStack removes the top item from the alternate data stack and pushes it onto the main data stack.
// Main data stack transformation: [... x1 x2 x3] -> [... x1 x2 x3 y3]
// Alt data stack transformation:  [... y1 y2 y3] -> [... y1 y2]
func opcodeFromAltStack(
	op *parsedOpcode, vm *Engine) error {

	so, err := vm.astack.PopByteArray()
	if err != nil {

		return err
	}
	vm.dstack.PushByteArray(so)
	return nil
}

// opcode2Drop removes the top 2 items from the data stack.
// Stack transformation: [... x1 x2 x3] -> [... x1]
func opcode2Drop(
	op *parsedOpcode, vm *Engine) error {

	return vm.dstack.DropN(2)
}

// opcode2Dup duplicates the top 2 items on the data stack.
// Stack transformation: [... x1 x2 x3] -> [... x1 x2 x3 x2 x3]
func opcode2Dup(
	op *parsedOpcode, vm *Engine) error {

	return vm.dstack.DupN(2)
}

// opcode3Dup duplicates the top 3 items on the data stack.
// Stack transformation: [... x1 x2 x3] -> [... x1 x2 x3 x1 x2 x3]
func opcode3Dup(
	op *parsedOpcode, vm *Engine) error {

	return vm.dstack.DupN(3)
}

// opcode2Over duplicates the 2 items before the top 2 items on the data stack.
// Stack transformation: [... x1 x2 x3 x4] -> [... x1 x2 x3 x4 x1 x2]
func opcode2Over(
	op *parsedOpcode, vm *Engine) error {

	return vm.dstack.OverN(2)
}

// opcode2Rot rotates the top 6 items on the data stack to the left twice.
// Stack transformation: [... x1 x2 x3 x4 x5 x6] -> [... x3 x4 x5 x6 x1 x2]
func opcode2Rot(
	op *parsedOpcode, vm *Engine) error {

	return vm.dstack.RotN(2)
}

// opcode2Swap swaps the top 2 items on the data stack with the 2 that come before them.
// Stack transformation: [... x1 x2 x3 x4] -> [... x3 x4 x1 x2]
func opcode2Swap(
	op *parsedOpcode, vm *Engine) error {

	return vm.dstack.SwapN(2)
}

// opcodeIfDup duplicates the top item of the stack if it is not zero.
// Stack transformation (x1==0): [... x1] -> [... x1]
// Stack transformation (x1!=0): [... x1] -> [... x1 x1]
func opcodeIfDup(
	op *parsedOpcode, vm *Engine) error {

	so, err := vm.dstack.PeekByteArray(0)
	if err != nil {

		return err
	}

	// Push copy of data iff it isn't zero
	if asBool(so) {

		vm.dstack.PushByteArray(so)
	}
	return nil
}

// opcodeDepth pushes the depth of the data stack prior to executing this opcode, encoded as a number, onto the data stack.
// Stack transformation: [...] -> [... <num of items on the stack>]
// Example with 2 items: [x1 x2] -> [x1 x2 2]
// Example with 3 items: [x1 x2 x3] -> [x1 x2 x3 3]
func opcodeDepth(
	op *parsedOpcode, vm *Engine) error {

	vm.dstack.PushInt(scriptNum(vm.dstack.Depth()))
	return nil
}

// opcodeDrop removes the top item from the data stack.
// Stack transformation: [... x1 x2 x3] -> [... x1 x2]
func opcodeDrop(
	op *parsedOpcode, vm *Engine) error {

	return vm.dstack.DropN(1)
}

// opcodeDup duplicates the top item on the data stack.
// Stack transformation: [... x1 x2 x3] -> [... x1 x2 x3 x3]
func opcodeDup(
	op *parsedOpcode, vm *Engine) error {

	return vm.dstack.DupN(1)
}

// opcodeNip removes the item before the top item on the data stack.
// Stack transformation: [... x1 x2 x3] -> [... x1 x3]
func opcodeNip(
	op *parsedOpcode, vm *Engine) error {

	return vm.dstack.NipN(1)
}

// opcodeOver duplicates the item before the top item on the data stack.
// Stack transformation: [... x1 x2 x3] -> [... x1 x2 x3 x2]
func opcodeOver(
	op *parsedOpcode, vm *Engine) error {

	return vm.dstack.OverN(1)
}

// opcodePick treats the top item on the data stack as an integer and duplicates the item on the stack that number of items back to the top.
// Stack transformation: [xn ... x2 x1 x0 n] -> [xn ... x2 x1 x0 xn]
// Example with n=1: [x2 x1 x0 1] -> [x2 x1 x0 x1]
// Example with n=2: [x2 x1 x0 2] -> [x2 x1 x0 x2]
func opcodePick(
	op *parsedOpcode, vm *Engine) error {

	val, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	return vm.dstack.PickN(val.Int32())
}

// opcodeRoll treats the top item on the data stack as an integer and moves the item on the stack that number of items back to the top.
// Stack transformation: [xn ... x2 x1 x0 n] -> [... x2 x1 x0 xn]
// Example with n=1: [x2 x1 x0 1] -> [x2 x0 x1]
// Example with n=2: [x2 x1 x0 2] -> [x1 x0 x2]
func opcodeRoll(
	op *parsedOpcode, vm *Engine) error {

	val, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	return vm.dstack.RollN(val.Int32())
}

// opcodeRot rotates the top 3 items on the data stack to the left.
// Stack transformation: [... x1 x2 x3] -> [... x2 x3 x1]
func opcodeRot(
	op *parsedOpcode, vm *Engine) error {

	return vm.dstack.RotN(1)
}

// opcodeSwap swaps the top two items on the stack.
// Stack transformation: [... x1 x2] -> [... x2 x1]
func opcodeSwap(
	op *parsedOpcode, vm *Engine) error {

	return vm.dstack.SwapN(1)
}

// opcodeTuck inserts a duplicate of the top item of the data stack before the second-to-top item.
// Stack transformation: [... x1 x2] -> [... x2 x1 x2]
func opcodeTuck(
	op *parsedOpcode, vm *Engine) error {

	return vm.dstack.Tuck()
}

// opcodeSize pushes the size of the top item of the data stack onto the data stack.
// Stack transformation: [... x1] -> [... x1 len(x1)]
func opcodeSize(
	op *parsedOpcode, vm *Engine) error {

	so, err := vm.dstack.PeekByteArray(0)
	if err != nil {

		return err
	}
	vm.dstack.PushInt(scriptNum(len(so)))
	return nil
}

// opcodeEqual removes the top 2 items of the data stack, compares them as raw bytes, and pushes the result, encoded as a boolean, back to the stack.
// Stack transformation: [... x1 x2] -> [... bool]
func opcodeEqual(
	op *parsedOpcode, vm *Engine) error {

	a, err := vm.dstack.PopByteArray()
	if err != nil {

		return err
	}
	b, err := vm.dstack.PopByteArray()
	if err != nil {

		return err
	}
	vm.dstack.PushBool(bytes.Equal(a, b))
	return nil
}

// opcodeEqualVerify is a combination of opcodeEqual and opcodeVerify. Specifically, it removes the top 2 items of the data stack, compares them,
// and pushes the result, encoded as a boolean, back to the stack.  Then, it examines the top item on the data stack as a boolean value and verifies it evaluates to true.  An error is returned if it does not.
// Stack transformation: [... x1 x2] -> [... bool] -> [...]
func opcodeEqualVerify(
	op *parsedOpcode, vm *Engine) error {

	err := opcodeEqual(op, vm)
	if err == nil {

		err = abstractVerify(op, vm, ErrEqualVerify)
	}
	return err
}

// opcode1Add treats the top item on the data stack as an integer and replaces it with its incremented value (plus 1).
// Stack transformation: [... x1 x2] -> [... x1 x2+1]
func opcode1Add(
	op *parsedOpcode, vm *Engine) error {

	m, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	vm.dstack.PushInt(m + 1)
	return nil
}

// opcode1Sub treats the top item on the data stack as an integer and replaces it with its decremented value (minus 1).
// Stack transformation: [... x1 x2] -> [... x1 x2-1]
func opcode1Sub(
	op *parsedOpcode, vm *Engine) error {

	m, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	vm.dstack.PushInt(m - 1)
	return nil
}

// opcodeNegate treats the top item on the data stack as an integer and replaces it with its negation.
// Stack transformation: [... x1 x2] -> [... x1 -x2]
func opcodeNegate(
	op *parsedOpcode, vm *Engine) error {

	m, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	vm.dstack.PushInt(-m)
	return nil
}

// opcodeAbs treats the top item on the data stack as an integer and replaces it it with its absolute value.
// Stack transformation: [... x1 x2] -> [... x1 abs(x2)]
func opcodeAbs(
	op *parsedOpcode, vm *Engine) error {

	m, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	if m < 0 {

		m = -m
	}
	vm.dstack.PushInt(m)
	return nil
}

// opcodeNot treats the top item on the data stack as an integer and replaces it with its "inverted" value (0 becomes 1, non-zero becomes 0). NOTE: While it would probably make more sense to treat the top item as a boolean, and push the opposite, which is really what the intention of this opcode is, it is extremely important that is not done because integers are interpreted differently than booleans and the consensus rules for this opcode dictate the item is interpreted as an integer.
// Stack transformation (x2==0): [... x1 0] -> [... x1 1]
// Stack transformation (x2!=0): [... x1 1] -> [... x1 0]
// Stack transformation (x2!=0): [... x1 17] -> [... x1 0]
func opcodeNot(
	op *parsedOpcode, vm *Engine) error {

	m, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	if m == 0 {

		vm.dstack.PushInt(scriptNum(1))
	} else {

		vm.dstack.PushInt(scriptNum(0))
	}
	return nil
}

// opcode0NotEqual treats the top item on the data stack as an integer and replaces it with either a 0 if it is zero, or a 1 if it is not zero.
// Stack transformation (x2==0): [... x1 0] -> [... x1 0]
// Stack transformation (x2!=0): [... x1 1] -> [... x1 1]
// Stack transformation (x2!=0): [... x1 17] -> [... x1 1]
func opcode0NotEqual(
	op *parsedOpcode, vm *Engine) error {

	m, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	if m != 0 {

		m = 1
	}
	vm.dstack.PushInt(m)
	return nil
}

// opcodeAdd treats the top two items on the data stack as integers and replaces them with their sum.
// Stack transformation: [... x1 x2] -> [... x1+x2]
func opcodeAdd(
	op *parsedOpcode, vm *Engine) error {

	v0, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	v1, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	vm.dstack.PushInt(v0 + v1)
	return nil
}

// opcodeSub treats the top two items on the data stack as integers and replaces them with the result of subtracting the top entry from the second-to-top entry.
// Stack transformation: [... x1 x2] -> [... x1-x2]
func opcodeSub(
	op *parsedOpcode, vm *Engine) error {

	v0, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	v1, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	vm.dstack.PushInt(v1 - v0)
	return nil
}

// opcodeBoolAnd treats the top two items on the data stack as integers.  When both of them are not zero, they are replaced with a 1, otherwise a 0.
// Stack transformation (x1==0, x2==0): [... 0 0] -> [... 0]
// Stack transformation (x1!=0, x2==0): [... 5 0] -> [... 0]
// Stack transformation (x1==0, x2!=0): [... 0 7] -> [... 0]
// Stack transformation (x1!=0, x2!=0): [... 4 8] -> [... 1]
func opcodeBoolAnd(
	op *parsedOpcode, vm *Engine) error {

	v0, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	v1, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	if v0 != 0 && v1 != 0 {

		vm.dstack.PushInt(scriptNum(1))
	} else {

		vm.dstack.PushInt(scriptNum(0))
	}
	return nil
}

// opcodeBoolOr treats the top two items on the data stack as integers.  When either of them are not zero, they are replaced with a 1, otherwise a 0.
// Stack transformation (x1==0, x2==0): [... 0 0] -> [... 0]
// Stack transformation (x1!=0, x2==0): [... 5 0] -> [... 1]
// Stack transformation (x1==0, x2!=0): [... 0 7] -> [... 1]
// Stack transformation (x1!=0, x2!=0): [... 4 8] -> [... 1]
func opcodeBoolOr(
	op *parsedOpcode, vm *Engine) error {

	v0, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	v1, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	if v0 != 0 || v1 != 0 {

		vm.dstack.PushInt(scriptNum(1))
	} else {

		vm.dstack.PushInt(scriptNum(0))
	}
	return nil
}

// opcodeNumEqual treats the top two items on the data stack as integers.  When they are equal, they are replaced with a 1, otherwise a 0.
// Stack transformation (x1==x2): [... 5 5] -> [... 1]
// Stack transformation (x1!=x2): [... 5 7] -> [... 0]
func opcodeNumEqual(
	op *parsedOpcode, vm *Engine) error {

	v0, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	v1, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	if v0 == v1 {

		vm.dstack.PushInt(scriptNum(1))
	} else {

		vm.dstack.PushInt(scriptNum(0))
	}
	return nil
}

// opcodeNumEqualVerify is a combination of opcodeNumEqual and opcodeVerify. Specifically, treats the top two items on the data stack as integers.  When they are equal, they are replaced with a 1, otherwise a 0.  Then, it examines the top item on the data stack as a boolean value and verifies it evaluates to true.  An error is returned if it does not.
// Stack transformation: [... x1 x2] -> [... bool] -> [...]
func opcodeNumEqualVerify(
	op *parsedOpcode, vm *Engine) error {

	err := opcodeNumEqual(op, vm)
	if err == nil {

		err = abstractVerify(op, vm, ErrNumEqualVerify)
	}
	return err
}

// opcodeNumNotEqual treats the top two items on the data stack as integers. When they are NOT equal, they are replaced with a 1, otherwise a 0.
// Stack transformation (x1==x2): [... 5 5] -> [... 0]
// Stack transformation (x1!=x2): [... 5 7] -> [... 1]
func opcodeNumNotEqual(
	op *parsedOpcode, vm *Engine) error {

	v0, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	v1, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	if v0 != v1 {

		vm.dstack.PushInt(scriptNum(1))
	} else {

		vm.dstack.PushInt(scriptNum(0))
	}
	return nil
}

// opcodeLessThan treats the top two items on the data stack as integers.  When the second-to-top item is less than the top item, they are replaced with a 1, otherwise a 0.
// Stack transformation: [... x1 x2] -> [... bool]
func opcodeLessThan(
	op *parsedOpcode, vm *Engine) error {

	v0, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	v1, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	if v1 < v0 {

		vm.dstack.PushInt(scriptNum(1))
	} else {

		vm.dstack.PushInt(scriptNum(0))
	}
	return nil
}

// opcodeGreaterThan treats the top two items on the data stack as integers. When the second-to-top item is greater than the top item, they are replaced with a 1, otherwise a 0.
// Stack transformation: [... x1 x2] -> [... bool]
func opcodeGreaterThan(
	op *parsedOpcode, vm *Engine) error {

	v0, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	v1, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	if v1 > v0 {

		vm.dstack.PushInt(scriptNum(1))
	} else {

		vm.dstack.PushInt(scriptNum(0))
	}
	return nil
}

// opcodeLessThanOrEqual treats the top two items on the data stack as integers. When the second-to-top item is less than or equal to the top item, they are replaced with a 1, otherwise a 0.
// Stack transformation: [... x1 x2] -> [... bool]
func opcodeLessThanOrEqual(
	op *parsedOpcode, vm *Engine) error {

	v0, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	v1, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	if v1 <= v0 {

		vm.dstack.PushInt(scriptNum(1))
	} else {

		vm.dstack.PushInt(scriptNum(0))
	}
	return nil
}

// opcodeGreaterThanOrEqual treats the top two items on the data stack as integers.  When the second-to-top item is greater than or equal to the top item, they are replaced with a 1, otherwise a 0.
// Stack transformation: [... x1 x2] -> [... bool]
func opcodeGreaterThanOrEqual(
	op *parsedOpcode, vm *Engine) error {

	v0, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	v1, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	if v1 >= v0 {

		vm.dstack.PushInt(scriptNum(1))
	} else {

		vm.dstack.PushInt(scriptNum(0))
	}
	return nil
}

// opcodeMin treats the top two items on the data stack as integers and replaces with the minimum of the two.
// Stack transformation: [... x1 x2] -> [... min(x1, x2)]
func opcodeMin(
	op *parsedOpcode, vm *Engine) error {

	v0, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	v1, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	if v1 < v0 {

		vm.dstack.PushInt(v1)
	} else {

		vm.dstack.PushInt(v0)
	}
	return nil
}

// opcodeMax treats the top two items on the data stack as integers and replaces them with the maximum of the two.
// Stack transformation: [... x1 x2] -> [... max(x1, x2)]
func opcodeMax(
	op *parsedOpcode, vm *Engine) error {

	v0, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	v1, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	if v1 > v0 {

		vm.dstack.PushInt(v1)
	} else {

		vm.dstack.PushInt(v0)
	}
	return nil
}

// opcodeWithin treats the top 3 items on the data stack as integers.  When the value to test is within the specified range (left inclusive), they are replaced with a 1, otherwise a 0. The top item is the max value, the second-top-item is the minimum value, and the third-to-top item is the value to test.
// Stack transformation: [... x1 min max] -> [... bool]
func opcodeWithin(
	op *parsedOpcode, vm *Engine) error {

	maxVal, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	minVal, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	x, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	if x >= minVal && x < maxVal {

		vm.dstack.PushInt(scriptNum(1))
	} else {

		vm.dstack.PushInt(scriptNum(0))
	}
	return nil
}

// calcHash calculates the hash of hasher over buf.
func calcHash(
	buf []byte, hasher hash.Hash) []byte {

	hasher.Write(buf)
	return hasher.Sum(nil)
}

// opcodeRipemd160 treats the top item of the data stack as raw bytes and replaces it with ripemd160(data).
// Stack transformation: [... x1] -> [... ripemd160(x1)]
func opcodeRipemd160(
	op *parsedOpcode, vm *Engine) error {

	buf, err := vm.dstack.PopByteArray()
	if err != nil {

		return err
	}
	vm.dstack.PushByteArray(calcHash(buf, ripemd160.New()))
	return nil
}

// opcodeSha1 treats the top item of the data stack as raw bytes and replaces it with sha1(data).
// Stack transformation: [... x1] -> [... sha1(x1)]
func opcodeSha1(
	op *parsedOpcode, vm *Engine) error {

	buf, err := vm.dstack.PopByteArray()
	if err != nil {

		return err
	}
	hash := sha1.Sum(buf)
	vm.dstack.PushByteArray(hash[:])
	return nil
}

// opcodeSha256 treats the top item of the data stack as raw bytes and replaces it with sha256(data).
// Stack transformation: [... x1] -> [... sha256(x1)]
func opcodeSha256(
	op *parsedOpcode, vm *Engine) error {

	buf, err := vm.dstack.PopByteArray()
	if err != nil {

		return err
	}
	hash := sha256.Sum256(buf)
	vm.dstack.PushByteArray(hash[:])
	return nil
}

// opcodeHash160 treats the top item of the data stack as raw bytes and replaces it with ripemd160(sha256(data)).
// Stack transformation: [... x1] -> [... ripemd160(sha256(x1))]
func opcodeHash160(
	op *parsedOpcode, vm *Engine) error {

	buf, err := vm.dstack.PopByteArray()
	if err != nil {

		return err
	}
	hash := sha256.Sum256(buf)
	vm.dstack.PushByteArray(calcHash(hash[:], ripemd160.New()))
	return nil
}

// opcodeHash256 treats the top item of the data stack as raw bytes and replaces it with sha256(sha256(data)).
// Stack transformation: [... x1] -> [... sha256(sha256(x1))]
func opcodeHash256(
	op *parsedOpcode, vm *Engine) error {

	buf, err := vm.dstack.PopByteArray()
	if err != nil {

		return err
	}
	vm.dstack.PushByteArray(chainhash.DoubleHashB(buf))
	return nil
}

// opcodeCodeSeparator stores the current script offset as the most recently seen OpCodeSeparator which is used during signature checking. This opcode does not change the contents of the data stack.
func opcodeCodeSeparator(
	op *parsedOpcode, vm *Engine) error {

	vm.lastCodeSep = int(vm.scriptOff.Load())
	return nil
}

// opcodeCheckSig treats the top 2 items on the stack as a public key and a signature and replaces them with a bool which indicates if the signature was successfully verified.
// The process of verifying a signature requires calculating a signature hash in the same way the transaction signer did.  It involves hashing portions of the transaction based on the hash type byte (which is the final byte of the signature) and the portion of the script starting from the most recent OpCodeSeparator (or the beginning of the script if there are none) to the end of the script (with any other OP_CODESEPARATORs removed).  Once this "script hash" is calculated, the signature is checked using standard cryptographic methods against the provided public key.
// Stack transformation: [... signature pubkey] -> [... bool]
func opcodeCheckSig(
	op *parsedOpcode, vm *Engine) error {

	pkBytes, err := vm.dstack.PopByteArray()
	if err != nil {

		return err
	}
	fullSigBytes, err := vm.dstack.PopByteArray()
	if err != nil {

		return err
	}

	// The signature actually needs needs to be longer than this, but at 1 byte is needed for the hash type below.  The full length is checked depending on the script flags and upon parsing the signature.
	if len(fullSigBytes) < 1 {

		vm.dstack.PushBool(false)
		return nil
	}

	// Trim off hashtype from the signature string and check if the signature and pubkey conform to the strict encoding requirements depending on the flags.

	// NOTE: When the strict encoding flags are set, any errors in the signature or public encoding here result in an immediate script error (and thus no result bool is pushed to the data stack).  This differs from the logic below where any errors in parsing the signature is treated as the signature failure resulting in false being pushed to the data stack.  This is required because the more general script validation consensus rules do not have the new strict encoding requirements enabled by the flags.
	hashType := SigHashType(fullSigBytes[len(fullSigBytes)-1])
	sigBytes := fullSigBytes[:len(fullSigBytes)-1]
	if err := vm.checkHashTypeEncoding(hashType); err != nil {

		return err
	}
	if err := vm.checkSignatureEncoding(sigBytes); err != nil {

		return err
	}
	if err := vm.checkPubKeyEncoding(pkBytes); err != nil {

		return err
	}

	// Get script starting from the most recent OpCodeSeparator.
	subScript := vm.subScript()

	// Generate the signature hash based on the signature hash type.
	var hash []byte
	if vm.isWitnessVersionActive(0) {

		var sigHashes *TxSigHashes

		if vm.hashCache != nil {

			sigHashes = vm.hashCache
		} else {

			sigHashes = NewTxSigHashes(&vm.tx)
		}
		hash, err = calcWitnessSignatureHash(subScript, sigHashes, hashType,
			&vm.tx, vm.txIdx, vm.inputAmount)

		if err != nil {

			return err
		}
	} else {

		// Remove the signature since there is no way for a signature to sign itself.
		subScript = removeOpcodeByData(subScript, fullSigBytes)
		hash = calcSignatureHash(subScript, hashType, &vm.tx, vm.txIdx)
	}
	pubKey, err := ec.ParsePubKey(pkBytes, ec.S256())
	if err != nil {

		vm.dstack.PushBool(false)
		return nil
	}
	var signature *ec.Signature
	if vm.hasFlag(ScriptVerifyStrictEncoding) ||
		vm.hasFlag(ScriptVerifyDERSignatures) {

		signature, err = ec.ParseDERSignature(sigBytes, ec.S256())
	} else {

		signature, err = ec.ParseSignature(sigBytes, ec.S256())
	}
	if err != nil {

		vm.dstack.PushBool(false)
		return nil
	}
	var valid bool
	if vm.sigCache != nil {

		var sigHash chainhash.Hash
		copy(sigHash[:], hash)
		valid = vm.sigCache.Exists(sigHash, signature, pubKey)

		if !valid && signature.Verify(hash, pubKey) {

			vm.sigCache.Add(sigHash, signature, pubKey)
			valid = true
		}
	} else {

		valid = signature.Verify(hash, pubKey)
	}
	if !valid && vm.hasFlag(ScriptVerifyNullFail) && len(sigBytes) > 0 {

		str := "signature not empty on failed checksig"
		return scriptError(ErrNullFail, str)
	}
	vm.dstack.PushBool(valid)
	return nil
}

// opcodeCheckSigVerify is a combination of opcodeCheckSig and opcodeVerify. The opcodeCheckSig function is invoked followed by opcodeVerify.  See the documentation for each of those opcodes for more details.
// Stack transformation: signature pubkey] -> [... bool] -> [...]
func opcodeCheckSigVerify(
	op *parsedOpcode, vm *Engine) error {

	err := opcodeCheckSig(op, vm)
	if err == nil {

		err = abstractVerify(op, vm, ErrCheckSigVerify)
	}
	return err
}

// parsedSigInfo houses a raw signature along with its parsed form and a flag for whether or not it has already been parsed.  It is used to prevent parsing the same signature multiple times when verifying a multisig.

type parsedSigInfo struct {
	signature       []byte
	parsedSignature *ec.Signature
	parsed          bool
}

// opcodeCheckMultiSig treats the top item on the stack as an integer number of public keys, followed by that many entries as raw data representing the public keys, followed by the integer number of signatures, followed by that many entries as raw data representing the signatures. Due to a bug in the original Satoshi client implementation, an additional dummy argument is also required by the consensus rules, although it is not used.  The dummy value SHOULD be an OpZero, although that is not required by the consensus rules.  When the ScriptStrictMultiSig flag is set, it must be OpZero.
// All of the aforementioned stack items are replaced with a bool which indicates if the requisite number of signatures were successfully verified. See the opcodeCheckSigVerify documentation for more details about the process for verifying each signature.
// Stack transformation:
// [... dummy [sig ...] numsigs [pubkey ...] numpubkeys] -> [... bool]
func opcodeCheckMultiSig(
	op *parsedOpcode, vm *Engine) error {

	numKeys, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	numPubKeys := int(numKeys.Int32())
	if numPubKeys < 0 {

		str := fmt.Sprintf("number of pubkeys %d is negative",
			numPubKeys)
		return scriptError(ErrInvalidPubKeyCount, str)
	}
	if numPubKeys > MaxPubKeysPerMultiSig {

		str := fmt.Sprintf("too many pubkeys: %d > %d",
			numPubKeys, MaxPubKeysPerMultiSig)
		return scriptError(ErrInvalidPubKeyCount, str)
	}
	vm.numOps += numPubKeys
	if vm.numOps > MaxOpsPerScript {

		str := fmt.Sprintf("exceeded max operation limit of %d",
			MaxOpsPerScript)
		return scriptError(ErrTooManyOperations, str)
	}
	pubKeys := make([][]byte, 0, numPubKeys)

	for i := 0; i < numPubKeys; i++ {

		pubKey, err := vm.dstack.PopByteArray()

		if err != nil {

			return err
		}
		pubKeys = append(pubKeys, pubKey)
	}
	numSigs, err := vm.dstack.PopInt()
	if err != nil {

		return err
	}
	numSignatures := int(numSigs.Int32())
	if numSignatures < 0 {

		str := fmt.Sprintf("number of signatures %d is negative",
			numSignatures)
		return scriptError(ErrInvalidSignatureCount, str)
	}
	if numSignatures > numPubKeys {

		str := fmt.Sprintf("more signatures than pubkeys: %d > %d",
			numSignatures, numPubKeys)
		return scriptError(ErrInvalidSignatureCount, str)
	}
	signatures := make([]*parsedSigInfo, 0, numSignatures)

	for i := 0; i < numSignatures; i++ {

		signature, err := vm.dstack.PopByteArray()

		if err != nil {

			return err
		}
		sigInfo := &parsedSigInfo{signature: signature}
		signatures = append(signatures, sigInfo)
	}

	// A bug in the original Satoshi client implementation means one more stack value than should be used must be popped.  Unfortunately, this buggy behavior is now part of the consensus and a hard fork would be required to fix it.
	dummy, err := vm.dstack.PopByteArray()
	if err != nil {

		return err
	}

	// Since the dummy argument is otherwise not checked, it could be any value which unfortunately provides a source of malleability.  Thus, there is a script flag to force an error when the value is NOT 0.
	if vm.hasFlag(ScriptStrictMultiSig) && len(dummy) != 0 {

		str := fmt.Sprintf("multisig dummy argument has length %d "+
			"instead of 0", len(dummy))
		return scriptError(ErrSigNullDummy, str)
	}

	// Get script starting from the most recent OpCodeSeparator.
	script := vm.subScript()

	// Remove the signature in pre version 0 segwit scripts since there is no way for a signature to sign itself.
	if !vm.isWitnessVersionActive(0) {

		for _, sigInfo := range signatures {

			script = removeOpcodeByData(script, sigInfo.signature)
		}
	}
	success := true
	numPubKeys++
	pubKeyIdx := -1
	signatureIdx := 0

	for numSignatures > 0 {

		// When there are more signatures than public keys remaining, there is no way to succeed since too many signatures are invalid, so exit early.
		pubKeyIdx++
		numPubKeys--

		if numSignatures > numPubKeys {

			success = false
			break
		}
		sigInfo := signatures[signatureIdx]
		pubKey := pubKeys[pubKeyIdx]
		// The order of the signature and public key evaluation is important here since it can be distinguished by an OpCheckMultiSig NOT when the strict encoding flag is set.
		rawSig := sigInfo.signature

		if len(rawSig) == 0 {

			// Skip to the next pubkey if signature is empty.
			continue
		}
		// Split the signature into hash type and signature components.
		hashType := SigHashType(rawSig[len(rawSig)-1])
		signature := rawSig[:len(rawSig)-1]
		// Only parse and check the signature encoding once.
		var parsedSig *ec.Signature

		if !sigInfo.parsed {

			if err := vm.checkHashTypeEncoding(hashType); err != nil {

				return err
			}

			if err := vm.checkSignatureEncoding(signature); err != nil {

				return err
			}
			// Parse the signature.
			var err error

			if vm.hasFlag(ScriptVerifyStrictEncoding) ||

				vm.hasFlag(ScriptVerifyDERSignatures) {

				parsedSig, err = ec.ParseDERSignature(signature,
					ec.S256())
			} else {

				parsedSig, err = ec.ParseSignature(signature,
					ec.S256())
			}
			sigInfo.parsed = true

			if err != nil {

				continue
			}
			sigInfo.parsedSignature = parsedSig
		} else {

			// Skip to the next pubkey if the signature is invalid.

			if sigInfo.parsedSignature == nil {

				continue
			}
			// Use the already parsed signature.
			parsedSig = sigInfo.parsedSignature
		}

		if err := vm.checkPubKeyEncoding(pubKey); err != nil {

			return err
		}
		// Parse the pubkey.
		parsedPubKey, err := ec.ParsePubKey(pubKey, ec.S256())

		if err != nil {

			continue
		}
		// Generate the signature hash based on the signature hash type.
		var hash []byte

		if vm.isWitnessVersionActive(0) {

			var sigHashes *TxSigHashes

			if vm.hashCache != nil {

				sigHashes = vm.hashCache
			} else {

				sigHashes = NewTxSigHashes(&vm.tx)
			}
			hash, err = calcWitnessSignatureHash(script, sigHashes, hashType,
				&vm.tx, vm.txIdx, vm.inputAmount)

			if err != nil {

				return err
			}
		} else {

			hash = calcSignatureHash(script, hashType, &vm.tx, vm.txIdx)
		}
		var valid bool

		if vm.sigCache != nil {

			var sigHash chainhash.Hash
			copy(sigHash[:], hash)
			valid = vm.sigCache.Exists(sigHash, parsedSig, parsedPubKey)

			if !valid && parsedSig.Verify(hash, parsedPubKey) {

				vm.sigCache.Add(sigHash, parsedSig, parsedPubKey)
				valid = true
			}
		} else {

			valid = parsedSig.Verify(hash, parsedPubKey)
		}

		if valid {

			// PubKey verified, move on to the next signature.
			signatureIdx++
			numSignatures--
		}
	}
	if !success && vm.hasFlag(ScriptVerifyNullFail) {

		for _, sig := range signatures {

			if len(sig.signature) > 0 {

				str := "not all signatures empty on failed checkmultisig"
				return scriptError(ErrNullFail, str)
			}
		}
	}
	vm.dstack.PushBool(success)
	return nil
}

// opcodeCheckMultiSigVerify is a combination of opcodeCheckMultiSig and opcodeVerify.  The opcodeCheckMultiSig is invoked followed by opcodeVerify. See the documentation for each of those opcodes for more details.
// Stack transformation:
// [... dummy [sig ...] numsigs [pubkey ...] numpubkeys] -> [... bool] -> [...]
func opcodeCheckMultiSigVerify(
	op *parsedOpcode, vm *Engine) error {

	err := opcodeCheckMultiSig(op, vm)
	if err == nil {

		err = abstractVerify(op, vm, ErrCheckMultiSigVerify)
	}
	return err
}

// OpcodeByName is a map that can be used to lookup an opcode by its human-readable name (OpCheckMultiSig, OpCheckSig, etc).
var OpcodeByName = make(map[string]byte)

func init() {

	// Initialize the opcode name to value map using the contents of the opcode array.  Also add entries for "OpFalse", "OpTrue", and "OpNoOp2" since they are aliases for "OpZero", "Op1", and "OpCheckLockTimeVerify" respectively.

	for _, op := range opcodeArray {

		OpcodeByName[op.name] = op.value
	}
	OpcodeByName["OpFalse"] = OpFalse
	OpcodeByName["OpTrue"] = OpTrue
	OpcodeByName["OpNoOp2"] = OpCheckLockTimeVerify
	OpcodeByName["OpNoOp3"] = OpCheckSequenceVerify
}
