from solcx import compile_files
import sys

compiled_sol = compile_files(
    ['./DynamicArray.sol'],
    output_values = ['bin', 'abi']
)

contract = compiled_sol['./DynamicArray.sol:DynamicArrayTest']
# concurrent_queue = compiled_sol['./ConcurrentLibInterface.sol:ConcurrentQueue']

with open('darray.txt', 'w') as f:
    f.write('code = "{}"\n'.format(contract['bin']))
    # f.write('abi = "{}"\n'.format(concurrent_queue['abi']))
