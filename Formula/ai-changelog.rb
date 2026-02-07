# Homebrew formula for ai-changelog.
# To publish a new version:
# 1. Tag a release: git tag v1.0.0 && git push origin v1.0.0
# 2. Create a GitHub release for that tag (or use the tag alone)
# 3. Update the version and sha256 below. Get sha256 with:
#    curl -sL "https://github.com/lucasbrogni/ai-changelog/archive/refs/tags/v1.0.0.tar.gz" | shasum -a 256
# 4. Copy this file to your homebrew-ai-changelog tap repo and push.

class AiChangelog < Formula
  desc "Generates changelogs from git history using a local LLM via Ollama"
  homepage "https://github.com/lucasbrogni/ai-changelog"
  url "https://github.com/lucasbrogni/ai-changelog/archive/refs/tags/v1.0.0.tar.gz"
  # Replace with actual tarball sha256 after creating release v1.0.0:
  # curl -sL "https://github.com/lucasbrogni/ai-changelog/archive/refs/tags/v1.0.0.tar.gz" | shasum -a 256
  sha256 "0000000000000000000000000000000000000000000000000000000000000000"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"ai-changelog", "."
  end

  test do
    assert_match "Generates changelogs", shell_output("#{bin}/ai-changelog --help")
  end
end
