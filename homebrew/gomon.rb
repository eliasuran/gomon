class Gomon < Formula
  desc "Automatic restarting of go http server when change in file"
  homepage "https://gomon.eliasuran.dev"
  url "https://github.com/eliasuran/gomon/archive/refs/tags/latest.tar.gz"
  sha256 "checksum"

  depends_on "go" => :build

  def install
    ENV["GO111MODULE"] = "auto"
    system "go", "build", "-o", bin/"gomon", "-ldflags", "-s -w"
  end
end
