require 'rspec'
require 'yaml'


describe 'Certificate Duration ops-file' do

  # bosh int cfcr.yml + ops-file
  # iterate over /variables in resulting manifest
  # check if certificates have 'duration'
  # 

  it 'modifies all the certificates generated in the cfcr base manifest' do
    manifest_filepath = File.expand_path("../manifests", Dir.pwd)
    bosh_int_output = `bosh interpolate #{manifest_filepath}/cfcr.yml --ops-file \
                          #{manifest_filepath}/ops-files/set-certificate-duration.yml \
                          --var certificate-duration=1460`

    puts bosh_int_output
  end
end
