from solcx import compile_files
import sys
import os

targetPath = os.path.dirname(os.path.realpath(__file__))

compiled_sol = compile_files(
    [targetPath + '/NativeStorage.sol' ],
    output_values = ['bin', 'abi']
)

contract = compiled_sol[targetPath + '/NativeStorage.sol' + ':NativeStorage']

with open(targetPath + '/bytecode.txt', 'w') as f:
    f.write('{}'.format(contract['bin']))

