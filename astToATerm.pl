#!/usr/bin/perl -w
# astToATerm converts the output of ast.Fprint to ATerm
use v5.30;
use Data::Dumper;

$/ = '';

sub typeToATermName {
    $_ = $_[0];
    s/\[\](\w++)/${1}s/;
    return $_;
}

my $str = <>;
$str =~ s/\*ast\.//g;
$str =~ s/\bast\.//g;
# substite (len = 2) { ... } sections with [ ... ]
while ($str =~ s/(\S++) \(len = \d++\) (\{(?:[^{}]++|(?2))*+\})/typeToATermName($1) . '([' . (substr(($2 =~ s,\d++: ,,gr), 1, -1)) . '])'/ge) {};
while ($str =~ s/(\S++) (\{(?:[^{}]++|(?2))*+\})/$1 . '(' . substr($2, 1, -1) . ')'/ge) {}
# # remove keys
$str =~ s/(\w++): //g;
$str =~ s/\./ /g;
# remove line numbers
$str =~ s/^.{8}//gm;
# # add commas between elements
$str =~ s/^(\s++)([\)\]]|\]\))(\n\1\S)/$1$2,$3/gm;
print $str;
