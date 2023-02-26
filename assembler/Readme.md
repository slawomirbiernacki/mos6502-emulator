### How to compile a test

1. Install docker
2. Pull `functional-tests` git submodule with `git submodule init` & `git submodule update`. The tests come from this repo https://github.com/amb5l/6502_65C02_functional_tests
3. Build docker image containing the assembler `docker build -t=ca65 .`
4. Run the container mounting the sources `docker run -it -v $(pwd)/functional-tests:/assembler ca65 ./asm_cmp.sh 6502_functional_test`
5. Output files should be created in `functional-tests/ca65` directory.


