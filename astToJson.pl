#!/usr/bin/perl -w
# astToJson converts the output of ast.Fprint to JSON
use v5.30;
use Data::Dumper;

$/ = '';

my $str = <>;
# substite (len = 2) { ... } sections with [ ... ]
while ($str =~ s/\S++ \(len = \d++\) (\{(?:[^{}]++|(?1))*+\})/'[' . (substr(($1 =~ s,\d++: ,,gr), 1, -1)) . ']'/ge) {};
$str =~ s/\S++ \{/{/g;
# quote keys
$str =~ s/(\w++):/"$1":/g;
$str =~ s/\./ /g;
# remove line numbers
$str =~ s/^.{8}//gm;
# add commas between elements
$str =~ s/^(\s++)([\}\]])(\n\1\S)/$1$2,$3/gm;
print $str;
