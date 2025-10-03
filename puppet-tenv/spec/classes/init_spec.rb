require 'spec_helper'

describe 'tenv' do
  on_supported_os.each do |os, os_facts|
    context "on #{os}" do
      let(:facts) { os_facts }

      it { is_expected.to compile }

      context 'with default parameters' do
        it { is_expected.to contain_class('tenv') }
        it { is_expected.to contain_class('tenv::install') }
        it { is_expected.to contain_class('tenv::config') }
        it { is_expected.to contain_class('tenv::cosign') }
        it { is_expected.to contain_package('curl') }
        it { is_expected.to contain_package('jq') }
        it { is_expected.to contain_package('unzip') }
        it { is_expected.to contain_package('ca-certificates') }
      end

      context 'without cosign' do
        let(:params) { { install_cosign: false } }
        it { is_expected.not_to contain_class('tenv::cosign') }
      end

      context 'without shell configuration' do
        let(:params) { { configure_shell: false } }
        it { is_expected.not_to contain_class('tenv::config') }
      end

      context 'with custom version' do
        let(:params) { { version: 'v2.6.1' } }
        it { is_expected.to compile }
      end

      context 'with multiple users' do
        let(:params) { { users: ['user1', 'user2'] } }
        it { is_expected.to compile }
        it { is_expected.to contain_class('tenv::config') }
      end

      context 'with auto_install enabled' do
        let(:params) { { auto_install: true } }
        it { is_expected.to compile }
      end

      context 'without managing prerequisites' do
        let(:params) { { manage_prerequisites: false } }
        it { is_expected.not_to contain_package('curl') }
      end
    end
  end

  context 'on unsupported OS' do
    let(:facts) do
      {
        os: {
          family: 'Windows'
        }
      }
    end

    it { is_expected.to compile.and_raise_error(/Unsupported operating system/) }
  end
end
