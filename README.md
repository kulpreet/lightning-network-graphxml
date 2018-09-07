
# Lightning Network DescribeGraph to GraphXML

Converts lightning describegraph outputs

1. Directed graph with all attributes made available
2. To provide input to Boost Graph Library for analysis


for input in ~/projects/l2/lightning-network-graph-analysis/inputs/* ; do ./lightning-network-graphxml --filename $input; done
