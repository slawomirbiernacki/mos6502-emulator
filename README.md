![build](https://github.com/slawomirbiernacki/mos6502-emulator/actions/workflows/go.yml/badge.svg)


# MOS6502-emulator
![image](https://user-images.githubusercontent.com/10660820/213313997-6248858e-d8eb-4333-a951-ff458ad537dd.png)

This repo contains a 6502 microprocessor emulator written in Go.
I've implemented it mostly for the educational purposes, but it's intended to be a fully working emulator.

### Implementation status:
* All opcodes implemented - emulator passes Klaus Dormann's functional tests

TODO:

1. decimal mode
2. IRQ & other Interruptions
3. Ticking speed

### Useful material I used during the implementation

* Inspiration for implementing this emulator https://www.youtube.com/watch?v=m6l3Elk7-Hg&ab_channel=Computerphile
* 6502 fundamentals explanation https://medium.com/codeburst/an-introduction-to-6502-assembly-and-low-level-programming-7c11fa6b9cb9
* How Do I Write an Emulator? https://atarihq.com/danb/files/emu_vol1.txt
* Main references I used:
  * https://web.archive.org/web/20160406122905/http://homepage.ntlworld.com/cyborgsystems/CS_Main/6502/6502.htm (contains small errors here and there)
  * https://www.pagetable.com/c64ref/6502/?tab=2
  * https://www.nesdev.org/wiki/Status_flags
* Comprehensive tests https://github.com/Klaus2m5/6502_65C02_functional_tests
* Video I found useful to understand how to use the above tests https://www.youtube.com/watch?v=ywN4ABwmldQ
* Web based emulator great for testing expected outcomes of operations https://skilldrick.github.io/easy6502/

Other:
* https://github.com/topics/6502-emulation
* https://llx.com/Neil/a2/opcodes.html
* https://www.middle-engine.com/blog/posts/2020/06/23/programming-the-nes-the-6502-in-detail
* https://archive.org/details/mos_microcomputers_programming_manual/
