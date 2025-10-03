require 'spec_helper'

describe 'tenv::cosign' do
  on_supported_os.each do |os, os_facts|
    context "on #{os}" do
      let(:facts) { os_facts }
      let(:pre_condition) { 'include tenv' }

      it { is_expected.to compile }

      it { is_expected.to contain_exec('get_cosign_version') }

      case os_facts[:os][:family]
      when 'Debian'
        it { is_expected.to contain_exec('download_cosign') }
        it { is_expected.to contain_package('cosign').with_provider('dpkg') }
        it { is_expected.to contain_file('/tmp/cosign.deb').with_ensure('absent') }

      when 'RedHat'
        it { is_expected.to contain_exec('download_cosign') }
        it { is_expected.to contain_package('cosign').with_provider('rpm') }
        it { is_expected.to contain_file('/tmp/cosign.rpm').with_ensure('absent') }

      when 'Archlinux'
        it { is_expected.to contain_package('cosign') }
      end
    end
  end
end
