require 'benchmark'
require 'csv'
require 'pry'

num_pats = %i(10 100 1000 10000 100000)
num_text_lines = %i(10 100 1000 10000 50000)

class Output
  def initialize(biblio_data, grep_data, lines, pats)
    @biblio_data = biblio_data
    @grep_data = grep_data
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
end

outputs = []

Benchmark.bm do |x|
  num_text_lines.each do |lines|
    num_pats.each do |pats|
      grep_output = x.report("grep -F") { `grep -F -f pats-#{pats}.txt text-#{lines}.txt` }
      biblio_output = x.report("biblio") { `./cmd -patternsFile pats-#{pats}.txt -file text-#{lines}.txt` }

      outputs << Output.new(biblio_output, grep_output, lines, pats)
    end
  end
end

headers = ['patterns', 'lines', '', 'biblio', 'fgrep', ]

CSV.open(
  'benchmarks.csv',
  'w',
  write_headers: true,
  headers: headers
) do |writer|
  outputs.each do |output|
    writer << [output.pats, output.lines, '', *output.biblio_a, *output.grep_a]
  end
end
