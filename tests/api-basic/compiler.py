from solcx import compile_files
import sys
import os

targetPath = os.path.dirname(os.path.realpath(__file__))

compiled_sol = compile_files(
    [targetPath + '/ApiBase.sol' ],
    output_values = ['bin', 'abi']
)

contract = compiled_sol[targetPath + '/ApiBase.sol'  + ':Bytes']

with open(targetPath + '/bytecode.txt', 'w') as f:
    f.write('{}'.format(contract['bin']))

