from solcx import compile_files
import sys
import os

targetPath = os.path.dirname(os.path.realpath(__file__))

compiled_sol = compile_files(
    [targetPath + '/api_base.sol' ],
    output_values = ['bin', 'abi']
)

contract = compiled_sol[targetPath + '/api_base.sol'  + ':Bytes']

with open(targetPath + '/api_base_bytecode.txt', 'w') as f:
    f.write('{}'.format(contract['bin']))

