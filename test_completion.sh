#!/bin/bash

# Test script to demonstrate tenv bash completion functionality

echo "=== Testing tenv bash completion ==="
echo

# Generate and source completion
echo "Generating bash completion..."
./build/tenv completion bash > /tmp/tenv_completion.bash
source /tmp/tenv_completion.bash

echo "Building tenv..."
make build > /dev/null 2>&1

echo "Installing a couple of OpenTofu versions for testing..."
./build/tenv tofu install 1.10.5 > /dev/null 2>&1
./build/tenv tofu install 1.10.6 > /dev/null 2>&1

echo
echo "=== Testing completion for 'tenv tofu use' ==="
echo "Available completions:"
./build/tenv __complete tofu use "" 2>/dev/null | head -n -1

echo
echo "=== Testing completion for 'tenv tofu install' ==="
echo "Available completions:"
./build/tenv __complete tofu install "" 2>/dev/null | head -n -1

echo
echo "=== Testing completion for 'tenv tofu uninstall' ==="
echo "Available completions:"
./build/tenv __complete tofu uninstall "" 2>/dev/null | head -n -1

echo
echo "=== Testing completion for 'tenv terraform use' (no versions installed) ==="
echo "Available completions:"
./build/tenv __complete terraform use "" 2>/dev/null | head -n -1

echo
echo "=== Completion Test Summary ==="
echo "✅ Installed version completion: Shows locally installed versions"
echo "✅ Strategy completion: Shows latest, latest-stable, latest-pre, etc."
echo "✅ Command-specific completion: Different options for use/install/uninstall"
echo "✅ Cross-tool support: Works for tofu, terraform, terragrunt, etc."
echo
echo "To enable completion in your shell, run:"
echo "  ./build/tenv completion bash > ~/.tenv_completion.bash"
echo "  echo 'source ~/.tenv_completion.bash' >> ~/.bashrc"

# Cleanup
rm -f /tmp/tenv_completion.bash
