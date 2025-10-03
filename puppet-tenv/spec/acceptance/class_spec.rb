require 'spec_helper_acceptance'

describe 'tenv class' do
  context 'basic installation' do
    it 'works idempotently with no errors' do
      pp = <<-MANIFEST
        class { 'tenv':
          version => 'latest',
        }
      MANIFEST

      # Run it twice and test for idempotency
      apply_manifest(pp, catch_failures: true)
      apply_manifest(pp, catch_changes: true)
    end

    describe package('tenv') do
      it { is_expected.to be_installed }
    end

    describe command('tenv version') do
      its(:exit_status) { is_expected.to eq 0 }
      its(:stdout) { is_expected.to match(/tenv version/) }
    end
  end

  context 'with cosign' do
    it 'installs cosign' do
      pp = <<-MANIFEST
        class { 'tenv':
          install_cosign => true,
        }
      MANIFEST

      apply_manifest(pp, catch_failures: true)
    end

    describe command('cosign version') do
      its(:exit_status) { is_expected.to eq 0 }
    end
  end

  context 'shell configuration' do
    it 'configures bash environment' do
      pp = <<-MANIFEST
        class { 'tenv':
          shell => 'bash',
          auto_install => true,
        }
      MANIFEST

      apply_manifest(pp, catch_failures: true)
    end

    describe file('/root/.bashrc') do
      it { is_expected.to be_file }
      its(:content) { is_expected.to match(/TENV_ROOT/) }
      its(:content) { is_expected.to match(/TENV_AUTO_INSTALL/) }
    end

    describe file('/root/.tenv.completion.bash') do
      it { is_expected.to be_file }
    end
  end

  context 'tool installation' do
    it 'can install terraform' do
      pp = <<-MANIFEST
        class { 'tenv': }
      MANIFEST

      apply_manifest(pp, catch_failures: true)
    end

    describe command('tenv tf install 1.6.0') do
      its(:exit_status) { is_expected.to eq 0 }
    end

    describe command('tenv tf list') do
      its(:exit_status) { is_expected.to eq 0 }
      its(:stdout) { is_expected.to match(/1.6.0/) }
    end
  end
end
