require 'benchmark'
require 'csv'
require 'pry'

num_pats = %i(100000 10000 1000 100 10)
num_text_lines = %i(10000 1000 100 10)

class Output
  def initialize(biblio_data, grep_data, ripgrep_data, lines, pats)
    @biblio_data = biblio_data
    @grep_data = grep_data
    @ripgrep_data = ripgrep_data
    @lines = lines
    @pats = pats
  end

  def lines
    @lines
  end

  def pats
    @pats
  end

  def biblio_a
    [ sprintf("%0.05f", @biblio_data.real) ]
  end

  def grep_a
    [ sprintf("%0.05f", @grep_data.real) ]
  end

  def ripgrep_a
    [ sprintf("%0.05f", @ripgrep_data.real) ]
  end
end

outputs = []

Benchmark.bm do |x|
  num_text_lines.each do |lines|
    num_pats.each do |pats|
      grep_output = x.report("grep -F") { `grep -F -f pats-#{pats}.txt text-#{lines}.txt -c` }
      ripgrep_output = x.report("ripgrep -F") { `rg -F -f pats-#{pats}.txt text-#{lines}.txt -c` }
      biblio_output = x.report("biblio") { `./cmd -patternsFile pats-#{pats}.txt -file text-#{lines}.txt` }

      outputs << Output.new(biblio_output, grep_output, ripgrep_output, lines, pats)
    end
  end
end

headers = ['patterns', 'lines', '', 'biblio', 'grep -F', 'rg -F']

CSV.open(
  'benchmarks.csv',
  'w',
  write_headers: true,
  headers: headers
) do |writer|
  outputs.each do |output|
    writer << [output.pats, output.lines, '', *output.biblio_a, *output.grep_a, *output.ripgrep_a]
  end
end
