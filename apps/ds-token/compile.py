import os
from solcx import compile_files

def compile_contracts(dir):
    sources = []
    for root, _, files in os.walk(dir):
        for file in files:
            if file.endswith('.sol'):
                sources.append(os.path.join(root, file))

    # print(sources)
    return compile_files(sources, output_values = ['abi', 'bin'])

targetPath = os.path.dirname(os.path.realpath(__file__))
compiled_sol = compile_contracts(targetPath + '/contracts')
dstoken = compiled_sol[targetPath + '/contracts/token.sol:DSToken']

with open('bytecode.txt', 'w') as f:
    f.write('code = "{}"\n'.format(dstoken['bin']))