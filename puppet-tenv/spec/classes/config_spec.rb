require 'spec_helper'

describe 'tenv::config' do
  on_supported_os.each do |os, os_facts|
    context "on #{os}" do
      let(:facts) { os_facts }
      let(:pre_condition) { 'include tenv' }

      it { is_expected.to compile }

      context 'with default user (root)' do
        it { is_expected.to contain_file('/root/.tenv').with_ensure('directory') }
        it { is_expected.to contain_file_line('tenv_root_root') }
      end

      context 'with bash shell' do
        let(:pre_condition) do
          "class { 'tenv': shell => 'bash' }"
        end

        it { is_expected.to contain_file_line('tenv_root_root').with_path('/root/.bashrc') }
        it { is_expected.to contain_exec('tenv_completion_root_bash') }
      end

      context 'with zsh shell' do
        let(:pre_condition) do
          "class { 'tenv': shell => 'zsh' }"
        end

        it { is_expected.to contain_file_line('tenv_root_root').with_path('/root/.zshrc') }
        it { is_expected.to contain_exec('tenv_completion_root_zsh') }
      end

      context 'with fish shell' do
        let(:pre_condition) do
          "class { 'tenv': shell => 'fish' }"
        end

        it { is_expected.to contain_file('/root/.config/fish').with_ensure('directory') }
        it { is_expected.to contain_exec('tenv_completion_root_fish') }
      end

      context 'with auto_install enabled' do
        let(:pre_condition) do
          "class { 'tenv': auto_install => true }"
        end

        it { is_expected.to contain_file_line('tenv_auto_install_root') }
      end

      context 'with github_token' do
        let(:pre_condition) do
          "class { 'tenv': github_token => 'ghp_test123' }"
        end

        it { is_expected.to contain_file_line('tenv_github_token_root') }
      end

      context 'with multiple users' do
        let(:pre_condition) do
          "class { 'tenv': users => ['user1', 'user2'] }"
        end

        it { is_expected.to contain_file('/home/user1/.tenv') }
        it { is_expected.to contain_file('/home/user2/.tenv') }
        it { is_expected.to contain_file_line('tenv_root_user1') }
        it { is_expected.to contain_file_line('tenv_root_user2') }
      end

      context 'without completion' do
        let(:pre_condition) do
          "class { 'tenv': setup_completion => false }"
        end

        it { is_expected.not_to contain_exec('tenv_completion_root_bash') }
      end
    end
  end
end
