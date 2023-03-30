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

compiled_sol = compile_contracts('./ds_token_v2')

dstoken = compiled_sol['./ds_token_v2/token.sol:DSToken']

with open('ds_token_v2.txt', 'w') as f:
    f.write('code = "{}"\n'.format(dstoken['bin']))