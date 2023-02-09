package memory

/*
Rockwell Memory map: https://web.archive.org/web/20160406122905/http://homepage.ntlworld.com/cyborgsystems/CS_Main/6502/6502.htm#BY_PRP
------
|  $0000
|  Zero page
|  $00FF
------
|  $0100
|  Stack (page 1) - push(go up), pull (go down)
|  $01FF
------
|  $0200
|  RAM
|  $03FF
------
|  $0400
|  Memory mapped I/O
|  $7FFF(?)
------
|  $8000
|  ROM
|  $FFF9
 FFFA       - Vector address for NMI (low byte)
 FFFB       - Vector address for NMI (high byte)
 FFFC       - Vector address for RESET (low byte)
 FFFD       - Vector address for RESET (high byte)
 FFFE       - Vector address for IRQ & BRK (low byte)
 FFFF       - Vector address for IRQ & BRK  (high byte)
------
*/

type Memory struct {
	// ($0 - $FFFF in hex)
	Mem [65536]byte
}
