require_relative 'services_pb'
require_relative 'services_services_pb'

stub = Proto::Jjogaegi::Stub.new('localhost:5000', :this_channel_is_insecure)
resp = stub.run(Proto::RunRequest.new(
    config: Proto::RunConfig.new(
      parser:    "list",
      formatter: "csv",
      options: {
         'debug' => 'true'
      }
    ),
    input: "한국 yyy".b,
))

puts resp.output