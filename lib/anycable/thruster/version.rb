require_relative "../../thruster/version"

module AnyCable
  module Thruster
    VERSION = "0.1.16"

    # Very basic validation to ensure the versions are in sync
    if ::Thruster::VERSION.split(".").take(2) != VERSION.split(".").take(2)
      raise "Minor or major version mismatch"
    end

    if ::Thruster::VERSION.split(".")[2].to_i > VERSION.split(".")[2].to_i
      raise "Patch version mismatch: must be greater or equal to Thruster's patch version"
    end
  end
end
