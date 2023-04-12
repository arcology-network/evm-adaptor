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

compiled_sol = compile_contracts('./pk_v2')

sale_auction = compiled_sol['./pk_v2/Auction/SaleClockAuction.sol:SaleClockAuction']
siring_auction = compiled_sol['./pk_v2/Auction/SiringClockAuction.sol:SiringClockAuction']
gene_science = compiled_sol['./pk_v2/GeneScience.sol:GeneScience']
kitty_core = compiled_sol['./pk_v2/KittyCore.sol:KittyCore']

with open('pk_v2_code.txt', 'w') as f:
    f.write('sale = "{}"\n'.format(sale_auction['bin']))
    f.write('sire = "{}"\n'.format(siring_auction['bin']))
    f.write('gene = "{}"\n'.format(gene_science['bin']))
    f.write('core = "{}"\n'.format(kitty_core['bin']))