# Homebrew formula for ai-changelog (lives in this repo).
# To publish a new version: tag a release, then update url/sha256 below.
# Get sha256: curl -sL "https://github.com/brognilucas/ai-changelog/archive/refs/tags/vX.Y.Z.tar.gz" | shasum -a 256

class AiChangelog < Formula
  desc "Generates changelogs from git history using a local LLM via Ollama"
  homepage "https://github.com/brognilucas/ai-changelog"
  url "https://github.com/brognilucas/ai-changelog/archive/refs/tags/v1.0.0.tar.gz"
  # Replace with actual tarball sha256 after creating release v1.0.0:
  # curl -sL "https://github.com/brognilucas/ai-changelog/archive/refs/tags/v1.0.0.tar.gz" | shasum -a 256
  sha256 "b6b1679134281b61a2e4bc142f3f2d06cef3239e2bdb2b5560c9c99bde5ce218"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"ai-changelog", "."
  end

  test do
    assert_match "Generates changelogs", shell_output("#{bin}/ai-changelog --help")
  end
end
