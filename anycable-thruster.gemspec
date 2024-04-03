require_relative "lib/anycable/thruster/version"

Gem::Specification.new do |s|
  s.name        = "anycable-thruster"
  s.version     = AnyCable::Thruster::VERSION
  s.summary     = "Zero-config HTTP/2 proxy with embedded AnyCable"
  s.description = "A zero-config HTTP/2 proxy for lightweight production deployments with AnyCable real-time server included"
  s.authors     = [ "Kevin McConnell", "Vladimir Dementyev", "Igor Platonov" ]
  s.email       = "anycable@evilmartians.io"
  s.homepage    = "https://github.com/anycable/thruster"
  s.license     = "MIT"

  s.metadata = {
    "homepage_uri" => s.homepage,
    "rubygems_mfa_required" => "true"
  }

  s.files = Dir[ "{lib}/**/*", "MIT-LICENSE", "README.md" ]
  s.bindir = "exe"
  s.executables << "thrust"
end
