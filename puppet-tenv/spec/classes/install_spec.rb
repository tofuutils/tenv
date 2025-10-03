require 'spec_helper'

describe 'tenv::install' do
  on_supported_os.each do |os, os_facts|
    context "on #{os}" do
      let(:facts) { os_facts }
      let(:pre_condition) { 'include tenv' }

      it { is_expected.to compile }

      case os_facts[:os][:family]
      when 'Debian'
        it { is_expected.to contain_exec('download_tenv') }
        it { is_expected.to contain_package('tenv').with_provider('dpkg') }
        it { is_expected.to contain_exec('verify_tenv') }
        
        it 'downloads the correct package' do
          is_expected.to contain_exec('download_tenv').with(
            command: %r{curl -sL.*\.deb},
          )
        end

      when 'RedHat'
        it { is_expected.to contain_exec('download_tenv') }
        it { is_expected.to contain_package('tenv').with_provider('rpm') }
        it { is_expected.to contain_exec('verify_tenv') }
        
        it 'downloads the correct package' do
          is_expected.to contain_exec('download_tenv').with(
            command: %r{curl -sL.*\.rpm},
          )
        end

      when 'Archlinux'
        it { is_expected.to contain_exec('install_tenv_aur') }
        it { is_expected.to contain_exec('verify_tenv') }
      end
    end
  end
end
