require 'puppetlabs_spec_helper/module_spec_helper'
require 'rspec-puppet-facts'

include RspecPuppetFacts

RSpec.configure do |c|
  c.before :each do
    Puppet.settings[:strict] = :warning
  end
end
