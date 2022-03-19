from solcx import compile_files
import sys

compiled_sol = compile_files(
    ['./NativeStorage.sol'],
    output_values = ['bin', 'abi']
)

defer_perf = compiled_sol['./NativeStorage.sol:NativeStorage']
# concurrent_queue = compiled_sol['./ConcurrentLibInterface.sol:ConcurrentQueue']

with open('code.txt', 'w') as f:
    f.write('code = "{}"\n'.format(defer_perf['bin']))
    # f.write('abi = "{}"\n'.format(concurrent_queue['abi']))
